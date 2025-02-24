// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nodes

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

// INodesNode is an auto generated low-level Go binding around an user-defined struct.
type INodesNode struct {
	SigningKeyPub        []byte
	HttpAddress          string
	IsReplicationEnabled bool
	IsApiEnabled         bool
	IsActive             bool
	MinMonthlyFee        *big.Int
}

// INodesNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type INodesNodeWithId struct {
	NodeId *big.Int
	Node   INodesNode
}

// NodesMetaData contains all meta data concerning the Nodes contract.
var NodesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_initialAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchUpdateActive\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"isActive\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.Node[]\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveNodesIDs\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesIDs\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"allNodesList\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeIsActive\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxActiveNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeOperatorCommissionPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateActive\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateHttpAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateIsApiEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateIsReplicationEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxActiveNodes\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMinMonthlyFee\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateNodeOperatorCommissionPercent\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ApiEnabledUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BaseURIUpdated\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HttpAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newHttpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxActiveNodesUpdated\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinMonthlyFeeUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeActivateUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAdded\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeOperatorCommissionPercentUpdated\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeTransferred\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplicationEnabledUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCommissionPercent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHttpAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInputLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigningKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesBelowCurrentCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeAlreadyActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeAlreadyInactive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x6080604052600a805464ffffffffff19166014179055348015610020575f5ffd5b5060405161460e38038061460e83398101604081905261003f91610317565b60408051808201825260128152712c26aa28102737b2329027b832b930ba37b960711b602080830191909152825180840190935260048352630584d54560e41b90830152906202a300836001600160a01b0381166100b657604051636116401160e11b81525f600482015260240160405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100df5f8261018b565b50600391506100f0905083826103dc565b5060046100fd82826103dc565b5050506001600160a01b0381166101275760405163e6c4247b60e01b815260040160405180910390fd5b61013e5f5160206145ee5f395f51905f525f6101fa565b6101555f5160206145ce5f395f51905f525f6101fa565b61016c5f5160206145ee5f395f51905f528261018b565b506101845f5160206145ce5f395f51905f528261018b565b5050610496565b5f826101e7575f6101a46002546001600160a01b031690565b6001600160a01b0316146101cb57604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6101f18383610226565b90505b92915050565b8161021857604051631fe1e13d60e11b815260040160405180910390fd5b61022282826102cd565b5050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff166102c6575f838152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905561027e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016101f4565b505f6101f4565b5f82815260208190526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b5f60208284031215610327575f5ffd5b81516001600160a01b038116811461033d575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061036c57607f821691505b60208210810361038a57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156103d757805f5260205f20601f840160051c810160208510156103b55750805b601f840160051c820191505b818110156103d4575f81556001016103c1565b50505b505050565b81516001600160401b038111156103f5576103f5610344565b610409816104038454610358565b84610390565b6020601f82116001811461043b575f83156104245750848201515b5f19600385901b1c1916600184901b1784556103d4565b5f84815260208120601f198516915b8281101561046a578785015182556020948501946001909201910161044a565b508482101561048757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b61412b806104a35f395ff3fe608060405234801561000f575f5ffd5b506004361061033b575f3560e01c80637d8389fd116101b3578063b88d4fde116100f3578063d547741f1161009e578063dac8e0eb11610079578063dac8e0eb14610775578063e985e9c514610788578063f3194a39146107c3578063fd967f47146107cc575f5ffd5b8063d547741f14610733578063d59f9fe014610746578063d602b9fd1461076d575f5ffd5b8063cefc1429116100ce578063cefc1429146106d9578063cf6eefb7146106e1578063d1a706f014610720575f5ffd5b8063b88d4fde146106ab578063c87b56dd146106be578063cc8463c8146106d1575f5ffd5b80639d32f9ba1161015e578063a217fddf11610139578063a217fddf14610669578063a22cb46514610670578063b7a8e35f14610683578063b80758cd14610698575f5ffd5b80639d32f9ba1461060e578063a1174e7d1461062d578063a1eda53c14610642575f5ffd5b806391d148541161018e57806391d14854146105bd57806394f8a324146105f357806395d89b4114610606575f5ffd5b80637d8389fd1461059c57806384ef8ffc146105a45780638da5cb5b146105b5575f5ffd5b806342842e0e1161027e578063634e93da116102295780636b51e919116102045780636b51e9191461053a5780636ec97bfc1461054f57806370a082311461056257806375b238fc14610575575f5ffd5b8063634e93da146105015780636352211e14610514578063649a5ec714610527575f5ffd5b806350d809931161025957806350d80993146104c857806355f804b3146104db5780635d04ef1c146104ee575f5ffd5b806342842e0e146104825780634f0f4aa91461049557806350d0215f146104b5575f5ffd5b8063095ea7b3116102e9578063248a9ca3116102c4578063248a9ca3146104195780632f2ff15d1461044957806335f31cfd1461045c57806336568abe1461046f575f5ffd5b8063095ea7b3146103eb5780630aa6220b146103fe57806323b872dd14610406575f5ffd5b80630549c152116103195780630549c1521461039857806306fdde03146103ab578063081812fc146103c0575f5ffd5b806301ffc9a71461033f578063022d63fb1461036757806304bb1e3d14610383575b5f5ffd5b61035261034d3660046135b6565b6107d5565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff909116815260200161035e565b6103966103913660046135e5565b6107e5565b005b6103966103a6366004613657565b6109b4565b6103b3610a74565b60405161035e91906136f1565b6103d36103ce366004613703565b610b04565b6040516001600160a01b03909116815260200161035e565b6103966103f9366004613730565b610b2b565b610396610b3a565b610396610414366004613758565b610b4f565b61043b610427366004613703565b5f9081526020819052604090206001015490565b60405190815260200161035e565b610396610457366004613792565b610bd3565b61039661046a3660046135e5565b610c14565b61039661047d366004613792565b610cfb565b610396610490366004613758565b610deb565b6104a86104a3366004613703565b610e0a565b60405161035e919061381b565b600a54610100900463ffffffff1661043b565b6103966104d636600461386b565b610ffd565b6103966104e93660046138b3565b611107565b6103526104fc366004613703565b611238565b61039661050f3660046138f2565b611244565b6103d3610522366004613703565b611257565b61039661053536600461390b565b611261565b610542611274565b60405161035e9190613930565b61043b61055d3660046139b1565b6114d9565b61043b6105703660046138f2565b6117fd565b61043b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b61043b61185b565b6002546001600160a01b03166103d3565b6103d361186b565b6103526105cb366004613792565b5f918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b610396610601366004613a39565b61187e565b6103b361193c565b600a5461061b9060ff1681565b60405160ff909116815260200161035e565b61063561194b565b60405161035e9190613a59565b61064a611bd1565b6040805165ffffffffffff93841681529290911660208301520161035e565b61043b5f81565b61039661067e366004613ae5565b611c4b565b61068b611c56565b60405161035e9190613b0d565b6103966106a6366004613703565b611c62565b6103966106b9366004613b7c565b611d05565b6103b36106cc366004613703565b611d23565b61036c611d88565b610396611e25565b600154604080516001600160a01b03831681527401000000000000000000000000000000000000000090920465ffffffffffff1660208301520161035e565b61039661072e366004613c5a565b611e74565b610396610741366004613792565b611f43565b61043b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a81565b610396611f84565b610396610783366004613703565b611f96565b610352610796366004613c7a565b6001600160a01b039182165f90815260086020908152604080832093909416825291909152205460ff1690565b61043b600e5481565b61043b61271081565b5f6107df82612074565b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561080f81612115565b5f838152600560205260409020546001600160a01b031661085c576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81156108ef57600a5460ff16610872600c61211f565b106108a9576040517f950be9a500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108b4600c84612128565b6108ea576040517f9288a2ec00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610930565b6108fa600c84612133565b610930576040517fd623cf5300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838152600b602052604090819020600201805484151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff9091161790555183907f4a518a6a9ee77b4498e418883bc42338213163021cf974718d9fe36511d6010f906109a790851515815260200190565b60405180910390a2505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756109de81612115565b838214610a17576040517f7db491eb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5b84811015610a6c57610a64868683818110610a3657610a36613ca2565b90506020020135858584818110610a4f57610a4f613ca2565b90506020020160208101906103919190613ccf565b600101610a19565b505050505050565b606060038054610a8390613ce8565b80601f0160208091040260200160405190810160405280929190818152602001828054610aaf90613ce8565b8015610afa5780601f10610ad157610100808354040283529160200191610afa565b820191905f5260205f20905b815481529060010190602001808311610add57829003601f168201915b5050505050905090565b5f610b0e8261213e565b505f828152600760205260409020546001600160a01b03166107df565b610b3682823361218f565b5050565b5f610b4481612115565b610b4c61219c565b50565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a610b7981612115565b610b82826121a8565b610b8d84848461222b565b826001600160a01b0316846001600160a01b0316837e80108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be7260405160405180910390a450505050565b81610c0a576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b3682826122e0565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a610c3e81612115565b5f838152600560205260409020546001600160a01b0316610c8b576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838152600b602090815260409182902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016851515908117909155915191825284917fda2a657fb74ca331fb64eaabeca91a4ab0c68fd7ce7a8938a1a709903cf9be1e91016109a7565b81158015610d1657506002546001600160a01b038281169116145b15610de1576001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610d5d575065ffffffffffff8116155b80610d7057504265ffffffffffff821610155b15610db6576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024015b60405180910390fd5b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b610b368282612304565b610e0583838360405180602001604052805f815250611d05565b505050565b6040805160c081018252606080825260208083018290525f8385018190529183018290526080830182905260a083018290528482526005905291909120546001600160a01b0316610e87576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f828152600b602052604090819020815160c08101909252805482908290610eae90613ce8565b80601f0160208091040260200160405190810160405280929190818152602001828054610eda90613ce8565b8015610f255780601f10610efc57610100808354040283529160200191610f25565b820191905f5260205f20905b815481529060010190602001808311610f0857829003601f168201915b50505050508152602001600182018054610f3e90613ce8565b80601f0160208091040260200160405190810160405280929190818152602001828054610f6a90613ce8565b8015610fb55780601f10610f8c57610100808354040283529160200191610fb5565b820191905f5260205f20905b815481529060010190602001808311610f9857829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015292915050565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a61102781612115565b5f848152600560205260409020546001600160a01b0316611074576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816110ab576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f848152600b602052604090206001016110c6838583613d7d565b50837f15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed84846040516110f9929190613e60565b60405180910390a250505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561113181612115565b81611168576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8282611175600182613ea0565b81811061118457611184613ca2565b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916602f60f81b146111ec576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60096111f9838583613d7d565b507f6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad838360405161122b929190613e60565b60405180910390a1505050565b5f6107df600c83612350565b5f61124e81612115565b610b3682612367565b5f6107df8261213e565b5f61126b81612115565b610b36826123d9565b6060611280600c61211f565b67ffffffffffffffff81111561129857611298613b4f565b6040519080825280602002602001820160405280156112fa57816020015b6040805160c0810182526060808252602082018190525f92820183905281018290526080810182905260a08101919091528152602001906001900390816112b65790505b5090505f5b611309600c61211f565b8163ffffffff1610156114d557600b5f61132d600c63ffffffff8086169061244116565b81526020019081526020015f206040518060c00160405290815f8201805461135490613ce8565b80601f016020809104026020016040519081016040528092919081815260200182805461138090613ce8565b80156113cb5780601f106113a2576101008083540402835291602001916113cb565b820191905f5260205f20905b8154815290600101906020018083116113ae57829003601f168201915b505050505081526020016001820180546113e490613ce8565b80601f016020809104026020016040519081016040528092919081815260200182805461141090613ce8565b801561145b5780601f106114325761010080835404028352916020019161145b565b820191905f5260205f20905b81548152906001019060200180831161143e57829003601f168201915b5050509183525050600282015460ff8082161515602084015261010082048116151560408401526201000090910416151560608201526003909101546080909101528251839063ffffffff84169081106114b7576114b7613ca2565b602002602001018190525080806114cd90613eb3565b9150506112ff565b5090565b5f7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561150481612115565b6001600160a01b038816611544576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8561157b576040517f8125403000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836115b2576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a8054610100900463ffffffff169060016115cd83613eb3565b91906101000a81548163ffffffff021916908363ffffffff160217905550505f6064600a60019054906101000a900463ffffffff1661160c9190613ed7565b905061161e898263ffffffff1661244c565b6040518060c0016040528089898080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250505090825250604080516020601f8a01819004810282018101909252888152918101919089908990819084018382808284375f9201829052509385525050506020808301829052604080840183905260608401839052608090930188905263ffffffff85168252600b905220815181906116d59082613ef6565b50602082015160018201906116ea9082613ef6565b5060408281015160028301805460608601516080870151151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff951515959095167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090931692909217939093179290921691909117905560a090920151600390910155516001600160a01b038a169063ffffffff8316907f663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c906117e3908c908c908c908c908c90613fb1565b60405180910390a363ffffffff1698975050505050505050565b5f6001600160a01b038216611840576040517f89c62b640000000000000000000000000000000000000000000000000000000081525f6004820152602401610dad565b506001600160a01b03165f9081526006602052604090205490565b5f611866600c61211f565b905090565b5f6118666002546001600160a01b031690565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a6118a881612115565b5f838152600560205260409020546001600160a01b03166118f5576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838152600b6020526040908190206003018390555183907f27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a906109a79085815260200190565b606060048054610a8390613ce8565b600a54606090610100900463ffffffff1667ffffffffffffffff81111561197457611974613b4f565b6040519080825280602002602001820160405280156119ad57816020015b61199a613538565b8152602001906001900390816119925790505b5090505f5b600a5463ffffffff610100909104811690821610156114d5575f6119d7826001613fea565b6119e2906064613ed7565b9050611a0a8163ffffffff165f908152600560205260409020546001600160a01b0316151590565b15611bc85760405180604001604052808263ffffffff168152602001600b5f8463ffffffff1681526020019081526020015f206040518060c00160405290815f82018054611a5790613ce8565b80601f0160208091040260200160405190810160405280929190818152602001828054611a8390613ce8565b8015611ace5780601f10611aa557610100808354040283529160200191611ace565b820191905f5260205f20905b815481529060010190602001808311611ab157829003601f168201915b50505050508152602001600182018054611ae790613ce8565b80601f0160208091040260200160405190810160405280929190818152602001828054611b1390613ce8565b8015611b5e5780601f10611b3557610100808354040283529160200191611b5e565b820191905f5260205f20905b815481529060010190602001808311611b4157829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff8516908110611bbc57611bbc613ca2565b60200260200101819052505b506001016119b2565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015611c1357504265ffffffffffff821610155b611c1e575f5f611c43565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b610b363383836124df565b6060611866600c6125b4565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611c8c81612115565b612710821115611cc8576040517f47d3b04600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e8290556040518281527f6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0906020015b60405180910390a15050565b611d10848484610b4f565b611d1d33858585856125c0565b50505050565b6060611d2e8261213e565b505f611d38612762565b90505f815111611d565760405180602001604052805f815250611d81565b80611d6084612771565b604051602001611d7192919061401d565b6040516020818303038152906040525b9392505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015611dc957504265ffffffffffff8216105b611dfb576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16611e1f565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b6001546001600160a01b0316338114611e6c576040517fc22c8022000000000000000000000000000000000000000000000000000000008152336004820152602401610dad565b610b4c61280e565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611e9e81612115565b611ea8600c61211f565b8260ff1611611ee3576040517f39beadee00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff84169081179091556040519081527f6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d82190602001611cf9565b81611f7a576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b3682826128e5565b5f611f8e81612115565b610b4c612909565b5f818152600560205260409020546001600160a01b03163314611fe5576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f818152600b602052604090819020600201805460ff61010080830482161581027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff9093169290921792839055925184937f9ee9d27515d2b2213be541ebd78d4500491c27c656ec4bc67eec934a57d3876893612069939004161515815260200190565b60405180910390a250565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd00000000000000000000000000000000000000000000000000000000148061210657507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b5e139f00000000000000000000000000000000000000000000000000000000145b806107df57506107df82612913565b610b4c8133612968565b5f6107df825490565b5f611d8183836129d3565b5f611d818383612a1f565b5f818152600560205260408120546001600160a01b0316806107df576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101849052602401610dad565b610e058383836001612b09565b6121a65f5f612c5c565b565b6121b3600c82612350565b15610b4c576121c3600c82612133565b505f818152600b6020908152604080832060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff1690555191825282917f4a518a6a9ee77b4498e418883bc42338213163021cf974718d9fe36511d6010f9101612069565b6001600160a01b03821661226d576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610dad565b5f612279838333612da8565b9050836001600160a01b0316816001600160a01b031614611d1d576040517f64283d7b0000000000000000000000000000000000000000000000000000000081526001600160a01b0380861660048301526024820184905282166044820152606401610dad565b5f828152602081905260409020600101546122fa81612115565b611d1d8383612eb2565b6001600160a01b0381163314612346576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e058282612f49565b5f8181526001830160205260408120541515611d81565b5f612370611d88565b61237942612f9d565b6123839190614031565b905061238f8282612fe8565b60405165ffffffffffff821681526001600160a01b038316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b5f6123e382613076565b6123ec42612f9d565b6123f69190614031565b90506124028282612c5c565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b9101611cf9565b5f611d8183836130bd565b6001600160a01b03821661248e576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610dad565b5f61249a83835f612da8565b90506001600160a01b03811615610e05576040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610dad565b6001600160a01b03821661252a576040517f5b08ba180000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610dad565b6001600160a01b038381165f8181526008602090815260408083209487168084529482529182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b60605f611d81836130e3565b6001600160a01b0383163b1561275b576040517f150b7a020000000000000000000000000000000000000000000000000000000081526001600160a01b0384169063150b7a029061261b90889088908790879060040161404f565b6020604051808303815f875af1925050508015612655575060408051601f3d908101601f191682019092526126529181019061408f565b60015b6126d5573d808015612682576040519150601f19603f3d011682016040523d82523d5f602084013e612687565b606091505b5080515f036126cd576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610dad565b805181602001fd5b7fffffffff0000000000000000000000000000000000000000000000000000000081167f150b7a020000000000000000000000000000000000000000000000000000000014610a6c576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610dad565b5050505050565b606060098054610a8390613ce8565b60605f61277d8361313c565b60010190505f8167ffffffffffffffff81111561279c5761279c613b4f565b6040519080825280601f01601f1916602001820160405280156127c6576020820181803683370190505b5090508181016020015b5f19017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85049450846127d057509392505050565b6001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff1680158061285157504265ffffffffffff821610155b15612892576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610dad565b6128ad5f6128a86002546001600160a01b031690565b612f49565b506128b85f83612eb2565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f828152602081905260409020600101546128ff81612115565b611d1d8383612f49565b6121a65f5f612fe8565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f314987860000000000000000000000000000000000000000000000000000000014806107df57506107df8261321d565b5f828152602081815260408083206001600160a01b038516845290915290205460ff16610b36576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b038216600482015260248101839052604401610dad565b5f818152600183016020526040812054612a1857508154600181810184555f8481526020808220909301849055845484825282860190935260409020919091556107df565b505f6107df565b5f8181526001830160205260408120548015612af9575f612a41600183613ea0565b85549091505f90612a5490600190613ea0565b9050808214612ab3575f865f018281548110612a7257612a72613ca2565b905f5260205f200154905080875f018481548110612a9257612a92613ca2565b5f918252602080832090910192909255918252600188019052604090208390555b8554869080612ac457612ac46140aa565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f9055600193505050506107df565b5f9150506107df565b5092915050565b8080612b1d57506001600160a01b03821615155b15612c15575f612b2c8461213e565b90506001600160a01b03831615801590612b585750826001600160a01b0316816001600160a01b031614155b8015612b8957506001600160a01b038082165f9081526008602090815260408083209387168352929052205460ff16155b15612bcb576040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526001600160a01b0384166004820152602401610dad565b8115612c135783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260076020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015612d30574265ffffffffffff82161015612d07576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055612d30565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b50600280546001600160a01b03167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f828152600560205260408120546001600160a01b0390811690831615612dd457612dd48184866132b3565b6001600160a01b03811615612e0e57612def5f855f5f612b09565b6001600160a01b0381165f90815260066020526040902080545f190190555b6001600160a01b03851615612e3c576001600160a01b0385165f908152600660205260409020805460010190555b5f8481526005602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b5f82612f3f575f612ecb6002546001600160a01b031690565b6001600160a01b031614612f0b576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384161790555b611d818383613349565b5f82158015612f6557506002546001600160a01b038381169116145b15612f9357600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b611d818383613407565b5f65ffffffffffff8211156114d5576040517f6dfcc6500000000000000000000000000000000000000000000000000000000081526030600482015260248101839052604401610dad565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff000000000000000000000000000000000000000000000000000084166001600160a01b03881617179093559004168015610e05576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f613080611d88565b90508065ffffffffffff168365ffffffffffff16116130a8576130a383826140d7565b611d81565b611d8165ffffffffffff8416620697806134a6565b5f825f0182815481106130d2576130d2613ca2565b905f5260205f200154905092915050565b6060815f0180548060200260200160405190810160405280929190818152602001828054801561313057602002820191905f5260205f20905b81548152602001906001019080831161311c575b50505050509050919050565b5f807a184f03e93ff9f4daa797ed6e38ed64bf6a1f0100000000000000008310613184577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef810000000083106131b0576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc1000083106131ce57662386f26fc10000830492506010015b6305f5e10083106131e6576305f5e100830492506008015b61271083106131fa57612710830492506004015b6064831061320c576064830492506002015b600a83106107df5760010192915050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806107df57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316146107df565b6132be8383836134b5565b610e05576001600160a01b038316613305576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101829052602401610dad565b6040517f177e802f0000000000000000000000000000000000000000000000000000000081526001600160a01b038316600482015260248101829052604401610dad565b5f828152602081815260408083206001600160a01b038516845290915281205460ff16612a18575f838152602081815260408083206001600160a01b0386168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556133bf3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016107df565b5f828152602081815260408083206001600160a01b038516845290915281205460ff1615612a18575f838152602081815260408083206001600160a01b038616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016107df565b5f828218828410028218611d81565b5f6001600160a01b038316158015906135305750826001600160a01b0316846001600160a01b0316148061350d57506001600160a01b038085165f9081526008602090815260408083209387168352929052205460ff165b8061353057505f828152600760205260409020546001600160a01b038481169116145b949350505050565b60405180604001604052805f81526020016135846040518060c0016040528060608152602001606081526020015f151581526020015f151581526020015f151581526020015f81525090565b905290565b7fffffffff0000000000000000000000000000000000000000000000000000000081168114610b4c575f5ffd5b5f602082840312156135c6575f5ffd5b8135611d8181613589565b803580151581146135e0575f5ffd5b919050565b5f5f604083850312156135f6575f5ffd5b82359150613606602084016135d1565b90509250929050565b5f5f83601f84011261361f575f5ffd5b50813567ffffffffffffffff811115613636575f5ffd5b6020830191508360208260051b8501011115613650575f5ffd5b9250929050565b5f5f5f5f6040858703121561366a575f5ffd5b843567ffffffffffffffff811115613680575f5ffd5b61368c8782880161360f565b909550935050602085013567ffffffffffffffff8111156136ab575f5ffd5b6136b78782880161360f565b95989497509550505050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f611d8160208301846136c3565b5f60208284031215613713575f5ffd5b5035919050565b80356001600160a01b03811681146135e0575f5ffd5b5f5f60408385031215613741575f5ffd5b61374a8361371a565b946020939093013593505050565b5f5f5f6060848603121561376a575f5ffd5b6137738461371a565b92506137816020850161371a565b929592945050506040919091013590565b5f5f604083850312156137a3575f5ffd5b823591506136066020840161371a565b5f815160c084526137c760c08501826136c3565b9050602083015184820360208601526137e082826136c3565b91505060408301511515604085015260608301511515606085015260808301511515608085015260a083015160a08501528091505092915050565b602081525f611d8160208301846137b3565b5f5f83601f84011261383d575f5ffd5b50813567ffffffffffffffff811115613854575f5ffd5b602083019150836020828501011115613650575f5ffd5b5f5f5f6040848603121561387d575f5ffd5b83359250602084013567ffffffffffffffff81111561389a575f5ffd5b6138a68682870161382d565b9497909650939450505050565b5f5f602083850312156138c4575f5ffd5b823567ffffffffffffffff8111156138da575f5ffd5b6138e68582860161382d565b90969095509350505050565b5f60208284031215613902575f5ffd5b611d818261371a565b5f6020828403121561391b575f5ffd5b813565ffffffffffff81168114611d81575f5ffd5b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b828110156139a5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08786030184526139908583516137b3565b94506020938401939190910190600101613956565b50929695505050505050565b5f5f5f5f5f5f608087890312156139c6575f5ffd5b6139cf8761371a565b9550602087013567ffffffffffffffff8111156139ea575f5ffd5b6139f689828a0161382d565b909650945050604087013567ffffffffffffffff811115613a15575f5ffd5b613a2189828a0161382d565b979a9699509497949695606090950135949350505050565b5f5f60408385031215613a4a575f5ffd5b50508035926020909101359150565b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b828110156139a5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08786030184528151805186526020810151905060406020870152613acf60408701826137b3565b9550506020938401939190910190600101613a7f565b5f5f60408385031215613af6575f5ffd5b613aff8361371a565b9150613606602084016135d1565b602080825282518282018190525f918401906040840190835b81811015613b44578351835260209384019390920191600101613b26565b509095945050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f5f5f60808587031215613b8f575f5ffd5b613b988561371a565b9350613ba66020860161371a565b925060408501359150606085013567ffffffffffffffff811115613bc8575f5ffd5b8501601f81018713613bd8575f5ffd5b803567ffffffffffffffff811115613bf257613bf2613b4f565b604051601f19603f601f19601f8501160116810181811067ffffffffffffffff82111715613c2257613c22613b4f565b604052818152828201602001891015613c39575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f60208284031215613c6a575f5ffd5b813560ff81168114611d81575f5ffd5b5f5f60408385031215613c8b575f5ffd5b613c948361371a565b91506136066020840161371a565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f60208284031215613cdf575f5ffd5b611d81826135d1565b600181811c90821680613cfc57607f821691505b602082108103613d33577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b601f821115610e0557805f5260205f20601f840160051c81016020851015613d5e5750805b601f840160051c820191505b8181101561275b575f8155600101613d6a565b67ffffffffffffffff831115613d9557613d95613b4f565b613da983613da38354613ce8565b83613d39565b5f601f841160018114613dda575f8515613dc35750838201355b5f19600387901b1c1916600186901b17835561275b565b5f83815260208120601f198716915b82811015613e095786850135825560209485019460019092019101613de9565b5086821015613e25575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b81835281816020850137505f602082840101525f6020601f19601f840116840101905092915050565b602081525f613530602083018486613e37565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b818103818111156107df576107df613e73565b5f63ffffffff821663ffffffff8103613ece57613ece613e73565b60010192915050565b63ffffffff8181168382160290811690818114612b0257612b02613e73565b815167ffffffffffffffff811115613f1057613f10613b4f565b613f2481613f1e8454613ce8565b84613d39565b6020601f821160018114613f56575f8315613f3f5750848201515b5f19600385901b1c1916600184901b17845561275b565b5f84815260208120601f198516915b82811015613f855787850151825560209485019460019092019101613f65565b5084821015613fa257868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b606081525f613fc4606083018789613e37565b8281036020840152613fd7818688613e37565b9150508260408301529695505050505050565b63ffffffff81811683821601908111156107df576107df613e73565b5f81518060208401855e5f93019283525090919050565b5f61353061402b8386614006565b84614006565b65ffffffffffff81811683821601908111156107df576107df613e73565b6001600160a01b03851681526001600160a01b0384166020820152826040820152608060608201525f61408560808301846136c3565b9695505050505050565b5f6020828403121561409f575f5ffd5b8151611d8181613589565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffd5b65ffffffffffff82811682821603908111156107df576107df613e7356fea26469706673582212201dae73aeb9e1f66e20a46167ce65b28be563e6475543ca975f18353107af4bad64736f6c634300081c0033daf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56aa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775",
}

// NodesABI is the input ABI used to generate the binding from.
// Deprecated: Use NodesMetaData.ABI instead.
var NodesABI = NodesMetaData.ABI

// NodesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodesMetaData.Bin instead.
var NodesBin = NodesMetaData.Bin

// DeployNodes deploys a new Ethereum contract, binding an instance of Nodes to it.
func DeployNodes(auth *bind.TransactOpts, backend bind.ContractBackend, _initialAdmin common.Address) (common.Address, *types.Transaction, *Nodes, error) {
	parsed, err := NodesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodesBin), backend, _initialAdmin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Nodes{NodesCaller: NodesCaller{contract: contract}, NodesTransactor: NodesTransactor{contract: contract}, NodesFilterer: NodesFilterer{contract: contract}}, nil
}

// Nodes is an auto generated Go binding around an Ethereum contract.
type Nodes struct {
	NodesCaller     // Read-only binding to the contract
	NodesTransactor // Write-only binding to the contract
	NodesFilterer   // Log filterer for contract events
}

// NodesCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodesSession struct {
	Contract     *Nodes            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodesCallerSession struct {
	Contract *NodesCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// NodesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodesTransactorSession struct {
	Contract     *NodesTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodesRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodesRaw struct {
	Contract *Nodes // Generic contract binding to access the raw methods on
}

// NodesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodesCallerRaw struct {
	Contract *NodesCaller // Generic read-only contract binding to access the raw methods on
}

// NodesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodesTransactorRaw struct {
	Contract *NodesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodes creates a new instance of Nodes, bound to a specific deployed contract.
func NewNodes(address common.Address, backend bind.ContractBackend) (*Nodes, error) {
	contract, err := bindNodes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Nodes{NodesCaller: NodesCaller{contract: contract}, NodesTransactor: NodesTransactor{contract: contract}, NodesFilterer: NodesFilterer{contract: contract}}, nil
}

// NewNodesCaller creates a new read-only instance of Nodes, bound to a specific deployed contract.
func NewNodesCaller(address common.Address, caller bind.ContractCaller) (*NodesCaller, error) {
	contract, err := bindNodes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodesCaller{contract: contract}, nil
}

// NewNodesTransactor creates a new write-only instance of Nodes, bound to a specific deployed contract.
func NewNodesTransactor(address common.Address, transactor bind.ContractTransactor) (*NodesTransactor, error) {
	contract, err := bindNodes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodesTransactor{contract: contract}, nil
}

// NewNodesFilterer creates a new log filterer instance of Nodes, bound to a specific deployed contract.
func NewNodesFilterer(address common.Address, filterer bind.ContractFilterer) (*NodesFilterer, error) {
	contract, err := bindNodes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodesFilterer{contract: contract}, nil
}

// bindNodes binds a generic wrapper to an already deployed contract.
func bindNodes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nodes *NodesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nodes.Contract.NodesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nodes *NodesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.Contract.NodesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nodes *NodesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nodes.Contract.NodesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nodes *NodesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nodes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nodes *NodesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nodes *NodesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nodes.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Nodes *NodesCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Nodes *NodesSession) ADMINROLE() ([32]byte, error) {
	return _Nodes.Contract.ADMINROLE(&_Nodes.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Nodes *NodesCallerSession) ADMINROLE() ([32]byte, error) {
	return _Nodes.Contract.ADMINROLE(&_Nodes.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Nodes *NodesCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Nodes *NodesSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Nodes.Contract.DEFAULTADMINROLE(&_Nodes.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Nodes *NodesCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Nodes.Contract.DEFAULTADMINROLE(&_Nodes.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_Nodes *NodesCaller) MAXBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "MAX_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_Nodes *NodesSession) MAXBPS() (*big.Int, error) {
	return _Nodes.Contract.MAXBPS(&_Nodes.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_Nodes *NodesCallerSession) MAXBPS() (*big.Int, error) {
	return _Nodes.Contract.MAXBPS(&_Nodes.CallOpts)
}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_Nodes *NodesCaller) NODEMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "NODE_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_Nodes *NodesSession) NODEMANAGERROLE() ([32]byte, error) {
	return _Nodes.Contract.NODEMANAGERROLE(&_Nodes.CallOpts)
}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_Nodes *NodesCallerSession) NODEMANAGERROLE() ([32]byte, error) {
	return _Nodes.Contract.NODEMANAGERROLE(&_Nodes.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Nodes *NodesCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Nodes *NodesSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Nodes.Contract.BalanceOf(&_Nodes.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Nodes *NodesCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Nodes.Contract.BalanceOf(&_Nodes.CallOpts, owner)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_Nodes *NodesCaller) DefaultAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "defaultAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_Nodes *NodesSession) DefaultAdmin() (common.Address, error) {
	return _Nodes.Contract.DefaultAdmin(&_Nodes.CallOpts)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_Nodes *NodesCallerSession) DefaultAdmin() (common.Address, error) {
	return _Nodes.Contract.DefaultAdmin(&_Nodes.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_Nodes *NodesCaller) DefaultAdminDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "defaultAdminDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_Nodes *NodesSession) DefaultAdminDelay() (*big.Int, error) {
	return _Nodes.Contract.DefaultAdminDelay(&_Nodes.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_Nodes *NodesCallerSession) DefaultAdminDelay() (*big.Int, error) {
	return _Nodes.Contract.DefaultAdminDelay(&_Nodes.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_Nodes *NodesCaller) DefaultAdminDelayIncreaseWait(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "defaultAdminDelayIncreaseWait")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_Nodes *NodesSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _Nodes.Contract.DefaultAdminDelayIncreaseWait(&_Nodes.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_Nodes *NodesCallerSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _Nodes.Contract.DefaultAdminDelayIncreaseWait(&_Nodes.CallOpts)
}

// GetActiveNodes is a free data retrieval call binding the contract method 0x6b51e919.
//
// Solidity: function getActiveNodes() view returns((bytes,string,bool,bool,bool,uint256)[] activeNodes)
func (_Nodes *NodesCaller) GetActiveNodes(opts *bind.CallOpts) ([]INodesNode, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveNodes")

	if err != nil {
		return *new([]INodesNode), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNode)).(*[]INodesNode)

	return out0, err

}

// GetActiveNodes is a free data retrieval call binding the contract method 0x6b51e919.
//
// Solidity: function getActiveNodes() view returns((bytes,string,bool,bool,bool,uint256)[] activeNodes)
func (_Nodes *NodesSession) GetActiveNodes() ([]INodesNode, error) {
	return _Nodes.Contract.GetActiveNodes(&_Nodes.CallOpts)
}

// GetActiveNodes is a free data retrieval call binding the contract method 0x6b51e919.
//
// Solidity: function getActiveNodes() view returns((bytes,string,bool,bool,bool,uint256)[] activeNodes)
func (_Nodes *NodesCallerSession) GetActiveNodes() ([]INodesNode, error) {
	return _Nodes.Contract.GetActiveNodes(&_Nodes.CallOpts)
}

// GetActiveNodesCount is a free data retrieval call binding the contract method 0x7d8389fd.
//
// Solidity: function getActiveNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesCaller) GetActiveNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveNodesCount is a free data retrieval call binding the contract method 0x7d8389fd.
//
// Solidity: function getActiveNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesSession) GetActiveNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetActiveNodesCount(&_Nodes.CallOpts)
}

// GetActiveNodesCount is a free data retrieval call binding the contract method 0x7d8389fd.
//
// Solidity: function getActiveNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesCallerSession) GetActiveNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetActiveNodesCount(&_Nodes.CallOpts)
}

// GetActiveNodesIDs is a free data retrieval call binding the contract method 0xb7a8e35f.
//
// Solidity: function getActiveNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesCaller) GetActiveNodesIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveNodesIDs")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveNodesIDs is a free data retrieval call binding the contract method 0xb7a8e35f.
//
// Solidity: function getActiveNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesSession) GetActiveNodesIDs() ([]*big.Int, error) {
	return _Nodes.Contract.GetActiveNodesIDs(&_Nodes.CallOpts)
}

// GetActiveNodesIDs is a free data retrieval call binding the contract method 0xb7a8e35f.
//
// Solidity: function getActiveNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesCallerSession) GetActiveNodesIDs() ([]*big.Int, error) {
	return _Nodes.Contract.GetActiveNodesIDs(&_Nodes.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodesList)
func (_Nodes *NodesCaller) GetAllNodes(opts *bind.CallOpts) ([]INodesNodeWithId, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getAllNodes")

	if err != nil {
		return *new([]INodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNodeWithId)).(*[]INodesNodeWithId)

	return out0, err

}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodesList)
func (_Nodes *NodesSession) GetAllNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetAllNodes(&_Nodes.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodesList)
func (_Nodes *NodesCallerSession) GetAllNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetAllNodes(&_Nodes.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_Nodes *NodesCaller) GetAllNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getAllNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_Nodes *NodesSession) GetAllNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetAllNodesCount(&_Nodes.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_Nodes *NodesCallerSession) GetAllNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetAllNodesCount(&_Nodes.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Nodes *NodesCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Nodes *NodesSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Nodes.Contract.GetApproved(&_Nodes.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Nodes *NodesCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Nodes.Contract.GetApproved(&_Nodes.CallOpts, tokenId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,bool,bool,uint256) node)
func (_Nodes *NodesCaller) GetNode(opts *bind.CallOpts, nodeId *big.Int) (INodesNode, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getNode", nodeId)

	if err != nil {
		return *new(INodesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(INodesNode)).(*INodesNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,bool,bool,uint256) node)
func (_Nodes *NodesSession) GetNode(nodeId *big.Int) (INodesNode, error) {
	return _Nodes.Contract.GetNode(&_Nodes.CallOpts, nodeId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,bool,bool,uint256) node)
func (_Nodes *NodesCallerSession) GetNode(nodeId *big.Int) (INodesNode, error) {
	return _Nodes.Contract.GetNode(&_Nodes.CallOpts, nodeId)
}

// GetNodeIsActive is a free data retrieval call binding the contract method 0x5d04ef1c.
//
// Solidity: function getNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesCaller) GetNodeIsActive(opts *bind.CallOpts, nodeId *big.Int) (bool, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getNodeIsActive", nodeId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetNodeIsActive is a free data retrieval call binding the contract method 0x5d04ef1c.
//
// Solidity: function getNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesSession) GetNodeIsActive(nodeId *big.Int) (bool, error) {
	return _Nodes.Contract.GetNodeIsActive(&_Nodes.CallOpts, nodeId)
}

// GetNodeIsActive is a free data retrieval call binding the contract method 0x5d04ef1c.
//
// Solidity: function getNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesCallerSession) GetNodeIsActive(nodeId *big.Int) (bool, error) {
	return _Nodes.Contract.GetNodeIsActive(&_Nodes.CallOpts, nodeId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Nodes *NodesCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Nodes *NodesSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Nodes.Contract.GetRoleAdmin(&_Nodes.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Nodes *NodesCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Nodes.Contract.GetRoleAdmin(&_Nodes.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Nodes *NodesCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Nodes *NodesSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Nodes.Contract.HasRole(&_Nodes.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Nodes *NodesCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Nodes.Contract.HasRole(&_Nodes.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Nodes *NodesCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Nodes *NodesSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Nodes.Contract.IsApprovedForAll(&_Nodes.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Nodes *NodesCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Nodes.Contract.IsApprovedForAll(&_Nodes.CallOpts, owner, operator)
}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_Nodes *NodesCaller) MaxActiveNodes(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "maxActiveNodes")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_Nodes *NodesSession) MaxActiveNodes() (uint8, error) {
	return _Nodes.Contract.MaxActiveNodes(&_Nodes.CallOpts)
}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_Nodes *NodesCallerSession) MaxActiveNodes() (uint8, error) {
	return _Nodes.Contract.MaxActiveNodes(&_Nodes.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nodes *NodesCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nodes *NodesSession) Name() (string, error) {
	return _Nodes.Contract.Name(&_Nodes.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nodes *NodesCallerSession) Name() (string, error) {
	return _Nodes.Contract.Name(&_Nodes.CallOpts)
}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_Nodes *NodesCaller) NodeOperatorCommissionPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "nodeOperatorCommissionPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_Nodes *NodesSession) NodeOperatorCommissionPercent() (*big.Int, error) {
	return _Nodes.Contract.NodeOperatorCommissionPercent(&_Nodes.CallOpts)
}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_Nodes *NodesCallerSession) NodeOperatorCommissionPercent() (*big.Int, error) {
	return _Nodes.Contract.NodeOperatorCommissionPercent(&_Nodes.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nodes *NodesCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nodes *NodesSession) Owner() (common.Address, error) {
	return _Nodes.Contract.Owner(&_Nodes.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nodes *NodesCallerSession) Owner() (common.Address, error) {
	return _Nodes.Contract.Owner(&_Nodes.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Nodes *NodesCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Nodes *NodesSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Nodes.Contract.OwnerOf(&_Nodes.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Nodes *NodesCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Nodes.Contract.OwnerOf(&_Nodes.CallOpts, tokenId)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_Nodes *NodesCaller) PendingDefaultAdmin(opts *bind.CallOpts) (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "pendingDefaultAdmin")

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
func (_Nodes *NodesSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _Nodes.Contract.PendingDefaultAdmin(&_Nodes.CallOpts)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_Nodes *NodesCallerSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _Nodes.Contract.PendingDefaultAdmin(&_Nodes.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_Nodes *NodesCaller) PendingDefaultAdminDelay(opts *bind.CallOpts) (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "pendingDefaultAdminDelay")

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
func (_Nodes *NodesSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _Nodes.Contract.PendingDefaultAdminDelay(&_Nodes.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_Nodes *NodesCallerSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _Nodes.Contract.PendingDefaultAdminDelay(&_Nodes.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nodes *NodesCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nodes *NodesSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nodes.Contract.SupportsInterface(&_Nodes.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nodes *NodesCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nodes.Contract.SupportsInterface(&_Nodes.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Nodes *NodesCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Nodes *NodesSession) Symbol() (string, error) {
	return _Nodes.Contract.Symbol(&_Nodes.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Nodes *NodesCallerSession) Symbol() (string, error) {
	return _Nodes.Contract.Symbol(&_Nodes.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Nodes *NodesCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Nodes *NodesSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Nodes.Contract.TokenURI(&_Nodes.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Nodes *NodesCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Nodes.Contract.TokenURI(&_Nodes.CallOpts, tokenId)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_Nodes *NodesTransactor) AcceptDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "acceptDefaultAdminTransfer")
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_Nodes *NodesSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _Nodes.Contract.AcceptDefaultAdminTransfer(&_Nodes.TransactOpts)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_Nodes *NodesTransactorSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _Nodes.Contract.AcceptDefaultAdminTransfer(&_Nodes.TransactOpts)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee) returns(uint256)
func (_Nodes *NodesTransactor) AddNode(opts *bind.TransactOpts, to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "addNode", to, signingKeyPub, httpAddress, minMonthlyFee)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee) returns(uint256)
func (_Nodes *NodesSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFee)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee) returns(uint256)
func (_Nodes *NodesTransactorSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFee)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Nodes *NodesTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Nodes *NodesSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.Approve(&_Nodes.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Nodes *NodesTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.Approve(&_Nodes.TransactOpts, to, tokenId)
}

// BatchUpdateActive is a paid mutator transaction binding the contract method 0x0549c152.
//
// Solidity: function batchUpdateActive(uint256[] nodeIds, bool[] isActive) returns()
func (_Nodes *NodesTransactor) BatchUpdateActive(opts *bind.TransactOpts, nodeIds []*big.Int, isActive []bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "batchUpdateActive", nodeIds, isActive)
}

// BatchUpdateActive is a paid mutator transaction binding the contract method 0x0549c152.
//
// Solidity: function batchUpdateActive(uint256[] nodeIds, bool[] isActive) returns()
func (_Nodes *NodesSession) BatchUpdateActive(nodeIds []*big.Int, isActive []bool) (*types.Transaction, error) {
	return _Nodes.Contract.BatchUpdateActive(&_Nodes.TransactOpts, nodeIds, isActive)
}

// BatchUpdateActive is a paid mutator transaction binding the contract method 0x0549c152.
//
// Solidity: function batchUpdateActive(uint256[] nodeIds, bool[] isActive) returns()
func (_Nodes *NodesTransactorSession) BatchUpdateActive(nodeIds []*big.Int, isActive []bool) (*types.Transaction, error) {
	return _Nodes.Contract.BatchUpdateActive(&_Nodes.TransactOpts, nodeIds, isActive)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_Nodes *NodesTransactor) BeginDefaultAdminTransfer(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "beginDefaultAdminTransfer", newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_Nodes *NodesSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.BeginDefaultAdminTransfer(&_Nodes.TransactOpts, newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_Nodes *NodesTransactorSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.BeginDefaultAdminTransfer(&_Nodes.TransactOpts, newAdmin)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_Nodes *NodesTransactor) CancelDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "cancelDefaultAdminTransfer")
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_Nodes *NodesSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _Nodes.Contract.CancelDefaultAdminTransfer(&_Nodes.TransactOpts)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_Nodes *NodesTransactorSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _Nodes.Contract.CancelDefaultAdminTransfer(&_Nodes.TransactOpts)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_Nodes *NodesTransactor) ChangeDefaultAdminDelay(opts *bind.TransactOpts, newDelay *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "changeDefaultAdminDelay", newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_Nodes *NodesSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.ChangeDefaultAdminDelay(&_Nodes.TransactOpts, newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_Nodes *NodesTransactorSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.ChangeDefaultAdminDelay(&_Nodes.TransactOpts, newDelay)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Nodes *NodesTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Nodes *NodesSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.GrantRole(&_Nodes.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Nodes *NodesTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.GrantRole(&_Nodes.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Nodes *NodesTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Nodes *NodesSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.RenounceRole(&_Nodes.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Nodes *NodesTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.RenounceRole(&_Nodes.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Nodes *NodesTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Nodes *NodesSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.RevokeRole(&_Nodes.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Nodes *NodesTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.RevokeRole(&_Nodes.TransactOpts, role, account)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_Nodes *NodesTransactor) RollbackDefaultAdminDelay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "rollbackDefaultAdminDelay")
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_Nodes *NodesSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _Nodes.Contract.RollbackDefaultAdminDelay(&_Nodes.TransactOpts)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_Nodes *NodesTransactorSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _Nodes.Contract.RollbackDefaultAdminDelay(&_Nodes.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Nodes *NodesTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Nodes *NodesSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.SafeTransferFrom(&_Nodes.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Nodes *NodesTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.SafeTransferFrom(&_Nodes.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Nodes *NodesTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Nodes *NodesSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Nodes.Contract.SafeTransferFrom0(&_Nodes.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Nodes *NodesTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Nodes.Contract.SafeTransferFrom0(&_Nodes.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nodes *NodesTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nodes *NodesSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nodes.Contract.SetApprovalForAll(&_Nodes.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nodes *NodesTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nodes.Contract.SetApprovalForAll(&_Nodes.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_Nodes *NodesTransactor) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setBaseURI", newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_Nodes *NodesSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _Nodes.Contract.SetBaseURI(&_Nodes.TransactOpts, newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_Nodes *NodesTransactorSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _Nodes.Contract.SetBaseURI(&_Nodes.TransactOpts, newBaseURI)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_Nodes *NodesTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "transferFrom", from, to, nodeId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_Nodes *NodesSession) TransferFrom(from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.TransferFrom(&_Nodes.TransactOpts, from, to, nodeId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_Nodes *NodesTransactorSession) TransferFrom(from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.TransferFrom(&_Nodes.TransactOpts, from, to, nodeId)
}

// UpdateActive is a paid mutator transaction binding the contract method 0x04bb1e3d.
//
// Solidity: function updateActive(uint256 nodeId, bool isActive) returns()
func (_Nodes *NodesTransactor) UpdateActive(opts *bind.TransactOpts, nodeId *big.Int, isActive bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateActive", nodeId, isActive)
}

// UpdateActive is a paid mutator transaction binding the contract method 0x04bb1e3d.
//
// Solidity: function updateActive(uint256 nodeId, bool isActive) returns()
func (_Nodes *NodesSession) UpdateActive(nodeId *big.Int, isActive bool) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateActive(&_Nodes.TransactOpts, nodeId, isActive)
}

// UpdateActive is a paid mutator transaction binding the contract method 0x04bb1e3d.
//
// Solidity: function updateActive(uint256 nodeId, bool isActive) returns()
func (_Nodes *NodesTransactorSession) UpdateActive(nodeId *big.Int, isActive bool) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateActive(&_Nodes.TransactOpts, nodeId, isActive)
}

// UpdateHttpAddress is a paid mutator transaction binding the contract method 0x50d80993.
//
// Solidity: function updateHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_Nodes *NodesTransactor) UpdateHttpAddress(opts *bind.TransactOpts, nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateHttpAddress", nodeId, httpAddress)
}

// UpdateHttpAddress is a paid mutator transaction binding the contract method 0x50d80993.
//
// Solidity: function updateHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_Nodes *NodesSession) UpdateHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateHttpAddress(&_Nodes.TransactOpts, nodeId, httpAddress)
}

// UpdateHttpAddress is a paid mutator transaction binding the contract method 0x50d80993.
//
// Solidity: function updateHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_Nodes *NodesTransactorSession) UpdateHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateHttpAddress(&_Nodes.TransactOpts, nodeId, httpAddress)
}

// UpdateIsApiEnabled is a paid mutator transaction binding the contract method 0xdac8e0eb.
//
// Solidity: function updateIsApiEnabled(uint256 nodeId) returns()
func (_Nodes *NodesTransactor) UpdateIsApiEnabled(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateIsApiEnabled", nodeId)
}

// UpdateIsApiEnabled is a paid mutator transaction binding the contract method 0xdac8e0eb.
//
// Solidity: function updateIsApiEnabled(uint256 nodeId) returns()
func (_Nodes *NodesSession) UpdateIsApiEnabled(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateIsApiEnabled(&_Nodes.TransactOpts, nodeId)
}

// UpdateIsApiEnabled is a paid mutator transaction binding the contract method 0xdac8e0eb.
//
// Solidity: function updateIsApiEnabled(uint256 nodeId) returns()
func (_Nodes *NodesTransactorSession) UpdateIsApiEnabled(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateIsApiEnabled(&_Nodes.TransactOpts, nodeId)
}

// UpdateIsReplicationEnabled is a paid mutator transaction binding the contract method 0x35f31cfd.
//
// Solidity: function updateIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_Nodes *NodesTransactor) UpdateIsReplicationEnabled(opts *bind.TransactOpts, nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateIsReplicationEnabled", nodeId, isReplicationEnabled)
}

// UpdateIsReplicationEnabled is a paid mutator transaction binding the contract method 0x35f31cfd.
//
// Solidity: function updateIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_Nodes *NodesSession) UpdateIsReplicationEnabled(nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateIsReplicationEnabled(&_Nodes.TransactOpts, nodeId, isReplicationEnabled)
}

// UpdateIsReplicationEnabled is a paid mutator transaction binding the contract method 0x35f31cfd.
//
// Solidity: function updateIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_Nodes *NodesTransactorSession) UpdateIsReplicationEnabled(nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateIsReplicationEnabled(&_Nodes.TransactOpts, nodeId, isReplicationEnabled)
}

// UpdateMaxActiveNodes is a paid mutator transaction binding the contract method 0xd1a706f0.
//
// Solidity: function updateMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_Nodes *NodesTransactor) UpdateMaxActiveNodes(opts *bind.TransactOpts, newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateMaxActiveNodes", newMaxActiveNodes)
}

// UpdateMaxActiveNodes is a paid mutator transaction binding the contract method 0xd1a706f0.
//
// Solidity: function updateMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_Nodes *NodesSession) UpdateMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateMaxActiveNodes(&_Nodes.TransactOpts, newMaxActiveNodes)
}

// UpdateMaxActiveNodes is a paid mutator transaction binding the contract method 0xd1a706f0.
//
// Solidity: function updateMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_Nodes *NodesTransactorSession) UpdateMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateMaxActiveNodes(&_Nodes.TransactOpts, newMaxActiveNodes)
}

// UpdateMinMonthlyFee is a paid mutator transaction binding the contract method 0x94f8a324.
//
// Solidity: function updateMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) returns()
func (_Nodes *NodesTransactor) UpdateMinMonthlyFee(opts *bind.TransactOpts, nodeId *big.Int, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateMinMonthlyFee", nodeId, minMonthlyFee)
}

// UpdateMinMonthlyFee is a paid mutator transaction binding the contract method 0x94f8a324.
//
// Solidity: function updateMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) returns()
func (_Nodes *NodesSession) UpdateMinMonthlyFee(nodeId *big.Int, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateMinMonthlyFee(&_Nodes.TransactOpts, nodeId, minMonthlyFee)
}

// UpdateMinMonthlyFee is a paid mutator transaction binding the contract method 0x94f8a324.
//
// Solidity: function updateMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) returns()
func (_Nodes *NodesTransactorSession) UpdateMinMonthlyFee(nodeId *big.Int, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateMinMonthlyFee(&_Nodes.TransactOpts, nodeId, minMonthlyFee)
}

// UpdateNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xb80758cd.
//
// Solidity: function updateNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_Nodes *NodesTransactor) UpdateNodeOperatorCommissionPercent(opts *bind.TransactOpts, newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateNodeOperatorCommissionPercent", newCommissionPercent)
}

// UpdateNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xb80758cd.
//
// Solidity: function updateNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_Nodes *NodesSession) UpdateNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateNodeOperatorCommissionPercent(&_Nodes.TransactOpts, newCommissionPercent)
}

// UpdateNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xb80758cd.
//
// Solidity: function updateNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_Nodes *NodesTransactorSession) UpdateNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateNodeOperatorCommissionPercent(&_Nodes.TransactOpts, newCommissionPercent)
}

// NodesApiEnabledUpdatedIterator is returned from FilterApiEnabledUpdated and is used to iterate over the raw logs and unpacked data for ApiEnabledUpdated events raised by the Nodes contract.
type NodesApiEnabledUpdatedIterator struct {
	Event *NodesApiEnabledUpdated // Event containing the contract specifics and raw log

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
func (it *NodesApiEnabledUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesApiEnabledUpdated)
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
		it.Event = new(NodesApiEnabledUpdated)
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
func (it *NodesApiEnabledUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesApiEnabledUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesApiEnabledUpdated represents a ApiEnabledUpdated event raised by the Nodes contract.
type NodesApiEnabledUpdated struct {
	NodeId       *big.Int
	IsApiEnabled bool
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterApiEnabledUpdated is a free log retrieval operation binding the contract event 0x9ee9d27515d2b2213be541ebd78d4500491c27c656ec4bc67eec934a57d38768.
//
// Solidity: event ApiEnabledUpdated(uint256 indexed nodeId, bool isApiEnabled)
func (_Nodes *NodesFilterer) FilterApiEnabledUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesApiEnabledUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ApiEnabledUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesApiEnabledUpdatedIterator{contract: _Nodes.contract, event: "ApiEnabledUpdated", logs: logs, sub: sub}, nil
}

// WatchApiEnabledUpdated is a free log subscription operation binding the contract event 0x9ee9d27515d2b2213be541ebd78d4500491c27c656ec4bc67eec934a57d38768.
//
// Solidity: event ApiEnabledUpdated(uint256 indexed nodeId, bool isApiEnabled)
func (_Nodes *NodesFilterer) WatchApiEnabledUpdated(opts *bind.WatchOpts, sink chan<- *NodesApiEnabledUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ApiEnabledUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesApiEnabledUpdated)
				if err := _Nodes.contract.UnpackLog(event, "ApiEnabledUpdated", log); err != nil {
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

// ParseApiEnabledUpdated is a log parse operation binding the contract event 0x9ee9d27515d2b2213be541ebd78d4500491c27c656ec4bc67eec934a57d38768.
//
// Solidity: event ApiEnabledUpdated(uint256 indexed nodeId, bool isApiEnabled)
func (_Nodes *NodesFilterer) ParseApiEnabledUpdated(log types.Log) (*NodesApiEnabledUpdated, error) {
	event := new(NodesApiEnabledUpdated)
	if err := _Nodes.contract.UnpackLog(event, "ApiEnabledUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Nodes contract.
type NodesApprovalIterator struct {
	Event *NodesApproval // Event containing the contract specifics and raw log

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
func (it *NodesApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesApproval)
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
		it.Event = new(NodesApproval)
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
func (it *NodesApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesApproval represents a Approval event raised by the Nodes contract.
type NodesApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Nodes *NodesFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NodesApprovalIterator, error) {

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

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesApprovalIterator{contract: _Nodes.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Nodes *NodesFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NodesApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesApproval)
				if err := _Nodes.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseApproval(log types.Log) (*NodesApproval, error) {
	event := new(NodesApproval)
	if err := _Nodes.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Nodes contract.
type NodesApprovalForAllIterator struct {
	Event *NodesApprovalForAll // Event containing the contract specifics and raw log

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
func (it *NodesApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesApprovalForAll)
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
		it.Event = new(NodesApprovalForAll)
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
func (it *NodesApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesApprovalForAll represents a ApprovalForAll event raised by the Nodes contract.
type NodesApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Nodes *NodesFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*NodesApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NodesApprovalForAllIterator{contract: _Nodes.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Nodes *NodesFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NodesApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesApprovalForAll)
				if err := _Nodes.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseApprovalForAll(log types.Log) (*NodesApprovalForAll, error) {
	event := new(NodesApprovalForAll)
	if err := _Nodes.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesBaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the Nodes contract.
type NodesBaseURIUpdatedIterator struct {
	Event *NodesBaseURIUpdated // Event containing the contract specifics and raw log

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
func (it *NodesBaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesBaseURIUpdated)
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
		it.Event = new(NodesBaseURIUpdated)
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
func (it *NodesBaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesBaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesBaseURIUpdated represents a BaseURIUpdated event raised by the Nodes contract.
type NodesBaseURIUpdated struct {
	NewBaseURI string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_Nodes *NodesFilterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*NodesBaseURIUpdatedIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesBaseURIUpdatedIterator{contract: _Nodes.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_Nodes *NodesFilterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *NodesBaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesBaseURIUpdated)
				if err := _Nodes.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseBaseURIUpdated(log types.Log) (*NodesBaseURIUpdated, error) {
	event := new(NodesBaseURIUpdated)
	if err := _Nodes.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesDefaultAdminDelayChangeCanceledIterator is returned from FilterDefaultAdminDelayChangeCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeCanceled events raised by the Nodes contract.
type NodesDefaultAdminDelayChangeCanceledIterator struct {
	Event *NodesDefaultAdminDelayChangeCanceled // Event containing the contract specifics and raw log

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
func (it *NodesDefaultAdminDelayChangeCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesDefaultAdminDelayChangeCanceled)
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
		it.Event = new(NodesDefaultAdminDelayChangeCanceled)
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
func (it *NodesDefaultAdminDelayChangeCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesDefaultAdminDelayChangeCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the Nodes contract.
type NodesDefaultAdminDelayChangeCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeCanceled is a free log retrieval operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_Nodes *NodesFilterer) FilterDefaultAdminDelayChangeCanceled(opts *bind.FilterOpts) (*NodesDefaultAdminDelayChangeCanceledIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return &NodesDefaultAdminDelayChangeCanceledIterator{contract: _Nodes.contract, event: "DefaultAdminDelayChangeCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeCanceled is a free log subscription operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_Nodes *NodesFilterer) WatchDefaultAdminDelayChangeCanceled(opts *bind.WatchOpts, sink chan<- *NodesDefaultAdminDelayChangeCanceled) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesDefaultAdminDelayChangeCanceled)
				if err := _Nodes.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseDefaultAdminDelayChangeCanceled(log types.Log) (*NodesDefaultAdminDelayChangeCanceled, error) {
	event := new(NodesDefaultAdminDelayChangeCanceled)
	if err := _Nodes.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesDefaultAdminDelayChangeScheduledIterator is returned from FilterDefaultAdminDelayChangeScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeScheduled events raised by the Nodes contract.
type NodesDefaultAdminDelayChangeScheduledIterator struct {
	Event *NodesDefaultAdminDelayChangeScheduled // Event containing the contract specifics and raw log

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
func (it *NodesDefaultAdminDelayChangeScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesDefaultAdminDelayChangeScheduled)
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
		it.Event = new(NodesDefaultAdminDelayChangeScheduled)
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
func (it *NodesDefaultAdminDelayChangeScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesDefaultAdminDelayChangeScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the Nodes contract.
type NodesDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeScheduled is a free log retrieval operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_Nodes *NodesFilterer) FilterDefaultAdminDelayChangeScheduled(opts *bind.FilterOpts) (*NodesDefaultAdminDelayChangeScheduledIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return &NodesDefaultAdminDelayChangeScheduledIterator{contract: _Nodes.contract, event: "DefaultAdminDelayChangeScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeScheduled is a free log subscription operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_Nodes *NodesFilterer) WatchDefaultAdminDelayChangeScheduled(opts *bind.WatchOpts, sink chan<- *NodesDefaultAdminDelayChangeScheduled) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesDefaultAdminDelayChangeScheduled)
				if err := _Nodes.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseDefaultAdminDelayChangeScheduled(log types.Log) (*NodesDefaultAdminDelayChangeScheduled, error) {
	event := new(NodesDefaultAdminDelayChangeScheduled)
	if err := _Nodes.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesDefaultAdminTransferCanceledIterator is returned from FilterDefaultAdminTransferCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferCanceled events raised by the Nodes contract.
type NodesDefaultAdminTransferCanceledIterator struct {
	Event *NodesDefaultAdminTransferCanceled // Event containing the contract specifics and raw log

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
func (it *NodesDefaultAdminTransferCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesDefaultAdminTransferCanceled)
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
		it.Event = new(NodesDefaultAdminTransferCanceled)
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
func (it *NodesDefaultAdminTransferCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesDefaultAdminTransferCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the Nodes contract.
type NodesDefaultAdminTransferCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferCanceled is a free log retrieval operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_Nodes *NodesFilterer) FilterDefaultAdminTransferCanceled(opts *bind.FilterOpts) (*NodesDefaultAdminTransferCanceledIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return &NodesDefaultAdminTransferCanceledIterator{contract: _Nodes.contract, event: "DefaultAdminTransferCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferCanceled is a free log subscription operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_Nodes *NodesFilterer) WatchDefaultAdminTransferCanceled(opts *bind.WatchOpts, sink chan<- *NodesDefaultAdminTransferCanceled) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesDefaultAdminTransferCanceled)
				if err := _Nodes.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseDefaultAdminTransferCanceled(log types.Log) (*NodesDefaultAdminTransferCanceled, error) {
	event := new(NodesDefaultAdminTransferCanceled)
	if err := _Nodes.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesDefaultAdminTransferScheduledIterator is returned from FilterDefaultAdminTransferScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferScheduled events raised by the Nodes contract.
type NodesDefaultAdminTransferScheduledIterator struct {
	Event *NodesDefaultAdminTransferScheduled // Event containing the contract specifics and raw log

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
func (it *NodesDefaultAdminTransferScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesDefaultAdminTransferScheduled)
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
		it.Event = new(NodesDefaultAdminTransferScheduled)
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
func (it *NodesDefaultAdminTransferScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesDefaultAdminTransferScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the Nodes contract.
type NodesDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferScheduled is a free log retrieval operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_Nodes *NodesFilterer) FilterDefaultAdminTransferScheduled(opts *bind.FilterOpts, newAdmin []common.Address) (*NodesDefaultAdminTransferScheduledIterator, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return &NodesDefaultAdminTransferScheduledIterator{contract: _Nodes.contract, event: "DefaultAdminTransferScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferScheduled is a free log subscription operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_Nodes *NodesFilterer) WatchDefaultAdminTransferScheduled(opts *bind.WatchOpts, sink chan<- *NodesDefaultAdminTransferScheduled, newAdmin []common.Address) (event.Subscription, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesDefaultAdminTransferScheduled)
				if err := _Nodes.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseDefaultAdminTransferScheduled(log types.Log) (*NodesDefaultAdminTransferScheduled, error) {
	event := new(NodesDefaultAdminTransferScheduled)
	if err := _Nodes.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesHttpAddressUpdatedIterator is returned from FilterHttpAddressUpdated and is used to iterate over the raw logs and unpacked data for HttpAddressUpdated events raised by the Nodes contract.
type NodesHttpAddressUpdatedIterator struct {
	Event *NodesHttpAddressUpdated // Event containing the contract specifics and raw log

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
func (it *NodesHttpAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesHttpAddressUpdated)
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
		it.Event = new(NodesHttpAddressUpdated)
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
func (it *NodesHttpAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesHttpAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesHttpAddressUpdated represents a HttpAddressUpdated event raised by the Nodes contract.
type NodesHttpAddressUpdated struct {
	NodeId         *big.Int
	NewHttpAddress string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterHttpAddressUpdated is a free log retrieval operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_Nodes *NodesFilterer) FilterHttpAddressUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesHttpAddressUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesHttpAddressUpdatedIterator{contract: _Nodes.contract, event: "HttpAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchHttpAddressUpdated is a free log subscription operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_Nodes *NodesFilterer) WatchHttpAddressUpdated(opts *bind.WatchOpts, sink chan<- *NodesHttpAddressUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesHttpAddressUpdated)
				if err := _Nodes.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseHttpAddressUpdated(log types.Log) (*NodesHttpAddressUpdated, error) {
	event := new(NodesHttpAddressUpdated)
	if err := _Nodes.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesMaxActiveNodesUpdatedIterator is returned from FilterMaxActiveNodesUpdated and is used to iterate over the raw logs and unpacked data for MaxActiveNodesUpdated events raised by the Nodes contract.
type NodesMaxActiveNodesUpdatedIterator struct {
	Event *NodesMaxActiveNodesUpdated // Event containing the contract specifics and raw log

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
func (it *NodesMaxActiveNodesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesMaxActiveNodesUpdated)
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
		it.Event = new(NodesMaxActiveNodesUpdated)
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
func (it *NodesMaxActiveNodesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesMaxActiveNodesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesMaxActiveNodesUpdated represents a MaxActiveNodesUpdated event raised by the Nodes contract.
type NodesMaxActiveNodesUpdated struct {
	NewMaxActiveNodes uint8
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterMaxActiveNodesUpdated is a free log retrieval operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_Nodes *NodesFilterer) FilterMaxActiveNodesUpdated(opts *bind.FilterOpts) (*NodesMaxActiveNodesUpdatedIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "MaxActiveNodesUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesMaxActiveNodesUpdatedIterator{contract: _Nodes.contract, event: "MaxActiveNodesUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxActiveNodesUpdated is a free log subscription operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_Nodes *NodesFilterer) WatchMaxActiveNodesUpdated(opts *bind.WatchOpts, sink chan<- *NodesMaxActiveNodesUpdated) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "MaxActiveNodesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesMaxActiveNodesUpdated)
				if err := _Nodes.contract.UnpackLog(event, "MaxActiveNodesUpdated", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseMaxActiveNodesUpdated(log types.Log) (*NodesMaxActiveNodesUpdated, error) {
	event := new(NodesMaxActiveNodesUpdated)
	if err := _Nodes.contract.UnpackLog(event, "MaxActiveNodesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesMinMonthlyFeeUpdatedIterator is returned from FilterMinMonthlyFeeUpdated and is used to iterate over the raw logs and unpacked data for MinMonthlyFeeUpdated events raised by the Nodes contract.
type NodesMinMonthlyFeeUpdatedIterator struct {
	Event *NodesMinMonthlyFeeUpdated // Event containing the contract specifics and raw log

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
func (it *NodesMinMonthlyFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesMinMonthlyFeeUpdated)
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
		it.Event = new(NodesMinMonthlyFeeUpdated)
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
func (it *NodesMinMonthlyFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesMinMonthlyFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesMinMonthlyFeeUpdated represents a MinMonthlyFeeUpdated event raised by the Nodes contract.
type NodesMinMonthlyFeeUpdated struct {
	NodeId        *big.Int
	MinMonthlyFee *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMinMonthlyFeeUpdated is a free log retrieval operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFee)
func (_Nodes *NodesFilterer) FilterMinMonthlyFeeUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesMinMonthlyFeeUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "MinMonthlyFeeUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesMinMonthlyFeeUpdatedIterator{contract: _Nodes.contract, event: "MinMonthlyFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinMonthlyFeeUpdated is a free log subscription operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFee)
func (_Nodes *NodesFilterer) WatchMinMonthlyFeeUpdated(opts *bind.WatchOpts, sink chan<- *NodesMinMonthlyFeeUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "MinMonthlyFeeUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesMinMonthlyFeeUpdated)
				if err := _Nodes.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
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
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFee)
func (_Nodes *NodesFilterer) ParseMinMonthlyFeeUpdated(log types.Log) (*NodesMinMonthlyFeeUpdated, error) {
	event := new(NodesMinMonthlyFeeUpdated)
	if err := _Nodes.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesNodeActivateUpdatedIterator is returned from FilterNodeActivateUpdated and is used to iterate over the raw logs and unpacked data for NodeActivateUpdated events raised by the Nodes contract.
type NodesNodeActivateUpdatedIterator struct {
	Event *NodesNodeActivateUpdated // Event containing the contract specifics and raw log

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
func (it *NodesNodeActivateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeActivateUpdated)
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
		it.Event = new(NodesNodeActivateUpdated)
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
func (it *NodesNodeActivateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeActivateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeActivateUpdated represents a NodeActivateUpdated event raised by the Nodes contract.
type NodesNodeActivateUpdated struct {
	NodeId   *big.Int
	IsActive bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNodeActivateUpdated is a free log retrieval operation binding the contract event 0x4a518a6a9ee77b4498e418883bc42338213163021cf974718d9fe36511d6010f.
//
// Solidity: event NodeActivateUpdated(uint256 indexed nodeId, bool isActive)
func (_Nodes *NodesFilterer) FilterNodeActivateUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesNodeActivateUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeActivateUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesNodeActivateUpdatedIterator{contract: _Nodes.contract, event: "NodeActivateUpdated", logs: logs, sub: sub}, nil
}

// WatchNodeActivateUpdated is a free log subscription operation binding the contract event 0x4a518a6a9ee77b4498e418883bc42338213163021cf974718d9fe36511d6010f.
//
// Solidity: event NodeActivateUpdated(uint256 indexed nodeId, bool isActive)
func (_Nodes *NodesFilterer) WatchNodeActivateUpdated(opts *bind.WatchOpts, sink chan<- *NodesNodeActivateUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeActivateUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeActivateUpdated)
				if err := _Nodes.contract.UnpackLog(event, "NodeActivateUpdated", log); err != nil {
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

// ParseNodeActivateUpdated is a log parse operation binding the contract event 0x4a518a6a9ee77b4498e418883bc42338213163021cf974718d9fe36511d6010f.
//
// Solidity: event NodeActivateUpdated(uint256 indexed nodeId, bool isActive)
func (_Nodes *NodesFilterer) ParseNodeActivateUpdated(log types.Log) (*NodesNodeActivateUpdated, error) {
	event := new(NodesNodeActivateUpdated)
	if err := _Nodes.contract.UnpackLog(event, "NodeActivateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesNodeAddedIterator is returned from FilterNodeAdded and is used to iterate over the raw logs and unpacked data for NodeAdded events raised by the Nodes contract.
type NodesNodeAddedIterator struct {
	Event *NodesNodeAdded // Event containing the contract specifics and raw log

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
func (it *NodesNodeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeAdded)
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
		it.Event = new(NodesNodeAdded)
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
func (it *NodesNodeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeAdded represents a NodeAdded event raised by the Nodes contract.
type NodesNodeAdded struct {
	NodeId        *big.Int
	Owner         common.Address
	SigningKeyPub []byte
	HttpAddress   string
	MinMonthlyFee *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNodeAdded is a free log retrieval operation binding the contract event 0x663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c.
//
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee)
func (_Nodes *NodesFilterer) FilterNodeAdded(opts *bind.FilterOpts, nodeId []*big.Int, owner []common.Address) (*NodesNodeAddedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeAdded", nodeIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &NodesNodeAddedIterator{contract: _Nodes.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

// WatchNodeAdded is a free log subscription operation binding the contract event 0x663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c.
//
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee)
func (_Nodes *NodesFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *NodesNodeAdded, nodeId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeAdded", nodeIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeAdded)
				if err := _Nodes.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee)
func (_Nodes *NodesFilterer) ParseNodeAdded(log types.Log) (*NodesNodeAdded, error) {
	event := new(NodesNodeAdded)
	if err := _Nodes.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesNodeOperatorCommissionPercentUpdatedIterator is returned from FilterNodeOperatorCommissionPercentUpdated and is used to iterate over the raw logs and unpacked data for NodeOperatorCommissionPercentUpdated events raised by the Nodes contract.
type NodesNodeOperatorCommissionPercentUpdatedIterator struct {
	Event *NodesNodeOperatorCommissionPercentUpdated // Event containing the contract specifics and raw log

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
func (it *NodesNodeOperatorCommissionPercentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeOperatorCommissionPercentUpdated)
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
		it.Event = new(NodesNodeOperatorCommissionPercentUpdated)
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
func (it *NodesNodeOperatorCommissionPercentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeOperatorCommissionPercentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeOperatorCommissionPercentUpdated represents a NodeOperatorCommissionPercentUpdated event raised by the Nodes contract.
type NodesNodeOperatorCommissionPercentUpdated struct {
	NewCommissionPercent *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterNodeOperatorCommissionPercentUpdated is a free log retrieval operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_Nodes *NodesFilterer) FilterNodeOperatorCommissionPercentUpdated(opts *bind.FilterOpts) (*NodesNodeOperatorCommissionPercentUpdatedIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeOperatorCommissionPercentUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesNodeOperatorCommissionPercentUpdatedIterator{contract: _Nodes.contract, event: "NodeOperatorCommissionPercentUpdated", logs: logs, sub: sub}, nil
}

// WatchNodeOperatorCommissionPercentUpdated is a free log subscription operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_Nodes *NodesFilterer) WatchNodeOperatorCommissionPercentUpdated(opts *bind.WatchOpts, sink chan<- *NodesNodeOperatorCommissionPercentUpdated) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeOperatorCommissionPercentUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeOperatorCommissionPercentUpdated)
				if err := _Nodes.contract.UnpackLog(event, "NodeOperatorCommissionPercentUpdated", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseNodeOperatorCommissionPercentUpdated(log types.Log) (*NodesNodeOperatorCommissionPercentUpdated, error) {
	event := new(NodesNodeOperatorCommissionPercentUpdated)
	if err := _Nodes.contract.UnpackLog(event, "NodeOperatorCommissionPercentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesNodeTransferredIterator is returned from FilterNodeTransferred and is used to iterate over the raw logs and unpacked data for NodeTransferred events raised by the Nodes contract.
type NodesNodeTransferredIterator struct {
	Event *NodesNodeTransferred // Event containing the contract specifics and raw log

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
func (it *NodesNodeTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeTransferred)
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
		it.Event = new(NodesNodeTransferred)
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
func (it *NodesNodeTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeTransferred represents a NodeTransferred event raised by the Nodes contract.
type NodesNodeTransferred struct {
	NodeId *big.Int
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeTransferred is a free log retrieval operation binding the contract event 0x0080108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be72.
//
// Solidity: event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to)
func (_Nodes *NodesFilterer) FilterNodeTransferred(opts *bind.FilterOpts, nodeId []*big.Int, from []common.Address, to []common.Address) (*NodesNodeTransferredIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeTransferred", nodeIdRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NodesNodeTransferredIterator{contract: _Nodes.contract, event: "NodeTransferred", logs: logs, sub: sub}, nil
}

// WatchNodeTransferred is a free log subscription operation binding the contract event 0x0080108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be72.
//
// Solidity: event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to)
func (_Nodes *NodesFilterer) WatchNodeTransferred(opts *bind.WatchOpts, sink chan<- *NodesNodeTransferred, nodeId []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeTransferred", nodeIdRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeTransferred)
				if err := _Nodes.contract.UnpackLog(event, "NodeTransferred", log); err != nil {
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

// ParseNodeTransferred is a log parse operation binding the contract event 0x0080108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be72.
//
// Solidity: event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to)
func (_Nodes *NodesFilterer) ParseNodeTransferred(log types.Log) (*NodesNodeTransferred, error) {
	event := new(NodesNodeTransferred)
	if err := _Nodes.contract.UnpackLog(event, "NodeTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesReplicationEnabledUpdatedIterator is returned from FilterReplicationEnabledUpdated and is used to iterate over the raw logs and unpacked data for ReplicationEnabledUpdated events raised by the Nodes contract.
type NodesReplicationEnabledUpdatedIterator struct {
	Event *NodesReplicationEnabledUpdated // Event containing the contract specifics and raw log

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
func (it *NodesReplicationEnabledUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesReplicationEnabledUpdated)
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
		it.Event = new(NodesReplicationEnabledUpdated)
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
func (it *NodesReplicationEnabledUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesReplicationEnabledUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesReplicationEnabledUpdated represents a ReplicationEnabledUpdated event raised by the Nodes contract.
type NodesReplicationEnabledUpdated struct {
	NodeId               *big.Int
	IsReplicationEnabled bool
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterReplicationEnabledUpdated is a free log retrieval operation binding the contract event 0xda2a657fb74ca331fb64eaabeca91a4ab0c68fd7ce7a8938a1a709903cf9be1e.
//
// Solidity: event ReplicationEnabledUpdated(uint256 indexed nodeId, bool isReplicationEnabled)
func (_Nodes *NodesFilterer) FilterReplicationEnabledUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesReplicationEnabledUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ReplicationEnabledUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesReplicationEnabledUpdatedIterator{contract: _Nodes.contract, event: "ReplicationEnabledUpdated", logs: logs, sub: sub}, nil
}

// WatchReplicationEnabledUpdated is a free log subscription operation binding the contract event 0xda2a657fb74ca331fb64eaabeca91a4ab0c68fd7ce7a8938a1a709903cf9be1e.
//
// Solidity: event ReplicationEnabledUpdated(uint256 indexed nodeId, bool isReplicationEnabled)
func (_Nodes *NodesFilterer) WatchReplicationEnabledUpdated(opts *bind.WatchOpts, sink chan<- *NodesReplicationEnabledUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ReplicationEnabledUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesReplicationEnabledUpdated)
				if err := _Nodes.contract.UnpackLog(event, "ReplicationEnabledUpdated", log); err != nil {
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

// ParseReplicationEnabledUpdated is a log parse operation binding the contract event 0xda2a657fb74ca331fb64eaabeca91a4ab0c68fd7ce7a8938a1a709903cf9be1e.
//
// Solidity: event ReplicationEnabledUpdated(uint256 indexed nodeId, bool isReplicationEnabled)
func (_Nodes *NodesFilterer) ParseReplicationEnabledUpdated(log types.Log) (*NodesReplicationEnabledUpdated, error) {
	event := new(NodesReplicationEnabledUpdated)
	if err := _Nodes.contract.UnpackLog(event, "ReplicationEnabledUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Nodes contract.
type NodesRoleAdminChangedIterator struct {
	Event *NodesRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *NodesRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesRoleAdminChanged)
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
		it.Event = new(NodesRoleAdminChanged)
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
func (it *NodesRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesRoleAdminChanged represents a RoleAdminChanged event raised by the Nodes contract.
type NodesRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Nodes *NodesFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*NodesRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &NodesRoleAdminChangedIterator{contract: _Nodes.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Nodes *NodesFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *NodesRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesRoleAdminChanged)
				if err := _Nodes.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseRoleAdminChanged(log types.Log) (*NodesRoleAdminChanged, error) {
	event := new(NodesRoleAdminChanged)
	if err := _Nodes.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Nodes contract.
type NodesRoleGrantedIterator struct {
	Event *NodesRoleGranted // Event containing the contract specifics and raw log

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
func (it *NodesRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesRoleGranted)
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
		it.Event = new(NodesRoleGranted)
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
func (it *NodesRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesRoleGranted represents a RoleGranted event raised by the Nodes contract.
type NodesRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Nodes *NodesFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*NodesRoleGrantedIterator, error) {

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

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NodesRoleGrantedIterator{contract: _Nodes.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Nodes *NodesFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *NodesRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesRoleGranted)
				if err := _Nodes.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseRoleGranted(log types.Log) (*NodesRoleGranted, error) {
	event := new(NodesRoleGranted)
	if err := _Nodes.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Nodes contract.
type NodesRoleRevokedIterator struct {
	Event *NodesRoleRevoked // Event containing the contract specifics and raw log

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
func (it *NodesRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesRoleRevoked)
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
		it.Event = new(NodesRoleRevoked)
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
func (it *NodesRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesRoleRevoked represents a RoleRevoked event raised by the Nodes contract.
type NodesRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Nodes *NodesFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*NodesRoleRevokedIterator, error) {

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

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NodesRoleRevokedIterator{contract: _Nodes.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Nodes *NodesFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *NodesRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesRoleRevoked)
				if err := _Nodes.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseRoleRevoked(log types.Log) (*NodesRoleRevoked, error) {
	event := new(NodesRoleRevoked)
	if err := _Nodes.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Nodes contract.
type NodesTransferIterator struct {
	Event *NodesTransfer // Event containing the contract specifics and raw log

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
func (it *NodesTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesTransfer)
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
		it.Event = new(NodesTransfer)
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
func (it *NodesTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesTransfer represents a Transfer event raised by the Nodes contract.
type NodesTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Nodes *NodesFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NodesTransferIterator, error) {

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

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesTransferIterator{contract: _Nodes.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Nodes *NodesFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NodesTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesTransfer)
				if err := _Nodes.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseTransfer(log types.Log) (*NodesTransfer, error) {
	event := new(NodesTransfer)
	if err := _Nodes.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
