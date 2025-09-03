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
	payerPrivateKey string,
	cfg *config.ContractsOptions,
) (*blockchain.BlockchainPublisher, error) {
	nonceManager := sqlmgr.NewSQLBackedNonceManager(db, log)

	signer, err := blockchain.NewPrivateKeySigner(
		payerPrivateKey,
		cfg.AppChain.ChainID,
	)
	if err != nil {
		return nil, err
	}

	appChainClient, err := blockchain.NewRPCClient(
		ctx,
		cfg.AppChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	return blockchain.NewBlockchainPublisher(
		ctx,
		log,
		appChainClient,
		signer,
		*cfg,
		nonceManager,
	)
}
