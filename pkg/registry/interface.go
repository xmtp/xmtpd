package registry

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
)

/*
*
A dumbed down interface of abis.NodesCaller for generating mocks
*/
type NodeRegistryContract interface {
	GetAllNodes(opts *bind.CallOpts) ([]noderegistry.INodeRegistryNodeWithId, error)
}

/*
*
The NodeRegistry is responsible for fetching the list of nodes from the registry contract
and notifying listeners when the list of nodes changes.
*/
type NodeRegistry interface {
	GetNodes() ([]Node, error)
	GetNode(uint32) (*Node, error)
	OnNewNodes() <-chan []Node
	OnChangedNode(uint32) <-chan Node
	Stop()
}
