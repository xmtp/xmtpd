package migrator

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type SerializableNode struct {
	NodeID             uint32 `json:"node_id"`
	OwnerAddress       string `json:"owner_address"`
	SigningKeyPub      string `json:"signing_key_pub"`
	HTTPAddress        string `json:"http_address"`
	InCanonicalNetwork bool   `json:"in_canonical_network"`
}

func ReadFromRegistry(chainCaller blockchain.INodeRegistryCaller) ([]SerializableNode, error) {
	nodes, err := chainCaller.GetAllNodes(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve nodes from registry")
	}

	serializableNodes := make([]SerializableNode, len(nodes))
	for i, node := range nodes {
		owner, err := chainCaller.OwnerOf(context.Background(), node.NodeId)
		if err != nil {
			return nil, errors.Wrapf(err, "could not retrieve owner for node %d", node.NodeId)
		}

		pubKey, err := crypto.UnmarshalPubkey(node.Node.SigningPublicKey)
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal node signing public key")
		}

		serializableNodes[i] = SerializableNode{
			NodeID:             node.NodeId,
			OwnerAddress:       owner.Hex(),
			SigningKeyPub:      utils.EcdsaPublicKeyToString(pubKey),
			HTTPAddress:        node.Node.HttpAddress,
			InCanonicalNetwork: node.Node.IsCanonical,
		}
	}

	return serializableNodes, nil
}

func findNodeWithPubKey(oldNodes []SerializableNode, pubKey string) *uint32 {
	for _, node := range oldNodes {
		if node.SigningKeyPub == pubKey {
			return &node.NodeID
		}
	}
	return nil
}

func WriteToRegistry(
	ctx context.Context,
	newNodes []SerializableNode,
	oldNodes []SerializableNode,
	chainAdmin blockchain.INodeRegistryAdmin,
) error {
	for _, node := range newNodes {
		alreadyRegisteredNodeID := findNodeWithPubKey(oldNodes, node.SigningKeyPub)

		if alreadyRegisteredNodeID != nil {
			if node.InCanonicalNetwork {
				err := chainAdmin.AddToNetwork(ctx, *alreadyRegisteredNodeID)
				if err != nil {
					return err
				}
			}
		} else {
			signingKey, err := utils.ParseEcdsaPublicKey(node.SigningKeyPub)
			if err != nil {
				return err
			}

			ownerAddress := common.HexToAddress(node.OwnerAddress)

			nodeID, err := chainAdmin.AddNode(
				ctx,
				ownerAddress,
				signingKey,
				node.HTTPAddress,
			)
			if err != nil {
				return err
			}

			if node.InCanonicalNetwork {
				err = chainAdmin.AddToNetwork(ctx, nodeID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func DumpNodesToFile(nodes []SerializableNode, outFile string) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return json.NewEncoder(file).Encode(nodes)
}

func ImportNodesFromFile(filePath string) ([]SerializableNode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var nodes []SerializableNode
	err = json.NewDecoder(file).Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
