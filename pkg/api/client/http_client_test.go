package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/client"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func Test_Client(t *testing.T) {
	var header http.Header
	var host string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header = r.Header
		host = r.Host
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()
	ctx := test.NewContext(t)
	client := client.NewHTTPClient(ctx.Logger(),
		server.URL, "test", "test client",
		client.WithHeader("Host", "override.localhost"),
		client.WithHeader("x-special", "something"))
	_, err := client.Query(ctx, api.NewQuery("topic"))
	assert.NoError(t, err)
	assert.Equal(t, "override.localhost", host)
	assert.Equal(t, "something", header.Get("X-Special"))
	assert.Equal(t, "test client", header.Get("X-App-Version"))
	assert.Equal(t, "xmtpd/test", header.Get("X-Client-Version"))
}
