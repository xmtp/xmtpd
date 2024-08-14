package registrant

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"database/sql"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type Registrant struct {
	record     registry.Node
	privateKey *ecdsa.PrivateKey
}

func NewRegistrant(
	ctx context.Context,
	db *queries.Queries,
	nodeRegistry registry.NodeRegistry,
	privateKeyString string,
) (*Registrant, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	records, err := nodeRegistry.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("unable to get nodes from registry: %v", err)
	}

	i := slices.IndexFunc(records, func(e registry.Node) bool {
		return e.SigningKey.Equal(&privateKey.PublicKey)
	})
	if i == -1 {
		return nil, fmt.Errorf("no matching public key found in registry")
	}
	record := records[i]

	_, err = db.InsertNodeInfo(
		ctx,
		queries.InsertNodeInfoParams{NodeID: int32(record.NodeID), PublicKey: crypto.FromECDSAPub(record.SigningKey)},
	)
	if err == sql.ErrNoRows {
		// Node info already exists in database - verify it matches
		// the record
		nodeInfo, err := db.SelectNodeInfo(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve node info from database: %v", err)
		}
		if nodeInfo.NodeID != int32(record.NodeID) {
			return nil, fmt.Errorf("registry node ID does not match ID in database")
		}
		if !bytes.Equal(nodeInfo.PublicKey, crypto.FromECDSAPub(record.SigningKey)) {
			return nil, fmt.Errorf("registry public key does not match public key in database")
		}
	}

	return &Registrant{
		record:     record,
		privateKey: privateKey,
	}, nil
}

func (r *Registrant) sid(localID int64) uint64 {
	if !utils.IsValidLocalID(localID) {
		// Either indicates ID exhaustion or developer error -
		// the service should not continue running either way
		panic(fmt.Sprintf("invalid local ID %d", localID))
	}
	return utils.SID(r.record.NodeID, localID)
}

func (r *Registrant) sign(data []byte) ([]byte, error) {
	hash := crypto.Keccak256(data)
	return crypto.Sign(hash, r.privateKey)
}

func (r *Registrant) SignStagedEnvelope(
	stagedEnv queries.StagedOriginatorEnvelope,
) (*message_api.OriginatorEnvelope, error) {
	payerEnv := &message_api.PayerEnvelope{}
	if err := proto.Unmarshal(stagedEnv.PayerEnvelope, payerEnv); err != nil {
		return nil, err
	}
	unsignedEnv := message_api.UnsignedOriginatorEnvelope{
		OriginatorSid: r.sid(stagedEnv.ID),
		OriginatorNs:  stagedEnv.OriginatorTime.UnixNano(),
		PayerEnvelope: payerEnv,
	}
	unsignedBytes, err := proto.Marshal(&unsignedEnv)
	if err != nil {
		return nil, err
	}

	sig, err := r.sign(unsignedBytes)
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
