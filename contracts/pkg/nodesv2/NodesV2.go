// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nodesv2

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
	IsDisabled           bool
	MinMonthlyFee        *big.Int
}

// INodesNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type INodesNodeWithId struct {
	NodeId *big.Int
	Node   INodesNode
}

// NodesV2MetaData contains all meta data concerning the NodesV2 contract.
var NodesV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_initialAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disableNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enableNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveApiNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveApiNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveApiNodesIDs\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesIDs\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveReplicationNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveReplicationNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveReplicationNodesIDs\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesIDs\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"allNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApiNodeIsActive\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeOperatorCommissionPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"commissionPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReplicationNodeIsActive\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxActiveNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeOperatorCommissionPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeFromApiNodes\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeFromReplicationNodes\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setHttpAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIsApiEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIsReplicationEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxActiveNodes\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinMonthlyFee\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNodeOperatorCommissionPercent\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ApiDisabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApiEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BaseURIUpdated\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HttpAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newHttpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxActiveNodesUpdated\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinMonthlyFeeUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAdded\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"minMonthlyFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeDisabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeOperatorCommissionPercentUpdated\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeTransferred\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplicationDisabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplicationEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCommissionPercent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHttpAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInputLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNodeConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigningKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesBelowCurrentCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeIsDisabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x6080604052600a805464ffffffffff19166014179055348015610020575f5ffd5b50604051614a17380380614a1783398101604081905261003f91610317565b60408051808201825260128152712c26aa28102737b2329027b832b930ba37b960711b602080830191909152825180840190935260048352630584d54560e41b90830152906202a300836001600160a01b0381166100b657604051636116401160e11b81525f600482015260240160405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100df5f8261018b565b50600391506100f0905083826103dc565b5060046100fd82826103dc565b5050506001600160a01b0381166101275760405163e6c4247b60e01b815260040160405180910390fd5b61013e5f5160206149f75f395f51905f525f6101fa565b6101555f5160206149d75f395f51905f525f6101fa565b61016c5f5160206149f75f395f51905f528261018b565b506101845f5160206149d75f395f51905f528261018b565b5050610496565b5f826101e7575f6101a46002546001600160a01b031690565b6001600160a01b0316146101cb57604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6101f18383610226565b90505b92915050565b8161021857604051631fe1e13d60e11b815260040160405180910390fd5b61022282826102cd565b5050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff166102c6575f838152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905561027e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016101f4565b505f6101f4565b5f82815260208190526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b5f60208284031215610327575f5ffd5b81516001600160a01b038116811461033d575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061036c57607f821691505b60208210810361038a57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156103d757805f5260205f20601f840160051c810160208510156103b55750805b601f840160051c820191505b818110156103d4575f81556001016103c1565b50505b505050565b81516001600160401b038111156103f5576103f5610344565b610409816104038454610358565b84610390565b6020601f82116001811461043b575f83156104245750848201515b5f19600385901b1c1916600184901b1784556103d4565b5f84815260208120601f198516915b8281101561046a578785015182556020948501946001909201910161044a565b508482101561048757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b614534806104a35f395ff3fe608060405234801561000f575f5ffd5b5060043610610388575f3560e01c80638da5cb5b116101df578063cc8463c811610109578063d74a2a50116100a9578063f3194a3911610079578063f3194a391461084c578063f579d7e114610855578063fb1120e21461085d578063fd967f4714610865575f5ffd5b8063d74a2a50146107e3578063e18cb254146107f6578063e985e9c514610809578063ebe487bf14610844575f5ffd5b8063cf6eefb7116100e4578063cf6eefb714610762578063d547741f146107a1578063d59f9fe0146107b4578063d602b9fd146107db575f5ffd5b8063cc8463c81461073f578063ce99948914610747578063cefc14291461075a575f5ffd5b8063a1eda53c1161017f578063b88d4fde1161014f578063b88d4fde146106fe578063b9b140d614610711578063c4741f3114610719578063c87b56dd1461072c575f5ffd5b8063a1eda53c146106aa578063a217fddf146106d1578063a22cb465146106d8578063a835f88e146106eb575f5ffd5b806391d14854116101ba57806391d148541461064557806395d89b411461067b5780639d32f9ba14610683578063a1174e7d146106a2575f5ffd5b80638da5cb5b146106155780638ed9ea341461061d5780638fbbf62314610630575f5ffd5b806342842e0e116102c0578063646453ba1161026057806375b238fc1161023057806375b238fc146105b757806379e0d58c146105de57806384ef8ffc146105f1578063895620b714610602575f5ffd5b8063646453ba14610569578063649a5ec71461057e5780636ec97bfc1461059157806370a08231146105a4575f5ffd5b806350d0215f1161029b57806350d0215f1461051d57806355f804b314610530578063634e93da146105435780636352211e14610556575f5ffd5b806342842e0e146104d757806344ff624e146104ea5780634f0f4aa9146104fd575f5ffd5b8063203ede771161032b578063248a9ca311610306578063248a9ca31461047c5780632f2ff15d1461049e57806336568abe146104b15780633d2853fb146104c4575f5ffd5b8063203ede771461044357806321fbd7cb1461045657806323b872dd14610469575f5ffd5b8063081812fc11610366578063081812fc146103e5578063095ea7b3146104105780630aa6220b1461042557806317e3b3a91461042d575f5ffd5b806301ffc9a71461038c578063022d63fb146103b457806306fdde03146103d0575b5f5ffd5b61039f61039a366004613afa565b61086e565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff90911681526020016103ab565b6103d861087e565b6040516103ab9190613b43565b6103f86103f3366004613b55565b61090e565b6040516001600160a01b0390911681526020016103ab565b61042361041e366004613b87565b610935565b005b610423610944565b610435610959565b6040519081526020016103ab565b610423610451366004613b55565b610969565b61039f610464366004613b55565b6109e9565b610423610477366004613baf565b6109f5565b61043561048a366004613b55565b5f9081526020819052604090206001015490565b6104236104ac366004613be9565b610a82565b6104236104bf366004613be9565b610ac3565b6104236104d2366004613b55565b610bb3565b6104236104e5366004613baf565b610c8b565b61039f6104f8366004613b55565b610caa565b61051061050b366004613b55565b610cb6565b6040516103ab9190613c7b565b600a54610100900463ffffffff16610435565b61042361053e366004613cd2565b610ea9565b610423610551366004613d11565b610fda565b6103f8610564366004613b55565b610fed565b610571610ff7565b6040516103ab9190613d2a565b61042361058c366004613d6c565b611003565b61043561059f366004613d91565b611016565b6104356105b2366004613d11565b61133a565b6104357fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b6104236105ec366004613b55565b611398565b6002546001600160a01b03166103f8565b610423610610366004613e28565b611418565b6103f861151f565b61042361062b366004613e49565b611532565b610638611602565b6040516103ab9190613e69565b61039f610653366004613be9565b5f918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b6103d8611873565b600a546106909060ff1681565b60405160ff90911681526020016103ab565b610638611882565b6106b2611b08565b6040805165ffffffffffff9384168152929091166020830152016103ab565b6104355f81565b6104236106e6366004613f01565b611b82565b6104236106f9366004613b55565b611b8d565b61042361070c366004613f56565b611c7f565b601054610435565b610423610727366004613b55565b611c9d565b6103d861073a366004613b55565b611d38565b6103b9611d9d565b610423610755366004614034565b611e3a565b610423611f05565b600154604080516001600160a01b03831681527401000000000000000000000000000000000000000090920465ffffffffffff166020830152016103ab565b6104236107af366004613be9565b611f54565b6104357fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a81565b610423611f95565b6104236107f1366004614054565b611fa7565b610423610804366004613e28565b6120b1565b61039f61081736600461409c565b6001600160a01b039182165f90815260086020908152604080832093909416825291909152205460ff1690565b6106386121b8565b61043560105481565b610435612425565b610571612430565b61043561271081565b5f6108788261243c565b92915050565b60606003805461088d906140c4565b80601f01602080910402602001604051908101604052809291908181526020018280546108b9906140c4565b80156109045780601f106108db57610100808354040283529160200191610904565b820191905f5260205f20905b8154815290600101906020018083116108e757829003601f168201915b5050505050905090565b5f610918826124dd565b505f828152600760205260409020546001600160a01b0316610878565b61094082823361252e565b5050565b5f61094e8161253b565b610956612545565b50565b5f610964600e612551565b905090565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756109938161253b565b5f828152600560205260409020546001600160a01b03166109e0576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109408261255a565b5f610878600c836125bc565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a610a1f8161253b565b610a28826125d3565b610a318261255a565b610a3c848484612653565b826001600160a01b0316846001600160a01b0316837e80108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be7260405160405180910390a450505050565b81610ab9576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109408282612708565b81158015610ade57506002546001600160a01b038281169116145b15610ba9576001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610b25575065ffffffffffff8116155b80610b3857504265ffffffffffff821610155b15610b7e576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024015b60405180910390fd5b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b610940828261272c565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610bdd8161253b565b5f828152600560205260409020546001600160a01b0316610c2a576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f828152600b602052604080822060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff1690555183917ff044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a91a25050565b610ca583838360405180602001604052805f815250611c7f565b505050565b5f610878600e836125bc565b6040805160c081018252606080825260208083018290525f8385018190529183018290526080830182905260a083018290528482526005905291909120546001600160a01b0316610d33576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f828152600b602052604090819020815160c08101909252805482908290610d5a906140c4565b80601f0160208091040260200160405190810160405280929190818152602001828054610d86906140c4565b8015610dd15780601f10610da857610100808354040283529160200191610dd1565b820191905f5260205f20905b815481529060010190602001808311610db457829003601f168201915b50505050508152602001600182018054610dea906140c4565b80601f0160208091040260200160405190810160405280929190818152602001828054610e16906140c4565b8015610e615780601f10610e3857610100808354040283529160200191610e61565b820191905f5260205f20905b815481529060010190602001808311610e4457829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015292915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610ed38161253b565b81610f0a576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8282610f17600182614142565b818110610f2657610f26614155565b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916602f60f81b14610f8e576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6009610f9b8385836141c6565b507f6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad8383604051610fcd9291906142a9565b60405180910390a1505050565b5f610fe48161253b565b61094082612778565b5f610878826124dd565b6060610964600c6127ea565b5f61100d8161253b565b610940826127f6565b5f7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756110418161253b565b6001600160a01b038816611081576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b856110b8576040517f8125403000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836110ef576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a8054610100900463ffffffff1690600161110a836142bc565b91906101000a81548163ffffffff021916908363ffffffff160217905550505f6064600a60019054906101000a900463ffffffff1661114991906142e0565b905061115b898263ffffffff1661285e565b6040518060c0016040528089898080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250505090825250604080516020601f8a01819004810282018101909252888152918101919089908990819084018382808284375f9201829052509385525050506020808301829052604080840183905260608401839052608090930188905263ffffffff85168252600b9052208151819061121290826142ff565b506020820151600182019061122790826142ff565b5060408281015160028301805460608601516080870151151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff951515959095167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090931692909217939093179290921691909117905560a090920151600390910155516001600160a01b038a169063ffffffff8316907f663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c90611320908c908c908c908c908c906143ba565b60405180910390a363ffffffff1698975050505050505050565b5f6001600160a01b03821661137d576040517f89c62b640000000000000000000000000000000000000000000000000000000081525f6004820152602401610b75565b506001600160a01b03165f9081526006602052604090205490565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756113c28161253b565b5f828152600560205260409020546001600160a01b031661140f576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610940826125d3565b5f8281526005602052604090205482906001600160a01b0316611467576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f818152600b602052604090206002015462010000900460ff16156114b8576040517fc40c6f6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838152600560205260409020546001600160a01b03163314611507576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b811561151657610ca5836128f1565b610ca5836125d3565b5f6109646002546001600160a01b031690565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561155c8161253b565b611566600c612551565b8260ff161080611581575061157b600e612551565b8260ff16105b156115b8576040517f39beadee00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a805460ff191660ff84169081179091556040519081527f6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821906020015b60405180910390a15050565b606061160e600e612551565b67ffffffffffffffff81111561162657611626613f29565b60405190808252806020026020018201604052801561165f57816020015b61164c613a7c565b8152602001906001900390816116445790505b5090505f5b61166e600e612551565b8163ffffffff16101561186f575f611690600e63ffffffff808516906129bb16565b5f818152600560205260409020549091506001600160a01b03161561185c576040518060400160405280828152602001600b5f8481526020019081526020015f206040518060c00160405290815f820180546116eb906140c4565b80601f0160208091040260200160405190810160405280929190818152602001828054611717906140c4565b80156117625780601f1061173957610100808354040283529160200191611762565b820191905f5260205f20905b81548152906001019060200180831161174557829003601f168201915b5050505050815260200160018201805461177b906140c4565b80601f01602080910402602001604051908101604052809291908181526020018280546117a7906140c4565b80156117f25780601f106117c9576101008083540402835291602001916117f2565b820191905f5260205f20905b8154815290600101906020018083116117d557829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff851690811061185057611850614155565b60200260200101819052505b5080611867816142bc565b915050611664565b5090565b60606004805461088d906140c4565b600a54606090610100900463ffffffff1667ffffffffffffffff8111156118ab576118ab613f29565b6040519080825280602002602001820160405280156118e457816020015b6118d1613a7c565b8152602001906001900390816118c95790505b5090505f5b600a5463ffffffff6101009091048116908216101561186f575f61190e8260016143f3565b6119199060646142e0565b90506119418163ffffffff165f908152600560205260409020546001600160a01b0316151590565b15611aff5760405180604001604052808263ffffffff168152602001600b5f8463ffffffff1681526020019081526020015f206040518060c00160405290815f8201805461198e906140c4565b80601f01602080910402602001604051908101604052809291908181526020018280546119ba906140c4565b8015611a055780601f106119dc57610100808354040283529160200191611a05565b820191905f5260205f20905b8154815290600101906020018083116119e857829003601f168201915b50505050508152602001600182018054611a1e906140c4565b80601f0160208091040260200160405190810160405280929190818152602001828054611a4a906140c4565b8015611a955780601f10611a6c57610100808354040283529160200191611a95565b820191905f5260205f20905b815481529060010190602001808311611a7857829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff8516908110611af357611af3614155565b60200260200101819052505b506001016118e9565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015611b4a57504265ffffffffffff821610155b611b55575f5f611b7a565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b6109403383836129c6565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611bb78161253b565b5f828152600560205260409020546001600160a01b0316611c04576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f828152600b6020526040902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff1662010000179055611c48826125d3565b611c518261255a565b60405182907fa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8905f90a25050565b611c8a8484846109f5565b611c973385858585612a7d565b50505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611cc78161253b565b612710821115611d03576040517f47d3b04600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60108290556040518281527f6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0906020016115f6565b6060611d43826124dd565b505f611d4d612c21565b90505f815111611d6b5760405180602001604052805f815250611d96565b80611d7584612c30565b604051602001611d86929190614426565b6040516020818303038152906040525b9392505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015611dde57504265ffffffffffff8216105b611e10576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16611e34565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a611e648161253b565b5f838152600560205260409020546001600160a01b0316611eb1576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838152600b6020526040908190206003018390555183907f27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a90611ef89085815260200190565b60405180910390a2505050565b6001546001600160a01b0316338114611f4c576040517fc22c8022000000000000000000000000000000000000000000000000000000008152336004820152602401610b75565b610956612ccd565b81611f8b576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109408282612da4565b5f611f9f8161253b565b610956612dc8565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a611fd18161253b565b5f848152600560205260409020546001600160a01b031661201e576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81612055576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f848152600b602052604090206001016120708385836141c6565b50837f15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed84846040516120a39291906142a9565b60405180910390a250505050565b5f8281526005602052604090205482906001600160a01b0316612100576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f818152600b602052604090206002015462010000900460ff1615612151576040517fc40c6f6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838152600560205260409020546001600160a01b031633146121a0576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81156121af57610ca583612dd2565b610ca58361255a565b60606121c4600c612551565b67ffffffffffffffff8111156121dc576121dc613f29565b60405190808252806020026020018201604052801561221557816020015b612202613a7c565b8152602001906001900390816121fa5790505b5090505f5b612224600c612551565b8163ffffffff16101561186f575f612246600c63ffffffff808516906129bb16565b5f818152600560205260409020549091506001600160a01b031615612412576040518060400160405280828152602001600b5f8481526020019081526020015f206040518060c00160405290815f820180546122a1906140c4565b80601f01602080910402602001604051908101604052809291908181526020018280546122cd906140c4565b80156123185780601f106122ef57610100808354040283529160200191612318565b820191905f5260205f20905b8154815290600101906020018083116122fb57829003601f168201915b50505050508152602001600182018054612331906140c4565b80601f016020809104026020016040519081016040528092919081815260200182805461235d906140c4565b80156123a85780601f1061237f576101008083540402835291602001916123a8565b820191905f5260205f20905b81548152906001019060200180831161238b57829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff851690811061240657612406614155565b60200260200101819052505b508061241d816142bc565b91505061221a565b5f610964600c612551565b6060610964600e6127ea565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd0000000000000000000000000000000000000000000000000000000014806124ce57507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b5e139f00000000000000000000000000000000000000000000000000000000145b80610878575061087882612e7d565b5f818152600560205260408120546001600160a01b031680610878576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101849052602401610b75565b610ca58383836001612ed2565b6109568133613025565b61254f5f5f613090565b565b5f610878825490565b5f818152600b60205260409020600201805460ff1916905561257d600e826125bc565b1561258f5761258d600e826131dc565b505b60405181907fa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce280761808905f90a250565b5f8181526001830160205260408120541515611d96565b5f818152600b6020526040902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055612614600c826125bc565b1561262657612624600c826131dc565b505b60405181907f9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f3905f90a250565b6001600160a01b038216612695576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610b75565b5f6126a18383336131e7565b9050836001600160a01b0316816001600160a01b031614611c97576040517f64283d7b0000000000000000000000000000000000000000000000000000000081526001600160a01b0380861660048301526024820184905282166044820152606401610b75565b5f828152602081905260409020600101546127228161253b565b611c9783836132f1565b6001600160a01b038116331461276e576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ca58282613388565b5f612781611d9d565b61278a426133dc565b612794919061443a565b90506127a08282613427565b60405165ffffffffffff821681526001600160a01b038316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b60605f611d96836134b5565b5f6128008261350e565b612809426133dc565b612813919061443a565b905061281f8282613090565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b91016115f6565b6001600160a01b0382166128a0576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610b75565b5f6128ac83835f6131e7565b90506001600160a01b03811615610ca5576040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610b75565b600a5460ff16612901600c612551565b10612938576040517f950be9a500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f818152600b6020526040902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010017905561297d600c826125bc565b61298e5761298c600c82613555565b505b60405181907f3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da173905f90a250565b5f611d968383613560565b6001600160a01b038216612a11576040517f5b08ba180000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610b75565b6001600160a01b038381165f81815260086020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b15612c1a576040517f150b7a020000000000000000000000000000000000000000000000000000000081526001600160a01b0384169063150b7a0290612ad8908890889087908790600401614458565b6020604051808303815f875af1925050508015612b12575060408051601f3d908101601f19168201909252612b0f91810190614498565b60015b612b92573d808015612b3f576040519150601f19603f3d011682016040523d82523d5f602084013e612b44565b606091505b5080515f03612b8a576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610b75565b805181602001fd5b7fffffffff0000000000000000000000000000000000000000000000000000000081167f150b7a020000000000000000000000000000000000000000000000000000000014612c18576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610b75565b505b5050505050565b60606009805461088d906140c4565b60605f612c3c83613586565b60010190505f8167ffffffffffffffff811115612c5b57612c5b613f29565b6040519080825280601f01601f191660200182016040528015612c85576020820181803683370190505b5090508181016020015b5f19017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a8504945084612c8f57509392505050565b6001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff16801580612d1057504265ffffffffffff821610155b15612d51576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610b75565b612d6c5f612d676002546001600160a01b031690565b613388565b50612d775f836132f1565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f82815260208190526040902060010154612dbe8161253b565b611c978383613388565b61254f5f5f613427565b600a5460ff16612de2600e612551565b10612e19576040517f950be9a500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f818152600b60205260409020600201805460ff19166001179055612e3f600e826125bc565b612e5057612e4e600e82613555565b505b60405181907fd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d0905f90a250565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f31498786000000000000000000000000000000000000000000000000000000001480610878575061087882613667565b8080612ee657506001600160a01b03821615155b15612fde575f612ef5846124dd565b90506001600160a01b03831615801590612f215750826001600160a01b0316816001600160a01b031614155b8015612f5257506001600160a01b038082165f9081526008602090815260408083209387168352929052205460ff16155b15612f94576040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526001600160a01b0384166004820152602401610b75565b8115612fdc5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260076020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b5f828152602081815260408083206001600160a01b038516845290915290205460ff16610940576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b038216600482015260248101839052604401610b75565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015613164574265ffffffffffff8216101561313b576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055613164565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b50600280546001600160a01b03167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f611d9683836136fd565b5f828152600560205260408120546001600160a01b0390811690831615613213576132138184866137e7565b6001600160a01b0381161561324d5761322e5f855f5f612ed2565b6001600160a01b0381165f90815260066020526040902080545f190190555b6001600160a01b0385161561327b576001600160a01b0385165f908152600660205260409020805460010190555b5f8481526005602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b5f8261337e575f61330a6002546001600160a01b031690565b6001600160a01b03161461334a576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384161790555b611d96838361387d565b5f821580156133a457506002546001600160a01b038381169116145b156133d257600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b611d968383613924565b5f65ffffffffffff82111561186f576040517f6dfcc6500000000000000000000000000000000000000000000000000000000081526030600482015260248101839052604401610b75565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff000000000000000000000000000000000000000000000000000084166001600160a01b03881617179093559004168015610ca5576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b6060815f0180548060200260200160405190810160405280929190818152602001828054801561350257602002820191905f5260205f20905b8154815260200190600101908083116134ee575b50505050509050919050565b5f5f613518611d9d565b90508065ffffffffffff168365ffffffffffff16116135405761353b83826144b3565b611d96565b611d9665ffffffffffff8416620697806139a5565b5f611d9683836139b4565b5f825f01828154811061357557613575614155565b905f5260205f200154905092915050565b5f807a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083106135ce577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef810000000083106135fa576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc10000831061361857662386f26fc10000830492506010015b6305f5e1008310613630576305f5e100830492506008015b612710831061364457612710830492506004015b60648310613656576064830492506002015b600a83106108785760010192915050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061087857507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000831614610878565b5f81815260018301602052604081205480156137d7575f61371f600183614142565b85549091505f9061373290600190614142565b9050808214613791575f865f01828154811061375057613750614155565b905f5260205f200154905080875f01848154811061377057613770614155565b5f918252602080832090910192909255918252600188019052604090208390555b85548690806137a2576137a26144d1565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f905560019350505050610878565b5f915050610878565b5092915050565b6137f28383836139f9565b610ca5576001600160a01b038316613839576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101829052602401610b75565b6040517f177e802f0000000000000000000000000000000000000000000000000000000081526001600160a01b038316600482015260248101829052604401610b75565b5f828152602081815260408083206001600160a01b038516845290915281205460ff1661391d575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556138d53390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610878565b505f610878565b5f828152602081815260408083206001600160a01b038516845290915281205460ff161561391d575f838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4506001610878565b5f828218828410028218611d96565b5f81815260018301602052604081205461391d57508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155610878565b5f6001600160a01b03831615801590613a745750826001600160a01b0316846001600160a01b03161480613a5157506001600160a01b038085165f9081526008602090815260408083209387168352929052205460ff165b80613a7457505f828152600760205260409020546001600160a01b038481169116145b949350505050565b60405180604001604052805f8152602001613ac86040518060c0016040528060608152602001606081526020015f151581526020015f151581526020015f151581526020015f81525090565b905290565b7fffffffff0000000000000000000000000000000000000000000000000000000081168114610956575f5ffd5b5f60208284031215613b0a575f5ffd5b8135611d9681613acd565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f611d966020830184613b15565b5f60208284031215613b65575f5ffd5b5035919050565b80356001600160a01b0381168114613b82575f5ffd5b919050565b5f5f60408385031215613b98575f5ffd5b613ba183613b6c565b946020939093013593505050565b5f5f5f60608486031215613bc1575f5ffd5b613bca84613b6c565b9250613bd860208501613b6c565b929592945050506040919091013590565b5f5f60408385031215613bfa575f5ffd5b82359150613c0a60208401613b6c565b90509250929050565b5f815160c08452613c2760c0850182613b15565b905060208301518482036020860152613c408282613b15565b91505060408301511515604085015260608301511515606085015260808301511515608085015260a083015160a08501528091505092915050565b602081525f611d966020830184613c13565b5f5f83601f840112613c9d575f5ffd5b50813567ffffffffffffffff811115613cb4575f5ffd5b602083019150836020828501011115613ccb575f5ffd5b9250929050565b5f5f60208385031215613ce3575f5ffd5b823567ffffffffffffffff811115613cf9575f5ffd5b613d0585828601613c8d565b90969095509350505050565b5f60208284031215613d21575f5ffd5b611d9682613b6c565b602080825282518282018190525f918401906040840190835b81811015613d61578351835260209384019390920191600101613d43565b509095945050505050565b5f60208284031215613d7c575f5ffd5b813565ffffffffffff81168114611d96575f5ffd5b5f5f5f5f5f5f60808789031215613da6575f5ffd5b613daf87613b6c565b9550602087013567ffffffffffffffff811115613dca575f5ffd5b613dd689828a01613c8d565b909650945050604087013567ffffffffffffffff811115613df5575f5ffd5b613e0189828a01613c8d565b979a9699509497949695606090950135949350505050565b80358015158114613b82575f5ffd5b5f5f60408385031215613e39575f5ffd5b82359150613c0a60208401613e19565b5f60208284031215613e59575f5ffd5b813560ff81168114611d96575f5ffd5b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b82811015613ef5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08786030184528151805186526020810151905060406020870152613edf6040870182613c13565b9550506020938401939190910190600101613e8f565b50929695505050505050565b5f5f60408385031215613f12575f5ffd5b613f1b83613b6c565b9150613c0a60208401613e19565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f5f5f60808587031215613f69575f5ffd5b613f7285613b6c565b9350613f8060208601613b6c565b925060408501359150606085013567ffffffffffffffff811115613fa2575f5ffd5b8501601f81018713613fb2575f5ffd5b803567ffffffffffffffff811115613fcc57613fcc613f29565b604051601f19603f601f19601f8501160116810181811067ffffffffffffffff82111715613ffc57613ffc613f29565b604052818152828201602001891015614013575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f60408385031215614045575f5ffd5b50508035926020909101359150565b5f5f5f60408486031215614066575f5ffd5b83359250602084013567ffffffffffffffff811115614083575f5ffd5b61408f86828701613c8d565b9497909650939450505050565b5f5f604083850312156140ad575f5ffd5b6140b683613b6c565b9150613c0a60208401613b6c565b600181811c908216806140d857607f821691505b60208210810361410f577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8181038181111561087857610878614115565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b601f821115610ca557805f5260205f20601f840160051c810160208510156141a75750805b601f840160051c820191505b81811015612c1a575f81556001016141b3565b67ffffffffffffffff8311156141de576141de613f29565b6141f2836141ec83546140c4565b83614182565b5f601f841160018114614223575f851561420c5750838201355b5f19600387901b1c1916600186901b178355612c1a565b5f83815260208120601f198716915b828110156142525786850135825560209485019460019092019101614232565b508682101561426e575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b81835281816020850137505f602082840101525f6020601f19601f840116840101905092915050565b602081525f613a74602083018486614280565b5f63ffffffff821663ffffffff81036142d7576142d7614115565b60010192915050565b63ffffffff81811683821602908116908181146137e0576137e0614115565b815167ffffffffffffffff81111561431957614319613f29565b61432d8161432784546140c4565b84614182565b6020601f82116001811461435f575f83156143485750848201515b5f19600385901b1c1916600184901b178455612c1a565b5f84815260208120601f198516915b8281101561438e578785015182556020948501946001909201910161436e565b50848210156143ab57868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b606081525f6143cd606083018789614280565b82810360208401526143e0818688614280565b9150508260408301529695505050505050565b63ffffffff818116838216019081111561087857610878614115565b5f81518060208401855e5f93019283525090919050565b5f613a74614434838661440f565b8461440f565b65ffffffffffff818116838216019081111561087857610878614115565b6001600160a01b03851681526001600160a01b0384166020820152826040820152608060608201525f61448e6080830184613b15565b9695505050505050565b5f602082840312156144a8575f5ffd5b8151611d9681613acd565b65ffffffffffff828116828216039081111561087857610878614115565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffdfea26469706673582212209ce3f6f1c799c54a8de010296209cfc74378194439e7e7721705b1d5a3b8993464736f6c634300081c0033daf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56aa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775",
}

// NodesV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use NodesV2MetaData.ABI instead.
var NodesV2ABI = NodesV2MetaData.ABI

// NodesV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodesV2MetaData.Bin instead.
var NodesV2Bin = NodesV2MetaData.Bin

// DeployNodesV2 deploys a new Ethereum contract, binding an instance of NodesV2 to it.
func DeployNodesV2(auth *bind.TransactOpts, backend bind.ContractBackend, _initialAdmin common.Address) (common.Address, *types.Transaction, *NodesV2, error) {
	parsed, err := NodesV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodesV2Bin), backend, _initialAdmin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodesV2{NodesV2Caller: NodesV2Caller{contract: contract}, NodesV2Transactor: NodesV2Transactor{contract: contract}, NodesV2Filterer: NodesV2Filterer{contract: contract}}, nil
}

// NodesV2 is an auto generated Go binding around an Ethereum contract.
type NodesV2 struct {
	NodesV2Caller     // Read-only binding to the contract
	NodesV2Transactor // Write-only binding to the contract
	NodesV2Filterer   // Log filterer for contract events
}

// NodesV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type NodesV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type NodesV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodesV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodesV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodesV2Session struct {
	Contract     *NodesV2          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodesV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodesV2CallerSession struct {
	Contract *NodesV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// NodesV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodesV2TransactorSession struct {
	Contract     *NodesV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// NodesV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type NodesV2Raw struct {
	Contract *NodesV2 // Generic contract binding to access the raw methods on
}

// NodesV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodesV2CallerRaw struct {
	Contract *NodesV2Caller // Generic read-only contract binding to access the raw methods on
}

// NodesV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodesV2TransactorRaw struct {
	Contract *NodesV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewNodesV2 creates a new instance of NodesV2, bound to a specific deployed contract.
func NewNodesV2(address common.Address, backend bind.ContractBackend) (*NodesV2, error) {
	contract, err := bindNodesV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodesV2{NodesV2Caller: NodesV2Caller{contract: contract}, NodesV2Transactor: NodesV2Transactor{contract: contract}, NodesV2Filterer: NodesV2Filterer{contract: contract}}, nil
}

// NewNodesV2Caller creates a new read-only instance of NodesV2, bound to a specific deployed contract.
func NewNodesV2Caller(address common.Address, caller bind.ContractCaller) (*NodesV2Caller, error) {
	contract, err := bindNodesV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodesV2Caller{contract: contract}, nil
}

// NewNodesV2Transactor creates a new write-only instance of NodesV2, bound to a specific deployed contract.
func NewNodesV2Transactor(address common.Address, transactor bind.ContractTransactor) (*NodesV2Transactor, error) {
	contract, err := bindNodesV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodesV2Transactor{contract: contract}, nil
}

// NewNodesV2Filterer creates a new log filterer instance of NodesV2, bound to a specific deployed contract.
func NewNodesV2Filterer(address common.Address, filterer bind.ContractFilterer) (*NodesV2Filterer, error) {
	contract, err := bindNodesV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodesV2Filterer{contract: contract}, nil
}

// bindNodesV2 binds a generic wrapper to an already deployed contract.
func bindNodesV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodesV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodesV2 *NodesV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodesV2.Contract.NodesV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodesV2 *NodesV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodesV2.Contract.NodesV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodesV2 *NodesV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodesV2.Contract.NodesV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodesV2 *NodesV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodesV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodesV2 *NodesV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodesV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodesV2 *NodesV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodesV2.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2Caller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2Session) ADMINROLE() ([32]byte, error) {
	return _NodesV2.Contract.ADMINROLE(&_NodesV2.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2CallerSession) ADMINROLE() ([32]byte, error) {
	return _NodesV2.Contract.ADMINROLE(&_NodesV2.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _NodesV2.Contract.DEFAULTADMINROLE(&_NodesV2.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _NodesV2.Contract.DEFAULTADMINROLE(&_NodesV2.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_NodesV2 *NodesV2Caller) MAXBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "MAX_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_NodesV2 *NodesV2Session) MAXBPS() (*big.Int, error) {
	return _NodesV2.Contract.MAXBPS(&_NodesV2.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint256)
func (_NodesV2 *NodesV2CallerSession) MAXBPS() (*big.Int, error) {
	return _NodesV2.Contract.MAXBPS(&_NodesV2.CallOpts)
}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2Caller) NODEMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "NODE_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2Session) NODEMANAGERROLE() ([32]byte, error) {
	return _NodesV2.Contract.NODEMANAGERROLE(&_NodesV2.CallOpts)
}

// NODEMANAGERROLE is a free data retrieval call binding the contract method 0xd59f9fe0.
//
// Solidity: function NODE_MANAGER_ROLE() view returns(bytes32)
func (_NodesV2 *NodesV2CallerSession) NODEMANAGERROLE() ([32]byte, error) {
	return _NodesV2.Contract.NODEMANAGERROLE(&_NodesV2.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodesV2 *NodesV2Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodesV2 *NodesV2Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodesV2.Contract.BalanceOf(&_NodesV2.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodesV2 *NodesV2CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodesV2.Contract.BalanceOf(&_NodesV2.CallOpts, owner)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_NodesV2 *NodesV2Caller) DefaultAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "defaultAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_NodesV2 *NodesV2Session) DefaultAdmin() (common.Address, error) {
	return _NodesV2.Contract.DefaultAdmin(&_NodesV2.CallOpts)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_NodesV2 *NodesV2CallerSession) DefaultAdmin() (common.Address, error) {
	return _NodesV2.Contract.DefaultAdmin(&_NodesV2.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_NodesV2 *NodesV2Caller) DefaultAdminDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "defaultAdminDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_NodesV2 *NodesV2Session) DefaultAdminDelay() (*big.Int, error) {
	return _NodesV2.Contract.DefaultAdminDelay(&_NodesV2.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_NodesV2 *NodesV2CallerSession) DefaultAdminDelay() (*big.Int, error) {
	return _NodesV2.Contract.DefaultAdminDelay(&_NodesV2.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_NodesV2 *NodesV2Caller) DefaultAdminDelayIncreaseWait(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "defaultAdminDelayIncreaseWait")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_NodesV2 *NodesV2Session) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _NodesV2.Contract.DefaultAdminDelayIncreaseWait(&_NodesV2.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_NodesV2 *NodesV2CallerSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _NodesV2.Contract.DefaultAdminDelayIncreaseWait(&_NodesV2.CallOpts)
}

// GetActiveApiNodes is a free data retrieval call binding the contract method 0xebe487bf.
//
// Solidity: function getActiveApiNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_NodesV2 *NodesV2Caller) GetActiveApiNodes(opts *bind.CallOpts) ([]INodesNodeWithId, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getActiveApiNodes")

	if err != nil {
		return *new([]INodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNodeWithId)).(*[]INodesNodeWithId)

	return out0, err

}

// GetActiveApiNodes is a free data retrieval call binding the contract method 0xebe487bf.
//
// Solidity: function getActiveApiNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_NodesV2 *NodesV2Session) GetActiveApiNodes() ([]INodesNodeWithId, error) {
	return _NodesV2.Contract.GetActiveApiNodes(&_NodesV2.CallOpts)
}

// GetActiveApiNodes is a free data retrieval call binding the contract method 0xebe487bf.
//
// Solidity: function getActiveApiNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_NodesV2 *NodesV2CallerSession) GetActiveApiNodes() ([]INodesNodeWithId, error) {
	return _NodesV2.Contract.GetActiveApiNodes(&_NodesV2.CallOpts)
}

// GetActiveApiNodesCount is a free data retrieval call binding the contract method 0xf579d7e1.
//
// Solidity: function getActiveApiNodesCount() view returns(uint256 activeNodesCount)
func (_NodesV2 *NodesV2Caller) GetActiveApiNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getActiveApiNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveApiNodesCount is a free data retrieval call binding the contract method 0xf579d7e1.
//
// Solidity: function getActiveApiNodesCount() view returns(uint256 activeNodesCount)
func (_NodesV2 *NodesV2Session) GetActiveApiNodesCount() (*big.Int, error) {
	return _NodesV2.Contract.GetActiveApiNodesCount(&_NodesV2.CallOpts)
}

// GetActiveApiNodesCount is a free data retrieval call binding the contract method 0xf579d7e1.
//
// Solidity: function getActiveApiNodesCount() view returns(uint256 activeNodesCount)
func (_NodesV2 *NodesV2CallerSession) GetActiveApiNodesCount() (*big.Int, error) {
	return _NodesV2.Contract.GetActiveApiNodesCount(&_NodesV2.CallOpts)
}

// GetActiveApiNodesIDs is a free data retrieval call binding the contract method 0x646453ba.
//
// Solidity: function getActiveApiNodesIDs() view returns(uint256[] activeNodesIDs)
func (_NodesV2 *NodesV2Caller) GetActiveApiNodesIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getActiveApiNodesIDs")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveApiNodesIDs is a free data retrieval call binding the contract method 0x646453ba.
//
// Solidity: function getActiveApiNodesIDs() view returns(uint256[] activeNodesIDs)
func (_NodesV2 *NodesV2Session) GetActiveApiNodesIDs() ([]*big.Int, error) {
	return _NodesV2.Contract.GetActiveApiNodesIDs(&_NodesV2.CallOpts)
}

// GetActiveApiNodesIDs is a free data retrieval call binding the contract method 0x646453ba.
//
// Solidity: function getActiveApiNodesIDs() view returns(uint256[] activeNodesIDs)
func (_NodesV2 *NodesV2CallerSession) GetActiveApiNodesIDs() ([]*big.Int, error) {
	return _NodesV2.Contract.GetActiveApiNodesIDs(&_NodesV2.CallOpts)
}

// GetActiveReplicationNodes is a free data retrieval call binding the contract method 0x8fbbf623.
//
// Solidity: function getActiveReplicationNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_NodesV2 *NodesV2Caller) GetActiveReplicationNodes(opts *bind.CallOpts) ([]INodesNodeWithId, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getActiveReplicationNodes")

	if err != nil {
		return *new([]INodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNodeWithId)).(*[]INodesNodeWithId)

	return out0, err

}

// GetActiveReplicationNodes is a free data retrieval call binding the contract method 0x8fbbf623.
//
// Solidity: function getActiveReplicationNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_NodesV2 *NodesV2Session) GetActiveReplicationNodes() ([]INodesNodeWithId, error) {
	return _NodesV2.Contract.GetActiveReplicationNodes(&_NodesV2.CallOpts)
}

// GetActiveReplicationNodes is a free data retrieval call binding the contract method 0x8fbbf623.
//
// Solidity: function getActiveReplicationNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_NodesV2 *NodesV2CallerSession) GetActiveReplicationNodes() ([]INodesNodeWithId, error) {
	return _NodesV2.Contract.GetActiveReplicationNodes(&_NodesV2.CallOpts)
}

// GetActiveReplicationNodesCount is a free data retrieval call binding the contract method 0x17e3b3a9.
//
// Solidity: function getActiveReplicationNodesCount() view returns(uint256 activeNodesCount)
func (_NodesV2 *NodesV2Caller) GetActiveReplicationNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getActiveReplicationNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveReplicationNodesCount is a free data retrieval call binding the contract method 0x17e3b3a9.
//
// Solidity: function getActiveReplicationNodesCount() view returns(uint256 activeNodesCount)
func (_NodesV2 *NodesV2Session) GetActiveReplicationNodesCount() (*big.Int, error) {
	return _NodesV2.Contract.GetActiveReplicationNodesCount(&_NodesV2.CallOpts)
}

// GetActiveReplicationNodesCount is a free data retrieval call binding the contract method 0x17e3b3a9.
//
// Solidity: function getActiveReplicationNodesCount() view returns(uint256 activeNodesCount)
func (_NodesV2 *NodesV2CallerSession) GetActiveReplicationNodesCount() (*big.Int, error) {
	return _NodesV2.Contract.GetActiveReplicationNodesCount(&_NodesV2.CallOpts)
}

// GetActiveReplicationNodesIDs is a free data retrieval call binding the contract method 0xfb1120e2.
//
// Solidity: function getActiveReplicationNodesIDs() view returns(uint256[] activeNodesIDs)
func (_NodesV2 *NodesV2Caller) GetActiveReplicationNodesIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getActiveReplicationNodesIDs")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveReplicationNodesIDs is a free data retrieval call binding the contract method 0xfb1120e2.
//
// Solidity: function getActiveReplicationNodesIDs() view returns(uint256[] activeNodesIDs)
func (_NodesV2 *NodesV2Session) GetActiveReplicationNodesIDs() ([]*big.Int, error) {
	return _NodesV2.Contract.GetActiveReplicationNodesIDs(&_NodesV2.CallOpts)
}

// GetActiveReplicationNodesIDs is a free data retrieval call binding the contract method 0xfb1120e2.
//
// Solidity: function getActiveReplicationNodesIDs() view returns(uint256[] activeNodesIDs)
func (_NodesV2 *NodesV2CallerSession) GetActiveReplicationNodesIDs() ([]*big.Int, error) {
	return _NodesV2.Contract.GetActiveReplicationNodesIDs(&_NodesV2.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodes)
func (_NodesV2 *NodesV2Caller) GetAllNodes(opts *bind.CallOpts) ([]INodesNodeWithId, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getAllNodes")

	if err != nil {
		return *new([]INodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNodeWithId)).(*[]INodesNodeWithId)

	return out0, err

}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodes)
func (_NodesV2 *NodesV2Session) GetAllNodes() ([]INodesNodeWithId, error) {
	return _NodesV2.Contract.GetAllNodes(&_NodesV2.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodes)
func (_NodesV2 *NodesV2CallerSession) GetAllNodes() ([]INodesNodeWithId, error) {
	return _NodesV2.Contract.GetAllNodes(&_NodesV2.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_NodesV2 *NodesV2Caller) GetAllNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getAllNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_NodesV2 *NodesV2Session) GetAllNodesCount() (*big.Int, error) {
	return _NodesV2.Contract.GetAllNodesCount(&_NodesV2.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint256 nodeCount)
func (_NodesV2 *NodesV2CallerSession) GetAllNodesCount() (*big.Int, error) {
	return _NodesV2.Contract.GetAllNodesCount(&_NodesV2.CallOpts)
}

// GetApiNodeIsActive is a free data retrieval call binding the contract method 0x21fbd7cb.
//
// Solidity: function getApiNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_NodesV2 *NodesV2Caller) GetApiNodeIsActive(opts *bind.CallOpts, nodeId *big.Int) (bool, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getApiNodeIsActive", nodeId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetApiNodeIsActive is a free data retrieval call binding the contract method 0x21fbd7cb.
//
// Solidity: function getApiNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_NodesV2 *NodesV2Session) GetApiNodeIsActive(nodeId *big.Int) (bool, error) {
	return _NodesV2.Contract.GetApiNodeIsActive(&_NodesV2.CallOpts, nodeId)
}

// GetApiNodeIsActive is a free data retrieval call binding the contract method 0x21fbd7cb.
//
// Solidity: function getApiNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_NodesV2 *NodesV2CallerSession) GetApiNodeIsActive(nodeId *big.Int) (bool, error) {
	return _NodesV2.Contract.GetApiNodeIsActive(&_NodesV2.CallOpts, nodeId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodesV2 *NodesV2Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodesV2 *NodesV2Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodesV2.Contract.GetApproved(&_NodesV2.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodesV2 *NodesV2CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodesV2.Contract.GetApproved(&_NodesV2.CallOpts, tokenId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,bool,bool,uint256) node)
func (_NodesV2 *NodesV2Caller) GetNode(opts *bind.CallOpts, nodeId *big.Int) (INodesNode, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getNode", nodeId)

	if err != nil {
		return *new(INodesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(INodesNode)).(*INodesNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,bool,bool,uint256) node)
func (_NodesV2 *NodesV2Session) GetNode(nodeId *big.Int) (INodesNode, error) {
	return _NodesV2.Contract.GetNode(&_NodesV2.CallOpts, nodeId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 nodeId) view returns((bytes,string,bool,bool,bool,uint256) node)
func (_NodesV2 *NodesV2CallerSession) GetNode(nodeId *big.Int) (INodesNode, error) {
	return _NodesV2.Contract.GetNode(&_NodesV2.CallOpts, nodeId)
}

// GetNodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xb9b140d6.
//
// Solidity: function getNodeOperatorCommissionPercent() view returns(uint256 commissionPercent)
func (_NodesV2 *NodesV2Caller) GetNodeOperatorCommissionPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getNodeOperatorCommissionPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xb9b140d6.
//
// Solidity: function getNodeOperatorCommissionPercent() view returns(uint256 commissionPercent)
func (_NodesV2 *NodesV2Session) GetNodeOperatorCommissionPercent() (*big.Int, error) {
	return _NodesV2.Contract.GetNodeOperatorCommissionPercent(&_NodesV2.CallOpts)
}

// GetNodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xb9b140d6.
//
// Solidity: function getNodeOperatorCommissionPercent() view returns(uint256 commissionPercent)
func (_NodesV2 *NodesV2CallerSession) GetNodeOperatorCommissionPercent() (*big.Int, error) {
	return _NodesV2.Contract.GetNodeOperatorCommissionPercent(&_NodesV2.CallOpts)
}

// GetReplicationNodeIsActive is a free data retrieval call binding the contract method 0x44ff624e.
//
// Solidity: function getReplicationNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_NodesV2 *NodesV2Caller) GetReplicationNodeIsActive(opts *bind.CallOpts, nodeId *big.Int) (bool, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getReplicationNodeIsActive", nodeId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetReplicationNodeIsActive is a free data retrieval call binding the contract method 0x44ff624e.
//
// Solidity: function getReplicationNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_NodesV2 *NodesV2Session) GetReplicationNodeIsActive(nodeId *big.Int) (bool, error) {
	return _NodesV2.Contract.GetReplicationNodeIsActive(&_NodesV2.CallOpts, nodeId)
}

// GetReplicationNodeIsActive is a free data retrieval call binding the contract method 0x44ff624e.
//
// Solidity: function getReplicationNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_NodesV2 *NodesV2CallerSession) GetReplicationNodeIsActive(nodeId *big.Int) (bool, error) {
	return _NodesV2.Contract.GetReplicationNodeIsActive(&_NodesV2.CallOpts, nodeId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_NodesV2 *NodesV2Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_NodesV2 *NodesV2Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _NodesV2.Contract.GetRoleAdmin(&_NodesV2.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_NodesV2 *NodesV2CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _NodesV2.Contract.GetRoleAdmin(&_NodesV2.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_NodesV2 *NodesV2Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_NodesV2 *NodesV2Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _NodesV2.Contract.HasRole(&_NodesV2.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_NodesV2 *NodesV2CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _NodesV2.Contract.HasRole(&_NodesV2.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodesV2 *NodesV2Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodesV2 *NodesV2Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodesV2.Contract.IsApprovedForAll(&_NodesV2.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodesV2 *NodesV2CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodesV2.Contract.IsApprovedForAll(&_NodesV2.CallOpts, owner, operator)
}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_NodesV2 *NodesV2Caller) MaxActiveNodes(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "maxActiveNodes")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_NodesV2 *NodesV2Session) MaxActiveNodes() (uint8, error) {
	return _NodesV2.Contract.MaxActiveNodes(&_NodesV2.CallOpts)
}

// MaxActiveNodes is a free data retrieval call binding the contract method 0x9d32f9ba.
//
// Solidity: function maxActiveNodes() view returns(uint8)
func (_NodesV2 *NodesV2CallerSession) MaxActiveNodes() (uint8, error) {
	return _NodesV2.Contract.MaxActiveNodes(&_NodesV2.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodesV2 *NodesV2Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodesV2 *NodesV2Session) Name() (string, error) {
	return _NodesV2.Contract.Name(&_NodesV2.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodesV2 *NodesV2CallerSession) Name() (string, error) {
	return _NodesV2.Contract.Name(&_NodesV2.CallOpts)
}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_NodesV2 *NodesV2Caller) NodeOperatorCommissionPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "nodeOperatorCommissionPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_NodesV2 *NodesV2Session) NodeOperatorCommissionPercent() (*big.Int, error) {
	return _NodesV2.Contract.NodeOperatorCommissionPercent(&_NodesV2.CallOpts)
}

// NodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xf3194a39.
//
// Solidity: function nodeOperatorCommissionPercent() view returns(uint256)
func (_NodesV2 *NodesV2CallerSession) NodeOperatorCommissionPercent() (*big.Int, error) {
	return _NodesV2.Contract.NodeOperatorCommissionPercent(&_NodesV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodesV2 *NodesV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodesV2 *NodesV2Session) Owner() (common.Address, error) {
	return _NodesV2.Contract.Owner(&_NodesV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NodesV2 *NodesV2CallerSession) Owner() (common.Address, error) {
	return _NodesV2.Contract.Owner(&_NodesV2.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodesV2 *NodesV2Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodesV2 *NodesV2Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodesV2.Contract.OwnerOf(&_NodesV2.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodesV2 *NodesV2CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodesV2.Contract.OwnerOf(&_NodesV2.CallOpts, tokenId)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_NodesV2 *NodesV2Caller) PendingDefaultAdmin(opts *bind.CallOpts) (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "pendingDefaultAdmin")

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
func (_NodesV2 *NodesV2Session) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _NodesV2.Contract.PendingDefaultAdmin(&_NodesV2.CallOpts)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_NodesV2 *NodesV2CallerSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _NodesV2.Contract.PendingDefaultAdmin(&_NodesV2.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_NodesV2 *NodesV2Caller) PendingDefaultAdminDelay(opts *bind.CallOpts) (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "pendingDefaultAdminDelay")

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
func (_NodesV2 *NodesV2Session) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _NodesV2.Contract.PendingDefaultAdminDelay(&_NodesV2.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_NodesV2 *NodesV2CallerSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _NodesV2.Contract.PendingDefaultAdminDelay(&_NodesV2.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodesV2 *NodesV2Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodesV2 *NodesV2Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodesV2.Contract.SupportsInterface(&_NodesV2.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodesV2 *NodesV2CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodesV2.Contract.SupportsInterface(&_NodesV2.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodesV2 *NodesV2Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodesV2 *NodesV2Session) Symbol() (string, error) {
	return _NodesV2.Contract.Symbol(&_NodesV2.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodesV2 *NodesV2CallerSession) Symbol() (string, error) {
	return _NodesV2.Contract.Symbol(&_NodesV2.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodesV2 *NodesV2Caller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NodesV2.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodesV2 *NodesV2Session) TokenURI(tokenId *big.Int) (string, error) {
	return _NodesV2.Contract.TokenURI(&_NodesV2.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodesV2 *NodesV2CallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodesV2.Contract.TokenURI(&_NodesV2.CallOpts, tokenId)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_NodesV2 *NodesV2Transactor) AcceptDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "acceptDefaultAdminTransfer")
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_NodesV2 *NodesV2Session) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodesV2.Contract.AcceptDefaultAdminTransfer(&_NodesV2.TransactOpts)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_NodesV2 *NodesV2TransactorSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodesV2.Contract.AcceptDefaultAdminTransfer(&_NodesV2.TransactOpts)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee) returns(uint256)
func (_NodesV2 *NodesV2Transactor) AddNode(opts *bind.TransactOpts, to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "addNode", to, signingKeyPub, httpAddress, minMonthlyFee)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee) returns(uint256)
func (_NodesV2 *NodesV2Session) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.AddNode(&_NodesV2.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFee)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee) returns(uint256)
func (_NodesV2 *NodesV2TransactorSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.AddNode(&_NodesV2.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFee)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodesV2 *NodesV2Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodesV2 *NodesV2Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.Approve(&_NodesV2.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodesV2 *NodesV2TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.Approve(&_NodesV2.TransactOpts, to, tokenId)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_NodesV2 *NodesV2Transactor) BeginDefaultAdminTransfer(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "beginDefaultAdminTransfer", newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_NodesV2 *NodesV2Session) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.BeginDefaultAdminTransfer(&_NodesV2.TransactOpts, newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_NodesV2 *NodesV2TransactorSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.BeginDefaultAdminTransfer(&_NodesV2.TransactOpts, newAdmin)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_NodesV2 *NodesV2Transactor) CancelDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "cancelDefaultAdminTransfer")
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_NodesV2 *NodesV2Session) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodesV2.Contract.CancelDefaultAdminTransfer(&_NodesV2.TransactOpts)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_NodesV2 *NodesV2TransactorSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _NodesV2.Contract.CancelDefaultAdminTransfer(&_NodesV2.TransactOpts)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_NodesV2 *NodesV2Transactor) ChangeDefaultAdminDelay(opts *bind.TransactOpts, newDelay *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "changeDefaultAdminDelay", newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_NodesV2 *NodesV2Session) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.ChangeDefaultAdminDelay(&_NodesV2.TransactOpts, newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_NodesV2 *NodesV2TransactorSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.ChangeDefaultAdminDelay(&_NodesV2.TransactOpts, newDelay)
}

// DisableNode is a paid mutator transaction binding the contract method 0xa835f88e.
//
// Solidity: function disableNode(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Transactor) DisableNode(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "disableNode", nodeId)
}

// DisableNode is a paid mutator transaction binding the contract method 0xa835f88e.
//
// Solidity: function disableNode(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Session) DisableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.DisableNode(&_NodesV2.TransactOpts, nodeId)
}

// DisableNode is a paid mutator transaction binding the contract method 0xa835f88e.
//
// Solidity: function disableNode(uint256 nodeId) returns()
func (_NodesV2 *NodesV2TransactorSession) DisableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.DisableNode(&_NodesV2.TransactOpts, nodeId)
}

// EnableNode is a paid mutator transaction binding the contract method 0x3d2853fb.
//
// Solidity: function enableNode(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Transactor) EnableNode(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "enableNode", nodeId)
}

// EnableNode is a paid mutator transaction binding the contract method 0x3d2853fb.
//
// Solidity: function enableNode(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Session) EnableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.EnableNode(&_NodesV2.TransactOpts, nodeId)
}

// EnableNode is a paid mutator transaction binding the contract method 0x3d2853fb.
//
// Solidity: function enableNode(uint256 nodeId) returns()
func (_NodesV2 *NodesV2TransactorSession) EnableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.EnableNode(&_NodesV2.TransactOpts, nodeId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.GrantRole(&_NodesV2.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.GrantRole(&_NodesV2.TransactOpts, role, account)
}

// RemoveFromApiNodes is a paid mutator transaction binding the contract method 0x79e0d58c.
//
// Solidity: function removeFromApiNodes(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Transactor) RemoveFromApiNodes(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "removeFromApiNodes", nodeId)
}

// RemoveFromApiNodes is a paid mutator transaction binding the contract method 0x79e0d58c.
//
// Solidity: function removeFromApiNodes(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Session) RemoveFromApiNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.RemoveFromApiNodes(&_NodesV2.TransactOpts, nodeId)
}

// RemoveFromApiNodes is a paid mutator transaction binding the contract method 0x79e0d58c.
//
// Solidity: function removeFromApiNodes(uint256 nodeId) returns()
func (_NodesV2 *NodesV2TransactorSession) RemoveFromApiNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.RemoveFromApiNodes(&_NodesV2.TransactOpts, nodeId)
}

// RemoveFromReplicationNodes is a paid mutator transaction binding the contract method 0x203ede77.
//
// Solidity: function removeFromReplicationNodes(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Transactor) RemoveFromReplicationNodes(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "removeFromReplicationNodes", nodeId)
}

// RemoveFromReplicationNodes is a paid mutator transaction binding the contract method 0x203ede77.
//
// Solidity: function removeFromReplicationNodes(uint256 nodeId) returns()
func (_NodesV2 *NodesV2Session) RemoveFromReplicationNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.RemoveFromReplicationNodes(&_NodesV2.TransactOpts, nodeId)
}

// RemoveFromReplicationNodes is a paid mutator transaction binding the contract method 0x203ede77.
//
// Solidity: function removeFromReplicationNodes(uint256 nodeId) returns()
func (_NodesV2 *NodesV2TransactorSession) RemoveFromReplicationNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.RemoveFromReplicationNodes(&_NodesV2.TransactOpts, nodeId)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2Session) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.RenounceRole(&_NodesV2.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2TransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.RenounceRole(&_NodesV2.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.RevokeRole(&_NodesV2.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_NodesV2 *NodesV2TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _NodesV2.Contract.RevokeRole(&_NodesV2.TransactOpts, role, account)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_NodesV2 *NodesV2Transactor) RollbackDefaultAdminDelay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "rollbackDefaultAdminDelay")
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_NodesV2 *NodesV2Session) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _NodesV2.Contract.RollbackDefaultAdminDelay(&_NodesV2.TransactOpts)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_NodesV2 *NodesV2TransactorSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _NodesV2.Contract.RollbackDefaultAdminDelay(&_NodesV2.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodesV2 *NodesV2Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodesV2 *NodesV2Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.SafeTransferFrom(&_NodesV2.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodesV2 *NodesV2TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.SafeTransferFrom(&_NodesV2.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodesV2 *NodesV2Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodesV2 *NodesV2Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodesV2.Contract.SafeTransferFrom0(&_NodesV2.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodesV2 *NodesV2TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodesV2.Contract.SafeTransferFrom0(&_NodesV2.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodesV2 *NodesV2Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodesV2 *NodesV2Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodesV2.Contract.SetApprovalForAll(&_NodesV2.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodesV2 *NodesV2TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodesV2.Contract.SetApprovalForAll(&_NodesV2.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_NodesV2 *NodesV2Transactor) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setBaseURI", newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_NodesV2 *NodesV2Session) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _NodesV2.Contract.SetBaseURI(&_NodesV2.TransactOpts, newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_NodesV2 *NodesV2TransactorSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _NodesV2.Contract.SetBaseURI(&_NodesV2.TransactOpts, newBaseURI)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_NodesV2 *NodesV2Transactor) SetHttpAddress(opts *bind.TransactOpts, nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setHttpAddress", nodeId, httpAddress)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_NodesV2 *NodesV2Session) SetHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _NodesV2.Contract.SetHttpAddress(&_NodesV2.TransactOpts, nodeId, httpAddress)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_NodesV2 *NodesV2TransactorSession) SetHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _NodesV2.Contract.SetHttpAddress(&_NodesV2.TransactOpts, nodeId, httpAddress)
}

// SetIsApiEnabled is a paid mutator transaction binding the contract method 0x895620b7.
//
// Solidity: function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) returns()
func (_NodesV2 *NodesV2Transactor) SetIsApiEnabled(opts *bind.TransactOpts, nodeId *big.Int, isApiEnabled bool) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setIsApiEnabled", nodeId, isApiEnabled)
}

// SetIsApiEnabled is a paid mutator transaction binding the contract method 0x895620b7.
//
// Solidity: function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) returns()
func (_NodesV2 *NodesV2Session) SetIsApiEnabled(nodeId *big.Int, isApiEnabled bool) (*types.Transaction, error) {
	return _NodesV2.Contract.SetIsApiEnabled(&_NodesV2.TransactOpts, nodeId, isApiEnabled)
}

// SetIsApiEnabled is a paid mutator transaction binding the contract method 0x895620b7.
//
// Solidity: function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) returns()
func (_NodesV2 *NodesV2TransactorSession) SetIsApiEnabled(nodeId *big.Int, isApiEnabled bool) (*types.Transaction, error) {
	return _NodesV2.Contract.SetIsApiEnabled(&_NodesV2.TransactOpts, nodeId, isApiEnabled)
}

// SetIsReplicationEnabled is a paid mutator transaction binding the contract method 0xe18cb254.
//
// Solidity: function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_NodesV2 *NodesV2Transactor) SetIsReplicationEnabled(opts *bind.TransactOpts, nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setIsReplicationEnabled", nodeId, isReplicationEnabled)
}

// SetIsReplicationEnabled is a paid mutator transaction binding the contract method 0xe18cb254.
//
// Solidity: function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_NodesV2 *NodesV2Session) SetIsReplicationEnabled(nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _NodesV2.Contract.SetIsReplicationEnabled(&_NodesV2.TransactOpts, nodeId, isReplicationEnabled)
}

// SetIsReplicationEnabled is a paid mutator transaction binding the contract method 0xe18cb254.
//
// Solidity: function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_NodesV2 *NodesV2TransactorSession) SetIsReplicationEnabled(nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _NodesV2.Contract.SetIsReplicationEnabled(&_NodesV2.TransactOpts, nodeId, isReplicationEnabled)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_NodesV2 *NodesV2Transactor) SetMaxActiveNodes(opts *bind.TransactOpts, newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setMaxActiveNodes", newMaxActiveNodes)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_NodesV2 *NodesV2Session) SetMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _NodesV2.Contract.SetMaxActiveNodes(&_NodesV2.TransactOpts, newMaxActiveNodes)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_NodesV2 *NodesV2TransactorSession) SetMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _NodesV2.Contract.SetMaxActiveNodes(&_NodesV2.TransactOpts, newMaxActiveNodes)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) returns()
func (_NodesV2 *NodesV2Transactor) SetMinMonthlyFee(opts *bind.TransactOpts, nodeId *big.Int, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setMinMonthlyFee", nodeId, minMonthlyFee)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) returns()
func (_NodesV2 *NodesV2Session) SetMinMonthlyFee(nodeId *big.Int, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.SetMinMonthlyFee(&_NodesV2.TransactOpts, nodeId, minMonthlyFee)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) returns()
func (_NodesV2 *NodesV2TransactorSession) SetMinMonthlyFee(nodeId *big.Int, minMonthlyFee *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.SetMinMonthlyFee(&_NodesV2.TransactOpts, nodeId, minMonthlyFee)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_NodesV2 *NodesV2Transactor) SetNodeOperatorCommissionPercent(opts *bind.TransactOpts, newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "setNodeOperatorCommissionPercent", newCommissionPercent)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_NodesV2 *NodesV2Session) SetNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.SetNodeOperatorCommissionPercent(&_NodesV2.TransactOpts, newCommissionPercent)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_NodesV2 *NodesV2TransactorSession) SetNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.SetNodeOperatorCommissionPercent(&_NodesV2.TransactOpts, newCommissionPercent)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_NodesV2 *NodesV2Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.contract.Transact(opts, "transferFrom", from, to, nodeId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_NodesV2 *NodesV2Session) TransferFrom(from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.TransferFrom(&_NodesV2.TransactOpts, from, to, nodeId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 nodeId) returns()
func (_NodesV2 *NodesV2TransactorSession) TransferFrom(from common.Address, to common.Address, nodeId *big.Int) (*types.Transaction, error) {
	return _NodesV2.Contract.TransferFrom(&_NodesV2.TransactOpts, from, to, nodeId)
}

// NodesV2ApiDisabledIterator is returned from FilterApiDisabled and is used to iterate over the raw logs and unpacked data for ApiDisabled events raised by the NodesV2 contract.
type NodesV2ApiDisabledIterator struct {
	Event *NodesV2ApiDisabled // Event containing the contract specifics and raw log

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
func (it *NodesV2ApiDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2ApiDisabled)
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
		it.Event = new(NodesV2ApiDisabled)
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
func (it *NodesV2ApiDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2ApiDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2ApiDisabled represents a ApiDisabled event raised by the NodesV2 contract.
type NodesV2ApiDisabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterApiDisabled is a free log retrieval operation binding the contract event 0x9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f3.
//
// Solidity: event ApiDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) FilterApiDisabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2ApiDisabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "ApiDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2ApiDisabledIterator{contract: _NodesV2.contract, event: "ApiDisabled", logs: logs, sub: sub}, nil
}

// WatchApiDisabled is a free log subscription operation binding the contract event 0x9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f3.
//
// Solidity: event ApiDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) WatchApiDisabled(opts *bind.WatchOpts, sink chan<- *NodesV2ApiDisabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "ApiDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2ApiDisabled)
				if err := _NodesV2.contract.UnpackLog(event, "ApiDisabled", log); err != nil {
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

// ParseApiDisabled is a log parse operation binding the contract event 0x9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f3.
//
// Solidity: event ApiDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) ParseApiDisabled(log types.Log) (*NodesV2ApiDisabled, error) {
	event := new(NodesV2ApiDisabled)
	if err := _NodesV2.contract.UnpackLog(event, "ApiDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2ApiEnabledIterator is returned from FilterApiEnabled and is used to iterate over the raw logs and unpacked data for ApiEnabled events raised by the NodesV2 contract.
type NodesV2ApiEnabledIterator struct {
	Event *NodesV2ApiEnabled // Event containing the contract specifics and raw log

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
func (it *NodesV2ApiEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2ApiEnabled)
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
		it.Event = new(NodesV2ApiEnabled)
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
func (it *NodesV2ApiEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2ApiEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2ApiEnabled represents a ApiEnabled event raised by the NodesV2 contract.
type NodesV2ApiEnabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterApiEnabled is a free log retrieval operation binding the contract event 0x3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da173.
//
// Solidity: event ApiEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) FilterApiEnabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2ApiEnabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "ApiEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2ApiEnabledIterator{contract: _NodesV2.contract, event: "ApiEnabled", logs: logs, sub: sub}, nil
}

// WatchApiEnabled is a free log subscription operation binding the contract event 0x3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da173.
//
// Solidity: event ApiEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) WatchApiEnabled(opts *bind.WatchOpts, sink chan<- *NodesV2ApiEnabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "ApiEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2ApiEnabled)
				if err := _NodesV2.contract.UnpackLog(event, "ApiEnabled", log); err != nil {
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

// ParseApiEnabled is a log parse operation binding the contract event 0x3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da173.
//
// Solidity: event ApiEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) ParseApiEnabled(log types.Log) (*NodesV2ApiEnabled, error) {
	event := new(NodesV2ApiEnabled)
	if err := _NodesV2.contract.UnpackLog(event, "ApiEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the NodesV2 contract.
type NodesV2ApprovalIterator struct {
	Event *NodesV2Approval // Event containing the contract specifics and raw log

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
func (it *NodesV2ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2Approval)
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
		it.Event = new(NodesV2Approval)
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
func (it *NodesV2ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2Approval represents a Approval event raised by the NodesV2 contract.
type NodesV2Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodesV2 *NodesV2Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NodesV2ApprovalIterator, error) {

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

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2ApprovalIterator{contract: _NodesV2.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodesV2 *NodesV2Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NodesV2Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2Approval)
				if err := _NodesV2.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseApproval(log types.Log) (*NodesV2Approval, error) {
	event := new(NodesV2Approval)
	if err := _NodesV2.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the NodesV2 contract.
type NodesV2ApprovalForAllIterator struct {
	Event *NodesV2ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *NodesV2ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2ApprovalForAll)
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
		it.Event = new(NodesV2ApprovalForAll)
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
func (it *NodesV2ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2ApprovalForAll represents a ApprovalForAll event raised by the NodesV2 contract.
type NodesV2ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodesV2 *NodesV2Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*NodesV2ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2ApprovalForAllIterator{contract: _NodesV2.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodesV2 *NodesV2Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NodesV2ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2ApprovalForAll)
				if err := _NodesV2.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseApprovalForAll(log types.Log) (*NodesV2ApprovalForAll, error) {
	event := new(NodesV2ApprovalForAll)
	if err := _NodesV2.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2BaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the NodesV2 contract.
type NodesV2BaseURIUpdatedIterator struct {
	Event *NodesV2BaseURIUpdated // Event containing the contract specifics and raw log

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
func (it *NodesV2BaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2BaseURIUpdated)
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
		it.Event = new(NodesV2BaseURIUpdated)
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
func (it *NodesV2BaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2BaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2BaseURIUpdated represents a BaseURIUpdated event raised by the NodesV2 contract.
type NodesV2BaseURIUpdated struct {
	NewBaseURI string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_NodesV2 *NodesV2Filterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*NodesV2BaseURIUpdatedIterator, error) {

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesV2BaseURIUpdatedIterator{contract: _NodesV2.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_NodesV2 *NodesV2Filterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *NodesV2BaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2BaseURIUpdated)
				if err := _NodesV2.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseBaseURIUpdated(log types.Log) (*NodesV2BaseURIUpdated, error) {
	event := new(NodesV2BaseURIUpdated)
	if err := _NodesV2.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2DefaultAdminDelayChangeCanceledIterator is returned from FilterDefaultAdminDelayChangeCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeCanceled events raised by the NodesV2 contract.
type NodesV2DefaultAdminDelayChangeCanceledIterator struct {
	Event *NodesV2DefaultAdminDelayChangeCanceled // Event containing the contract specifics and raw log

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
func (it *NodesV2DefaultAdminDelayChangeCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2DefaultAdminDelayChangeCanceled)
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
		it.Event = new(NodesV2DefaultAdminDelayChangeCanceled)
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
func (it *NodesV2DefaultAdminDelayChangeCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2DefaultAdminDelayChangeCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2DefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the NodesV2 contract.
type NodesV2DefaultAdminDelayChangeCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeCanceled is a free log retrieval operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_NodesV2 *NodesV2Filterer) FilterDefaultAdminDelayChangeCanceled(opts *bind.FilterOpts) (*NodesV2DefaultAdminDelayChangeCanceledIterator, error) {

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return &NodesV2DefaultAdminDelayChangeCanceledIterator{contract: _NodesV2.contract, event: "DefaultAdminDelayChangeCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeCanceled is a free log subscription operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_NodesV2 *NodesV2Filterer) WatchDefaultAdminDelayChangeCanceled(opts *bind.WatchOpts, sink chan<- *NodesV2DefaultAdminDelayChangeCanceled) (event.Subscription, error) {

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2DefaultAdminDelayChangeCanceled)
				if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseDefaultAdminDelayChangeCanceled(log types.Log) (*NodesV2DefaultAdminDelayChangeCanceled, error) {
	event := new(NodesV2DefaultAdminDelayChangeCanceled)
	if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2DefaultAdminDelayChangeScheduledIterator is returned from FilterDefaultAdminDelayChangeScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeScheduled events raised by the NodesV2 contract.
type NodesV2DefaultAdminDelayChangeScheduledIterator struct {
	Event *NodesV2DefaultAdminDelayChangeScheduled // Event containing the contract specifics and raw log

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
func (it *NodesV2DefaultAdminDelayChangeScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2DefaultAdminDelayChangeScheduled)
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
		it.Event = new(NodesV2DefaultAdminDelayChangeScheduled)
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
func (it *NodesV2DefaultAdminDelayChangeScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2DefaultAdminDelayChangeScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2DefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the NodesV2 contract.
type NodesV2DefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeScheduled is a free log retrieval operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_NodesV2 *NodesV2Filterer) FilterDefaultAdminDelayChangeScheduled(opts *bind.FilterOpts) (*NodesV2DefaultAdminDelayChangeScheduledIterator, error) {

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return &NodesV2DefaultAdminDelayChangeScheduledIterator{contract: _NodesV2.contract, event: "DefaultAdminDelayChangeScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeScheduled is a free log subscription operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_NodesV2 *NodesV2Filterer) WatchDefaultAdminDelayChangeScheduled(opts *bind.WatchOpts, sink chan<- *NodesV2DefaultAdminDelayChangeScheduled) (event.Subscription, error) {

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2DefaultAdminDelayChangeScheduled)
				if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseDefaultAdminDelayChangeScheduled(log types.Log) (*NodesV2DefaultAdminDelayChangeScheduled, error) {
	event := new(NodesV2DefaultAdminDelayChangeScheduled)
	if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2DefaultAdminTransferCanceledIterator is returned from FilterDefaultAdminTransferCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferCanceled events raised by the NodesV2 contract.
type NodesV2DefaultAdminTransferCanceledIterator struct {
	Event *NodesV2DefaultAdminTransferCanceled // Event containing the contract specifics and raw log

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
func (it *NodesV2DefaultAdminTransferCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2DefaultAdminTransferCanceled)
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
		it.Event = new(NodesV2DefaultAdminTransferCanceled)
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
func (it *NodesV2DefaultAdminTransferCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2DefaultAdminTransferCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2DefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the NodesV2 contract.
type NodesV2DefaultAdminTransferCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferCanceled is a free log retrieval operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_NodesV2 *NodesV2Filterer) FilterDefaultAdminTransferCanceled(opts *bind.FilterOpts) (*NodesV2DefaultAdminTransferCanceledIterator, error) {

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return &NodesV2DefaultAdminTransferCanceledIterator{contract: _NodesV2.contract, event: "DefaultAdminTransferCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferCanceled is a free log subscription operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_NodesV2 *NodesV2Filterer) WatchDefaultAdminTransferCanceled(opts *bind.WatchOpts, sink chan<- *NodesV2DefaultAdminTransferCanceled) (event.Subscription, error) {

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2DefaultAdminTransferCanceled)
				if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseDefaultAdminTransferCanceled(log types.Log) (*NodesV2DefaultAdminTransferCanceled, error) {
	event := new(NodesV2DefaultAdminTransferCanceled)
	if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2DefaultAdminTransferScheduledIterator is returned from FilterDefaultAdminTransferScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferScheduled events raised by the NodesV2 contract.
type NodesV2DefaultAdminTransferScheduledIterator struct {
	Event *NodesV2DefaultAdminTransferScheduled // Event containing the contract specifics and raw log

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
func (it *NodesV2DefaultAdminTransferScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2DefaultAdminTransferScheduled)
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
		it.Event = new(NodesV2DefaultAdminTransferScheduled)
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
func (it *NodesV2DefaultAdminTransferScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2DefaultAdminTransferScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2DefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the NodesV2 contract.
type NodesV2DefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferScheduled is a free log retrieval operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_NodesV2 *NodesV2Filterer) FilterDefaultAdminTransferScheduled(opts *bind.FilterOpts, newAdmin []common.Address) (*NodesV2DefaultAdminTransferScheduledIterator, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2DefaultAdminTransferScheduledIterator{contract: _NodesV2.contract, event: "DefaultAdminTransferScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferScheduled is a free log subscription operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_NodesV2 *NodesV2Filterer) WatchDefaultAdminTransferScheduled(opts *bind.WatchOpts, sink chan<- *NodesV2DefaultAdminTransferScheduled, newAdmin []common.Address) (event.Subscription, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2DefaultAdminTransferScheduled)
				if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseDefaultAdminTransferScheduled(log types.Log) (*NodesV2DefaultAdminTransferScheduled, error) {
	event := new(NodesV2DefaultAdminTransferScheduled)
	if err := _NodesV2.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2HttpAddressUpdatedIterator is returned from FilterHttpAddressUpdated and is used to iterate over the raw logs and unpacked data for HttpAddressUpdated events raised by the NodesV2 contract.
type NodesV2HttpAddressUpdatedIterator struct {
	Event *NodesV2HttpAddressUpdated // Event containing the contract specifics and raw log

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
func (it *NodesV2HttpAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2HttpAddressUpdated)
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
		it.Event = new(NodesV2HttpAddressUpdated)
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
func (it *NodesV2HttpAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2HttpAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2HttpAddressUpdated represents a HttpAddressUpdated event raised by the NodesV2 contract.
type NodesV2HttpAddressUpdated struct {
	NodeId         *big.Int
	NewHttpAddress string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterHttpAddressUpdated is a free log retrieval operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_NodesV2 *NodesV2Filterer) FilterHttpAddressUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2HttpAddressUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2HttpAddressUpdatedIterator{contract: _NodesV2.contract, event: "HttpAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchHttpAddressUpdated is a free log subscription operation binding the contract event 0x15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed.
//
// Solidity: event HttpAddressUpdated(uint256 indexed nodeId, string newHttpAddress)
func (_NodesV2 *NodesV2Filterer) WatchHttpAddressUpdated(opts *bind.WatchOpts, sink chan<- *NodesV2HttpAddressUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2HttpAddressUpdated)
				if err := _NodesV2.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseHttpAddressUpdated(log types.Log) (*NodesV2HttpAddressUpdated, error) {
	event := new(NodesV2HttpAddressUpdated)
	if err := _NodesV2.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2MaxActiveNodesUpdatedIterator is returned from FilterMaxActiveNodesUpdated and is used to iterate over the raw logs and unpacked data for MaxActiveNodesUpdated events raised by the NodesV2 contract.
type NodesV2MaxActiveNodesUpdatedIterator struct {
	Event *NodesV2MaxActiveNodesUpdated // Event containing the contract specifics and raw log

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
func (it *NodesV2MaxActiveNodesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2MaxActiveNodesUpdated)
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
		it.Event = new(NodesV2MaxActiveNodesUpdated)
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
func (it *NodesV2MaxActiveNodesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2MaxActiveNodesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2MaxActiveNodesUpdated represents a MaxActiveNodesUpdated event raised by the NodesV2 contract.
type NodesV2MaxActiveNodesUpdated struct {
	NewMaxActiveNodes uint8
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterMaxActiveNodesUpdated is a free log retrieval operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_NodesV2 *NodesV2Filterer) FilterMaxActiveNodesUpdated(opts *bind.FilterOpts) (*NodesV2MaxActiveNodesUpdatedIterator, error) {

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "MaxActiveNodesUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesV2MaxActiveNodesUpdatedIterator{contract: _NodesV2.contract, event: "MaxActiveNodesUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxActiveNodesUpdated is a free log subscription operation binding the contract event 0x6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821.
//
// Solidity: event MaxActiveNodesUpdated(uint8 newMaxActiveNodes)
func (_NodesV2 *NodesV2Filterer) WatchMaxActiveNodesUpdated(opts *bind.WatchOpts, sink chan<- *NodesV2MaxActiveNodesUpdated) (event.Subscription, error) {

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "MaxActiveNodesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2MaxActiveNodesUpdated)
				if err := _NodesV2.contract.UnpackLog(event, "MaxActiveNodesUpdated", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseMaxActiveNodesUpdated(log types.Log) (*NodesV2MaxActiveNodesUpdated, error) {
	event := new(NodesV2MaxActiveNodesUpdated)
	if err := _NodesV2.contract.UnpackLog(event, "MaxActiveNodesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2MinMonthlyFeeUpdatedIterator is returned from FilterMinMonthlyFeeUpdated and is used to iterate over the raw logs and unpacked data for MinMonthlyFeeUpdated events raised by the NodesV2 contract.
type NodesV2MinMonthlyFeeUpdatedIterator struct {
	Event *NodesV2MinMonthlyFeeUpdated // Event containing the contract specifics and raw log

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
func (it *NodesV2MinMonthlyFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2MinMonthlyFeeUpdated)
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
		it.Event = new(NodesV2MinMonthlyFeeUpdated)
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
func (it *NodesV2MinMonthlyFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2MinMonthlyFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2MinMonthlyFeeUpdated represents a MinMonthlyFeeUpdated event raised by the NodesV2 contract.
type NodesV2MinMonthlyFeeUpdated struct {
	NodeId        *big.Int
	MinMonthlyFee *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMinMonthlyFeeUpdated is a free log retrieval operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFee)
func (_NodesV2 *NodesV2Filterer) FilterMinMonthlyFeeUpdated(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2MinMonthlyFeeUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "MinMonthlyFeeUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2MinMonthlyFeeUpdatedIterator{contract: _NodesV2.contract, event: "MinMonthlyFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinMonthlyFeeUpdated is a free log subscription operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFee)
func (_NodesV2 *NodesV2Filterer) WatchMinMonthlyFeeUpdated(opts *bind.WatchOpts, sink chan<- *NodesV2MinMonthlyFeeUpdated, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "MinMonthlyFeeUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2MinMonthlyFeeUpdated)
				if err := _NodesV2.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseMinMonthlyFeeUpdated(log types.Log) (*NodesV2MinMonthlyFeeUpdated, error) {
	event := new(NodesV2MinMonthlyFeeUpdated)
	if err := _NodesV2.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2NodeAddedIterator is returned from FilterNodeAdded and is used to iterate over the raw logs and unpacked data for NodeAdded events raised by the NodesV2 contract.
type NodesV2NodeAddedIterator struct {
	Event *NodesV2NodeAdded // Event containing the contract specifics and raw log

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
func (it *NodesV2NodeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2NodeAdded)
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
		it.Event = new(NodesV2NodeAdded)
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
func (it *NodesV2NodeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2NodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2NodeAdded represents a NodeAdded event raised by the NodesV2 contract.
type NodesV2NodeAdded struct {
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
func (_NodesV2 *NodesV2Filterer) FilterNodeAdded(opts *bind.FilterOpts, nodeId []*big.Int, owner []common.Address) (*NodesV2NodeAddedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "NodeAdded", nodeIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2NodeAddedIterator{contract: _NodesV2.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

// WatchNodeAdded is a free log subscription operation binding the contract event 0x663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c.
//
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFee)
func (_NodesV2 *NodesV2Filterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *NodesV2NodeAdded, nodeId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "NodeAdded", nodeIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2NodeAdded)
				if err := _NodesV2.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseNodeAdded(log types.Log) (*NodesV2NodeAdded, error) {
	event := new(NodesV2NodeAdded)
	if err := _NodesV2.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2NodeDisabledIterator is returned from FilterNodeDisabled and is used to iterate over the raw logs and unpacked data for NodeDisabled events raised by the NodesV2 contract.
type NodesV2NodeDisabledIterator struct {
	Event *NodesV2NodeDisabled // Event containing the contract specifics and raw log

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
func (it *NodesV2NodeDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2NodeDisabled)
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
		it.Event = new(NodesV2NodeDisabled)
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
func (it *NodesV2NodeDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2NodeDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2NodeDisabled represents a NodeDisabled event raised by the NodesV2 contract.
type NodesV2NodeDisabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeDisabled is a free log retrieval operation binding the contract event 0xa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8.
//
// Solidity: event NodeDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) FilterNodeDisabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2NodeDisabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "NodeDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2NodeDisabledIterator{contract: _NodesV2.contract, event: "NodeDisabled", logs: logs, sub: sub}, nil
}

// WatchNodeDisabled is a free log subscription operation binding the contract event 0xa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8.
//
// Solidity: event NodeDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) WatchNodeDisabled(opts *bind.WatchOpts, sink chan<- *NodesV2NodeDisabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "NodeDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2NodeDisabled)
				if err := _NodesV2.contract.UnpackLog(event, "NodeDisabled", log); err != nil {
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

// ParseNodeDisabled is a log parse operation binding the contract event 0xa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8.
//
// Solidity: event NodeDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) ParseNodeDisabled(log types.Log) (*NodesV2NodeDisabled, error) {
	event := new(NodesV2NodeDisabled)
	if err := _NodesV2.contract.UnpackLog(event, "NodeDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2NodeEnabledIterator is returned from FilterNodeEnabled and is used to iterate over the raw logs and unpacked data for NodeEnabled events raised by the NodesV2 contract.
type NodesV2NodeEnabledIterator struct {
	Event *NodesV2NodeEnabled // Event containing the contract specifics and raw log

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
func (it *NodesV2NodeEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2NodeEnabled)
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
		it.Event = new(NodesV2NodeEnabled)
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
func (it *NodesV2NodeEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2NodeEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2NodeEnabled represents a NodeEnabled event raised by the NodesV2 contract.
type NodesV2NodeEnabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeEnabled is a free log retrieval operation binding the contract event 0xf044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a.
//
// Solidity: event NodeEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) FilterNodeEnabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2NodeEnabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "NodeEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2NodeEnabledIterator{contract: _NodesV2.contract, event: "NodeEnabled", logs: logs, sub: sub}, nil
}

// WatchNodeEnabled is a free log subscription operation binding the contract event 0xf044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a.
//
// Solidity: event NodeEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) WatchNodeEnabled(opts *bind.WatchOpts, sink chan<- *NodesV2NodeEnabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "NodeEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2NodeEnabled)
				if err := _NodesV2.contract.UnpackLog(event, "NodeEnabled", log); err != nil {
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

// ParseNodeEnabled is a log parse operation binding the contract event 0xf044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a.
//
// Solidity: event NodeEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) ParseNodeEnabled(log types.Log) (*NodesV2NodeEnabled, error) {
	event := new(NodesV2NodeEnabled)
	if err := _NodesV2.contract.UnpackLog(event, "NodeEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2NodeOperatorCommissionPercentUpdatedIterator is returned from FilterNodeOperatorCommissionPercentUpdated and is used to iterate over the raw logs and unpacked data for NodeOperatorCommissionPercentUpdated events raised by the NodesV2 contract.
type NodesV2NodeOperatorCommissionPercentUpdatedIterator struct {
	Event *NodesV2NodeOperatorCommissionPercentUpdated // Event containing the contract specifics and raw log

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
func (it *NodesV2NodeOperatorCommissionPercentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2NodeOperatorCommissionPercentUpdated)
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
		it.Event = new(NodesV2NodeOperatorCommissionPercentUpdated)
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
func (it *NodesV2NodeOperatorCommissionPercentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2NodeOperatorCommissionPercentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2NodeOperatorCommissionPercentUpdated represents a NodeOperatorCommissionPercentUpdated event raised by the NodesV2 contract.
type NodesV2NodeOperatorCommissionPercentUpdated struct {
	NewCommissionPercent *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterNodeOperatorCommissionPercentUpdated is a free log retrieval operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_NodesV2 *NodesV2Filterer) FilterNodeOperatorCommissionPercentUpdated(opts *bind.FilterOpts) (*NodesV2NodeOperatorCommissionPercentUpdatedIterator, error) {

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "NodeOperatorCommissionPercentUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesV2NodeOperatorCommissionPercentUpdatedIterator{contract: _NodesV2.contract, event: "NodeOperatorCommissionPercentUpdated", logs: logs, sub: sub}, nil
}

// WatchNodeOperatorCommissionPercentUpdated is a free log subscription operation binding the contract event 0x6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb0.
//
// Solidity: event NodeOperatorCommissionPercentUpdated(uint256 newCommissionPercent)
func (_NodesV2 *NodesV2Filterer) WatchNodeOperatorCommissionPercentUpdated(opts *bind.WatchOpts, sink chan<- *NodesV2NodeOperatorCommissionPercentUpdated) (event.Subscription, error) {

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "NodeOperatorCommissionPercentUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2NodeOperatorCommissionPercentUpdated)
				if err := _NodesV2.contract.UnpackLog(event, "NodeOperatorCommissionPercentUpdated", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseNodeOperatorCommissionPercentUpdated(log types.Log) (*NodesV2NodeOperatorCommissionPercentUpdated, error) {
	event := new(NodesV2NodeOperatorCommissionPercentUpdated)
	if err := _NodesV2.contract.UnpackLog(event, "NodeOperatorCommissionPercentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2NodeTransferredIterator is returned from FilterNodeTransferred and is used to iterate over the raw logs and unpacked data for NodeTransferred events raised by the NodesV2 contract.
type NodesV2NodeTransferredIterator struct {
	Event *NodesV2NodeTransferred // Event containing the contract specifics and raw log

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
func (it *NodesV2NodeTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2NodeTransferred)
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
		it.Event = new(NodesV2NodeTransferred)
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
func (it *NodesV2NodeTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2NodeTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2NodeTransferred represents a NodeTransferred event raised by the NodesV2 contract.
type NodesV2NodeTransferred struct {
	NodeId *big.Int
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeTransferred is a free log retrieval operation binding the contract event 0x0080108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be72.
//
// Solidity: event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to)
func (_NodesV2 *NodesV2Filterer) FilterNodeTransferred(opts *bind.FilterOpts, nodeId []*big.Int, from []common.Address, to []common.Address) (*NodesV2NodeTransferredIterator, error) {

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

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "NodeTransferred", nodeIdRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2NodeTransferredIterator{contract: _NodesV2.contract, event: "NodeTransferred", logs: logs, sub: sub}, nil
}

// WatchNodeTransferred is a free log subscription operation binding the contract event 0x0080108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be72.
//
// Solidity: event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to)
func (_NodesV2 *NodesV2Filterer) WatchNodeTransferred(opts *bind.WatchOpts, sink chan<- *NodesV2NodeTransferred, nodeId []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "NodeTransferred", nodeIdRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2NodeTransferred)
				if err := _NodesV2.contract.UnpackLog(event, "NodeTransferred", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseNodeTransferred(log types.Log) (*NodesV2NodeTransferred, error) {
	event := new(NodesV2NodeTransferred)
	if err := _NodesV2.contract.UnpackLog(event, "NodeTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2ReplicationDisabledIterator is returned from FilterReplicationDisabled and is used to iterate over the raw logs and unpacked data for ReplicationDisabled events raised by the NodesV2 contract.
type NodesV2ReplicationDisabledIterator struct {
	Event *NodesV2ReplicationDisabled // Event containing the contract specifics and raw log

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
func (it *NodesV2ReplicationDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2ReplicationDisabled)
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
		it.Event = new(NodesV2ReplicationDisabled)
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
func (it *NodesV2ReplicationDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2ReplicationDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2ReplicationDisabled represents a ReplicationDisabled event raised by the NodesV2 contract.
type NodesV2ReplicationDisabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReplicationDisabled is a free log retrieval operation binding the contract event 0xa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce280761808.
//
// Solidity: event ReplicationDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) FilterReplicationDisabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2ReplicationDisabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "ReplicationDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2ReplicationDisabledIterator{contract: _NodesV2.contract, event: "ReplicationDisabled", logs: logs, sub: sub}, nil
}

// WatchReplicationDisabled is a free log subscription operation binding the contract event 0xa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce280761808.
//
// Solidity: event ReplicationDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) WatchReplicationDisabled(opts *bind.WatchOpts, sink chan<- *NodesV2ReplicationDisabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "ReplicationDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2ReplicationDisabled)
				if err := _NodesV2.contract.UnpackLog(event, "ReplicationDisabled", log); err != nil {
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

// ParseReplicationDisabled is a log parse operation binding the contract event 0xa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce280761808.
//
// Solidity: event ReplicationDisabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) ParseReplicationDisabled(log types.Log) (*NodesV2ReplicationDisabled, error) {
	event := new(NodesV2ReplicationDisabled)
	if err := _NodesV2.contract.UnpackLog(event, "ReplicationDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2ReplicationEnabledIterator is returned from FilterReplicationEnabled and is used to iterate over the raw logs and unpacked data for ReplicationEnabled events raised by the NodesV2 contract.
type NodesV2ReplicationEnabledIterator struct {
	Event *NodesV2ReplicationEnabled // Event containing the contract specifics and raw log

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
func (it *NodesV2ReplicationEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2ReplicationEnabled)
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
		it.Event = new(NodesV2ReplicationEnabled)
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
func (it *NodesV2ReplicationEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2ReplicationEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2ReplicationEnabled represents a ReplicationEnabled event raised by the NodesV2 contract.
type NodesV2ReplicationEnabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReplicationEnabled is a free log retrieval operation binding the contract event 0xd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d0.
//
// Solidity: event ReplicationEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) FilterReplicationEnabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesV2ReplicationEnabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "ReplicationEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2ReplicationEnabledIterator{contract: _NodesV2.contract, event: "ReplicationEnabled", logs: logs, sub: sub}, nil
}

// WatchReplicationEnabled is a free log subscription operation binding the contract event 0xd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d0.
//
// Solidity: event ReplicationEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) WatchReplicationEnabled(opts *bind.WatchOpts, sink chan<- *NodesV2ReplicationEnabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "ReplicationEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2ReplicationEnabled)
				if err := _NodesV2.contract.UnpackLog(event, "ReplicationEnabled", log); err != nil {
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

// ParseReplicationEnabled is a log parse operation binding the contract event 0xd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d0.
//
// Solidity: event ReplicationEnabled(uint256 indexed nodeId)
func (_NodesV2 *NodesV2Filterer) ParseReplicationEnabled(log types.Log) (*NodesV2ReplicationEnabled, error) {
	event := new(NodesV2ReplicationEnabled)
	if err := _NodesV2.contract.UnpackLog(event, "ReplicationEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the NodesV2 contract.
type NodesV2RoleAdminChangedIterator struct {
	Event *NodesV2RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *NodesV2RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2RoleAdminChanged)
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
		it.Event = new(NodesV2RoleAdminChanged)
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
func (it *NodesV2RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2RoleAdminChanged represents a RoleAdminChanged event raised by the NodesV2 contract.
type NodesV2RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_NodesV2 *NodesV2Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*NodesV2RoleAdminChangedIterator, error) {

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

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2RoleAdminChangedIterator{contract: _NodesV2.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_NodesV2 *NodesV2Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *NodesV2RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2RoleAdminChanged)
				if err := _NodesV2.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseRoleAdminChanged(log types.Log) (*NodesV2RoleAdminChanged, error) {
	event := new(NodesV2RoleAdminChanged)
	if err := _NodesV2.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the NodesV2 contract.
type NodesV2RoleGrantedIterator struct {
	Event *NodesV2RoleGranted // Event containing the contract specifics and raw log

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
func (it *NodesV2RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2RoleGranted)
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
		it.Event = new(NodesV2RoleGranted)
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
func (it *NodesV2RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2RoleGranted represents a RoleGranted event raised by the NodesV2 contract.
type NodesV2RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodesV2 *NodesV2Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*NodesV2RoleGrantedIterator, error) {

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

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2RoleGrantedIterator{contract: _NodesV2.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodesV2 *NodesV2Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *NodesV2RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2RoleGranted)
				if err := _NodesV2.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseRoleGranted(log types.Log) (*NodesV2RoleGranted, error) {
	event := new(NodesV2RoleGranted)
	if err := _NodesV2.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the NodesV2 contract.
type NodesV2RoleRevokedIterator struct {
	Event *NodesV2RoleRevoked // Event containing the contract specifics and raw log

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
func (it *NodesV2RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2RoleRevoked)
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
		it.Event = new(NodesV2RoleRevoked)
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
func (it *NodesV2RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2RoleRevoked represents a RoleRevoked event raised by the NodesV2 contract.
type NodesV2RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodesV2 *NodesV2Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*NodesV2RoleRevokedIterator, error) {

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

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2RoleRevokedIterator{contract: _NodesV2.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_NodesV2 *NodesV2Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *NodesV2RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2RoleRevoked)
				if err := _NodesV2.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseRoleRevoked(log types.Log) (*NodesV2RoleRevoked, error) {
	event := new(NodesV2RoleRevoked)
	if err := _NodesV2.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesV2TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the NodesV2 contract.
type NodesV2TransferIterator struct {
	Event *NodesV2Transfer // Event containing the contract specifics and raw log

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
func (it *NodesV2TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesV2Transfer)
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
		it.Event = new(NodesV2Transfer)
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
func (it *NodesV2TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesV2TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesV2Transfer represents a Transfer event raised by the NodesV2 contract.
type NodesV2Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodesV2 *NodesV2Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NodesV2TransferIterator, error) {

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

	logs, sub, err := _NodesV2.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesV2TransferIterator{contract: _NodesV2.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodesV2 *NodesV2Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NodesV2Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _NodesV2.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesV2Transfer)
				if err := _NodesV2.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_NodesV2 *NodesV2Filterer) ParseTransfer(log types.Log) (*NodesV2Transfer, error) {
	event := new(NodesV2Transfer)
	if err := _NodesV2.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
