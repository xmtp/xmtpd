package commands

import (
	"encoding/hex"
	"errors"
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
		return nil, errors.New("value exceeds 96 bits")
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
	blockchain.AppChainGatewayPausedKey:           ParamBool,
	blockchain.SettlementChainGatewayPausedKey:    ParamBool,
	blockchain.DistributionManagerPausedKey:       ParamBool,
	blockchain.GroupMessageBroadcasterPausedKey:   ParamBool,
	blockchain.IdentityUpdateBroadcasterPausedKey: ParamBool,
	blockchain.PayerRegistryPausedKey:             ParamBool,

	// Addresses
	blockchain.DistributionManagerProtocolFeesRecipientKey: ParamAddress,
	blockchain.GroupMessagePayloadBootstrapperKey:          ParamAddress,
	blockchain.IdentityUpdatePayloadBootstrapperKey:        ParamAddress,
	blockchain.NodeRegistryAdminKey:                        ParamAddress,
	blockchain.RateRegistryMigratorKey:                     ParamAddress,

	// Sizes / counts
	blockchain.NodeRegistryMaxCanonicalNodesKey:           ParamUint8,
	blockchain.GroupMessageBroadcasterMaxPayloadSizeKey:   ParamUint32,
	blockchain.GroupMessageBroadcasterMinPayloadSizeKey:   ParamUint32,
	blockchain.IdentityUpdateBroadcasterMaxPayloadSizeKey: ParamUint32,
	blockchain.IdentityUpdateBroadcasterMinPayloadSizeKey: ParamUint32,

	// Monetary / big ints
	blockchain.PayerRegistryMinimumDepositKey: ParamUint96,

	// Durations
	blockchain.PayerRegistryWithdrawLockPeriodKey: ParamUint32,

	// Fees / rates
	blockchain.PayerReportManagerProtocolFeeRateKey: ParamUint16,
	blockchain.RateRegistryCongestionFeeKey:         ParamUint64,
	blockchain.RateRegistryMessageFeeKey:            ParamUint64,
	blockchain.RateRegistryStorageFeeKey:            ParamUint64,
	blockchain.RateRegistryTargetRatePerMinuteKey:   ParamUint64,
}

func paramType(key string) ParamType {
	t, ok := paramTypeByKey[key]
	if !ok {
		return ParamUnknown
	}
	return t
}
