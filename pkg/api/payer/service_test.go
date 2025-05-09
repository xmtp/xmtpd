package payer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetReaderNode(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []registry.Node
		registryError error
		wantErr       bool
		checkResponse func(t *testing.T, resp *payer_api.GetReaderNodeResponse)
	}{
		{
			name: "success with multiple nodes",
			nodes: []registry.Node{
				testutils.GetHealthyNode(100),
				testutils.GetHealthyNode(101),
				testutils.GetHealthyNode(102),
				testutils.GetHealthyNode(103),
			},
			wantErr: false,
			checkResponse: func(t *testing.T, resp *payer_api.GetReaderNodeResponse) {
				require.NotEmpty(t, resp.ReaderNodeUrl)
				require.Len(t, resp.BackupNodeUrls, 3)

				for _, backup := range resp.BackupNodeUrls {
					require.NotEqual(t, resp.ReaderNodeUrl, backup)
				}
			},
		},
		{
			name:    "no nodes available",
			nodes:   []registry.Node{},
			wantErr: true,
		},
		{
			name:          "registry error",
			nodes:         nil,
			registryError: status.Errorf(codes.Unavailable, "registry unavailable"),
			wantErr:       true,
		},
		{
			name: "single node",
			nodes: []registry.Node{
				testutils.GetHealthyNode(100),
			},
			wantErr: false,
			checkResponse: func(t *testing.T, resp *payer_api.GetReaderNodeResponse) {
				require.NotEmpty(t, resp.ReaderNodeUrl)
				require.Empty(t, resp.BackupNodeUrls)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc, _, registryMocks, _, cleanup := buildPayerService(t)
			defer cleanup()

			registryMocks.On("GetNodes").Return(tt.nodes, tt.registryError)

			resp, err := svc.GetReaderNode(ctx, &payer_api.GetReaderNodeRequest{})

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
		})
	}
}
