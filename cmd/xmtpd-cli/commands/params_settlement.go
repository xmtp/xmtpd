package commands

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ---------- opts ----------

type SettlementSetOpts struct {
	KVs        []string
	Raw        bool
	TimeoutSec int
}

type SettlementGetOpts struct {
	Keys       []string
	Raw        bool
	TimeoutSec int
}

// ---------- root (params settlement) ----------

func paramsSettlementCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "settlement",
		Short:        "Operate on Settlement chain parameters",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		settlementSetCmd(),
		settlementGetCmd(),
	)
	return cmd
}

// ---------- set ----------

func settlementSetCmd() *cobra.Command {
	var opts SettlementSetOpts

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set parameter(s) in the Settlement Parameter Registry (generic key/value)",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return settlementSetHandler(opts)
		},
		Example: `
xmtpd-cli params settlement set \
  --kv xmtp.rateRegistry.messageFee=100 \
  --kv xmtp.settlementChainGateway.paused=true`,
	}

	cmd.Flags().StringArrayVar(&opts.KVs, "kv", nil, "key=value")
	cmd.Flags().IntVar(&opts.TimeoutSec, "timeout", 120, "timeout (seconds)")
	cmd.Flags().BoolVar(&opts.Raw, "raw", false, "treat value as 0x-prefixed 32-byte hex")

	_ = cmd.MarkFlagRequired("kv")

	return cmd
}

func settlementSetHandler(opts SettlementSetOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}

	if len(opts.KVs) == 0 {
		return errors.New("at least one --kv is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(opts.TimeoutSec)*time.Second,
	)
	defer cancel()

	paramAdmin, _, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup parameter admin", zap.Error(err))
		return err
	}

	type kv struct {
		key    string
		valB32 *[32]byte

		// typed
		pt    ParamType
		boolV *bool
		addrV *common.Address
		u8V   *uint8
		u16V  *uint16
		u32V  *uint32
		u64V  *uint64
		u96V  *big.Int
	}
	var items []kv

	for _, kvs := range opts.KVs {
		key, rawVal, perr := splitKV(kvs)
		if perr != nil {
			return perr
		}
		key = strings.TrimSpace(key)

		if opts.Raw {
			b32, perr := parseBytes32(rawVal)
			if perr != nil {
				return fmt.Errorf("invalid value for key %s: %w", key, perr)
			}
			items = append(items, kv{key: key, valB32: &b32})
			continue
		}

		item := kv{key: key, pt: paramType(key)}

		switch item.pt {
		case ParamBool:
			v, err := strconv.ParseBool(rawVal)
			if err != nil {
				return fmt.Errorf("parse bool for %s: %w", key, err)
			}
			item.boolV = &v
		case ParamAddress:
			v, err := parseAddressString(rawVal)
			if err != nil {
				return fmt.Errorf("parse address for %s: %w", key, err)
			}
			item.addrV = &v
		case ParamUint8:
			v, err := parseUintFit[uint8](rawVal)
			if err != nil {
				return fmt.Errorf("parse uint8 for %s: %w", key, err)
			}
			item.u8V = &v
		case ParamUint16:
			v, err := parseUintFit[uint16](rawVal)
			if err != nil {
				return fmt.Errorf("parse uint16 for %s: %w", key, err)
			}
			item.u16V = &v
		case ParamUint32:
			v, err := parseUintFit[uint32](rawVal)
			if err != nil {
				return fmt.Errorf("parse uint32 for %s: %w", key, err)
			}
			item.u32V = &v
		case ParamUint64:
			v, err := parseUintFit[uint64](rawVal)
			if err != nil {
				return fmt.Errorf("parse uint64 for %s: %w", key, err)
			}
			item.u64V = &v
		case ParamUint96:
			v, err := parseUint96(rawVal)
			if err != nil {
				return fmt.Errorf("parse uint96 for %s: %w", key, err)
			}
			item.u96V = v
		default:
			return fmt.Errorf("unsupported/unknown param type for key %s. Use --raw", key)
		}

		items = append(items, item)
	}

	for _, it := range items {
		if opts.Raw && it.valB32 != nil {
			err = paramAdmin.SetRawParameter(ctx, it.key, *it.valB32)
		} else {
			switch it.pt {
			case ParamBool:
				err = paramAdmin.SetBoolParameter(ctx, it.key, *it.boolV)
			case ParamAddress:
				err = paramAdmin.SetAddressParameter(ctx, it.key, *it.addrV)
			case ParamUint8:
				err = paramAdmin.SetUint8Parameter(ctx, it.key, *it.u8V)
			case ParamUint16:
				err = paramAdmin.SetUint16Parameter(ctx, it.key, *it.u16V)
			case ParamUint32:
				err = paramAdmin.SetUint32Parameter(ctx, it.key, *it.u32V)
			case ParamUint64:
				err = paramAdmin.SetUint64Parameter(ctx, it.key, *it.u64V)
			case ParamUint96:
				err = paramAdmin.SetUint96Parameter(ctx, it.key, it.u96V)
			default:
				return fmt.Errorf("unhandled param type for key %s", it.key)
			}
		}

		if err != nil {
			logger.Error("set parameter failed", zap.String("key", it.key), zap.Error(err))
			return err
		}
		logger.Info("parameter set", zap.String("key", it.key))
	}

	logger.Info("all parameters set successfully")
	return nil
}

// ---------- get ----------

func settlementGetCmd() *cobra.Command {
	var opts SettlementGetOpts

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get parameter(s) from the Settlement Parameter Registry (generic)",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return settlementGetHandler(opts)
		},
		Example: `
xmtpd-cli params settlement get \
  --key xmtp.rateRegistry.messageFee \
  --key xmtp.settlementChainGateway.paused`,
	}

	cmd.Flags().StringArrayVar(&opts.Keys, "key", nil, "parameter key (repeatable)")
	cmd.Flags().IntVar(&opts.TimeoutSec, "timeout", 60, "timeout (seconds)")
	cmd.Flags().BoolVar(&opts.Raw, "raw", false, "treat value as 0x-prefixed 32-byte hex")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

func settlementGetHandler(opts SettlementGetOpts) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}

	if len(opts.Keys) == 0 {
		return errors.New("at least one --key is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(opts.TimeoutSec)*time.Second,
	)
	defer cancel()

	paramAdmin, _, err := setupSettlementChainAdmin(ctx, logger)
	if err != nil {
		logger.Error("could not setup parameter admin", zap.Error(err))
		return err
	}

	for _, k := range opts.Keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}

		if opts.Raw {
			val, gerr := paramAdmin.GetRawParameter(ctx, k)
			if gerr != nil {
				logger.Error("get parameter failed", zap.String("key", k), zap.Error(gerr))
				return gerr
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("bytes32", "0x"+common.Bytes2Hex(val[:])),
			)
			continue
		}

		switch paramType(k) {
		case ParamBool:
			v, err := paramAdmin.GetParameterBool(ctx, k)
			if err != nil {
				logger.Error("get bool failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.Bool("bool", v),
			)
		case ParamAddress:
			v, err := paramAdmin.GetParameterAddress(ctx, k)
			if err != nil {
				logger.Error("get address failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("address", v.Hex()),
			)
		case ParamUint8:
			v, err := paramAdmin.GetParameterUint8(ctx, k)
			if err != nil {
				logger.Error("get uint8 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint8("uint8", v))
		case ParamUint16:
			v, err := paramAdmin.GetParameterUint16(ctx, k)
			if err != nil {
				logger.Error("get uint16 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint16("uint16", v))
		case ParamUint32:
			v, err := paramAdmin.GetParameterUint32(ctx, k)
			if err != nil {
				logger.Error("get uint32 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint32("uint32", v))
		case ParamUint64:
			v, err := paramAdmin.GetParameterUint64(ctx, k)
			if err != nil {
				logger.Error("get uint64 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter", zap.String("key", k), zap.Uint64("uint64", v))
		case ParamUint96:
			v, err := paramAdmin.GetParameterUint96(ctx, k)
			if err != nil {
				logger.Error("get uint96 failed", zap.String("key", k), zap.Error(err))
				return err
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("uint96", v.String()),
			)
		default:
			// Fallback: raw
			val, gerr := paramAdmin.GetRawParameter(ctx, k)
			if gerr != nil {
				logger.Error("get parameter failed", zap.String("key", k), zap.Error(gerr))
				return gerr
			}
			logger.Info("parameter",
				zap.String("key", k),
				zap.String("bytes32", "0x"+common.Bytes2Hex(val[:])),
			)
		}
	}

	return nil
}
