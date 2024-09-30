package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateEcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func ParseEcdsaPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(key, "0x"))
}

func EcdsaPrivateKeyToString(key *ecdsa.PrivateKey) string {
	keyBytes := crypto.FromECDSA(key)
	return "0x" + hex.EncodeToString(keyBytes)
}

// Take the stringified form of an ECDSA public key and return the *ecdsa.PublicKey
func ParseEcdsaPublicKey(key string) (*ecdsa.PublicKey, error) {
	return crypto.DecompressPubkey(common.FromHex(key))
}

func EcdsaPublicKeyToString(key *ecdsa.PublicKey) string {
	return "0x" + hex.EncodeToString(crypto.FromECDSAPub(key))
}
