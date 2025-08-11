package migrator

import (
	"context"
	"database/sql"

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

	appChainClient, err := blockchain.NewRPCClient(
		ctx,
		cfg.Contracts.AppChain.RPCURL,
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
