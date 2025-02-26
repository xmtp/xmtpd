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
	NodeID        uint32 `json:"node_id"`
	OwnerAddress  string `json:"owner_address"`
	SigningKeyPub string `json:"signing_key_pub"`
	HttpAddress   string `json:"http_address"`
	IsHealthy     bool   `json:"is_healthy"`
}

func ReadFromRegistry(
	chainCaller *blockchain.NodeRegistryCaller,
) ([]SerializableNode, error) {
	nodes, err := chainCaller.GetAllNodes(context.Background())
	if err != nil {
		return nil, err
	}

	serializableNodes := make([]SerializableNode, len(nodes))
	for i, node := range nodes {
		owner, err := chainCaller.OwnerOf(context.Background(), int64(node.NodeId))
		if err != nil {
			return nil, err
		}

		pubKey, err := crypto.UnmarshalPubkey(node.Node.SigningKeyPub)
		if err != nil {
			return nil, err
		}

		serializableNodes[i] = SerializableNode{
			NodeID:        node.NodeId,
			OwnerAddress:  owner.Hex(),
			SigningKeyPub: utils.EcdsaPublicKeyToString(pubKey),
			HttpAddress:   node.Node.HttpAddress,
			IsHealthy:     node.Node.IsHealthy,
		}
	}

	return serializableNodes, nil
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

func WriteToRegistry(
	logger *zap.Logger,
	nodes []SerializableNode,
	chainCaller *blockchain.NodeRegistryAdmin,
) error {
	ctx := context.Background()
	for _, node := range nodes {
		signingKey, err := utils.ParseEcdsaPublicKey(node.SigningKeyPub)
		if err != nil {
			return err
		}
		err = chainCaller.AddNode(
			ctx,
			node.OwnerAddress,
			signingKey,
			node.HttpAddress,
		)
		if err != nil {
			return err
		}

		if !node.IsHealthy {
			err = chainCaller.UpdateHealth(
				ctx,
				int64(node.NodeID),
				false,
			)
			if err != nil {
				return err
			}
		}

		logger.Info("wrote node to registry", zap.Any("node", node))
	}

	return nil
}
