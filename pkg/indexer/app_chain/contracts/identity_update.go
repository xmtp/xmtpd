package contracts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	identityUpdateName  = "identity-update-broadcaster"
	identityUpdateTopic = "IdentityUpdateCreated"
)

type IdentityUpdateBroadcaster struct {
	address common.Address
	topics  []common.Hash
	logger  *zap.Logger
	c.IBlockTracker
	c.IReorgHandler
	c.ILogStorer
}

var _ c.IContract = &IdentityUpdateBroadcaster{}

func NewIdentityUpdateBroadcaster(
	ctx context.Context,
	client *ethclient.Client,
	db *sql.DB,
	logger *zap.Logger,
	validationService mlsvalidate.MLSValidationService,
	address common.Address,
	chainID int64,
	startBlock uint64,
) (*IdentityUpdateBroadcaster, error) {
	contract, err := identityUpdateBroadcasterContract(address, client)
	if err != nil {
		return nil, err
	}

	querier := queries.New(db)

	identityUpdatesTracker, err := c.NewBlockTracker(
		ctx,
		client,
		address,
		querier,
		startBlock,
	)
	if err != nil {
		return nil, err
	}

	topics, err := identityUpdateBroadcasterTopic()
	if err != nil {
		return nil, err
	}

	logger = logger.Named(utils.IdentityUpdateBroadcasterLoggerName).
		With(utils.ContractAddressField(address.Hex()))

	identityUpdateStorer := NewIdentityUpdateStorer(db, logger, contract, validationService)

	reorgHandler := NewIdentityUpdateReorgHandler(logger)

	return &IdentityUpdateBroadcaster{
		address:       address,
		topics:        []common.Hash{topics},
		logger:        logger,
		IBlockTracker: identityUpdatesTracker,
		IReorgHandler: reorgHandler,
		ILogStorer:    identityUpdateStorer,
	}, nil
}

func (iu *IdentityUpdateBroadcaster) Address() common.Address {
	return iu.address
}

func (iu *IdentityUpdateBroadcaster) Topics() []common.Hash {
	return iu.topics
}

func (iu *IdentityUpdateBroadcaster) Logger() *zap.Logger {
	return iu.logger
}

func identityUpdateBroadcasterContract(
	address common.Address,
	client *ethclient.Client,
) (*iu.IdentityUpdateBroadcaster, error) {
	return iu.NewIdentityUpdateBroadcaster(
		address,
		client,
	)
}

func IdentityUpdateBroadcasterName(chainID int64) string {
	return fmt.Sprintf("%s-%v", identityUpdateName, chainID)
}

func identityUpdateBroadcasterTopic() (common.Hash, error) {
	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, identityUpdateTopic)
}
