package stress

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	IDENTITY_UPDATES_SIGNATURE = "addIdentityUpdate(bytes32,bytes)"
	IDENTITY_UPDATES_INBOX_ID  = "0xdee4b6c2d041591ff89952fb57487e594ec42d07063a94cd4535e80af25c8e2e"
	IDENTITY_UPDATES_PAYLOAD   = "0a23122102d3be290e677cce90cb6300f584dc29b201b5fa136aac59a5a1eec4ac4d001f592aa0030a790a770a2a30783231383164633965343032633731326566653161633130336564356430373730343965373133333910011a450a430a41ded55d591b4007c54463b3bcdb6ccb7d508b98b6a029e2b898fb364edd0a0a516a700701caa15da918d80dba740777cae1159cae83ab316fb45443fe26922e851c20010ad60112d3010a221220a0cd41533d41444017a9e4fc401b66f4207b9d9eb16d7d9b747aa26fb76b953912450a430a41ded55d591b4007c54463b3bcdb6ccb7d508b98b6a029e2b898fb364edd0a0a516a700701caa15da918d80dba740777cae1159cae83ab316fb45443fe26922e851c1a661a640a40540884bcc933b83cd1089af67818c7c673a4b76ad31aa823d4921a37013ce868617cac8b208cd0be1022733e7186b7e378ce1852544200a3998971e732d8510a1220a0cd41533d41444017a9e4fc401b66f4207b9d9eb16d7d9b747aa26fb76b953910b8cab6c5efa9bd9c181a4064336265323930653637376363653930636236333030663538346463323962323031623566613133366161633539613561316565633461633464303031663539"
)

func StressIdentityUpdates(
	ctx context.Context,
	logger *zap.Logger,
	n int,
	contractAddress, rpc, privateKey string,
) error {
	var wg sync.WaitGroup

	startingNonce, err := getCurrentNonce(ctx, privateKey, rpc)
	if err != nil {
		return fmt.Errorf("failed to get starting nonce: %s", err)
	}

	for i := 0; i < n; i++ {
		wg.Add(1)

		nonce := startingNonce + i

		cs := &CastSendCommand{
			ContractAddress: contractAddress,
			Function:        IDENTITY_UPDATES_SIGNATURE,
			FunctionArgs:    []string{IDENTITY_UPDATES_INBOX_ID, IDENTITY_UPDATES_PAYLOAD},
			Rpc:             rpc,
			PrivateKey:      privateKey,
			Nonce:           &nonce,
		}

		go func(idx int) {
			defer wg.Done()

			logger.Info("starting transaction", zap.Int("idx", idx))
			if err := cs.Run(ctx); err != nil {
				logger.Error("error", zap.Int("idx", idx), zap.Error(err))
			} else {
				logger.Info("completed transaction", zap.Int("idx", idx))
			}
		}(i)
	}

	wg.Wait()

	return nil
}

func getCurrentNonce(ctx context.Context, privateKey, rpcUrl string) (int, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum node: %s", err)
	}
	defer client.Close()

	address, err := getAddressFromPrivateKey(privateKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get address from private key: %w", err)
	}

	nonce, err := client.NonceAt(ctx, address, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce for address %s: %s", privateKey, err)
	}

	return int(nonce), nil
}

func getAddressFromPrivateKey(privateKeyHex string) (common.Address, error) {
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	} else {
		return common.Address{}, fmt.Errorf("private key must start with 0x")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address, nil
}
