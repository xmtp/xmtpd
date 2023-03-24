package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type httpClient struct {
	log        *zap.Logger
	url        string
	http       *retryablehttp.Client
	version    string
	appVersion string
	headers    map[string]string
}

type option func(*httpClient)

func WithHeader(name, value string) option {
	return func(c *httpClient) {
		c.headers[name] = value
	}
}

const (
	clientVersionHeaderKey = "x-client-version"
	appVersionHeaderKey    = "x-app-version"
)

func NewHTTPClient(log *zap.Logger, serverAddr string, gitCommit string, appVersion string, opts ...option) *httpClient {
	version := "xmtpd/" + shortGitCommit(gitCommit)
	http := retryablehttp.NewClient()
	http.CheckRetry = retryPolicy
	http.Logger = &logger{log}
	client := &httpClient{
		log:        log,
		http:       http,
		url:        serverAddr,
		version:    version,
		appVersion: appVersion,
		headers:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *httpClient) Close() error {
	c.http.HTTPClient.CloseIdleConnections()
	return nil
}

func (c *httpClient) Publish(ctx context.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	res, err := c.rawPublish(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *httpClient) Subscribe(ctx context.Context, req *messagev1.SubscribeRequest) (Stream, error) {
	stream, err := newHTTPStream(c.log, func() (*http.Response, error) {
		return c.post(ctx, "/message/v1/subscribe", req)
	})
	if err != nil {
		return nil, err
	}

	// Wait for subscribe confirmation.
	env, err := stream.Next(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "waiting for subscribe confirmation")
	}
	if !proto.Equal(env, &messagev1.Envelope{}) {
		return nil, fmt.Errorf("invalid subscribe confirmation: %s", env)
	}

	return stream, nil
}

func (c *httpClient) SubscribeAll(ctx context.Context) (Stream, error) {
	stream, err := newHTTPStream(c.log, func() (*http.Response, error) {
		return c.post(ctx, "/message/v1/subscribe-all", &messagev1.SubscribeAllRequest{})
	})
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (c *httpClient) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	res, err := c.rawQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *httpClient) BatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	res, err := c.rawBatchQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *httpClient) rawPublish(ctx context.Context, req *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	var res messagev1.PublishResponse
	resp, err := c.post(ctx, "/message/v1/publish", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	err = protojson.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *httpClient) rawQuery(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	var res messagev1.QueryResponse
	resp, err := c.post(ctx, "/message/v1/query", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	err = protojson.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, string(body))
	}
	return &res, nil
}

func (c *httpClient) rawBatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	var res messagev1.BatchQueryResponse
	resp, err := c.post(ctx, "/message/v1/batch-query", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	err = protojson.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, string(body))
	}
	return &res, nil
}

func (c *httpClient) post(ctx context.Context, path string, req interface{}) (*http.Response, error) {
	var reqJSON []byte
	var err error
	switch req := req.(type) {
	case proto.Message:
		reqJSON, err = protojson.Marshal(req)
		if err != nil {
			return nil, err
		}
	default:
		reqJSON, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}

	url := c.url + path

	post, err := retryablehttp.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	post = post.WithContext(ctx)
	post.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		if strings.ToLower(key) == "host" {
			post.Host = value
			continue
		}
		post.Header.Set(key, value)
	}
	md, _ := metadata.FromOutgoingContext(ctx)
	for key, vals := range md {
		if len(vals) > 0 {
			post.Header.Set(key, vals[0])
		}
	}
	clientVersion := post.Header.Get(clientVersionHeaderKey)
	if clientVersion == "" {
		post.Header.Set(clientVersionHeaderKey, c.version)
	}
	appVersion := post.Header.Get(appVersionHeaderKey)
	if appVersion == "" {
		post.Header.Set(appVersionHeaderKey, c.appVersion)
	}
	resp, err := c.http.Do(post)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if resp == nil {
		return false, nil
	}

	// Avoid conflicting with grpc-gateway max message size error.
	if resp.StatusCode == http.StatusTooManyRequests {
		return false, err
	}

	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

func shortGitCommit(val string) string {
	if len(val) > 0 {
		if len(val) < 7 {
			return val
		}
		return val[:7]
	}
	return ""
}

type logger struct {
	zap *zap.Logger
}

func (l *logger) Error(msg string, keysAndValues ...interface{}) {
	l.zap.Error(msg, zap.AnyFields(keysAndValues...)...)
}

func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	l.zap.Info(msg, zap.AnyFields(keysAndValues...)...)
}

func (l *logger) Debug(msg string, keysAndValues ...interface{}) {
	l.zap.Debug(msg, zap.AnyFields(keysAndValues...)...)
}

func (l *logger) Warn(msg string, keysAndValues ...interface{}) {
	l.zap.Warn(msg, zap.AnyFields(keysAndValues...)...)
}
