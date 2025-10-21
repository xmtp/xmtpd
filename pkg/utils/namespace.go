package utils

import (
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

const namespacePrefix = "xmtpd_"

func BuildNamespace(privateKey string, nodesAddress string) string {
	hash := ethcrypto.Keccak256(
		[]byte(privateKey),
		[]byte(nodesAddress),
	)

	return fmt.Sprintf("%s%s", namespacePrefix, HexEncode(hash)[:12])
}
