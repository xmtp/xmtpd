package payerreport

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBuildMerkle(t *testing.T) {
	cases := []struct {
		name      string
		payerMap  payerMap
		expectErr bool
	}{
		{
			name: "success",
			payerMap: payerMap{
				common.HexToAddress("0x1"): 100,
				common.HexToAddress("0x2"): 200,
			},
			expectErr: false,
		},
		{
			name: "negative fee",
			payerMap: payerMap{
				common.HexToAddress("0x1"): -100,
			},
			expectErr: true,
		},
		{
			name: "zero fee",
			payerMap: payerMap{
				common.HexToAddress("0x1"): 0,
			},
			expectErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := generateMerkleTree(c.payerMap)
			if c.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
