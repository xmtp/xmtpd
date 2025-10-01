package commands

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/blockchain"
)

func splitKV(s string) (key, val string, err error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid --kv %q (expected key=value)", s)
	}
	key = strings.TrimSpace(parts[0])
	val = strings.TrimSpace(parts[1])
	if key == "" || val == "" {
		return "", "", fmt.Errorf("invalid --kv %q (empty key or value)", s)
	}
	return key, val, nil
}

func parseBytes32(hexStr string) ([32]byte, error) {
	var out [32]byte
	h := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(hexStr)), "0x")
	if len(h) != 64 {
		return out, fmt.Errorf("want 32 bytes (64 hex chars), got %d", len(h))
	}
	b, err := hex.DecodeString(h)
	if err != nil {
		return out, fmt.Errorf("decode hex: %w", err)
	}
	copy(out[:], b)
	return out, nil
}

func parseAddressString(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("invalid address %q", s)
	}
	return common.HexToAddress(s), nil
}

func parseUintFit[T ~uint8 | ~uint16 | ~uint32 | ~uint64](s string) (T, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, err
	}
	maxSize := ^T(0)
	if u > uint64(maxSize) {
		return 0, fmt.Errorf("value %d overflows %T", u, maxSize)
	}
	return T(u), nil
}

func parseUint96(s string) (*big.Int, error) {
	// accepts decimal or 0x-hex
	ss := strings.TrimSpace(s)
	bi := new(big.Int)
	var ok bool
	if strings.HasPrefix(ss, "0x") || strings.HasPrefix(ss, "0X") {
		// hex
		_, ok = bi.SetString(ss[2:], 16)
	} else {
		_, ok = bi.SetString(ss, 10)
	}
	if !ok {
		return nil, fmt.Errorf("invalid uint96: %q", s)
	}
	// Ensure bi fits within 96 bits (<= 2^96-1)
	if bi.BitLen() > 96 {
		return nil, fmt.Errorf("value exceeds 96 bits")
	}
	return bi, nil
}

type ParamType int

const (
	ParamUnknown ParamType = iota
	ParamBool
	ParamAddress
	ParamUint8
	ParamUint16
	ParamUint32
	ParamUint64
	ParamUint96
)

// Known keys → types (fill with what we know today; easy to extend)
var paramTypeByKey = map[string]ParamType{
	// Paused flags (bool)
	blockchain.APP_CHAIN_GATEWAY_PAUSED_KEY:           ParamBool,
	blockchain.SETTLEMENT_CHAIN_GATEWAY_PAUSED_KEY:    ParamBool,
	blockchain.DISTRIBUTION_MANAGER_PAUSED_KEY:        ParamBool,
	blockchain.GROUP_MESSAGE_BROADCASTER_PAUSED_KEY:   ParamBool,
	blockchain.IDENTITY_UPDATE_BROADCASTER_PAUSED_KEY: ParamBool,
	blockchain.PAYER_REGISTRY_PAUSED_KEY:              ParamBool,

	// Addresses
	blockchain.DISTRIBUTION_MANAGER_PROTOCOL_FEES_RECIPIENT_KEY: ParamAddress,
	blockchain.GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY:           ParamAddress,
	blockchain.IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY:         ParamAddress,
	blockchain.NODE_REGISTRY_ADMIN_KEY:                          ParamAddress,
	blockchain.RATE_REGISTRY_MIGRATOR_KEY:                       ParamAddress,

	// Sizes / counts
	blockchain.NODE_REGISTRY_MAX_CANONICAL_NODES_KEY:            ParamUint8,
	blockchain.GROUP_MESSAGE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY:   ParamUint32,
	blockchain.GROUP_MESSAGE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY:   ParamUint32,
	blockchain.IDENTITY_UPDATE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY: ParamUint32,
	blockchain.IDENTITY_UPDATE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY: ParamUint32,

	// Monetary / big ints
	blockchain.PAYER_REGISTRY_MINIMUM_DEPOSIT_KEY: ParamUint96,

	// Durations
	blockchain.PAYER_REGISTRY_WITHDRAW_LOCK_PERIOD_KEY: ParamUint32,

	// Fees / rates
	blockchain.PAYER_REPORT_MANAGER_PROTOCOL_FEE_RATE_KEY: ParamUint16,
	blockchain.RATE_REGISTRY_CONGESTION_FEE_KEY:           ParamUint64,
	blockchain.RATE_REGISTRY_MESSAGE_FEE_KEY:              ParamUint64,
	blockchain.RATE_REGISTRY_STORAGE_FEE_KEY:              ParamUint64,
	blockchain.RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY:   ParamUint64,
}

func paramType(key string) ParamType {
	t, ok := paramTypeByKey[key]
	if !ok {
		return ParamUnknown
	}
	return t
}
