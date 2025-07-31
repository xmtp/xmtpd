package migrator

import (
	"context"
	"net"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	redisnoncemanager "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/redis"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/gateway"
	"go.uber.org/zap"
)

func setupNonceManager(
	ctx context.Context,
	log *zap.Logger,
	cfg *config.MigrationServerOptions,
) (noncemanager.NonceManager, error) {
	redisClient, err := gateway.SetupRedisClient(ctx, cfg.Redis.RedisUrl, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return redisnoncemanager.NewRedisBackedNonceManager(redisClient, log, "xmtpd:migrator")
}

func setupBlockchainPublisher(
	ctx context.Context,
	log *zap.Logger,
	cfg *config.MigrationServerOptions,
	nonceManager noncemanager.NonceManager,
) (*blockchain.BlockchainPublisher, error) {
	signer, err := blockchain.NewPrivateKeySigner(
		cfg.PayerPrivateKey,
		cfg.Contracts.AppChain.ChainID,
	)
	if err != nil {
		return nil, err
	}

	appChainClient, err := blockchain.NewClient(
		ctx,
		blockchain.WithWebSocketURL(cfg.Contracts.AppChain.WssURL),
		blockchain.WithKeepAliveConfig(net.KeepAliveConfig{
			Enable: false,
		}),
	)
	if err != nil {
		return nil, err
	}

	return blockchain.NewBlockchainPublisher(
		ctx,
		log,
		appChainClient,
		signer,
		cfg.Contracts,
		nonceManager,
	)
}
