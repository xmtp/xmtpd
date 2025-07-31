package migrator

import (
	"context"
	"database/sql"
	"net"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	sqlmgr "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/sql"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

func setupBlockchainPublisher(
	ctx context.Context,
	log *zap.Logger,
	db *sql.DB,
	cfg *config.MigrationServerOptions,
) (*blockchain.BlockchainPublisher, error) {
	nonceManager := sqlmgr.NewSQLBackedNonceManager(db, log)

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
