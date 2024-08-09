package node

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
)

// Stores: public key, private key, id
// Performs: signing, SID generation
// Verify against registry
// Verify against DB

type Node struct {
	record     registry.Record
	privateKey *ecdsa.PrivateKey
}

func NewNode(
	ctx context.Context,
	db *queries.Queries,
	nodeRegistry registry.NodeRegistry,
	privateKeyString string,
) (*Node, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}
	publicKey := crypto.FromECDSAPub(&privateKey.PublicKey)

	records, err := nodeRegistry.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("unable to get nodes from registry: %v", err)
	}

	i := slices.IndexFunc(records, func(e registry.Record) bool {
		return bytes.Equal(e.PublicKey, publicKey)
	})
	if i == -1 {
		return nil, fmt.Errorf("no matching public key found in registry")
	}
	record := records[i]

	_, err = db.InsertNodeInfo(
		ctx,
		queries.InsertNodeInfoParams{NodeID: record.ID, PublicKey: record.PublicKey},
	)
	if err == sql.ErrNoRows {
		// Node info already exists in database - verify it matches
		// node info on initialization
		nodeInfo, err := db.SelectNodeInfo(ctx)
		if err != nil {
			panic("unable to select node info")
		}
		if nodeInfo.NodeID != record.ID {
			panic("registry node ID does not match database entry")
		}
		if !bytes.Equal(nodeInfo.PublicKey, publicKey) {
			panic("public key does not match database entry")
		}
	}

	return &Node{
		record:     record,
		privateKey: privateKey,
	}, nil
}
