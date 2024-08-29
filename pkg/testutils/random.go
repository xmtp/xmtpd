package testutils

import (
	cryptoRand "crypto/rand"
	"math/rand"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomStringLower(n int) string {
	return strings.ToLower(RandomString(n))
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = cryptoRand.Read(b)
	return b
}

func RandomInboxId() string {
	bytes := RandomBytes(32)

	return HexEncode(bytes)
}

func RandomAddress() common.Address {
	bytes := RandomBytes(20)
	return common.BytesToAddress(bytes)
}

func RandomLogTopic() common.Hash {
	bytes := RandomBytes(32)
	return common.BytesToHash(bytes)
}
