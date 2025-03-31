// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package noderegistry

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// INodeRegistryNode is an auto generated low-level Go binding around an user-defined struct.
type INodeRegistryNode struct {
	SigningKeyPub             []byte
	HttpAddress               string
	InCanonicalNetwork        bool
	MinMonthlyFeeMicroDollars *big.Int
}

// INodeRegistryNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type INodeRegistryNodeWithId struct {
	NodeId *big.Int
	Node   INodeRegistryNode
}

// NodeRegistryMetaData contains all meta data concerning the NodeRegistry contract.
var NodeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_INCREMENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addToNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"allNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodeRegistry.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodeRegistry.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"inCanonicalNetwork\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsCanonicalNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isCanonicalNode\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodeRegistry.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"inCanonicalNetwork\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxActiveNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeOperatorCommissionPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeFromNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setHttpAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxActiveNodes\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinMonthlyFee\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNodeOperatorCommissionPercent\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"supported\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BaseURIUpdated\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HttpAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newHttpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxActiveNodesUpdated\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinMonthlyFeeUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAdded\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAddedToCanonicalNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeOperatorCommissionPercentUpdated\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeRemovedFromCanonicalNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"FailedToAddNodeToCanonicalNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToRemoveNodeFromCanonicalNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCommissionPercent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHttpAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigningKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesBelowCurrentCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x6080604052600a805464ffffffffff19166014179055348015610020575f5ffd5b50604051613bb7380380613bb783398101604081905261003f91610317565b60408051808201825260128152712c26aa28102737b2329027b832b930ba37b960711b602080830191909152825180840190935260048352630584d54560e41b90830152906202a300836001600160a01b0381166100b657604051636116401160e11b81525f600482015260240160405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100df5f8261018b565b50600391506100f0905083826103dc565b5060046100fd82826103dc565b5050506001600160a01b0381166101275760405163e6c4247b60e01b815260040160405180910390fd5b61013e5f516020613b975f395f51905f525f6101fa565b6101555f516020613b775f395f51905f525f6101fa565b61016c5f516020613b975f395f51905f528261018b565b506101845f516020613b775f395f51905f528261018b565b5050610496565b5f826101e7575f6101a46002546001600160a01b031690565b6001600160a01b0316146101cb57604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6101f18383610226565b90505b92915050565b8161021857604051631fe1e13d60e11b815260040160405180910390fd5b61022282826102cd565b5050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff166102c6575f838152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905561027e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016101f4565b505f6101f4565b5f82815260208190526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b5f60208284031215610327575f5ffd5b81516001600160a01b038116811461033d575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061036c57607f821691505b60208210810361038a57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156103d757805f5260205f20601f840160051c810160208510156103b55750805b601f840160051c820191505b818110156103d4575f81556001016103c1565b50505b505050565b81516001600160401b038111156103f5576103f5610344565b610409816104038454610358565b84610390565b6020601f82116001811461043b575f83156104245750848201515b5f19600385901b1c1916600184901b1784556103d4565b5f84815260208120601f198516915b8281101561046a578785015182556020948501946001909201910161044a565b508482101561048757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b6136d4806104a35f395ff3fe608060405234801561000f575f5ffd5b506004361061030f575f3560e01c806395d89b411161019d578063cefc1429116100e8578063d881666b11610093578063f947c3d01161006e578063f947c3d014610735578063fd667d1e14610748578063fd967f4714610765575f5ffd5b8063d881666b146106de578063e985e9c5146106f1578063f3194a391461072c575f5ffd5b8063d59f9fe0116100c3578063d59f9fe01461069c578063d602b9fd146106c3578063d74a2a50146106cb575f5ffd5b8063cefc142914610642578063cf6eefb71461064a578063d547741f14610689575f5ffd5b8063afcf4ad311610148578063c87b56dd11610123578063c87b56dd14610614578063cc8463c814610627578063ce9994891461062f575f5ffd5b8063afcf4ad3146105db578063b88d4fde146105ee578063c4741f3114610601575f5ffd5b8063a1eda53c11610178578063a1eda53c1461059a578063a217fddf146105c1578063a22cb465146105c8575f5ffd5b806395d89b411461055e5780639d32f9ba14610566578063a1174e7d14610585575f5ffd5b806350d0215f1161025d57806370a08231116102085780638da5cb5b116101e35780638da5cb5b1461050d5780638ed9ea341461051557806391d1485414610528575f5ffd5b806370a08231146104c257806375b238fc146104d557806384ef8ffc146104fc575f5ffd5b80636352211e116102385780636352211e14610489578063649a5ec71461049c5780636ec97bfc146104af575f5ffd5b806350d0215f1461045057806355f804b314610463578063634e93da14610476575f5ffd5b806323b872dd116102bd57806336568abe1161029857806336568abe1461040a57806342842e0e1461041d5780634f0f4aa914610430575f5ffd5b806323b872dd146103b4578063248a9ca3146103c75780632f2ff15d146103f7575f5ffd5b8063081812fc116102ed578063081812fc1461036c578063095ea7b3146103975780630aa6220b146103ac575f5ffd5b806301ffc9a714610313578063022d63fb1461033b57806306fdde0314610357575b5f5ffd5b610326610321366004612d13565b61076e565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff9091168152602001610332565b61035f6107c9565b6040516103329190612d5c565b61037f61037a366004612d6e565b610859565b6040516001600160a01b039091168152602001610332565b6103aa6103a5366004612da0565b610880565b005b6103aa61088f565b6103aa6103c2366004612dc8565b6108a4565b6103e96103d5366004612d6e565b5f9081526020819052604090206001015490565b604051908152602001610332565b6103aa610405366004612e02565b6108df565b6103aa610418366004612e02565b610920565b6103aa61042b366004612dc8565b610a10565b61044361043e366004612d6e565b610a2f565b6040516103329190612e7c565b600a54610100900463ffffffff166103e9565b6103aa610471366004612ed3565b610bba565b6103aa610484366004612f12565b610d07565b61037f610497366004612d6e565b610d1a565b6103aa6104aa366004612f2b565b610d24565b6103e96104bd366004612f50565b610d37565b6103e96104d0366004612f12565b610f9a565b6103e97fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b6002546001600160a01b031661037f565b61037f610ff8565b6103aa610523366004612fd8565b611010565b610326610536366004612e02565b5f918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b61035f6110ca565b600a546105739060ff1681565b60405160ff9091168152602001610332565b61058d6110d9565b6040516103329190612ff8565b6105a2611318565b6040805165ffffffffffff938416815292909116602083015201610332565b6103e95f81565b6103aa6105d6366004613090565b611392565b6103266105e9366004612d6e565b61139d565b6103aa6105fc3660046130f6565b6113a9565b6103aa61060f366004612d6e565b6113c1565b61035f610622366004612d6e565b61145c565b6103406114c1565b6103aa61063d3660046131d4565b61155e565b6103aa6115e5565b600154604080516001600160a01b03831681527401000000000000000000000000000000000000000090920465ffffffffffff16602083015201610332565b6103aa610697366004612e02565b611634565b6103e97fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a81565b6103aa611675565b6103aa6106d93660046131f4565b611687565b6103aa6106ec366004612d6e565b61174d565b6103266106ff36600461323c565b6001600160a01b039182165f90815260086020908152604080832093909416825291909152205460ff1690565b6103e9600e5481565b6103aa610743366004612d6e565b611804565b610750606481565b60405163ffffffff9091168152602001610332565b6103e961271081565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167fb5c09a6b0000000000000000000000000000000000000000000000000000000014806107c357506107c382611905565b92915050565b6060600380546107d890613264565b80601f016020809104026020016040519081016040528092919081815260200182805461080490613264565b801561084f5780601f106108265761010080835404028352916020019161084f565b820191905f5260205f20905b81548152906001019060200180831161083257829003601f168201915b5050505050905090565b5f610863826119a6565b505f828152600760205260409020546001600160a01b03166107c3565b61088b8282336119f7565b5050565b5f61089981611a04565b6108a1611a0e565b50565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a6108ce81611a04565b6108d9848484611a1a565b50505050565b81610916576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61088b8282611acf565b8115801561093b57506002546001600160a01b038281169116145b15610a06576001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610982575065ffffffffffff8116155b8061099557504265ffffffffffff821610155b156109db576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024015b60405180910390fd5b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b61088b8282611af3565b610a2a83838360405180602001604052805f8152506113a9565b505050565b610a5a604051806080016040528060608152602001606081526020015f151581526020015f81525090565b610a6382611b3f565b5f828152600b602052604090819020815160808101909252805482908290610a8a90613264565b80601f0160208091040260200160405190810160405280929190818152602001828054610ab690613264565b8015610b015780601f10610ad857610100808354040283529160200191610b01565b820191905f5260205f20905b815481529060010190602001808311610ae457829003601f168201915b50505050508152602001600182018054610b1a90613264565b80601f0160208091040260200160405190810160405280929190818152602001828054610b4690613264565b8015610b915780601f10610b6857610100808354040283529160200191610b91565b820191905f5260205f20905b815481529060010190602001808311610b7457829003601f168201915b5050509183525050600282015460ff161515602082015260039091015460409091015292915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610be481611a04565b81610c1b576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f2f000000000000000000000000000000000000000000000000000000000000008383610c496001826132e2565b818110610c5857610c586132f5565b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610cbb576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6009610cc8838583613366565b507f6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad8383604051610cfa929190613449565b60405180910390a1505050565b5f610d1181611a04565b61088b82611b8c565b5f6107c3826119a6565b5f610d2e81611a04565b61088b82611bfe565b5f7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610d6281611a04565b6001600160a01b038816610da2576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85610dd9576040517f8125403000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83610e10576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6064600a600181819054906101000a900463ffffffff16610e309061345c565b91906101000a81548163ffffffff021916908363ffffffff1602179055610e579190613480565b63ffffffff169150604051806080016040528088888080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250505090825250604080516020601f8901819004810282018101909252878152918101919088908890819084018382808284375f92018290525093855250505060208083018290526040928301879052858252600b90522081518190610f01908261349f565b5060208201516001820190610f16908261349f565b50604082015160028201805460ff1916911515919091179055606090910151600390910155610f458883611c66565b876001600160a01b0316827f663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c8989898989604051610f8795949392919061355a565b60405180910390a3509695505050505050565b5f6001600160a01b038216610fdd576040517f89c62b640000000000000000000000000000000000000000000000000000000081525f60048201526024016109d2565b506001600160a01b03165f9081526006602052604090205490565b5f61100b6002546001600160a01b031690565b905090565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561103a81611a04565b611044600c611cf9565b8260ff161015611080576040517f39beadee00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a805460ff191660ff84169081179091556040519081527f6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821906020015b60405180910390a15050565b6060600480546107d890613264565b600a54606090610100900463ffffffff1667ffffffffffffffff811115611102576111026130c9565b60405190808252806020026020018201604052801561113b57816020015b611128612ca5565b8152602001906001900390816111205790505b5090505f5b600a5463ffffffff61010090910481169082161015611314575f611165826001613593565b611170906064613480565b905060405180604001604052808263ffffffff168152602001600b5f8463ffffffff1681526020019081526020015f206040518060800160405290815f820180546111ba90613264565b80601f01602080910402602001604051908101604052809291908181526020018280546111e690613264565b80156112315780601f1061120857610100808354040283529160200191611231565b820191905f5260205f20905b81548152906001019060200180831161121457829003601f168201915b5050505050815260200160018201805461124a90613264565b80601f016020809104026020016040519081016040528092919081815260200182805461127690613264565b80156112c15780601f10611298576101008083540402835291602001916112c1565b820191905f5260205f20905b8154815290600101906020018083116112a457829003601f168201915b5050509183525050600282015460ff161515602082015260039091015460409091015290528351849063ffffffff8516908110611300576113006132f5565b602090810291909101015250600101611140565b5090565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561135a57504265ffffffffffff821610155b611365575f5f61138a565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b61088b338383611d02565b5f6107c3600c83611db9565b6113b48484846108a4565b6108d93385858585611dd0565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756113eb81611a04565b612710821115611427576040517f47d3b04600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e8290556040518281527f6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0906020016110be565b6060611467826119a6565b505f611471611f74565b90505f81511161148f5760405180602001604052805f8152506114ba565b8061149984611f83565b6040516020016114aa9291906135c6565b6040516020818303038152906040525b9392505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561150257504265ffffffffffff8216105b611534576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16611558565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a61158881611a04565b61159183611b3f565b5f838152600b6020526040908190206003018390555183907f27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a906115d89085815260200190565b60405180910390a2505050565b6001546001600160a01b031633811461162c576040517fc22c80220000000000000000000000000000000000000000000000000000000081523360048201526024016109d2565b6108a1612020565b8161166b576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61088b82826120f7565b5f61167f81611a04565b6108a161211b565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a6116b181611a04565b6116ba84611b3f565b816116f1576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f848152600b6020526040902060010161170c838583613366565b50837f15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed848460405161173f929190613449565b60405180910390a250505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561177781611a04565b61178082611b3f565b61178b600c83612125565b6117c1576040517fe31ff23600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f828152600b6020526040808220600201805460ff191690555183917f1b3bca5c7af55f35aad90a6fb8fcd0be0f294c332d42a01d87d47fc75f93c70691a25050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561182e81611a04565b61183782611b3f565b600a5460ff16611847600c611cf9565b1061187e576040517f950be9a500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611889600c83612130565b6118bf576040517f4992486d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f828152600b6020526040808220600201805460ff191660011790555183917ff5c33a68e71e241f24116ddc5051ad86f3d18505d210b4fc6d8235f8185a101291a25050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd00000000000000000000000000000000000000000000000000000000148061199757507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b5e139f00000000000000000000000000000000000000000000000000000000145b806107c357506107c38261213b565b5f818152600560205260408120546001600160a01b0316806107c3576040517f7e273289000000000000000000000000000000000000000000000000000000008152600481018490526024016109d2565b610a2a8383836001612190565b6108a181336122e3565b611a185f5f61234e565b565b6001600160a01b038216611a5c576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f60048201526024016109d2565b5f611a6883833361249a565b9050836001600160a01b0316816001600160a01b0316146108d9576040517f64283d7b0000000000000000000000000000000000000000000000000000000081526001600160a01b03808616600483015260248201849052821660448201526064016109d2565b5f82815260208190526040902060010154611ae981611a04565b6108d983836125a4565b6001600160a01b0381163314611b35576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a2a828261263b565b5f818152600560205260409020546001600160a01b03166108a1576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f611b956114c1565b611b9e4261268f565b611ba891906135da565b9050611bb482826126da565b60405165ffffffffffff821681526001600160a01b038316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b5f611c0882612768565b611c114261268f565b611c1b91906135da565b9050611c27828261234e565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b91016110be565b6001600160a01b038216611ca8576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f60048201526024016109d2565b5f611cb483835f61249a565b90506001600160a01b03811615610a2a576040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081525f60048201526024016109d2565b5f6107c3825490565b6001600160a01b038216611d4d576040517f5b08ba180000000000000000000000000000000000000000000000000000000081526001600160a01b03831660048201526024016109d2565b6001600160a01b038381165f81815260086020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b5f81815260018301602052604081205415156114ba565b6001600160a01b0383163b15611f6d576040517f150b7a020000000000000000000000000000000000000000000000000000000081526001600160a01b0384169063150b7a0290611e2b9088908890879087906004016135f8565b6020604051808303815f875af1925050508015611e65575060408051601f3d908101601f19168201909252611e6291810190613638565b60015b611ee5573d808015611e92576040519150601f19603f3d011682016040523d82523d5f602084013e611e97565b606091505b5080515f03611edd576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b03851660048201526024016109d2565b805181602001fd5b7fffffffff0000000000000000000000000000000000000000000000000000000081167f150b7a020000000000000000000000000000000000000000000000000000000014611f6b576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b03851660048201526024016109d2565b505b5050505050565b6060600980546107d890613264565b60605f611f8f836127af565b60010190505f8167ffffffffffffffff811115611fae57611fae6130c9565b6040519080825280601f01601f191660200182016040528015611fd8576020820181803683370190505b5090508181016020015b5f19017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a8504945084611fe257509392505050565b6001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff1680158061206357504265ffffffffffff821610155b156120a4576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024016109d2565b6120bf5f6120ba6002546001600160a01b031690565b61263b565b506120ca5f836125a4565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f8281526020819052604090206001015461211181611a04565b6108d9838361263b565b611a185f5f6126da565b5f6114ba8383612890565b5f6114ba838361297a565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f314987860000000000000000000000000000000000000000000000000000000014806107c357506107c3826129c6565b80806121a457506001600160a01b03821615155b1561229c575f6121b3846119a6565b90506001600160a01b038316158015906121df5750826001600160a01b0316816001600160a01b031614155b801561221057506001600160a01b038082165f9081526008602090815260408083209387168352929052205460ff16155b15612252576040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526001600160a01b03841660048201526024016109d2565b811561229a5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260076020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b5f828152602081815260408083206001600160a01b038516845290915290205460ff1661088b576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602481018390526044016109d2565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015612422574265ffffffffffff821610156123f9576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055612422565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b50600280546001600160a01b03167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f828152600560205260408120546001600160a01b03908116908316156124c6576124c6818486612a5c565b6001600160a01b03811615612500576124e15f855f5f612190565b6001600160a01b0381165f90815260066020526040902080545f190190555b6001600160a01b0385161561252e576001600160a01b0385165f908152600660205260409020805460010190555b5f8481526005602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b5f82612631575f6125bd6002546001600160a01b031690565b6001600160a01b0316146125fd576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384161790555b6114ba8383612af2565b5f8215801561265757506002546001600160a01b038381169116145b1561268557600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b6114ba8383612b92565b5f65ffffffffffff821115611314576040517f6dfcc65000000000000000000000000000000000000000000000000000000000815260306004820152602481018390526044016109d2565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff000000000000000000000000000000000000000000000000000084166001600160a01b03881617179093559004168015610a2a576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f6127726114c1565b90508065ffffffffffff168365ffffffffffff161161279a576127958382613653565b6114ba565b6114ba65ffffffffffff841662069780612c13565b5f807a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083106127f7577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef81000000008310612823576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc10000831061284157662386f26fc10000830492506010015b6305f5e1008310612859576305f5e100830492506008015b612710831061286d57612710830492506004015b6064831061287f576064830492506002015b600a83106107c35760010192915050565b5f818152600183016020526040812054801561296a575f6128b26001836132e2565b85549091505f906128c5906001906132e2565b9050808214612924575f865f0182815481106128e3576128e36132f5565b905f5260205f200154905080875f018481548110612903576129036132f5565b5f918252602080832090910192909255918252600188019052604090208390555b855486908061293557612935613671565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f9055600193505050506107c3565b5f9150506107c3565b5092915050565b5f8181526001830160205260408120546129bf57508154600181810184555f8481526020808220909301849055845484825282860190935260409020919091556107c3565b505f6107c3565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806107c357507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316146107c3565b612a67838383612c22565b610a2a576001600160a01b038316612aae576040517f7e273289000000000000000000000000000000000000000000000000000000008152600481018290526024016109d2565b6040517f177e802f0000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602481018290526044016109d2565b5f828152602081815260408083206001600160a01b038516845290915281205460ff166129bf575f838152602081815260408083206001600160a01b03861684529091529020805460ff19166001179055612b4a3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016107c3565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16156129bf575f838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016107c3565b5f8282188284100282186114ba565b5f6001600160a01b03831615801590612c9d5750826001600160a01b0316846001600160a01b03161480612c7a57506001600160a01b038085165f9081526008602090815260408083209387168352929052205460ff165b80612c9d57505f828152600760205260409020546001600160a01b038481169116145b949350505050565b60405180604001604052805f8152602001612ce1604051806080016040528060608152602001606081526020015f151581526020015f81525090565b905290565b7fffffffff00000000000000000000000000000000000000000000000000000000811681146108a1575f5ffd5b5f60208284031215612d23575f5ffd5b81356114ba81612ce6565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f6114ba6020830184612d2e565b5f60208284031215612d7e575f5ffd5b5035919050565b80356001600160a01b0381168114612d9b575f5ffd5b919050565b5f5f60408385031215612db1575f5ffd5b612dba83612d85565b946020939093013593505050565b5f5f5f60608486031215612dda575f5ffd5b612de384612d85565b9250612df160208501612d85565b929592945050506040919091013590565b5f5f60408385031215612e13575f5ffd5b82359150612e2360208401612d85565b90509250929050565b5f815160808452612e406080850182612d2e565b905060208301518482036020860152612e598282612d2e565b915050604083015115156040850152606083015160608501528091505092915050565b602081525f6114ba6020830184612e2c565b5f5f83601f840112612e9e575f5ffd5b50813567ffffffffffffffff811115612eb5575f5ffd5b602083019150836020828501011115612ecc575f5ffd5b9250929050565b5f5f60208385031215612ee4575f5ffd5b823567ffffffffffffffff811115612efa575f5ffd5b612f0685828601612e8e565b90969095509350505050565b5f60208284031215612f22575f5ffd5b6114ba82612d85565b5f60208284031215612f3b575f5ffd5b813565ffffffffffff811681146114ba575f5ffd5b5f5f5f5f5f5f60808789031215612f65575f5ffd5b612f6e87612d85565b9550602087013567ffffffffffffffff811115612f89575f5ffd5b612f9589828a01612e8e565b909650945050604087013567ffffffffffffffff811115612fb4575f5ffd5b612fc089828a01612e8e565b979a9699509497949695606090950135949350505050565b5f60208284031215612fe8575f5ffd5b813560ff811681146114ba575f5ffd5b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b82811015613084577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0878603018452815180518652602081015190506040602087015261306e6040870182612e2c565b955050602093840193919091019060010161301e565b50929695505050505050565b5f5f604083850312156130a1575f5ffd5b6130aa83612d85565b9150602083013580151581146130be575f5ffd5b809150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f5f5f60808587031215613109575f5ffd5b61311285612d85565b935061312060208601612d85565b925060408501359150606085013567ffffffffffffffff811115613142575f5ffd5b8501601f81018713613152575f5ffd5b803567ffffffffffffffff81111561316c5761316c6130c9565b604051601f19603f601f19601f8501160116810181811067ffffffffffffffff8211171561319c5761319c6130c9565b6040528181528282016020018910156131b3575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f604083850312156131e5575f5ffd5b50508035926020909101359150565b5f5f5f60408486031215613206575f5ffd5b83359250602084013567ffffffffffffffff811115613223575f5ffd5b61322f86828701612e8e565b9497909650939450505050565b5f5f6040838503121561324d575f5ffd5b61325683612d85565b9150612e2360208401612d85565b600181811c9082168061327857607f821691505b6020821081036132af577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b818103818111156107c3576107c36132b5565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b601f821115610a2a57805f5260205f20601f840160051c810160208510156133475750805b601f840160051c820191505b81811015611f6d575f8155600101613353565b67ffffffffffffffff83111561337e5761337e6130c9565b6133928361338c8354613264565b83613322565b5f601f8411600181146133c3575f85156133ac5750838201355b5f19600387901b1c1916600186901b178355611f6d565b5f83815260208120601f198716915b828110156133f257868501358255602094850194600190920191016133d2565b508682101561340e575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b81835281816020850137505f602082840101525f6020601f19601f840116840101905092915050565b602081525f612c9d602083018486613420565b5f63ffffffff821663ffffffff8103613477576134776132b5565b60010192915050565b63ffffffff8181168382160290811690818114612973576129736132b5565b815167ffffffffffffffff8111156134b9576134b96130c9565b6134cd816134c78454613264565b84613322565b6020601f8211600181146134ff575f83156134e85750848201515b5f19600385901b1c1916600184901b178455611f6d565b5f84815260208120601f198516915b8281101561352e578785015182556020948501946001909201910161350e565b508482101561354b57868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b606081525f61356d606083018789613420565b8281036020840152613580818688613420565b9150508260408301529695505050505050565b63ffffffff81811683821601908111156107c3576107c36132b5565b5f81518060208401855e5f93019283525090919050565b5f612c9d6135d483866135af565b846135af565b65ffffffffffff81811683821601908111156107c3576107c36132b5565b6001600160a01b03851681526001600160a01b0384166020820152826040820152608060608201525f61362e6080830184612d2e565b9695505050505050565b5f60208284031215613648575f5ffd5b81516114ba81612ce6565b65ffffffffffff82811682821603908111156107c3576107c36132b5565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffdfea26469706673582212205b93df304993af11c31f1fd9745cb4a487768f244abe291e51aafbe1a131d59f64736f6c634300081c0033daf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56aa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775",
}

// NodeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeRegistryMetaData.ABI instead.
var NodeRegistryABI = NodeRegistryMetaData.ABI

// NodeRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodeRegistryMetaData.Bin instead.
var NodeRegistryBin = NodeRegistryMetaData.Bin

// DeployNodeRegistry deploys a new Ethereum contract, binding an instance of NodeRegistry to it.
func DeployNodeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, initialAdmin common.Address) (common.Address, *types.Transaction, *NodeRegistry, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodeRegistryBin), backend, initialAdmin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NodeRegistry is an auto generated Go binding around an Ethereum contract.
type NodeRegistry struct {
	NodeRegistryCaller     // Read-only binding to the contract
	NodeRegistryTransactor // Write-only binding to the contract
	NodeRegistryFilterer   // Log filterer for contract events
}

// NodeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeRegistrySession struct {
	Contract     *NodeRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeRegistryCallerSession struct {
	Contract *NodeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// NodeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeRegistryTransactorSession struct {
	Contract     *NodeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NodeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeRegistryRaw struct {
	Contract *NodeRegistry // Generic contract binding to access the raw methods on
}

// NodeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeRegistryCallerRaw struct {
	Contract *NodeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// NodeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeRegistryTransactorRaw struct {
	Contract *NodeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeRegistry creates a new instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistry(address common.Address, backend bind.ContractBackend) (*NodeRegistry, error) {
	contract, err := bindNodeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NewNodeRegistryCaller creates a new read-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryCaller(address common.Address, caller bind.ContractCaller) (*NodeRegistryCaller, error) {
	contract, err := bindNodeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryCaller{contract: contract}, nil
}

// NewNodeRegistryTransactor creates a new write-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeRegistryTransactor, error) {
	contract, err := bindNodeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransactor{contract: contract}, nil
}

// NewNodeRegistryFilterer creates a new log filterer instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeRegistryFilterer, error) {
	contract, err := bindNodeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryFilterer{contract: contract}, nil
}

// bindNodeRegistry binds a generic wrapper to an already deployed contract.
func bindNodeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.NodeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistryCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistrySession) ADMINROLE() ([32]byte, error) {
	return _NodeRegistry.Contract.ADMINROLE(&_NodeRegistry.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistryCallerSession) ADMINROLE() ([32]byte, error) {
	return _NodeRegistry.Contract.ADMINROLE(&_NodeRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _NodeRegistry.Contract.DEFAULTADMINROLE(&_NodeRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _NodeRegistry.Contract.DEFAULTADMINROLE(&_NodeRegistry.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_NodeRegistry *NodeRegistryCaller) MAXBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "MAX_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_NodeRegistry *NodeRegistrySession) MAXBPS() (*big.Int, error) {
	return _NodeRegistry.Contract.MAXBPS(&_NodeRegistry.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_NodeRegistry *NodeRegistryCallerSession) MAXBPS() (*big.Int, error) {
	return _NodeRegistry.Contract.MAXBPS(&_NodeRegistry.CallOpts)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistryCaller) NODEINCREMENT(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "NODE_INCREMENT")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistrySession) NODEINCREMENT() (uint32, error) {
	return _NodeRegistry.Contract.NODEINCREMENT(&_NodeRegistry.CallOpts)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistryCallerSession) NODEINCREMENT() (uint32, error) {
	return _NodeRegistry.Contract.NODEINCREMENT(&_NodeRegistry.CallOpts)
}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistryCaller) NODEMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "NODE_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistrySession) NODEMANAGERROLE() ([32]byte, error) {
	return _NodeRegistry.Contract.NODEMANAGERROLE(&_NodeRegistry.CallOpts)
}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_NodeRegistry *NodeRegistryCallerSession) NODEMANAGERROLE() ([32]byte, error) {
	return _NodeRegistry.Contract.NODEMANAGERROLE(&_NodeRegistry.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistryCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistrySession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeRegistry.Contract.BalanceOf(&_NodeRegistry.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistryCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeRegistry.Contract.BalanceOf(&_NodeRegistry.CallOpts, owner)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_NodeRegistry *NodeRegistryCaller) DefaultAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "defaultAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_NodeRegistry *NodeRegistrySession) DefaultAdmin() (common.Address, error) {
	return _NodeRegistry.Contract.DefaultAdmin(&_NodeRegistry.CallOpts)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) DefaultAdmin() (common.Address, error) {
	return _NodeRegistry.Contract.DefaultAdmin(&_NodeRegistry.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_NodeRegistry *NodeRegistryCaller) DefaultAdminDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "defaultAdminDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_NodeRegistry *NodeRegistrySession) DefaultAdminDelay() (*big.Int, error) {
	return _NodeRegistry.Contract.DefaultAdminDelay(&_NodeRegistry.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_NodeRegistry *NodeRegistryCallerSession) DefaultAdminDelay() (*big.Int, error) {
	return _NodeRegistry.Contract.DefaultAdminDelay(&_NodeRegistry.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_NodeRegistry *NodeRegistryCaller) DefaultAdminDelayIncreaseWait(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "defaultAdminDelayIncreaseWait")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_NodeRegistry *NodeRegistrySession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _NodeRegistry.Contract.DefaultAdminDelayIncreaseWait(&_NodeRegistry.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_NodeRegistry *NodeRegistryCallerSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _NodeRegistry.Contract.DefaultAdminDelayIncreaseWait(&_NodeRegistry.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,uint256))[] allNodes)
func (_NodeRegistry *NodeRegistryCaller) GetAllNodes(opts *bind.CallOpts) ([]INodeRegistryNodeWithId, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getAllNodes")

	if err != nil {
		return *new([]INodeRegistryNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodeRegistryNodeWithId)).(*[]INodeRegistryNodeWithId)

	return out0, err

}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,uint256))[] allNodes)
func (_NodeRegistry *NodeRegistrySession) GetAllNodes() ([]INodeRegistryNodeWithId, error) {
	return _NodeRegistry.Contract.GetAllNodes(&_NodeRegistry.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,uint256))[] allNodes)
func (_NodeRegistry *NodeRegistryCallerSession) GetAllNodes() ([]INodeRegistryNodeWithId, error) {
	return _NodeRegistry.Contract.GetAllNodes(&_NodeRegistry.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_NodeRegistry *NodeRegistryCaller) GetAllNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getAllNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_NodeRegistry *NodeRegistrySession) GetAllNodesCount() (*big.Int, error) {
	return _NodeRegistry.Contract.GetAllNodesCount(&_NodeRegistry.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_NodeRegistry *NodeRegistryCallerSession) GetAllNodesCount() (*big.Int, error) {
	return _NodeRegistry.Contract.GetAllNodesCount(&_NodeRegistry.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistrySession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.GetApproved(&_NodeRegistry.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.GetApproved(&_NodeRegistry.CallOpts, tokenId)
}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xafcf4ad3.
//
// Solidity: function getIsCanonicalNode(uint256 nodeId) view returns(bool isCanonicalNode)
func (_NodeRegistry *NodeRegistryCaller) GetIsCanonicalNode(opts *bind.CallOpts, nodeId *big.Int) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getIsCanonicalNode", nodeId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xafcf4ad3.
//
// Solidity: function getIsCanonicalNode(uint256 nodeId) view returns(bool isCanonicalNode)
func (_NodeRegistry *NodeRegistrySession) GetIsCanonicalNode(nodeId *big.Int) (bool, error) {
	return _NodeRegistry.Contract.GetIsCanonicalNode(&_NodeRegistry.CallOpts, nodeId)
}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xafcf4ad3.
//
// Solidity: function getIsCanonicalNode(uint256 nodeId) view returns(bool isCanonicalNode)
func (_NodeRegistry *NodeRegistryCallerSession) GetIsCanonicalNode(nodeId *big.Int) (bool, error) {
	return _NodeRegistry.Contract.GetIsCanonicalNode(&_NodeRegistry.CallOpts, nodeId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,uint256) node)
func (_NodeRegistry *NodeRegistryCaller) GetNode(opts *bind.CallOpts, nodeId *big.Int) (INodeRegistryNode, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getNode", nodeId)

	if err != nil {
		return *new(INodeRegistryNode), err
	}

	out0 := *abi.ConvertType(out[0], new(INodeRegistryNode)).(*INodeRegistryNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,uint256) node)
func (_NodeRegistry *NodeRegistrySession) GetNode(nodeId *big.Int) (INodeRegistryNode, error) {
	return _NodeRegistry.Contract.GetNode(&_NodeRegistry.CallOpts, nodeId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,uint256) node)
func (_NodeRegistry *NodeRegistryCallerSession) GetNode(nodeId *big.Int) (INodeRegistryNode, error) {
	return _NodeRegistry.Contract.GetNode(&_NodeRegistry.CallOpts, nodeId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_NodeRegistry *NodeRegistryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_NodeRegistry *NodeRegistrySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _NodeRegistry.Contract.GetRoleAdmin(&_NodeRegistry.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_NodeRegistry *NodeRegistryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _NodeRegistry.Contract.GetRoleAdmin(&_NodeRegistry.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_NodeRegistry *NodeRegistryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_NodeRegistry *NodeRegistrySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _NodeRegistry.Contract.HasRole(&_NodeRegistry.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_NodeRegistry *NodeRegistryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _NodeRegistry.Contract.HasRole(&_NodeRegistry.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistryCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistrySession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodeRegistry.Contract.IsApprovedForAll(&_NodeRegistry.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistryCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodeRegistry.Contract.IsApprovedForAll(&_NodeRegistry.CallOpts, owner, operator)
}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_NodeRegistry *NodeRegistryCaller) MaxActiveNodes(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "maxActiveNodes")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_NodeRegistry *NodeRegistrySession) MaxActiveNodes() (uint8, error) {
	return _NodeRegistry.Contract.MaxActiveNodes(&_NodeRegistry.CallOpts)
}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_NodeRegistry *NodeRegistryCallerSession) MaxActiveNodes() (uint8, error) {
	return _NodeRegistry.Contract.MaxActiveNodes(&_NodeRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistryCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistrySession) Name() (string, error) {
	return _NodeRegistry.Contract.Name(&_NodeRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) Name() (string, error) {
	return _NodeRegistry.Contract.Name(&_NodeRegistry.CallOpts)
}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_NodeRegistry *NodeRegistryCaller) NodeOperatorCommissionPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "nodeOperatorCommissionPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_NodeRegistry *NodeRegistrySession) NodeOperatorCommissionPercent() (*big.Int, error) {
	return _NodeRegistry.Contract.NodeOperatorCommissionPercent(&_NodeRegistry.CallOpts)
}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_NodeRegistry *NodeRegistryCallerSession) NodeOperatorCommissionPercent() (*big.Int, error) {
	return _NodeRegistry.Contract.NodeOperatorCommissionPercent(&_NodeRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeRegistry *NodeRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeRegistry *NodeRegistrySession) Owner() (common.Address, error) {
	return _NodeRegistry.Contract.Owner(&_NodeRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) Owner() (common.Address, error) {
	return _NodeRegistry.Contract.Owner(&_NodeRegistry.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistrySession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.OwnerOf(&_NodeRegistry.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.OwnerOf(&_NodeRegistry.CallOpts, tokenId)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_NodeRegistry *NodeRegistryCaller) PendingDefaultAdmin(opts *bind.CallOpts) (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "pendingDefaultAdmin")

	outstruct := new(struct {
		NewAdmin common.Address
		Schedule *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NewAdmin = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Schedule = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_NodeRegistry *NodeRegistrySession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _NodeRegistry.Contract.PendingDefaultAdmin(&_NodeRegistry.CallOpts)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_NodeRegistry *NodeRegistryCallerSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _NodeRegistry.Contract.PendingDefaultAdmin(&_NodeRegistry.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_NodeRegistry *NodeRegistryCaller) PendingDefaultAdminDelay(opts *bind.CallOpts) (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "pendingDefaultAdminDelay")

	outstruct := new(struct {
		NewDelay *big.Int
		Schedule *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NewDelay = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Schedule = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_NodeRegistry *NodeRegistrySession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _NodeRegistry.Contract.PendingDefaultAdminDelay(&_NodeRegistry.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_NodeRegistry *NodeRegistryCallerSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _NodeRegistry.Contract.PendingDefaultAdminDelay(&_NodeRegistry.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool supported)
func (_NodeRegistry *NodeRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool supported)
func (_NodeRegistry *NodeRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodeRegistry.Contract.SupportsInterface(&_NodeRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool supported)
func (_NodeRegistry *NodeRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodeRegistry.Contract.SupportsInterface(&_NodeRegistry.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistryCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistrySession) Symbol() (string, error) {
	return _NodeRegistry.Contract.Symbol(&_NodeRegistry.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) Symbol() (string, error) {
	return _NodeRegistry.Contract.Symbol(&_NodeRegistry.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistryCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistrySession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodeRegistry.Contract.TokenURI(&_NodeRegistry.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodeRegistry.Contract.TokenURI(&_NodeRegistry.CallOpts, tokenId)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_NodeRegistry *NodeRegistryTransactor) AcceptDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "acceptDefaultAdminTransfer")
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_NodeRegistry *NodeRegistrySession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodeRegistry.Contract.AcceptDefaultAdminTransfer(&_NodeRegistry.TransactOpts)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodeRegistry.Contract.AcceptDefaultAdminTransfer(&_NodeRegistry.TransactOpts)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars) returns(uint256 nodeId)
func (_NodeRegistry *NodeRegistryTransactor) AddNode(opts *bind.TransactOpts, to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "addNode", to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars) returns(uint256 nodeId)
func (_NodeRegistry *NodeRegistrySession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddNode(&_NodeRegistry.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars) returns(uint256 nodeId)
func (_NodeRegistry *NodeRegistryTransactorSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddNode(&_NodeRegistry.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0xf947c3d0.
//
// Solidity: function addToNetwork(uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistryTransactor) AddToNetwork(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "addToNetwork", nodeId)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0xf947c3d0.
//
// Solidity: function addToNetwork(uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistrySession) AddToNetwork(nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddToNetwork(&_NodeRegistry.TransactOpts, nodeId)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0xf947c3d0.
//
// Solidity: function addToNetwork(uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) AddToNetwork(nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddToNetwork(&_NodeRegistry.TransactOpts, nodeId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.Approve(&_NodeRegistry.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.Approve(&_NodeRegistry.TransactOpts, to, tokenId)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_NodeRegistry *NodeRegistryTransactor) BeginDefaultAdminTransfer(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "beginDefaultAdminTransfer", newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_NodeRegistry *NodeRegistrySession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.BeginDefaultAdminTransfer(&_NodeRegistry.TransactOpts, newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.BeginDefaultAdminTransfer(&_NodeRegistry.TransactOpts, newAdmin)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_NodeRegistry *NodeRegistryTransactor) CancelDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "cancelDefaultAdminTransfer")
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_NodeRegistry *NodeRegistrySession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodeRegistry.Contract.CancelDefaultAdminTransfer(&_NodeRegistry.TransactOpts)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodeRegistry.Contract.CancelDefaultAdminTransfer(&_NodeRegistry.TransactOpts)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_NodeRegistry *NodeRegistryTransactor) ChangeDefaultAdminDelay(opts *bind.TransactOpts, newDelay *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "changeDefaultAdminDelay", newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_NodeRegistry *NodeRegistrySession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.ChangeDefaultAdminDelay(&_NodeRegistry.TransactOpts, newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.ChangeDefaultAdminDelay(&_NodeRegistry.TransactOpts, newDelay)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistrySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.GrantRole(&_NodeRegistry.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.GrantRole(&_NodeRegistry.TransactOpts, role, account)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0xd881666b.
//
// Solidity: function removeFromNetwork(uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistryTransactor) RemoveFromNetwork(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "removeFromNetwork", nodeId)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0xd881666b.
//
// Solidity: function removeFromNetwork(uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistrySession) RemoveFromNetwork(nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RemoveFromNetwork(&_NodeRegistry.TransactOpts, nodeId)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0xd881666b.
//
// Solidity: function removeFromNetwork(uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RemoveFromNetwork(nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RemoveFromNetwork(&_NodeRegistry.TransactOpts, nodeId)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistrySession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RenounceRole(&_NodeRegistry.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RenounceRole(&_NodeRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistrySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RevokeRole(&_NodeRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RevokeRole(&_NodeRegistry.TransactOpts, role, account)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_NodeRegistry *NodeRegistryTransactor) RollbackDefaultAdminDelay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "rollbackDefaultAdminDelay")
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_NodeRegistry *NodeRegistrySession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _NodeRegistry.Contract.RollbackDefaultAdminDelay(&_NodeRegistry.TransactOpts)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _NodeRegistry.Contract.RollbackDefaultAdminDelay(&_NodeRegistry.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistryTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistrySession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom0(&_NodeRegistry.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom0(&_NodeRegistry.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistrySession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetApprovalForAll(&_NodeRegistry.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetApprovalForAll(&_NodeRegistry.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setBaseURI", newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_NodeRegistry *NodeRegistrySession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetBaseURI(&_NodeRegistry.TransactOpts, newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetBaseURI(&_NodeRegistry.TransactOpts, newBaseURI)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetHttpAddress(opts *bind.TransactOpts, nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setHttpAddress", nodeId, httpAddress)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_NodeRegistry *NodeRegistrySession) SetHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetHttpAddress(&_NodeRegistry.TransactOpts, nodeId, httpAddress)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetHttpAddress(&_NodeRegistry.TransactOpts, nodeId, httpAddress)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetMaxActiveNodes(opts *bind.TransactOpts, newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setMaxActiveNodes", newMaxActiveNodes)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_NodeRegistry *NodeRegistrySession) SetMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetMaxActiveNodes(&_NodeRegistry.TransactOpts, newMaxActiveNodes)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetMaxActiveNodes(&_NodeRegistry.TransactOpts, newMaxActiveNodes)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetMinMonthlyFee(opts *bind.TransactOpts, nodeId *big.Int, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setMinMonthlyFee", nodeId, minMonthlyFeeMicroDollars)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) returns()
func (_NodeRegistry *NodeRegistrySession) SetMinMonthlyFee(nodeId *big.Int, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetMinMonthlyFee(&_NodeRegistry.TransactOpts, nodeId, minMonthlyFeeMicroDollars)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetMinMonthlyFee(nodeId *big.Int, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetMinMonthlyFee(&_NodeRegistry.TransactOpts, nodeId, minMonthlyFeeMicroDollars)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetNodeOperatorCommissionPercent(opts *bind.TransactOpts, newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setNodeOperatorCommissionPercent", newCommissionPercent)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_NodeRegistry *NodeRegistrySession) SetNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetNodeOperatorCommissionPercent(&_NodeRegistry.TransactOpts, newCommissionPercent)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetNodeOperatorCommissionPercent(&_NodeRegistry.TransactOpts, newCommissionPercent)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistryTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "transferFrom", from, to, nodeId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistrySession) TransferFrom(from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.TransferFrom(&_NodeRegistry.TransactOpts, from, to, nodeId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) TransferFrom(from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.TransferFrom(&_NodeRegistry.TransactOpts, from, to, nodeId)
}

// NodeRegistryApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the NodeRegistry contract.
type NodeRegistryApprovalIterator struct {
	Event *NodeRegistryApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryApproval represents a Approval event raised by the NodeRegistry contract.
type NodeRegistryApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NodeRegistryApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryApprovalIterator{contract: _NodeRegistry.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NodeRegistryApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryApproval)
				if err := _NodeRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) ParseApproval(log types.Log) (*NodeRegistryApproval, error) {
	event := new(NodeRegistryApproval)
	if err := _NodeRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the NodeRegistry contract.
type NodeRegistryApprovalForAllIterator struct {
	Event *NodeRegistryApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryApprovalForAll represents a ApprovalForAll event raised by the NodeRegistry contract.
type NodeRegistryApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*NodeRegistryApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryApprovalForAllIterator{contract: _NodeRegistry.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NodeRegistryApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryApprovalForAll)
				if err := _NodeRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) ParseApprovalForAll(log types.Log) (*NodeRegistryApprovalForAll, error) {
	event := new(NodeRegistryApprovalForAll)
	if err := _NodeRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryBaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the NodeRegistry contract.
type NodeRegistryBaseURIUpdatedIterator struct {
	Event *NodeRegistryBaseURIUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryBaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryBaseURIUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryBaseURIUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryBaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryBaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryBaseURIUpdated represents a BaseURIUpdated event raised by the NodeRegistry contract.
type NodeRegistryBaseURIUpdated struct {
	NewBaseURI string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_NodeRegistry *NodeRegistryFilterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*NodeRegistryBaseURIUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryBaseURIUpdatedIterator{contract: _NodeRegistry.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_NodeRegistry *NodeRegistryFilterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryBaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryBaseURIUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBaseURIUpdated is a log parse operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_NodeRegistry *NodeRegistryFilterer) ParseBaseURIUpdated(log types.Log) (*NodeRegistryBaseURIUpdated, error) {
	event := new(NodeRegistryBaseURIUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryDefaultAdminDelayChangeCanceledIterator is returned from FilterDefaultAdminDelayChangeCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeCanceled events raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminDelayChangeCanceledIterator struct {
	Event *NodeRegistryDefaultAdminDelayChangeCanceled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryDefaultAdminDelayChangeCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryDefaultAdminDelayChangeCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryDefaultAdminDelayChangeCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryDefaultAdminDelayChangeCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryDefaultAdminDelayChangeCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminDelayChangeCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeCanceled is a free log retrieval operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_NodeRegistry *NodeRegistryFilterer) FilterDefaultAdminDelayChangeCanceled(opts *bind.FilterOpts) (*NodeRegistryDefaultAdminDelayChangeCanceledIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryDefaultAdminDelayChangeCanceledIterator{contract: _NodeRegistry.contract, event: "DefaultAdminDelayChangeCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeCanceled is a free log subscription operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_NodeRegistry *NodeRegistryFilterer) WatchDefaultAdminDelayChangeCanceled(opts *bind.WatchOpts, sink chan<- *NodeRegistryDefaultAdminDelayChangeCanceled) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryDefaultAdminDelayChangeCanceled)
				if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDefaultAdminDelayChangeCanceled is a log parse operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_NodeRegistry *NodeRegistryFilterer) ParseDefaultAdminDelayChangeCanceled(log types.Log) (*NodeRegistryDefaultAdminDelayChangeCanceled, error) {
	event := new(NodeRegistryDefaultAdminDelayChangeCanceled)
	if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryDefaultAdminDelayChangeScheduledIterator is returned from FilterDefaultAdminDelayChangeScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeScheduled events raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminDelayChangeScheduledIterator struct {
	Event *NodeRegistryDefaultAdminDelayChangeScheduled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryDefaultAdminDelayChangeScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryDefaultAdminDelayChangeScheduled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryDefaultAdminDelayChangeScheduled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryDefaultAdminDelayChangeScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryDefaultAdminDelayChangeScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeScheduled is a free log retrieval operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_NodeRegistry *NodeRegistryFilterer) FilterDefaultAdminDelayChangeScheduled(opts *bind.FilterOpts) (*NodeRegistryDefaultAdminDelayChangeScheduledIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryDefaultAdminDelayChangeScheduledIterator{contract: _NodeRegistry.contract, event: "DefaultAdminDelayChangeScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeScheduled is a free log subscription operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_NodeRegistry *NodeRegistryFilterer) WatchDefaultAdminDelayChangeScheduled(opts *bind.WatchOpts, sink chan<- *NodeRegistryDefaultAdminDelayChangeScheduled) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryDefaultAdminDelayChangeScheduled)
				if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDefaultAdminDelayChangeScheduled is a log parse operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_NodeRegistry *NodeRegistryFilterer) ParseDefaultAdminDelayChangeScheduled(log types.Log) (*NodeRegistryDefaultAdminDelayChangeScheduled, error) {
	event := new(NodeRegistryDefaultAdminDelayChangeScheduled)
	if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryDefaultAdminTransferCanceledIterator is returned from FilterDefaultAdminTransferCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferCanceled events raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminTransferCanceledIterator struct {
	Event *NodeRegistryDefaultAdminTransferCanceled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryDefaultAdminTransferCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryDefaultAdminTransferCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryDefaultAdminTransferCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryDefaultAdminTransferCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryDefaultAdminTransferCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminTransferCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferCanceled is a free log retrieval operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_NodeRegistry *NodeRegistryFilterer) FilterDefaultAdminTransferCanceled(opts *bind.FilterOpts) (*NodeRegistryDefaultAdminTransferCanceledIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryDefaultAdminTransferCanceledIterator{contract: _NodeRegistry.contract, event: "DefaultAdminTransferCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferCanceled is a free log subscription operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_NodeRegistry *NodeRegistryFilterer) WatchDefaultAdminTransferCanceled(opts *bind.WatchOpts, sink chan<- *NodeRegistryDefaultAdminTransferCanceled) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryDefaultAdminTransferCanceled)
				if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDefaultAdminTransferCanceled is a log parse operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_NodeRegistry *NodeRegistryFilterer) ParseDefaultAdminTransferCanceled(log types.Log) (*NodeRegistryDefaultAdminTransferCanceled, error) {
	event := new(NodeRegistryDefaultAdminTransferCanceled)
	if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryDefaultAdminTransferScheduledIterator is returned from FilterDefaultAdminTransferScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferScheduled events raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminTransferScheduledIterator struct {
	Event *NodeRegistryDefaultAdminTransferScheduled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryDefaultAdminTransferScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryDefaultAdminTransferScheduled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryDefaultAdminTransferScheduled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryDefaultAdminTransferScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryDefaultAdminTransferScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the NodeRegistry contract.
type NodeRegistryDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferScheduled is a free log retrieval operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_NodeRegistry *NodeRegistryFilterer) FilterDefaultAdminTransferScheduled(opts *bind.FilterOpts, newAdmin []common.Address) (*NodeRegistryDefaultAdminTransferScheduledIterator, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryDefaultAdminTransferScheduledIterator{contract: _NodeRegistry.contract, event: "DefaultAdminTransferScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferScheduled is a free log subscription operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_NodeRegistry *NodeRegistryFilterer) WatchDefaultAdminTransferScheduled(opts *bind.WatchOpts, sink chan<- *NodeRegistryDefaultAdminTransferScheduled, newAdmin []common.Address) (event.Subscription, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryDefaultAdminTransferScheduled)
				if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDefaultAdminTransferScheduled is a log parse operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_NodeRegistry *NodeRegistryFilterer) ParseDefaultAdminTransferScheduled(log types.Log) (*NodeRegistryDefaultAdminTransferScheduled, error) {
	event := new(NodeRegistryDefaultAdminTransferScheduled)
	if err := _NodeRegistry.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryHttpAddressUpdatedIterator is returned from FilterHttpAddressUpdated and is used to iterate over the raw logs and unpacked data for HttpAddressUpdated events raised by the NodeRegistry contract.
type NodeRegistryHttpAddressUpdatedIterator struct {
	Event *NodeRegistryHttpAddressUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryHttpAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryHttpAddressUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryHttpAddressUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryHttpAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryHttpAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryHttpAddressUpdated represents a HttpAddressUpdated event raised by the NodeRegistry contract.
type NodeRegistryHttpAddressUpdated struct {
	NodeId         *big.Int
	NewHttpAddress string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterHttpAddressUpdated is a free log retrieval operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_NodeRegistry *NodeRegistryFilterer) FilterHttpAddressUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodeRegistryHttpAddressUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryHttpAddressUpdatedIterator{contract: _NodeRegistry.contract, event: "HttpAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchHttpAddressUpdated is a free log subscription operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_NodeRegistry *NodeRegistryFilterer) WatchHttpAddressUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryHttpAddressUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryHttpAddressUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHttpAddressUpdated is a log parse operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_NodeRegistry *NodeRegistryFilterer) ParseHttpAddressUpdated(log types.Log) (*NodeRegistryHttpAddressUpdated, error) {
	event := new(NodeRegistryHttpAddressUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryMaxActiveNodesUpdatedIterator is returned from FilterMaxActiveNodesUpdated and is used to iterate over the raw logs and unpacked data for MaxActiveNodesUpdated events raised by the NodeRegistry contract.
type NodeRegistryMaxActiveNodesUpdatedIterator struct {
	Event *NodeRegistryMaxActiveNodesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryMaxActiveNodesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryMaxActiveNodesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryMaxActiveNodesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryMaxActiveNodesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryMaxActiveNodesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryMaxActiveNodesUpdated represents a MaxActiveNodesUpdated event raised by the NodeRegistry contract.
type NodeRegistryMaxActiveNodesUpdated struct {
	NewMaxActiveNodes uint8
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterMaxActiveNodesUpdated is a free log retrieval operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_NodeRegistry *NodeRegistryFilterer) FilterMaxActiveNodesUpdated(opts *bind.FilterOpts) (*NodeRegistryMaxActiveNodesUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "MaxActiveNodesUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryMaxActiveNodesUpdatedIterator{contract: _NodeRegistry.contract, event: "MaxActiveNodesUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxActiveNodesUpdated is a free log subscription operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_NodeRegistry *NodeRegistryFilterer) WatchMaxActiveNodesUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryMaxActiveNodesUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "MaxActiveNodesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryMaxActiveNodesUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "MaxActiveNodesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMaxActiveNodesUpdated is a log parse operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_NodeRegistry *NodeRegistryFilterer) ParseMaxActiveNodesUpdated(log types.Log) (*NodeRegistryMaxActiveNodesUpdated, error) {
	event := new(NodeRegistryMaxActiveNodesUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "MaxActiveNodesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryMinMonthlyFeeUpdatedIterator is returned from FilterMinMonthlyFeeUpdated and is used to iterate over the raw logs and unpacked data for MinMonthlyFeeUpdated events raised by the NodeRegistry contract.
type NodeRegistryMinMonthlyFeeUpdatedIterator struct {
	Event *NodeRegistryMinMonthlyFeeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryMinMonthlyFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryMinMonthlyFeeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryMinMonthlyFeeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryMinMonthlyFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryMinMonthlyFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryMinMonthlyFeeUpdated represents a MinMonthlyFeeUpdated event raised by the NodeRegistry contract.
type NodeRegistryMinMonthlyFeeUpdated struct {
	NodeId                    *big.Int
	MinMonthlyFeeMicroDollars *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterMinMonthlyFeeUpdated is a free log retrieval operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFeeMicroDollars)
func (_NodeRegistry *NodeRegistryFilterer) FilterMinMonthlyFeeUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodeRegistryMinMonthlyFeeUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "MinMonthlyFeeUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryMinMonthlyFeeUpdatedIterator{contract: _NodeRegistry.contract, event: "MinMonthlyFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinMonthlyFeeUpdated is a free log subscription operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFeeMicroDollars)
func (_NodeRegistry *NodeRegistryFilterer) WatchMinMonthlyFeeUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryMinMonthlyFeeUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "MinMonthlyFeeUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryMinMonthlyFeeUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMinMonthlyFeeUpdated is a log parse operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFeeMicroDollars)
func (_NodeRegistry *NodeRegistryFilterer) ParseMinMonthlyFeeUpdated(log types.Log) (*NodeRegistryMinMonthlyFeeUpdated, error) {
	event := new(NodeRegistryMinMonthlyFeeUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeAddedIterator is returned from FilterNodeAdded and is used to iterate over the raw logs and unpacked data for NodeAdded events raised by the NodeRegistry contract.
type NodeRegistryNodeAddedIterator struct {
	Event *NodeRegistryNodeAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryNodeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryNodeAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryNodeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeAdded represents a NodeAdded event raised by the NodeRegistry contract.
type NodeRegistryNodeAdded struct {
	NodeId                    *big.Int
	Owner                     common.Address
	SigningKeyPub             []byte
	HttpAddress               string
	MinMonthlyFeeMicroDollars *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterNodeAdded is a free log retrieval operation binding the contract event 0x663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c.
//
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeAdded(opts *bind.FilterOpts, nodeId []*big.Int, owner []common.Address) (*NodeRegistryNodeAddedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeAdded", nodeIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeAddedIterator{contract: _NodeRegistry.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

// WatchNodeAdded is a free log subscription operation binding the contract event 0x663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c.
//
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeAdded, nodeId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeAdded", nodeIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeAdded)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeAdded is a log parse operation binding the contract event 0x663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c.
//
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeAdded(log types.Log) (*NodeRegistryNodeAdded, error) {
	event := new(NodeRegistryNodeAdded)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeAddedToCanonicalNetworkIterator is returned from FilterNodeAddedToCanonicalNetwork and is used to iterate over the raw logs and unpacked data for NodeAddedToCanonicalNetwork events raised by the NodeRegistry contract.
type NodeRegistryNodeAddedToCanonicalNetworkIterator struct {
	Event *NodeRegistryNodeAddedToCanonicalNetwork // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeAddedToCanonicalNetwork)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryNodeAddedToCanonicalNetwork)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeAddedToCanonicalNetwork represents a NodeAddedToCanonicalNetwork event raised by the NodeRegistry contract.
type NodeRegistryNodeAddedToCanonicalNetwork struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeAddedToCanonicalNetwork is a free log retrieval operation binding the contract event 0xf5c33a68e71e241f24116ddc5051ad86f3d18505d210b4fc6d8235f8185a1012.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint256 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeAddedToCanonicalNetwork(opts *bind.FilterOpts, nodeId []*big.Int) (*NodeRegistryNodeAddedToCanonicalNetworkIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeAddedToCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeAddedToCanonicalNetworkIterator{contract: _NodeRegistry.contract, event: "NodeAddedToCanonicalNetwork", logs: logs, sub: sub}, nil
}

// WatchNodeAddedToCanonicalNetwork is a free log subscription operation binding the contract event 0xf5c33a68e71e241f24116ddc5051ad86f3d18505d210b4fc6d8235f8185a1012.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint256 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeAddedToCanonicalNetwork(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeAddedToCanonicalNetwork, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeAddedToCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeAddedToCanonicalNetwork)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeAddedToCanonicalNetwork", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeAddedToCanonicalNetwork is a log parse operation binding the contract event 0xf5c33a68e71e241f24116ddc5051ad86f3d18505d210b4fc6d8235f8185a1012.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint256 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeAddedToCanonicalNetwork(log types.Log) (*NodeRegistryNodeAddedToCanonicalNetwork, error) {
	event := new(NodeRegistryNodeAddedToCanonicalNetwork)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeAddedToCanonicalNetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeOperatorCommissionPercentUpdatedIterator is returned from FilterNodeOperatorCommissionPercentUpdated and is used to iterate over the raw logs and unpacked data for NodeOperatorCommissionPercentUpdated events raised by the NodeRegistry contract.
type NodeRegistryNodeOperatorCommissionPercentUpdatedIterator struct {
	Event *NodeRegistryNodeOperatorCommissionPercentUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryNodeOperatorCommissionPercentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeOperatorCommissionPercentUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryNodeOperatorCommissionPercentUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryNodeOperatorCommissionPercentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeOperatorCommissionPercentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeOperatorCommissionPercentUpdated represents a NodeOperatorCommissionPercentUpdated event raised by the NodeRegistry contract.
type NodeRegistryNodeOperatorCommissionPercentUpdated struct {
	NewCommissionPercent *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterNodeOperatorCommissionPercentUpdated is a free log retrieval operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeOperatorCommissionPercentUpdated(opts *bind.FilterOpts) (*NodeRegistryNodeOperatorCommissionPercentUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeOperatorCommissionPercentUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeOperatorCommissionPercentUpdatedIterator{contract: _NodeRegistry.contract, event: "NodeOperatorCommissionPercentUpdated", logs: logs, sub: sub}, nil
}

// WatchNodeOperatorCommissionPercentUpdated is a free log subscription operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeOperatorCommissionPercentUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeOperatorCommissionPercentUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeOperatorCommissionPercentUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeOperatorCommissionPercentUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeOperatorCommissionPercentUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeOperatorCommissionPercentUpdated is a log parse operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeOperatorCommissionPercentUpdated(log types.Log) (*NodeRegistryNodeOperatorCommissionPercentUpdated, error) {
	event := new(NodeRegistryNodeOperatorCommissionPercentUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeOperatorCommissionPercentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeRemovedFromCanonicalNetworkIterator is returned from FilterNodeRemovedFromCanonicalNetwork and is used to iterate over the raw logs and unpacked data for NodeRemovedFromCanonicalNetwork events raised by the NodeRegistry contract.
type NodeRegistryNodeRemovedFromCanonicalNetworkIterator struct {
	Event *NodeRegistryNodeRemovedFromCanonicalNetwork // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeRemovedFromCanonicalNetwork)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryNodeRemovedFromCanonicalNetwork)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeRemovedFromCanonicalNetwork represents a NodeRemovedFromCanonicalNetwork event raised by the NodeRegistry contract.
type NodeRegistryNodeRemovedFromCanonicalNetwork struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeRemovedFromCanonicalNetwork is a free log retrieval operation binding the contract event 0x1b3bca5c7af55f35aad90a6fb8fcd0be0f294c332d42a01d87d47fc75f93c706.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint256 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeRemovedFromCanonicalNetwork(opts *bind.FilterOpts, nodeId []*big.Int) (*NodeRegistryNodeRemovedFromCanonicalNetworkIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeRemovedFromCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeRemovedFromCanonicalNetworkIterator{contract: _NodeRegistry.contract, event: "NodeRemovedFromCanonicalNetwork", logs: logs, sub: sub}, nil
}

// WatchNodeRemovedFromCanonicalNetwork is a free log subscription operation binding the contract event 0x1b3bca5c7af55f35aad90a6fb8fcd0be0f294c332d42a01d87d47fc75f93c706.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint256 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeRemovedFromCanonicalNetwork(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeRemovedFromCanonicalNetwork, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeRemovedFromCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeRemovedFromCanonicalNetwork)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeRemovedFromCanonicalNetwork", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeRemovedFromCanonicalNetwork is a log parse operation binding the contract event 0x1b3bca5c7af55f35aad90a6fb8fcd0be0f294c332d42a01d87d47fc75f93c706.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint256 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeRemovedFromCanonicalNetwork(log types.Log) (*NodeRegistryNodeRemovedFromCanonicalNetwork, error) {
	event := new(NodeRegistryNodeRemovedFromCanonicalNetwork)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeRemovedFromCanonicalNetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the NodeRegistry contract.
type NodeRegistryRoleAdminChangedIterator struct {
	Event *NodeRegistryRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the NodeRegistry contract.
type NodeRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_NodeRegistry *NodeRegistryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*NodeRegistryRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryRoleAdminChangedIterator{contract: _NodeRegistry.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_NodeRegistry *NodeRegistryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *NodeRegistryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryRoleAdminChanged)
				if err := _NodeRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_NodeRegistry *NodeRegistryFilterer) ParseRoleAdminChanged(log types.Log) (*NodeRegistryRoleAdminChanged, error) {
	event := new(NodeRegistryRoleAdminChanged)
	if err := _NodeRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the NodeRegistry contract.
type NodeRegistryRoleGrantedIterator struct {
	Event *NodeRegistryRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryRoleGranted represents a RoleGranted event raised by the NodeRegistry contract.
type NodeRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodeRegistry *NodeRegistryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*NodeRegistryRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryRoleGrantedIterator{contract: _NodeRegistry.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodeRegistry *NodeRegistryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *NodeRegistryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryRoleGranted)
				if err := _NodeRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodeRegistry *NodeRegistryFilterer) ParseRoleGranted(log types.Log) (*NodeRegistryRoleGranted, error) {
	event := new(NodeRegistryRoleGranted)
	if err := _NodeRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the NodeRegistry contract.
type NodeRegistryRoleRevokedIterator struct {
	Event *NodeRegistryRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryRoleRevoked represents a RoleRevoked event raised by the NodeRegistry contract.
type NodeRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodeRegistry *NodeRegistryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*NodeRegistryRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryRoleRevokedIterator{contract: _NodeRegistry.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodeRegistry *NodeRegistryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *NodeRegistryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryRoleRevoked)
				if err := _NodeRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodeRegistry *NodeRegistryFilterer) ParseRoleRevoked(log types.Log) (*NodeRegistryRoleRevoked, error) {
	event := new(NodeRegistryRoleRevoked)
	if err := _NodeRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the NodeRegistry contract.
type NodeRegistryTransferIterator struct {
	Event *NodeRegistryTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NodeRegistryTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NodeRegistryTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NodeRegistryTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryTransfer represents a Transfer event raised by the NodeRegistry contract.
type NodeRegistryTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NodeRegistryTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransferIterator{contract: _NodeRegistry.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NodeRegistryTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryTransfer)
				if err := _NodeRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) ParseTransfer(log types.Log) (*NodeRegistryTransfer, error) {
	event := new(NodeRegistryTransfer)
	if err := _NodeRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
