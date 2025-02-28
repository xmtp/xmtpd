package migrator

import (
	"context"
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type SerializableNodeV1 struct {
	NodeID        uint32 `json:"node_id"`
	OwnerAddress  string `json:"owner_address"`
	SigningKeyPub string `json:"signing_key_pub"`
	HttpAddress   string `json:"http_address"`
	IsHealthy     bool   `json:"is_healthy"`
}

type SerializableNodeV2 struct {
	NodeID               uint32 `json:"node_id"`
	OwnerAddress         string `json:"owner_address"`
	SigningKeyPub        string `json:"signing_key_pub"`
	HttpAddress          string `json:"http_address"`
	MinMonthlyFee        string `json:"min_monthly_fee"`
	IsReplicationEnabled bool   `json:"is_replication_enabled"`
	IsApiEnabled         bool   `json:"is_api_enabled"`
}

func ReadFromRegistryV1(
	chainCaller *blockchain.NodeRegistryCaller,
) ([]SerializableNodeV1, error) {
	nodes, err := chainCaller.GetAllNodes(context.Background())
	if err != nil {
		return nil, err
	}

	serializableNodes := make([]SerializableNodeV1, len(nodes))
	for i, node := range nodes {
		owner, err := chainCaller.OwnerOf(context.Background(), int64(node.NodeId))
		if err != nil {
			return nil, err
		}

		pubKey, err := crypto.UnmarshalPubkey(node.Node.SigningKeyPub)
		if err != nil {
			return nil, err
		}

		serializableNodes[i] = SerializableNodeV1{
			NodeID:        node.NodeId,
			OwnerAddress:  owner.Hex(),
			SigningKeyPub: utils.EcdsaPublicKeyToString(pubKey),
			HttpAddress:   node.Node.HttpAddress,
			IsHealthy:     node.Node.IsHealthy,
		}
	}

	return serializableNodes, nil
}

func DumpNodesToFile(nodes []SerializableNodeV1, outFile string) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(nodes)
}

func ImportNodesFromFile(filePath string) ([]SerializableNodeV1, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodes []SerializableNodeV1
	err = json.NewDecoder(file).Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func ReadFromRegistryV2(
	chainCaller *blockchain.NodeRegistryCallerV2,
) ([]SerializableNodeV2, error) {
	nodes, err := chainCaller.GetAllNodes(context.Background())
	if err != nil {
		return nil, err
	}

	serializableNodes := make([]SerializableNodeV2, len(nodes))
	for i, node := range nodes {
		owner, err := chainCaller.OwnerOf(context.Background(), node.NodeId.Int64())
		if err != nil {
			return nil, err
		}

		pubKey, err := crypto.UnmarshalPubkey(node.Node.SigningKeyPub)
		if err != nil {
			return nil, err
		}

		serializableNodes[i] = SerializableNodeV2{
			NodeID:               uint32(node.NodeId.Uint64()),
			OwnerAddress:         owner.Hex(),
			SigningKeyPub:        utils.EcdsaPublicKeyToString(pubKey),
			HttpAddress:          node.Node.HttpAddress,
			MinMonthlyFee:        node.Node.MinMonthlyFee.String(),
			IsReplicationEnabled: node.Node.IsReplicationEnabled,
			IsApiEnabled:         node.Node.IsApiEnabled,
		}
	}

	return serializableNodes, nil
}

func WriteToRegistryV2(
	logger *zap.Logger,
	nodes []SerializableNodeV1,
	chainAdmin *blockchain.NodeRegistryAdminV2,
) error {
	ctx := context.Background()

	for _, node := range nodes {
		signingKey, err := utils.ParseEcdsaPublicKey(node.SigningKeyPub)
		if err != nil {
			return err
		}

		minMonthlyFee := big.NewInt(0)

		err = chainAdmin.AddNodeV2(
			ctx,
			node.OwnerAddress,
			signingKey,
			node.HttpAddress,
			minMonthlyFee,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
