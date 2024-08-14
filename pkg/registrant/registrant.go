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
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
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
		queries.InsertNodeInfoParams{NodeID: int32(record.NodeId), PublicKey: crypto.FromECDSAPub(record.SigningKey)},
	)
	if err == sql.ErrNoRows {
		// Node info already exists in database - verify it matches
		// the record
		nodeInfo, err := db.SelectNodeInfo(ctx)
		if err != nil {
			panic("unable to select node info")
		}
		if nodeInfo.NodeID != int32(record.NodeId) {
			panic("registry node ID does not match database entry")
		}
		if !bytes.Equal(nodeInfo.PublicKey, crypto.FromECDSAPub(record.SigningKey)) {
			panic("public key does not match database entry")
		}
	}

	return &Registrant{
		record:     record,
		privateKey: privateKey,
	}, nil
}

func (r *Registrant) SID(localID int64) uint64 {
	if !utils.IsValidLocalID(localID) {
		// Either indicates ID exhaustion or developer error -
		// the service should not continue running either way
		panic(fmt.Sprintf("invalid local ID %d", localID))
	}
	return utils.SID(int32(r.record.NodeId), localID)
}

func (r *Registrant) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, r.privateKey)
}
