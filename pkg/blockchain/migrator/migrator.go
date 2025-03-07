package migrator

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type SerializableNode struct {
	NodeID               uint32 `json:"node_id"`
	OwnerAddress         string `json:"owner_address"`
	SigningKeyPub        string `json:"signing_key_pub"`
	HttpAddress          string `json:"http_address"`
	MinMonthlyFee        int64  `json:"min_monthly_fee"`
	IsReplicationEnabled bool   `json:"is_replication_enabled"`
	IsApiEnabled         bool   `json:"is_api_enabled"`
}

func ReadFromRegistry(chainCaller blockchain.INodeRegistryCaller) ([]SerializableNode, error) {
	nodes, err := chainCaller.GetAllNodes(context.Background())
	if err != nil {
		return nil, err
	}

	serializableNodes := make([]SerializableNode, len(nodes))
	for i, node := range nodes {
		owner, err := chainCaller.OwnerOf(context.Background(), node.NodeId.Int64())
		if err != nil {
			return nil, err
		}

		pubKey, err := crypto.UnmarshalPubkey(node.Node.SigningKeyPub)
		if err != nil {
			return nil, err
		}

		serializableNodes[i] = SerializableNode{
			NodeID:               uint32(node.NodeId.Int64()),
			OwnerAddress:         owner.Hex(),
			SigningKeyPub:        utils.EcdsaPublicKeyToString(pubKey),
			HttpAddress:          node.Node.HttpAddress,
			MinMonthlyFee:        node.Node.MinMonthlyFee.Int64(),
			IsReplicationEnabled: node.Node.IsReplicationEnabled,
			IsApiEnabled:         node.Node.IsApiEnabled,
		}
	}

	return serializableNodes, nil
}

func WriteToRegistry(
	logger *zap.Logger,
	nodes []SerializableNode,
	chainAdmin blockchain.INodeRegistryAdmin,
) error {
	ctx := context.Background()

	for _, node := range nodes {
		signingKey, err := utils.ParseEcdsaPublicKey(node.SigningKeyPub)
		if err != nil {
			return err
		}

		err = chainAdmin.AddNode(
			ctx,
			node.OwnerAddress,
			signingKey,
			node.HttpAddress,
			node.MinMonthlyFee,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func DumpNodesToFile(nodes []SerializableNode, outFile string) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(nodes)
}

func ImportNodesFromFile(filePath string) ([]SerializableNode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodes []SerializableNode
	err = json.NewDecoder(file).Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
