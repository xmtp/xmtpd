package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type httpStream struct {
	respC      chan *http.Response
	errC       chan error
	bodyReader *bufio.Reader
	body       io.ReadCloser
	closed     bool
}

func newHTTPStream(log *zap.Logger, reqFn func() (*http.Response, error)) (*httpStream, error) {
	s := &httpStream{
		respC: make(chan *http.Response, 1),
		errC:  make(chan error, 1),
	}

	go func() {
		// Streaming requests block until the first byte is sent, at which
		// point we can consume from the response body reader as a stream.
		resp, err := reqFn()
		if err != nil {
			if !s.closed {
				s.errC <- err
			} else {
				log.Error("requesting", zap.Error(err))
			}
			return
		}

		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				err = fmt.Errorf("%s: %s", resp.Status, string(body))
			}
			if !s.closed {
				s.errC <- err
			} else {
				log.Error("reading body", zap.Error(err))
			}
			return
		}

		if s.closed {
			return
		}
		s.respC <- resp
	}()

	return s, nil
}

func (s *httpStream) reader(ctx context.Context) (*bufio.Reader, error) {
	if s.bodyReader != nil {
		return s.bodyReader, nil
	}
	select {
	case err := <-s.errC:
		return nil, err
	case resp := <-s.respC:
		s.body = resp.Body
		s.bodyReader = bufio.NewReader(s.body)
		return s.bodyReader, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *httpStream) Next(ctx context.Context) (*messagev1.Envelope, error) {
	reader, err := s.reader(ctx)
	if err != nil {
		return nil, err
	}
	if s.body == nil { // stream was closed
		return nil, io.EOF
	}
	lineC := make(chan []byte)
	errC := make(chan error)
	go func() {
		line, err := reader.ReadBytes('\n')
		if ctx.Err() != nil {
			// If the context has already closed, then just return out of this.
			return
		}
		if err != nil {
			errC <- err
			return
		}
		lineC <- line
	}()
	var wrapper struct {
		Result interface{}
	}
	var line []byte
	select {
	case v := <-lineC:
		line = v
	case err := <-errC:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	err = json.Unmarshal(line, &wrapper)
	if err != nil {
		return nil, err
	}
	envJSON, err := json.Marshal(wrapper.Result)
	if err != nil {
		return nil, err
	}

	var env messagev1.Envelope
	err = protojson.Unmarshal(envJSON, &env)
	return &env, err
}

func (s *httpStream) Close() error {
	if !s.closed {
		defer close(s.respC)
		defer close(s.errC)
	}
	s.closed = true
	if s.body == nil {
		return nil
	}
	err := s.body.Close()
	s.body = nil
	return err
}
