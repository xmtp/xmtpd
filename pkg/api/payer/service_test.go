package payer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetReaderNode(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []registry.Node
		registryError error
		wantErr       bool
		checkResponse func(t *testing.T, resp *payer_api.GetNodesResponse)
	}{
		{
			name: "success with multiple nodes",
			nodes: []registry.Node{
				nodeRegistry.GetHealthyNode(100),
				nodeRegistry.GetHealthyNode(101),
				nodeRegistry.GetHealthyNode(102),
				nodeRegistry.GetHealthyNode(103),
			},
			wantErr: false,
			checkResponse: func(t *testing.T, resp *payer_api.GetNodesResponse) {
				require.NotEmpty(t, resp.Nodes)
				require.Len(t, resp.Nodes, 4)
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
				nodeRegistry.GetHealthyNode(100),
			},
			wantErr: false,
			checkResponse: func(t *testing.T, resp *payer_api.GetNodesResponse) {
				require.NotEmpty(t, resp.Nodes)
				require.Len(t, resp.Nodes, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc, _, registryMocks, _ := buildPayerService(t)

			registryMocks.On("GetNodes").Return(tt.nodes, tt.registryError)

			resp, err := svc.GetNodes(ctx, &payer_api.GetNodesRequest{})

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
