package utils

import (
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/config"
)

func BuildNamespace(options config.ServerOptions) string {
	hash := ethcrypto.Keccak256(
		[]byte(options.Signer.PrivateKey),
		[]byte(options.Contracts.NodesContractAddress),
	)

	return fmt.Sprintf("xmtpd_%s", HexEncode(hash)[:12])
}
