// Package server implements the replication server test utils.
package server

import (
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	r "github.com/xmtp/xmtpd/pkg/registry"
	s "github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/fees"
)

type EnabledServices struct {
	API     bool
	Indexer bool
	Reports bool
	Sync    bool
}

type TestServerCfg struct {
	ContractsOptions *config.ContractsOptions
	DB               *sql.DB
	Port             int
	PrivateKey       *ecdsa.PrivateKey
	Registry         r.NodeRegistry
	Services         EnabledServices
}

func NewTestBaseServer(
	t *testing.T,
	cfg TestServerCfg,
) *s.BaseServer {
	log := testutils.NewLog(t)

	server, err := s.NewBaseServer(
		s.WithContext(t.Context()),
		s.WithLogger(log),
		s.WithDB(cfg.DB),
		s.WithNodeRegistry(cfg.Registry),
		s.WithServerVersion(testutils.GetLatestVersion(t)),
		s.WithFeeCalculator(fees.NewTestFeeCalculator()),
		s.WithPromReg(prometheus.NewRegistry()),
		s.WithServerOptions(&config.ServerOptions{
			API: config.APIOptions{
				Port:                  cfg.Port,
				Enable:                cfg.Services.API,
				SendKeepAliveInterval: 30 * time.Second,
			},
			Contracts: *cfg.ContractsOptions,
			MlsValidation: config.MlsValidationOptions{
				GrpcAddress: "http://localhost:60051",
			},
			PayerReport: config.PayerReportOptions{
				Enable:                        cfg.Services.Reports,
				AttestationWorkerPollInterval: 10 * time.Second,
				GenerateReportSelfPeriod:      10 * time.Second,
				GenerateReportOthersPeriod:    10 * time.Second,
				ExpirySelfPeriod:              20 * time.Second,
				ExpiryOthersPeriod:            20 * time.Second,
			},
			Signer: config.SignerOptions{
				PrivateKey: hex.EncodeToString(crypto.FromECDSA(cfg.PrivateKey)),
			},
			Sync: config.SyncOptions{
				Enable: cfg.Services.Sync,
			},
			Indexer: config.IndexerOptions{
				Enable: cfg.Services.Indexer,
			},
		}))
	require.NoError(t, err)

	return server
}
