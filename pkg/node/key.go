package node

import (
	"encoding/hex"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/pkg/errors"
)

func GeneratePrivateKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 1)
	return priv, err
}

func PrivateKeyToHex(key crypto.PrivKey) (string, error) {
	keyBytes, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}

func DecodePrivateKey(keyHex string) (crypto.PrivKey, error) {
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, errors.Wrap(err, "decoding private key")
	}
	return crypto.UnmarshalPrivateKey(keyBytes)
}

func getOrCreatePrivateKey(keyHex string) (crypto.PrivKey, error) {
	if keyHex == "" {
		priv, err := GeneratePrivateKey()
		if err != nil {
			return nil, err
		}

		return priv, nil
	}
	return DecodePrivateKey(keyHex)
}
