package blockchain

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	paramReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"
	"go.uber.org/zap"
)

const (
	// cast sig "NoChange()" = 0xa88ee577.
	NoChangeError = "0xa88ee577"
)

func SetParameterRegistryKeys(
	ctx context.Context,
	signer TransactionSigner,
	logger *zap.Logger,
	client *ethclient.Client,
	parameterContract *paramReg.SettlementChainParameterRegistry,
	keys []string,
	values [][32]byte,
) error {
	if len(keys) != len(values) {
		return fmt.Errorf("keys and values must have the same length")
	}

	if len(keys) == 0 {
		return fmt.Errorf("keys must not be empty")
	}

	expectedEvents := len(keys)

	err := ExecuteTransaction(
		ctx,
		signer,
		logger,
		client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			if len(keys) == 1 {
				return parameterContract.Set(opts, keys[0], values[0])
			}

			return parameterContract.Set0(opts, keys, values)
		},
		func(receipt *types.Receipt) error {
			gotEvents := 0

			for _, log := range receipt.Logs {
				event, err := parameterContract.ParseParameterSet(*log)
				if err != nil {
					continue
				}

				gotEvents++

				logger.Info("set parameter",
					zap.String("key", event.Key.String()),
					zap.Uint64("parameter", decodeBytes32ToUint64(event.Value)),
				)
			}

			if gotEvents != expectedEvents {
				logger.Warn(
					"unexpected number of events, double check logs and on-chain state",
					zap.Int("expected", expectedEvents),
					zap.Int("got", gotEvents),
				)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func isNoChangeError(err error) bool {
	return strings.Contains(err.Error(), NoChangeError)
}
