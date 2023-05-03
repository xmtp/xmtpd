package e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	apiclient "github.com/xmtp/xmtpd/pkg/api/client"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

func (e *E2E) testConvergence(name string) error {
	nodeHosts := make([]string, len(e.opts.APIURLs))
	for i, apiURL := range e.opts.APIURLs {
		url, err := url.Parse(apiURL)
		if err != nil {
			return err
		}
		nodeHosts[i] = url.Hostname()
	}

	// Initialize clients for each node.
	clients := make([]apiclient.Client, len(e.opts.APIURLs))
	for i, apiURL := range e.opts.APIURLs {
		appVersion := "xmtpd-e2e/"
		if len(e.opts.GitCommit) > 0 {
			appVersion += e.opts.GitCommit[:7]
		}
		apiURL, clientOpts, err := parseAPIURL(apiURL)
		if err != nil {
			return err
		}
		clients[i] = apiclient.NewHTTPClient(e.log, apiURL, e.opts.GitCommit, appVersion, clientOpts...)
		defer clients[i].Close()
	}

	topic := "test-" + e.randomStringLower(12)

	ctx := context.WithTimeout(e.ctx, 30*time.Second)
	defer ctx.Close()

	// Subscribe across nodes.
	subscribeStart := time.Now().UTC()
	subs := make([]apiclient.Stream, len(clients))
	for nodeIndex, client := range clients {
		sub, err := client.Subscribe(ctx, &messagev1.SubscribeRequest{
			ContentTopics: []string{
				topic,
			},
		})
		if err != nil {
			if err == context.Canceled {
				e.log.Debug("context canceled", zap.Error(err))
				return nil
			}
			duration := time.Since(subscribeStart)
			e.log.Error("error subscribing", zap.Error(err))
			e.metrics.recordSubscribe(e.ctx, name, nodeHosts[nodeIndex], "failed", duration)
			clients[nodeIndex] = nil
		} else {
			subs[nodeIndex] = sub
			defer sub.Close()
			duration := time.Since(subscribeStart)
			e.metrics.recordSubscribe(e.ctx, name, nodeHosts[nodeIndex], "passed", duration)
		}
	}

	// Publish messages.
	publishStart := time.Now().UTC()
	envs := []*messagev1.Envelope{}
	var envsLock sync.Mutex
	var publishGroup errgroup.Group
	for nodeIndex, client := range clients {
		if client == nil {
			continue
		}
		client := client
		nodeIndex := nodeIndex
		publishGroup.Go(func() error {
			clientEnvs := make([]*messagev1.Envelope, e.opts.MessagePerClient)
			for j := 0; j < e.opts.MessagePerClient; j++ {
				clientEnvs[j] = &messagev1.Envelope{
					ContentTopic: topic,
					TimestampNs:  uint64(j + 1),
					Message:      []byte(fmt.Sprintf("msg%d-%d", nodeIndex+1, j+1)),
				}
			}
			func() {
				envsLock.Lock()
				defer envsLock.Unlock()
				envs = append(envs, clientEnvs...)
			}()
			_, err := client.Publish(ctx, &messagev1.PublishRequest{
				Envelopes: clientEnvs,
			})
			if err != nil {
				duration := time.Since(publishStart)
				e.log.Error("error publishing", zap.Error(err), zap.Duration("duration", time.Since(publishStart)))
				e.metrics.recordPublish(e.ctx, name, nodeHosts[nodeIndex], "failed", duration)
				return nil
			}

			duration := time.Since(publishStart)
			e.log.Info("published", zap.Duration("duration", duration), zap.String("node", nodeHosts[nodeIndex]))
			e.metrics.recordPublish(e.ctx, name, nodeHosts[nodeIndex], "passed", duration)

			return nil
		})
	}
	err := publishGroup.Wait()
	if err != nil {
		return err
	}

	// Expect them to be relayed to each subscription.
	var subscribeGroup errgroup.Group
	for nodeIndex, sub := range subs {
		if sub == nil {
			continue
		}
		sub := sub
		nodeIndex := nodeIndex
		subscribeGroup.Go(func() error {
			envC := make(chan *messagev1.Envelope, 100)
			go func(sub apiclient.Stream) {
				for {
					env, err := sub.Next(ctx)
					if err != nil {
						if isErrClosedConnection(err) || err == context.Canceled {
							break
						}
						e.log.Error("getting next", zap.Error(err))
						break
					}
					if env == nil {
						continue
					}
					envC <- env
				}
			}(sub)
			err := subscribeExpect(envC, envs)
			if err != nil {
				duration := time.Since(publishStart)
				e.log.Error("error checking subscription", zap.Error(err), zap.Duration("duration", time.Since(publishStart)))
				e.metrics.recordSubscribeConvergence(e.ctx, name, nodeHosts[nodeIndex], "failed", duration)
				return nil
			}

			subscribeConvergenceDuration := time.Since(publishStart)
			e.log.Info("subscribe converged", zap.Duration("duration", subscribeConvergenceDuration), zap.String("node", nodeHosts[nodeIndex]))
			e.metrics.recordSubscribeConvergence(e.ctx, name, nodeHosts[nodeIndex], "passed", subscribeConvergenceDuration)

			return nil
		})
	}

	// Expect that they're stored.
	var queryGroup errgroup.Group
	for nodeIndex, client := range clients {
		if client == nil {
			continue
		}
		client := client
		nodeIndex := nodeIndex
		queryGroup.Go(func() error {
			err := expectQueryMessagesEventually(ctx, client, []string{topic}, envs)
			if err != nil {
				duration := time.Since(publishStart)
				e.log.Error("error querying", zap.Error(err), zap.Duration("duration", duration))
				e.metrics.recordQueryConvergence(e.ctx, name, nodeHosts[nodeIndex], "failed", duration)
				return nil
			}

			duration := time.Since(publishStart)
			e.log.Info("query converged", zap.Duration("duration", duration), zap.String("node", nodeHosts[nodeIndex]))
			e.metrics.recordQueryConvergence(e.ctx, name, nodeHosts[nodeIndex], "passed", duration)

			return nil
		})
	}

	err = subscribeGroup.Wait()
	if err != nil {
		return err
	}

	err = queryGroup.Wait()
	if err != nil {
		return err
	}

	return nil
}

func subscribeExpect(envC chan *messagev1.Envelope, envs []*messagev1.Envelope) error {
	receivedEnvs := []*messagev1.Envelope{}
	waitC := time.After(5 * time.Second)
	var done bool
	for !done {
		select {
		case env := <-envC:
			receivedEnvs = append(receivedEnvs, env)
			if len(receivedEnvs) == len(envs) {
				done = true
			}
		case <-waitC:
			done = true
		}
	}
	err := envsDiff(envs, receivedEnvs)
	if err != nil {
		return errors.Wrap(err, "expected subscribe envelopes")
	}
	return nil
}

func isErrClosedConnection(err error) bool {
	return errors.Is(err, io.EOF) || strings.Contains(err.Error(), "closed network connection") || strings.Contains(err.Error(), "response body closed")
}

func expectQueryMessagesEventually(ctx context.Context, client apiclient.Client, contentTopics []string, expectedEnvs []*messagev1.Envelope) error {
	timeout := 10 * time.Second
	delay := 10 * time.Millisecond
	// TODO: higher delay for production?
	started := time.Now()
	for {
		envs, err := query(ctx, client, contentTopics)
		if err != nil {
			return errors.Wrap(err, "querying")
		}
		if len(envs) == len(expectedEnvs) {
			err := envsDiff(envs, expectedEnvs)
			if err != nil {
				return errors.Wrap(err, "expected query envelopes")
			}
			break
		}
		if time.Since(started) > timeout {
			err := envsDiff(envs, expectedEnvs)
			if err != nil {
				return errors.Wrap(err, "expected query envelopes")
			}
			return fmt.Errorf("timeout waiting for query expectation with no diff")
		}
		time.Sleep(delay)
	}
	return nil
}

func query(ctx context.Context, client apiclient.Client, contentTopics []string) ([]*messagev1.Envelope, error) {
	var envs []*messagev1.Envelope
	var pagingInfo *messagev1.PagingInfo
	for {
		res, err := client.Query(ctx, &messagev1.QueryRequest{
			ContentTopics: contentTopics,
			PagingInfo:    pagingInfo,
		})
		if err != nil {
			return nil, err
		}
		envs = append(envs, res.Envelopes...)
		if len(res.Envelopes) == 0 || res.PagingInfo == nil || res.PagingInfo.Cursor == nil {
			break
		}
		pagingInfo = res.PagingInfo
	}
	return envs, nil
}

func envsDiff(a, b []*messagev1.Envelope) error {
	diff := cmp.Diff(a, b,
		cmpopts.SortSlices(func(a, b *messagev1.Envelope) bool {
			if a.ContentTopic != b.ContentTopic {
				return a.ContentTopic < b.ContentTopic
			}
			if a.TimestampNs != b.TimestampNs {
				return a.TimestampNs < b.TimestampNs
			}
			return bytes.Compare(a.Message, b.Message) < 0
		}),
		cmp.Comparer(proto.Equal),
	)
	if diff != "" {
		return fmt.Errorf("expected equal, diff: %s", diff)
	}
	return nil
}

func parseAPIURL(apiURL string) (string, []apiclient.Option, error) {
	// If the API URL is a subdomain of localhost, then replace it with
	// localhost and include a Host header, since Go doesn't resolve this
	// DNS properly.
	opts := []apiclient.Option{}
	url, err := url.Parse(apiURL)
	if err != nil {
		return "", nil, err
	}
	urlParts := strings.Split(url.Hostname(), ".")
	if len(urlParts) == 2 && urlParts[1] == "localhost" {
		opts = append(opts, apiclient.WithHeader("Host", url.Hostname()))
		url.Host = "localhost"
		port := url.Port()
		if port != "80" {
			url.Host += ":" + port
		}
		apiURL = url.String()
	}
	return apiURL, opts, nil
}
