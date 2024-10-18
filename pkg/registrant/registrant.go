package registrant

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"slices"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type Registrant struct {
	record       *registry.Node
	privateKey   *ecdsa.PrivateKey
	tokenFactory *authn.TokenFactory
}

func NewRegistrant(
	ctx context.Context,
	log *zap.Logger,
	db *queries.Queries,
	nodeRegistry registry.NodeRegistry,
	privateKeyString string,
) (*Registrant, error) {
	privateKey, err := utils.ParseEcdsaPrivateKey(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse signer.private-key: %v", err)
	}

	record, err := getRegistryRecord(nodeRegistry, privateKey)
	if err != nil {
		return nil, err
	}
	log.Debug(fmt.Sprintf("Running with public key %x", crypto.FromECDSAPub(record.SigningKey)))

	if err = ensureDatabaseMatches(ctx, db, record); err != nil {
		return nil, err
	}

	tokenFactory := authn.NewTokenFactory(privateKey, record.NodeID)

	return &Registrant{
		record:       record,
		privateKey:   privateKey,
		tokenFactory: tokenFactory,
	}, nil
}

func (r *Registrant) sign(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, r.privateKey)
}

func (r *Registrant) NodeID() uint32 {
	return r.record.NodeID
}

func (r *Registrant) TokenFactory() *authn.TokenFactory {
	return r.tokenFactory
}

func (r *Registrant) SignStagedEnvelope(
	stagedEnv queries.StagedOriginatorEnvelope,
) (*message_api.OriginatorEnvelope, error) {
	payerEnv := &message_api.PayerEnvelope{}
	if err := proto.Unmarshal(stagedEnv.PayerEnvelope, payerEnv); err != nil {
		return nil, fmt.Errorf("Could not unmarshal payer envelope: %v", err)
	}
	unsignedEnv := message_api.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     r.record.NodeID,
		OriginatorSequenceId: uint64(stagedEnv.ID),
		OriginatorNs:         stagedEnv.OriginatorTime.UnixNano(),
		PayerEnvelope:        payerEnv,
	}
	unsignedBytes, err := proto.Marshal(&unsignedEnv)
	if err != nil {
		return nil, err
	}

	sig, err := r.sign(utils.HashOriginatorSignatureInput(unsignedBytes))
	if err != nil {
		return nil, err
	}

	signedEnv := message_api.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof: &message_api.OriginatorEnvelope_OriginatorSignature{
			OriginatorSignature: &associations.RecoverableEcdsaSignature{
				Bytes: sig,
			},
		},
	}

	return &signedEnv, nil
}

func getRegistryRecord(
	nodeRegistry registry.NodeRegistry,
	privateKey *ecdsa.PrivateKey,
) (*registry.Node, error) {
	records, err := nodeRegistry.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("unable to get nodes from registry: %v", err)
	}
	i := slices.IndexFunc(records, func(e registry.Node) bool {
		if e.NodeID == 0 {
			return false
		}
		return e.SigningKey.Equal(&privateKey.PublicKey)
	})
	if i == -1 {
		return nil, fmt.Errorf("no matching public key found in registry")
	}

	return &records[i], nil
}

// Prevents mistakes such as:
// - Running multiple nodes with different private keys against the same DB
// - Changing a server's configuration while pointing to data in an existing DB
func ensureDatabaseMatches(ctx context.Context, db *queries.Queries, record *registry.Node) error {
	numRows, err := db.InsertNodeInfo(
		ctx,
		queries.InsertNodeInfoParams{
			NodeID:    int32(record.NodeID),
			PublicKey: crypto.FromECDSAPub(record.SigningKey),
		},
	)
	if err != nil {
		return fmt.Errorf("unable to insert node info into database: %v", err)
	}

	if numRows == 0 {
		nodeInfo, err := db.SelectNodeInfo(ctx)
		if err != nil {
			return fmt.Errorf("unable to retrieve node info from database: %v", err)
		}
		if nodeInfo.NodeID != int32(record.NodeID) {
			return fmt.Errorf("registry node ID does not match ID in database")
		}
		if !bytes.Equal(nodeInfo.PublicKey, crypto.FromECDSAPub(record.SigningKey)) {
			return fmt.Errorf("registry public key does not match public key in database")
		}
	}

	return nil
}
