// Package stress implements the stress test for the identity updates.
package stress

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	identityUpdatesSignature = "addIdentityUpdate(bytes32,bytes)"
	identityUpdatesInboxID   = "0xdee4b6c2d041591ff89952fb57487e594ec42d07063a94cd4535e80af25c8e2e"
	identityUpdatesPayload   = "0a23122102d3be290e677cce90cb6300f584dc29b201b5fa136aac59a5a1eec4ac4d001f592aa0030a790a770a2a30783231383164633965343032633731326566653161633130336564356430373730343965373133333910011a450a430a41ded55d591b4007c54463b3bcdb6ccb7d508b98b6a029e2b898fb364edd0a0a516a700701caa15da918d80dba740777cae1159cae83ab316fb45443fe26922e851c20010ad60112d3010a221220a0cd41533d41444017a9e4fc401b66f4207b9d9eb16d7d9b747aa26fb76b953912450a430a41ded55d591b4007c54463b3bcdb6ccb7d508b98b6a029e2b898fb364edd0a0a516a700701caa15da918d80dba740777cae1159cae83ab316fb45443fe26922e851c1a661a640a40540884bcc933b83cd1089af67818c7c673a4b76ad31aa823d4921a37013ce868617cac8b208cd0be1022733e7186b7e378ce1852544200a3998971e732d8510a1220a0cd41533d41444017a9e4fc401b66f4207b9d9eb16d7d9b747aa26fb76b953910b8cab6c5efa9bd9c181a4064336265323930653637376363653930636236333030663538346463323962323031623566613133366161633539613561316565633461633464303031663539"
)

func StressIdentityUpdates(
	ctx context.Context,
	logger *zap.Logger,
	n int,
	contractAddress common.Address,
	rpc string,
	privateKey string,
	async bool,
) error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	type Result struct {
		Idx      int
		Nonce    uint64
		Success  bool
		Duration time.Duration
	}
	results := make([]Result, 0, n)

	// concurrency limiter
	// if you see NONCE TOO HIGH errors it might mean that this guesstimate is too high
	semaphore := make(chan struct{}, 50)

	startingNonce, err := getCurrentNonce(ctx, privateKey, rpc)
	if err != nil {
		return fmt.Errorf("failed to get starting nonce: %w", err)
	}

	var nonceCounter atomic.Uint64
	nonceCounter.Store(startingNonce)

	totalTime := time.Now()

	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			nonce := nonceCounter.Add(1) - 1

			cs := &CastSendCommand{
				ContractAddress: contractAddress.String(),
				Function:        identityUpdatesSignature,
				FunctionArgs:    []string{identityUpdatesInboxID, identityUpdatesPayload},
				RPC:             rpc,
				PrivateKey:      privateKey,
				Nonce:           &nonce,
				Async:           async,
			}

			startTime := time.Now()
			logger.Info("starting transaction", zap.Int("idx", idx), zap.Uint64("nonce", nonce))

			err := cs.Run(ctx)
			duration := time.Since(startTime)

			mu.Lock()
			results = append(results, Result{
				Idx:      idx,
				Nonce:    nonce,
				Success:  err == nil,
				Duration: duration,
			})
			mu.Unlock()

			if err != nil {
				logger.Error("error", zap.Int("idx", idx), zap.Error(err))
			} else {
				logger.Info(
					"completed transaction",
					zap.Int("idx", idx),
					zap.Uint64("nonce", nonce),
					zap.Duration("duration", duration),
				)
			}
		}(i)
	}

	wg.Wait()

	elapsedSec := time.Since(totalTime).Seconds()

	// Analytics
	var totalDuration time.Duration
	var successCount int
	for _, res := range results {
		totalDuration += res.Duration
		if res.Success {
			successCount++
		}
	}

	tps := float64(n) / elapsedSec

	avgDuration := totalDuration / time.Duration(n)
	logger.Info("stress test summary",
		zap.Int("total_transactions", n),
		zap.Int("successful_transactions", successCount),
		zap.Float64("success_rate", float64(successCount)/float64(n)),
		zap.Duration("total_duration", totalDuration),
		zap.Duration("average_duration", avgDuration),
		zap.Float64("tps", tps),
	)

	return nil
}

func getCurrentNonce(ctx context.Context, privateKey, wsURL string) (uint64, error) {
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	defer client.Close()

	address, err := getAddressFromPrivateKey(privateKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get address from private key: %w", err)
	}

	nonce, err := client.NonceAt(ctx, address, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce for address %s: %w", address, err)
	}

	return nonce, nil
}

func getAddressFromPrivateKey(privateKeyHex string) (common.Address, error) {
	privateKey, err := utils.ParseEcdsaPrivateKey(privateKeyHex)
	if err != nil {
		return common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address, nil
}
