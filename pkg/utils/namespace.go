package utils

import (
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func BuildNamespace(privateKey string, nodesAddress string) string {
	hash := ethcrypto.Keccak256(
		[]byte(privateKey),
		[]byte(nodesAddress),
	)

	return fmt.Sprintf("xmtpd_%s", HexEncode(hash)[:12])
}
