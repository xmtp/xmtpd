package payerreport_test

import (
	"testing"

	"github.com/xmtp/xmtpd/pkg/payerreport"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBuildMerkle(t *testing.T) {
	cases := []struct {
		name      string
		payerMap  payerreport.PayerMap
		expectErr bool
	}{
		{
			name: "success",
			payerMap: payerreport.PayerMap{
				common.HexToAddress("0x1"): 100,
				common.HexToAddress("0x2"): 200,
			},
			expectErr: false,
		},
		{
			name: "negative fee",
			payerMap: payerreport.PayerMap{
				common.HexToAddress("0x1"): -100,
			},
			expectErr: true,
		},
		{
			name: "zero fee",
			payerMap: payerreport.PayerMap{
				common.HexToAddress("0x1"): 0,
			},
			expectErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := payerreport.GenerateMerkleTree(c.payerMap)
			if c.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
