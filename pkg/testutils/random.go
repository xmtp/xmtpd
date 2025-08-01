package testutils

import (
	"crypto/ecdsa"
	cryptoRand "crypto/rand"
	"math/rand"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/utils"
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

func RandomInboxIDString() string {
	bytes := RandomBytes(32)

	return utils.HexEncode(bytes)
}

func RandomInboxIDBytes() [32]byte {
	var inboxID [32]byte
	copy(inboxID[:], RandomBytes(32))

	return inboxID
}

func RandomAddress() common.Address {
	bytes := RandomBytes(20)
	return common.BytesToAddress(bytes)
}

func RandomLogTopic() common.Hash {
	bytes := RandomBytes(32)
	return common.BytesToHash(bytes)
}

func RandomGroupID() [16]byte {
	var groupID [16]byte
	copy(groupID[:], RandomBytes(16))

	return groupID
}

func RandomPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	return key
}

func RandomDomainSeparator() common.Hash {
	bytes := RandomBytes(32)
	return common.BytesToHash(bytes)
}

func RandomBlockHash() common.Hash {
	bytes := RandomBytes(32)
	return common.BytesToHash(bytes)
}

func RandomInt32() int32 {
	return rand.Int31()
}

func RandomInt64() int64 {
	return rand.Int63()
}
