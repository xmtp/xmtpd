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
	SigningKeyPub             []byte
	HttpAddress               string
	IsReplicationEnabled      bool
	IsApiEnabled              bool
	IsDisabled                bool
	MinMonthlyFeeMicroDollars *big.Int
}

// INodesNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type INodesNodeWithId struct {
	NodeId *big.Int
	Node   INodesNode
}

// NodesMetaData contains all meta data concerning the Nodes contract.
var NodesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_INCREMENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disableNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enableNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveApiNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveApiNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveApiNodesIDs\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesIDs\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveReplicationNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveReplicationNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveReplicationNodesIDs\",\"inputs\":[],\"outputs\":[{\"name\":\"activeNodesIDs\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"allNodes\",\"type\":\"tuple[]\",\"internalType\":\"structINodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApiNodeIsActive\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeOperatorCommissionPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"commissionPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReplicationNodeIsActive\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxActiveNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeOperatorCommissionPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeFromApiNodes\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeFromReplicationNodes\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setHttpAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIsApiEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isApiEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIsReplicationEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isReplicationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxActiveNodes\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinMonthlyFee\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNodeOperatorCommissionPercent\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"supported\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ApiDisabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApiEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BaseURIUpdated\",\"inputs\":[{\"name\":\"newBaseURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HttpAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newHttpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxActiveNodesUpdated\",\"inputs\":[{\"name\":\"newMaxActiveNodes\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinMonthlyFeeUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAdded\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"minMonthlyFeeMicroDollars\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeDisabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeOperatorCommissionPercentUpdated\",\"inputs\":[{\"name\":\"newCommissionPercent\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeTransferred\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplicationDisabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReplicationEnabled\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCommissionPercent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHttpAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInputLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNodeConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigningKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesBelowCurrentCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeIsDisabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x6080604052600a805464ffffffffff19166014179055348015610020575f5ffd5b5060405161473638038061473683398101604081905261003f91610317565b60408051808201825260128152712c26aa28102737b2329027b832b930ba37b960711b602080830191909152825180840190935260048352630584d54560e41b90830152906202a300836001600160a01b0381166100b657604051636116401160e11b81525f600482015260240160405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100df5f8261018b565b50600391506100f0905083826103dc565b5060046100fd82826103dc565b5050506001600160a01b0381166101275760405163e6c4247b60e01b815260040160405180910390fd5b61013e5f5160206147165f395f51905f525f6101fa565b6101555f5160206146f65f395f51905f525f6101fa565b61016c5f5160206147165f395f51905f528261018b565b506101845f5160206146f65f395f51905f528261018b565b5050610496565b5f826101e7575f6101a46002546001600160a01b031690565b6001600160a01b0316146101cb57604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6101f18383610226565b90505b92915050565b8161021857604051631fe1e13d60e11b815260040160405180910390fd5b61022282826102cd565b5050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff166102c6575f838152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905561027e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016101f4565b505f6101f4565b5f82815260208190526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b5f60208284031215610327575f5ffd5b81516001600160a01b038116811461033d575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061036c57607f821691505b60208210810361038a57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156103d757805f5260205f20601f840160051c810160208510156103b55750805b601f840160051c820191505b818110156103d4575f81556001016103c1565b50505b505050565b81516001600160401b038111156103f5576103f5610344565b610409816104038454610358565b84610390565b6020601f82116001811461043b575f83156104245750848201515b5f19600385901b1c1916600184901b1784556103d4565b5f84815260208120601f198516915b8281101561046a578785015182556020948501946001909201910161044a565b508482101561048757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b614253806104a35f395ff3fe608060405234801561000f575f5ffd5b5060043610610393575f3560e01c80638ed9ea34116101df578063ce99948911610109578063e18cb254116100a9578063f579d7e111610079578063f579d7e114610860578063fb1120e214610868578063fd667d1e14610870578063fd967f471461088d575f5ffd5b8063e18cb25414610801578063e985e9c514610814578063ebe487bf1461084f578063f3194a3914610857575f5ffd5b8063d547741f116100e4578063d547741f146107ac578063d59f9fe0146107bf578063d602b9fd146107e6578063d74a2a50146107ee575f5ffd5b8063ce99948914610752578063cefc142914610765578063cf6eefb71461076d575f5ffd5b8063a217fddf1161017f578063b9b140d61161014f578063b9b140d61461071c578063c4741f3114610724578063c87b56dd14610737578063cc8463c81461074a575f5ffd5b8063a217fddf146106dc578063a22cb465146106e3578063a835f88e146106f6578063b88d4fde14610709575f5ffd5b806395d89b41116101ba57806395d89b41146106865780639d32f9ba1461068e578063a1174e7d146106ad578063a1eda53c146106b5575f5ffd5b80638ed9ea34146106285780638fbbf6231461063b57806391d1485414610650575f5ffd5b806344ff624e116102c0578063649a5ec71161026057806379e0d58c1161023057806379e0d58c146105e957806384ef8ffc146105fc578063895620b71461060d5780638da5cb5b14610620575f5ffd5b8063649a5ec7146105895780636ec97bfc1461059c57806370a08231146105af57806375b238fc146105c2575f5ffd5b806355f804b31161029b57806355f804b31461053b578063634e93da1461054e5780636352211e14610561578063646453ba14610574575f5ffd5b806344ff624e146104f55780634f0f4aa91461050857806350d0215f14610528575f5ffd5b8063203ede77116103365780632f2ff15d116103065780632f2ff15d146104a957806336568abe146104bc5780633d2853fb146104cf57806342842e0e146104e2575f5ffd5b8063203ede771461044e57806321fbd7cb1461046157806323b872dd14610474578063248a9ca314610487575f5ffd5b8063081812fc11610371578063081812fc146103f0578063095ea7b31461041b5780630aa6220b1461043057806317e3b3a914610438575f5ffd5b806301ffc9a714610397578063022d63fb146103bf57806306fdde03146103db575b5f5ffd5b6103aa6103a5366004613819565b610896565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff90911681526020016103b6565b6103e36108f1565b6040516103b69190613862565b6104036103fe366004613874565b610981565b6040516001600160a01b0390911681526020016103b6565b61042e6104293660046138a6565b6109a8565b005b61042e6109b7565b6104406109cc565b6040519081526020016103b6565b61042e61045c366004613874565b6109dc565b6103aa61046f366004613874565b610a18565b61042e6104823660046138ce565b610a24565b610440610495366004613874565b5f9081526020819052604090206001015490565b61042e6104b7366004613908565b610ab1565b61042e6104ca366004613908565b610af2565b61042e6104dd366004613874565b610be2565b61042e6104f03660046138ce565b610c76565b6103aa610503366004613874565b610c95565b61051b610516366004613874565b610ca1565b6040516103b6919061399a565b600a54610100900463ffffffff16610440565b61042e6105493660046139f1565b610e51565b61042e61055c366004613a30565b610f9e565b61040361056f366004613874565b610fb1565b61057c610fbb565b6040516103b69190613a49565b61042e610597366004613a8b565b610fc7565b6104406105aa366004613ab0565b610fda565b6104406105bd366004613a30565b6112d6565b6104407fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b61042e6105f7366004613874565b611334565b6002546001600160a01b0316610403565b61042e61061b366004613b47565b611370565b61040361139a565b61042e610636366004613b68565b6113ad565b61064361147d565b6040516103b69190613b88565b6103aa61065e366004613908565b5f918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b6103e36116cf565b600a5461069b9060ff1681565b60405160ff90911681526020016103b6565b6106436116de565b6106bd611938565b6040805165ffffffffffff9384168152929091166020830152016103b6565b6104405f81565b61042e6106f1366004613c20565b6119b2565b61042e610704366004613874565b6119bd565b61042e610717366004613c75565b611a6b565b601054610440565b61042e610732366004613874565b611a89565b6103e3610745366004613874565b611b24565b6103c4611b89565b61042e610760366004613d53565b611c26565b61042e611cad565b600154604080516001600160a01b03831681527401000000000000000000000000000000000000000090920465ffffffffffff166020830152016103b6565b61042e6107ba366004613908565b611cfc565b6104407fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a81565b61042e611d3d565b61042e6107fc366004613d73565b611d4f565b61042e61080f366004613b47565b611e15565b6103aa610822366004613dbb565b6001600160a01b039182165f90815260086020908152604080832093909416825291909152205460ff1690565b610643611e3f565b61044060105481565b61044061208d565b61057c612098565b610878606481565b60405163ffffffff90911681526020016103b6565b61044061271081565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f5ad5d4b10000000000000000000000000000000000000000000000000000000014806108eb57506108eb826120a4565b92915050565b60606003805461090090613de3565b80601f016020809104026020016040519081016040528092919081815260200182805461092c90613de3565b80156109775780601f1061094e57610100808354040283529160200191610977565b820191905f5260205f20905b81548152906001019060200180831161095a57829003601f168201915b5050505050905090565b5f61098b82612145565b505f828152600760205260409020546001600160a01b03166108eb565b6109b3828233612196565b5050565b5f6109c1816121a3565b6109c96121ad565b50565b5f6109d7600e6121b9565b905090565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610a06816121a3565b610a0f826121c2565b6109b38261220f565b5f6108eb600c83612263565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a610a4e816121a3565b610a578261227a565b610a608261220f565b610a6b8484846122ec565b826001600160a01b0316846001600160a01b0316837e80108bb11ee8badd8a48ff0b4585853d721b6e5ac7e3415f99413dac52be7260405160405180910390a450505050565b81610ae8576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109b382826123a1565b81158015610b0d57506002546001600160a01b038281169116145b15610bd8576001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610b54575065ffffffffffff8116155b80610b6757504265ffffffffffff821610155b15610bad576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff821660048201526024015b60405180910390fd5b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b6109b382826123c5565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610c0c816121a3565b610c15826121c2565b5f828152600b602052604080822060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff1690555183917ff044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a91a25050565b610c9083838360405180602001604052805f815250611a6b565b505050565b5f6108eb600e83612263565b6040805160c0810182526060808252602082018190525f92820183905281018290526080810182905260a0810191909152610cdb826121c2565b5f828152600b602052604090819020815160c08101909252805482908290610d0290613de3565b80601f0160208091040260200160405190810160405280929190818152602001828054610d2e90613de3565b8015610d795780601f10610d5057610100808354040283529160200191610d79565b820191905f5260205f20905b815481529060010190602001808311610d5c57829003601f168201915b50505050508152602001600182018054610d9290613de3565b80601f0160208091040260200160405190810160405280929190818152602001828054610dbe90613de3565b8015610e095780601f10610de057610100808354040283529160200191610e09565b820191905f5260205f20905b815481529060010190602001808311610dec57829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015292915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610e7b816121a3565b81610eb2576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f2f000000000000000000000000000000000000000000000000000000000000008383610ee0600182613e61565b818110610eef57610eef613e74565b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610f52576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6009610f5f838583613ee5565b507f6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad8383604051610f91929190613fc8565b60405180910390a1505050565b5f610fa8816121a3565b6109b382612411565b5f6108eb82612145565b60606109d7600c612483565b5f610fd1816121a3565b6109b38261248f565b5f7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611005816121a3565b6001600160a01b038816611045576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8561107c576040517f8125403000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b836110b3576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6064600a600181819054906101000a900463ffffffff166110d390613fdb565b91906101000a81548163ffffffff021916908363ffffffff16021790556110fa9190613fff565b63ffffffff1691506040518060c0016040528088888080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250505090825250604080516020601f8901819004810282018101909252878152918101919088908890819084018382808284375f92018290525093855250505060208083018290526040808401839052606084018390526080909301879052858252600b905220815181906111b3908261401e565b50602082015160018201906111c8908261401e565b50604082015160028201805460608501516080860151151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff951515959095167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090931692909217939093179290921691909117905560a09091015160039091015561128188836124f7565b876001600160a01b0316827f663d98c1e2bdf874fcd4fadcdf16242719c434e099664a3eb574322b78bd7c5c89898989896040516112c39594939291906140d9565b60405180910390a3509695505050505050565b5f6001600160a01b038216611319576040517f89c62b640000000000000000000000000000000000000000000000000000000081525f6004820152602401610ba4565b506001600160a01b03165f9081526006602052604090205490565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561135e816121a3565b611367826121c2565b6109b38261227a565b611379826121c2565b6113828261258a565b61138b826125db565b8015611367576109b38261262a565b5f6109d76002546001600160a01b031690565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756113d7816121a3565b6113e1600c6121b9565b8260ff1610806113fc57506113f6600e6121b9565b8260ff16105b15611433576040517f39beadee00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a805460ff191660ff84169081179091556040519081527f6dd6623df488fb2b38fa153b12758a1b41c8e49e88025f8d9fb1eba1b8f1d821906020015b60405180910390a15050565b6060611489600e6121b9565b67ffffffffffffffff8111156114a1576114a1613c48565b6040519080825280602002602001820160405280156114da57816020015b6114c761379b565b8152602001906001900390816114bf5790505b5090505f5b6114e9600e6121b9565b8163ffffffff1610156116cb575f61150b600e63ffffffff808516906126e716565b90506040518060400160405280828152602001600b5f8481526020019081526020015f206040518060c00160405290815f8201805461154990613de3565b80601f016020809104026020016040519081016040528092919081815260200182805461157590613de3565b80156115c05780601f10611597576101008083540402835291602001916115c0565b820191905f5260205f20905b8154815290600101906020018083116115a357829003601f168201915b505050505081526020016001820180546115d990613de3565b80601f016020809104026020016040519081016040528092919081815260200182805461160590613de3565b80156116505780601f1061162757610100808354040283529160200191611650565b820191905f5260205f20905b81548152906001019060200180831161163357829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff85169081106116ae576116ae613e74565b602002602001018190525050806116c490613fdb565b90506114df565b5090565b60606004805461090090613de3565b600a54606090610100900463ffffffff1667ffffffffffffffff81111561170757611707613c48565b60405190808252806020026020018201604052801561174057816020015b61172d61379b565b8152602001906001900390816117255790505b5090505f5b600a5463ffffffff610100909104811690821610156116cb575f61176a826001614112565b611775906064613fff565b905060405180604001604052808263ffffffff168152602001600b5f8463ffffffff1681526020019081526020015f206040518060c00160405290815f820180546117bf90613de3565b80601f01602080910402602001604051908101604052809291908181526020018280546117eb90613de3565b80156118365780601f1061180d57610100808354040283529160200191611836565b820191905f5260205f20905b81548152906001019060200180831161181957829003601f168201915b5050505050815260200160018201805461184f90613de3565b80601f016020809104026020016040519081016040528092919081815260200182805461187b90613de3565b80156118c65780601f1061189d576101008083540402835291602001916118c6565b820191905f5260205f20905b8154815290600101906020018083116118a957829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff851690811061192457611924613e74565b602090810291909101015250600101611745565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff16801515801561197a57504265ffffffffffff821610155b611985575f5f6119aa565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b6109b33383836126f2565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756119e7816121a3565b6119f0826121c2565b5f828152600b6020526040902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff1662010000179055611a348261227a565b611a3d8261220f565b60405182907fa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8905f90a25050565b611a76848484610a24565b611a8333858585856127a9565b50505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611ab3816121a3565b612710821115611aef576040517f47d3b04600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60108290556040518281527f6367530104bc8677601bbb2f410055f5144865bf130b2c7bed1af5ff39185eb090602001611471565b6060611b2f82612145565b505f611b3961294d565b90505f815111611b575760405180602001604052805f815250611b82565b80611b618461295c565b604051602001611b72929190614145565b6040516020818303038152906040525b9392505050565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015611bca57504265ffffffffffff8216105b611bfc576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16611c20565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a611c50816121a3565b611c59836121c2565b5f838152600b6020526040908190206003018390555183907f27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a90611ca09085815260200190565b60405180910390a2505050565b6001546001600160a01b0316338114611cf4576040517fc22c8022000000000000000000000000000000000000000000000000000000008152336004820152602401610ba4565b6109c96129f9565b81611d33576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109b38282612ad0565b5f611d47816121a3565b6109c9612af4565b7fdaf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56a611d79816121a3565b611d82846121c2565b81611db9576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f848152600b60205260409020600101611dd4838583613ee5565b50837f15c3eac3b34037e402127abd35c3804f49d489c361f5bb8ff237544f0dfff4ed8484604051611e07929190613fc8565b60405180910390a250505050565b611e1e826121c2565b611e278261258a565b611e30826125db565b8015610a0f576109b382612afe565b6060611e4b600c6121b9565b67ffffffffffffffff811115611e6357611e63613c48565b604051908082528060200260200182016040528015611e9c57816020015b611e8961379b565b815260200190600190039081611e815790505b5090505f5b611eab600c6121b9565b8163ffffffff1610156116cb575f611ecd600c63ffffffff808516906126e716565b90506040518060400160405280828152602001600b5f8481526020019081526020015f206040518060c00160405290815f82018054611f0b90613de3565b80601f0160208091040260200160405190810160405280929190818152602001828054611f3790613de3565b8015611f825780601f10611f5957610100808354040283529160200191611f82565b820191905f5260205f20905b815481529060010190602001808311611f6557829003601f168201915b50505050508152602001600182018054611f9b90613de3565b80601f0160208091040260200160405190810160405280929190818152602001828054611fc790613de3565b80156120125780601f10611fe957610100808354040283529160200191612012565b820191905f5260205f20905b815481529060010190602001808311611ff557829003601f168201915b5050509183525050600282015460ff80821615156020840152610100820481161515604084015262010000909104161515606082015260039091015460809091015290528351849063ffffffff851690811061207057612070613e74565b6020026020010181905250508061208690613fdb565b9050611ea1565b5f6109d7600c6121b9565b60606109d7600e612483565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd00000000000000000000000000000000000000000000000000000000148061213657507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b5e139f00000000000000000000000000000000000000000000000000000000145b806108eb57506108eb82612b9c565b5f818152600560205260408120546001600160a01b0316806108eb576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101849052602401610ba4565b610c908383836001612bf1565b6109c98133612d44565b6121b75f5f612daf565b565b5f6108eb825490565b5f818152600560205260409020546001600160a01b03166109c9576040517f5e926f7100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61221a600e82612efb565b6122215750565b5f818152600b6020526040808220600201805460ff191690555182917fa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce28076180891a250565b5f8181526001830160205260408120541515611b82565b612285600c82612efb565b61228c5750565b5f818152600b602052604080822060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1690555182917f9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f391a250565b6001600160a01b03821661232e576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610ba4565b5f61233a838333612f06565b9050836001600160a01b0316816001600160a01b031614611a83576040517f64283d7b0000000000000000000000000000000000000000000000000000000081526001600160a01b0380861660048301526024820184905282166044820152606401610ba4565b5f828152602081905260409020600101546123bb816121a3565b611a838383613010565b6001600160a01b0381163314612407576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c9082826130a7565b5f61241a611b89565b612423426130fb565b61242d9190614159565b90506124398282613146565b60405165ffffffffffff821681526001600160a01b038316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b60605f611b82836131d4565b5f6124998261322d565b6124a2426130fb565b6124ac9190614159565b90506124b88282612daf565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b9101611471565b6001600160a01b038216612539576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610ba4565b5f61254583835f612f06565b90506001600160a01b03811615610c90576040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610ba4565b5f818152600b602052604090206002015462010000900460ff16156109c9576040517fc40c6f6300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f818152600560205260409020546001600160a01b031633146109c9576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a5460ff1661263a600c6121b9565b10612671576040517f950be9a500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61267c600c82613274565b6126835750565b5f818152600b602052604080822060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555182917f3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da17391a250565b5f611b82838361327f565b6001600160a01b03821661273d576040517f5b08ba180000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610ba4565b6001600160a01b038381165f81815260086020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b15612946576040517f150b7a020000000000000000000000000000000000000000000000000000000081526001600160a01b0384169063150b7a0290612804908890889087908790600401614177565b6020604051808303815f875af192505050801561283e575060408051601f3d908101601f1916820190925261283b918101906141b7565b60015b6128be573d80801561286b576040519150601f19603f3d011682016040523d82523d5f602084013e612870565b606091505b5080515f036128b6576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610ba4565b805181602001fd5b7fffffffff0000000000000000000000000000000000000000000000000000000081167f150b7a020000000000000000000000000000000000000000000000000000000014612944576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610ba4565b505b5050505050565b60606009805461090090613de3565b60605f612968836132a5565b60010190505f8167ffffffffffffffff81111561298757612987613c48565b6040519080825280601f01601f1916602001820160405280156129b1576020820181803683370190505b5090508181016020015b5f19017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85049450846129bb57509392505050565b6001546001600160a01b0381169074010000000000000000000000000000000000000000900465ffffffffffff16801580612a3c57504265ffffffffffff821610155b15612a7d576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610ba4565b612a985f612a936002546001600160a01b031690565b6130a7565b50612aa35f83613010565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f82815260208190526040902060010154612aea816121a3565b611a8383836130a7565b6121b75f5f613146565b600a5460ff16612b0e600e6121b9565b10612b45576040517f950be9a500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612b50600e82613274565b612b575750565b5f818152600b6020526040808220600201805460ff191660011790555182917fd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d091a250565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f314987860000000000000000000000000000000000000000000000000000000014806108eb57506108eb82613386565b8080612c0557506001600160a01b03821615155b15612cfd575f612c1484612145565b90506001600160a01b03831615801590612c405750826001600160a01b0316816001600160a01b031614155b8015612c7157506001600160a01b038082165f9081526008602090815260408083209387168352929052205460ff16155b15612cb3576040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526001600160a01b0384166004820152602401610ba4565b8115612cfb5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260076020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b5f828152602081815260408083206001600160a01b038516845290915290205460ff166109b3576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b038216600482015260248101839052604401610ba4565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015612e83574265ffffffffffff82161015612e5a576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055612e83565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b50600280546001600160a01b03167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f611b82838361341c565b5f828152600560205260408120546001600160a01b0390811690831615612f3257612f32818486613506565b6001600160a01b03811615612f6c57612f4d5f855f5f612bf1565b6001600160a01b0381165f90815260066020526040902080545f190190555b6001600160a01b03851615612f9a576001600160a01b0385165f908152600660205260409020805460010190555b5f8481526005602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b5f8261309d575f6130296002546001600160a01b031690565b6001600160a01b031614613069576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0384161790555b611b82838361359c565b5f821580156130c357506002546001600160a01b038381169116145b156130f157600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b611b828383613643565b5f65ffffffffffff8211156116cb576040517f6dfcc6500000000000000000000000000000000000000000000000000000000081526030600482015260248101839052604401610ba4565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff000000000000000000000000000000000000000000000000000084166001600160a01b03881617179093559004168015610c90576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b6060815f0180548060200260200160405190810160405280929190818152602001828054801561322157602002820191905f5260205f20905b81548152602001906001019080831161320d575b50505050509050919050565b5f5f613237611b89565b90508065ffffffffffff168365ffffffffffff161161325f5761325a83826141d2565b611b82565b611b8265ffffffffffff8416620697806136c4565b5f611b8283836136d3565b5f825f01828154811061329457613294613e74565b905f5260205f200154905092915050565b5f807a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083106132ed577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef81000000008310613319576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc10000831061333757662386f26fc10000830492506010015b6305f5e100831061334f576305f5e100830492506008015b612710831061336357612710830492506004015b60648310613375576064830492506002015b600a83106108eb5760010192915050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806108eb57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316146108eb565b5f81815260018301602052604081205480156134f6575f61343e600183613e61565b85549091505f9061345190600190613e61565b90508082146134b0575f865f01828154811061346f5761346f613e74565b905f5260205f200154905080875f01848154811061348f5761348f613e74565b5f918252602080832090910192909255918252600188019052604090208390555b85548690806134c1576134c16141f0565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f9055600193505050506108eb565b5f9150506108eb565b5092915050565b613511838383613718565b610c90576001600160a01b038316613558576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101829052602401610ba4565b6040517f177e802f0000000000000000000000000000000000000000000000000000000081526001600160a01b038316600482015260248101829052604401610ba4565b5f828152602081815260408083206001600160a01b038516845290915281205460ff1661363c575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556135f43390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016108eb565b505f6108eb565b5f828152602081815260408083206001600160a01b038516845290915281205460ff161561363c575f838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016108eb565b5f828218828410028218611b82565b5f81815260018301602052604081205461363c57508154600181810184555f8481526020808220909301849055845484825282860190935260409020919091556108eb565b5f6001600160a01b038316158015906137935750826001600160a01b0316846001600160a01b0316148061377057506001600160a01b038085165f9081526008602090815260408083209387168352929052205460ff165b8061379357505f828152600760205260409020546001600160a01b038481169116145b949350505050565b60405180604001604052805f81526020016137e76040518060c0016040528060608152602001606081526020015f151581526020015f151581526020015f151581526020015f81525090565b905290565b7fffffffff00000000000000000000000000000000000000000000000000000000811681146109c9575f5ffd5b5f60208284031215613829575f5ffd5b8135611b82816137ec565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f611b826020830184613834565b5f60208284031215613884575f5ffd5b5035919050565b80356001600160a01b03811681146138a1575f5ffd5b919050565b5f5f604083850312156138b7575f5ffd5b6138c08361388b565b946020939093013593505050565b5f5f5f606084860312156138e0575f5ffd5b6138e98461388b565b92506138f76020850161388b565b929592945050506040919091013590565b5f5f60408385031215613919575f5ffd5b823591506139296020840161388b565b90509250929050565b5f815160c0845261394660c0850182613834565b90506020830151848203602086015261395f8282613834565b91505060408301511515604085015260608301511515606085015260808301511515608085015260a083015160a08501528091505092915050565b602081525f611b826020830184613932565b5f5f83601f8401126139bc575f5ffd5b50813567ffffffffffffffff8111156139d3575f5ffd5b6020830191508360208285010111156139ea575f5ffd5b9250929050565b5f5f60208385031215613a02575f5ffd5b823567ffffffffffffffff811115613a18575f5ffd5b613a24858286016139ac565b90969095509350505050565b5f60208284031215613a40575f5ffd5b611b828261388b565b602080825282518282018190525f918401906040840190835b81811015613a80578351835260209384019390920191600101613a62565b509095945050505050565b5f60208284031215613a9b575f5ffd5b813565ffffffffffff81168114611b82575f5ffd5b5f5f5f5f5f5f60808789031215613ac5575f5ffd5b613ace8761388b565b9550602087013567ffffffffffffffff811115613ae9575f5ffd5b613af589828a016139ac565b909650945050604087013567ffffffffffffffff811115613b14575f5ffd5b613b2089828a016139ac565b979a9699509497949695606090950135949350505050565b803580151581146138a1575f5ffd5b5f5f60408385031215613b58575f5ffd5b8235915061392960208401613b38565b5f60208284031215613b78575f5ffd5b813560ff81168114611b82575f5ffd5b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b82811015613c14577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08786030184528151805186526020810151905060406020870152613bfe6040870182613932565b9550506020938401939190910190600101613bae565b50929695505050505050565b5f5f60408385031215613c31575f5ffd5b613c3a8361388b565b915061392960208401613b38565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f5f5f60808587031215613c88575f5ffd5b613c918561388b565b9350613c9f6020860161388b565b925060408501359150606085013567ffffffffffffffff811115613cc1575f5ffd5b8501601f81018713613cd1575f5ffd5b803567ffffffffffffffff811115613ceb57613ceb613c48565b604051601f19603f601f19601f8501160116810181811067ffffffffffffffff82111715613d1b57613d1b613c48565b604052818152828201602001891015613d32575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f60408385031215613d64575f5ffd5b50508035926020909101359150565b5f5f5f60408486031215613d85575f5ffd5b83359250602084013567ffffffffffffffff811115613da2575f5ffd5b613dae868287016139ac565b9497909650939450505050565b5f5f60408385031215613dcc575f5ffd5b613dd58361388b565b91506139296020840161388b565b600181811c90821680613df757607f821691505b602082108103613e2e577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b818103818111156108eb576108eb613e34565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b601f821115610c9057805f5260205f20601f840160051c81016020851015613ec65750805b601f840160051c820191505b81811015612946575f8155600101613ed2565b67ffffffffffffffff831115613efd57613efd613c48565b613f1183613f0b8354613de3565b83613ea1565b5f601f841160018114613f42575f8515613f2b5750838201355b5f19600387901b1c1916600186901b178355612946565b5f83815260208120601f198716915b82811015613f715786850135825560209485019460019092019101613f51565b5086821015613f8d575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b81835281816020850137505f602082840101525f6020601f19601f840116840101905092915050565b602081525f613793602083018486613f9f565b5f63ffffffff821663ffffffff8103613ff657613ff6613e34565b60010192915050565b63ffffffff81811683821602908116908181146134ff576134ff613e34565b815167ffffffffffffffff81111561403857614038613c48565b61404c816140468454613de3565b84613ea1565b6020601f82116001811461407e575f83156140675750848201515b5f19600385901b1c1916600184901b178455612946565b5f84815260208120601f198516915b828110156140ad578785015182556020948501946001909201910161408d565b50848210156140ca57868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b606081525f6140ec606083018789613f9f565b82810360208401526140ff818688613f9f565b9150508260408301529695505050505050565b63ffffffff81811683821601908111156108eb576108eb613e34565b5f81518060208401855e5f93019283525090919050565b5f613793614153838661412e565b8461412e565b65ffffffffffff81811683821601908111156108eb576108eb613e34565b6001600160a01b03851681526001600160a01b0384166020820152826040820152608060608201525f6141ad6080830184613834565b9695505050505050565b5f602082840312156141c7575f5ffd5b8151611b82816137ec565b65ffffffffffff82811682821603908111156108eb576108eb613e34565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffdfea2646970667358221220b6d571dd1526ffac8e29836c0c3ebcdb95ef57f1d5955f90f437e97b78dd9afe64736f6c634300081c0033daf9ac3a6308052428e8806fd908cf472318416ed7d78b3f35dd94bbbafde56aa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775",
}

// NodesABI is the input ABI used to generate the binding from.
// Deprecated: Use NodesMetaData.ABI instead.
var NodesABI = NodesMetaData.ABI

// NodesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodesMetaData.Bin instead.
var NodesBin = NodesMetaData.Bin

// DeployNodes deploys a new Ethereum contract, binding an instance of Nodes to it.
func DeployNodes(auth *bind.TransactOpts, backend bind.ContractBackend, initialAdmin common.Address) (common.Address, *types.Transaction, *Nodes, error) {
	parsed, err := NodesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodesBin), backend, initialAdmin)
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

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_Nodes *NodesCaller) NODEINCREMENT(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "NODE_INCREMENT")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_Nodes *NodesSession) NODEINCREMENT() (uint32, error) {
	return _Nodes.Contract.NODEINCREMENT(&_Nodes.CallOpts)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_Nodes *NodesCallerSession) NODEINCREMENT() (uint32, error) {
	return _Nodes.Contract.NODEINCREMENT(&_Nodes.CallOpts)
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

// GetActiveApiNodes is a free data retrieval call binding the contract method 0xebe487bf.
//
// Solidity: function getActiveApiNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_Nodes *NodesCaller) GetActiveApiNodes(opts *bind.CallOpts) ([]INodesNodeWithId, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveApiNodes")

	if err != nil {
		return *new([]INodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNodeWithId)).(*[]INodesNodeWithId)

	return out0, err

}

// GetActiveApiNodes is a free data retrieval call binding the contract method 0xebe487bf.
//
// Solidity: function getActiveApiNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_Nodes *NodesSession) GetActiveApiNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetActiveApiNodes(&_Nodes.CallOpts)
}

// GetActiveApiNodes is a free data retrieval call binding the contract method 0xebe487bf.
//
// Solidity: function getActiveApiNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_Nodes *NodesCallerSession) GetActiveApiNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetActiveApiNodes(&_Nodes.CallOpts)
}

// GetActiveApiNodesCount is a free data retrieval call binding the contract method 0xf579d7e1.
//
// Solidity: function getActiveApiNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesCaller) GetActiveApiNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveApiNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveApiNodesCount is a free data retrieval call binding the contract method 0xf579d7e1.
//
// Solidity: function getActiveApiNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesSession) GetActiveApiNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetActiveApiNodesCount(&_Nodes.CallOpts)
}

// GetActiveApiNodesCount is a free data retrieval call binding the contract method 0xf579d7e1.
//
// Solidity: function getActiveApiNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesCallerSession) GetActiveApiNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetActiveApiNodesCount(&_Nodes.CallOpts)
}

// GetActiveApiNodesIDs is a free data retrieval call binding the contract method 0x646453ba.
//
// Solidity: function getActiveApiNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesCaller) GetActiveApiNodesIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveApiNodesIDs")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveApiNodesIDs is a free data retrieval call binding the contract method 0x646453ba.
//
// Solidity: function getActiveApiNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesSession) GetActiveApiNodesIDs() ([]*big.Int, error) {
	return _Nodes.Contract.GetActiveApiNodesIDs(&_Nodes.CallOpts)
}

// GetActiveApiNodesIDs is a free data retrieval call binding the contract method 0x646453ba.
//
// Solidity: function getActiveApiNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesCallerSession) GetActiveApiNodesIDs() ([]*big.Int, error) {
	return _Nodes.Contract.GetActiveApiNodesIDs(&_Nodes.CallOpts)
}

// GetActiveReplicationNodes is a free data retrieval call binding the contract method 0x8fbbf623.
//
// Solidity: function getActiveReplicationNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_Nodes *NodesCaller) GetActiveReplicationNodes(opts *bind.CallOpts) ([]INodesNodeWithId, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveReplicationNodes")

	if err != nil {
		return *new([]INodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodesNodeWithId)).(*[]INodesNodeWithId)

	return out0, err

}

// GetActiveReplicationNodes is a free data retrieval call binding the contract method 0x8fbbf623.
//
// Solidity: function getActiveReplicationNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_Nodes *NodesSession) GetActiveReplicationNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetActiveReplicationNodes(&_Nodes.CallOpts)
}

// GetActiveReplicationNodes is a free data retrieval call binding the contract method 0x8fbbf623.
//
// Solidity: function getActiveReplicationNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] activeNodes)
func (_Nodes *NodesCallerSession) GetActiveReplicationNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetActiveReplicationNodes(&_Nodes.CallOpts)
}

// GetActiveReplicationNodesCount is a free data retrieval call binding the contract method 0x17e3b3a9.
//
// Solidity: function getActiveReplicationNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesCaller) GetActiveReplicationNodesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveReplicationNodesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveReplicationNodesCount is a free data retrieval call binding the contract method 0x17e3b3a9.
//
// Solidity: function getActiveReplicationNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesSession) GetActiveReplicationNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetActiveReplicationNodesCount(&_Nodes.CallOpts)
}

// GetActiveReplicationNodesCount is a free data retrieval call binding the contract method 0x17e3b3a9.
//
// Solidity: function getActiveReplicationNodesCount() view returns(uint256 activeNodesCount)
func (_Nodes *NodesCallerSession) GetActiveReplicationNodesCount() (*big.Int, error) {
	return _Nodes.Contract.GetActiveReplicationNodesCount(&_Nodes.CallOpts)
}

// GetActiveReplicationNodesIDs is a free data retrieval call binding the contract method 0xfb1120e2.
//
// Solidity: function getActiveReplicationNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesCaller) GetActiveReplicationNodesIDs(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getActiveReplicationNodesIDs")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveReplicationNodesIDs is a free data retrieval call binding the contract method 0xfb1120e2.
//
// Solidity: function getActiveReplicationNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesSession) GetActiveReplicationNodesIDs() ([]*big.Int, error) {
	return _Nodes.Contract.GetActiveReplicationNodesIDs(&_Nodes.CallOpts)
}

// GetActiveReplicationNodesIDs is a free data retrieval call binding the contract method 0xfb1120e2.
//
// Solidity: function getActiveReplicationNodesIDs() view returns(uint256[] activeNodesIDs)
func (_Nodes *NodesCallerSession) GetActiveReplicationNodesIDs() ([]*big.Int, error) {
	return _Nodes.Contract.GetActiveReplicationNodesIDs(&_Nodes.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodes)
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
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodes)
func (_Nodes *NodesSession) GetAllNodes() ([]INodesNodeWithId, error) {
	return _Nodes.Contract.GetAllNodes(&_Nodes.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint256,(bytes,string,bool,bool,bool,uint256))[] allNodes)
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

// GetApiNodeIsActive is a free data retrieval call binding the contract method 0x21fbd7cb.
//
// Solidity: function getApiNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesCaller) GetApiNodeIsActive(opts *bind.CallOpts, nodeId *big.Int) (bool, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getApiNodeIsActive", nodeId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetApiNodeIsActive is a free data retrieval call binding the contract method 0x21fbd7cb.
//
// Solidity: function getApiNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesSession) GetApiNodeIsActive(nodeId *big.Int) (bool, error) {
	return _Nodes.Contract.GetApiNodeIsActive(&_Nodes.CallOpts, nodeId)
}

// GetApiNodeIsActive is a free data retrieval call binding the contract method 0x21fbd7cb.
//
// Solidity: function getApiNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesCallerSession) GetApiNodeIsActive(nodeId *big.Int) (bool, error) {
	return _Nodes.Contract.GetApiNodeIsActive(&_Nodes.CallOpts, nodeId)
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

// GetNodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xb9b140d6.
//
// Solidity: function getNodeOperatorCommissionPercent() view returns(uint256 commissionPercent)
func (_Nodes *NodesCaller) GetNodeOperatorCommissionPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getNodeOperatorCommissionPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xb9b140d6.
//
// Solidity: function getNodeOperatorCommissionPercent() view returns(uint256 commissionPercent)
func (_Nodes *NodesSession) GetNodeOperatorCommissionPercent() (*big.Int, error) {
	return _Nodes.Contract.GetNodeOperatorCommissionPercent(&_Nodes.CallOpts)
}

// GetNodeOperatorCommissionPercent is a free data retrieval call binding the contract method 0xb9b140d6.
//
// Solidity: function getNodeOperatorCommissionPercent() view returns(uint256 commissionPercent)
func (_Nodes *NodesCallerSession) GetNodeOperatorCommissionPercent() (*big.Int, error) {
	return _Nodes.Contract.GetNodeOperatorCommissionPercent(&_Nodes.CallOpts)
}

// GetReplicationNodeIsActive is a free data retrieval call binding the contract method 0x44ff624e.
//
// Solidity: function getReplicationNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesCaller) GetReplicationNodeIsActive(opts *bind.CallOpts, nodeId *big.Int) (bool, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getReplicationNodeIsActive", nodeId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetReplicationNodeIsActive is a free data retrieval call binding the contract method 0x44ff624e.
//
// Solidity: function getReplicationNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesSession) GetReplicationNodeIsActive(nodeId *big.Int) (bool, error) {
	return _Nodes.Contract.GetReplicationNodeIsActive(&_Nodes.CallOpts, nodeId)
}

// GetReplicationNodeIsActive is a free data retrieval call binding the contract method 0x44ff624e.
//
// Solidity: function getReplicationNodeIsActive(uint256 nodeId) view returns(bool isActive)
func (_Nodes *NodesCallerSession) GetReplicationNodeIsActive(nodeId *big.Int) (bool, error) {
	return _Nodes.Contract.GetReplicationNodeIsActive(&_Nodes.CallOpts, nodeId)
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
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool supported)
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
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool supported)
func (_Nodes *NodesSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nodes.Contract.SupportsInterface(&_Nodes.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool supported)
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
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars) returns(uint256 nodeId)
func (_Nodes *NodesTransactor) AddNode(opts *bind.TransactOpts, to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "addNode", to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars) returns(uint256 nodeId)
func (_Nodes *NodesSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars)
}

// AddNode is a paid mutator transaction binding the contract method 0x6ec97bfc.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars) returns(uint256 nodeId)
func (_Nodes *NodesTransactorSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars)
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

// DisableNode is a paid mutator transaction binding the contract method 0xa835f88e.
//
// Solidity: function disableNode(uint256 nodeId) returns()
func (_Nodes *NodesTransactor) DisableNode(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "disableNode", nodeId)
}

// DisableNode is a paid mutator transaction binding the contract method 0xa835f88e.
//
// Solidity: function disableNode(uint256 nodeId) returns()
func (_Nodes *NodesSession) DisableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.DisableNode(&_Nodes.TransactOpts, nodeId)
}

// DisableNode is a paid mutator transaction binding the contract method 0xa835f88e.
//
// Solidity: function disableNode(uint256 nodeId) returns()
func (_Nodes *NodesTransactorSession) DisableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.DisableNode(&_Nodes.TransactOpts, nodeId)
}

// EnableNode is a paid mutator transaction binding the contract method 0x3d2853fb.
//
// Solidity: function enableNode(uint256 nodeId) returns()
func (_Nodes *NodesTransactor) EnableNode(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "enableNode", nodeId)
}

// EnableNode is a paid mutator transaction binding the contract method 0x3d2853fb.
//
// Solidity: function enableNode(uint256 nodeId) returns()
func (_Nodes *NodesSession) EnableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.EnableNode(&_Nodes.TransactOpts, nodeId)
}

// EnableNode is a paid mutator transaction binding the contract method 0x3d2853fb.
//
// Solidity: function enableNode(uint256 nodeId) returns()
func (_Nodes *NodesTransactorSession) EnableNode(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.EnableNode(&_Nodes.TransactOpts, nodeId)
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

// RemoveFromApiNodes is a paid mutator transaction binding the contract method 0x79e0d58c.
//
// Solidity: function removeFromApiNodes(uint256 nodeId) returns()
func (_Nodes *NodesTransactor) RemoveFromApiNodes(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "removeFromApiNodes", nodeId)
}

// RemoveFromApiNodes is a paid mutator transaction binding the contract method 0x79e0d58c.
//
// Solidity: function removeFromApiNodes(uint256 nodeId) returns()
func (_Nodes *NodesSession) RemoveFromApiNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.RemoveFromApiNodes(&_Nodes.TransactOpts, nodeId)
}

// RemoveFromApiNodes is a paid mutator transaction binding the contract method 0x79e0d58c.
//
// Solidity: function removeFromApiNodes(uint256 nodeId) returns()
func (_Nodes *NodesTransactorSession) RemoveFromApiNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.RemoveFromApiNodes(&_Nodes.TransactOpts, nodeId)
}

// RemoveFromReplicationNodes is a paid mutator transaction binding the contract method 0x203ede77.
//
// Solidity: function removeFromReplicationNodes(uint256 nodeId) returns()
func (_Nodes *NodesTransactor) RemoveFromReplicationNodes(opts *bind.TransactOpts, nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "removeFromReplicationNodes", nodeId)
}

// RemoveFromReplicationNodes is a paid mutator transaction binding the contract method 0x203ede77.
//
// Solidity: function removeFromReplicationNodes(uint256 nodeId) returns()
func (_Nodes *NodesSession) RemoveFromReplicationNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.RemoveFromReplicationNodes(&_Nodes.TransactOpts, nodeId)
}

// RemoveFromReplicationNodes is a paid mutator transaction binding the contract method 0x203ede77.
//
// Solidity: function removeFromReplicationNodes(uint256 nodeId) returns()
func (_Nodes *NodesTransactorSession) RemoveFromReplicationNodes(nodeId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.RemoveFromReplicationNodes(&_Nodes.TransactOpts, nodeId)
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

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_Nodes *NodesTransactor) SetHttpAddress(opts *bind.TransactOpts, nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setHttpAddress", nodeId, httpAddress)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_Nodes *NodesSession) SetHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.SetHttpAddress(&_Nodes.TransactOpts, nodeId, httpAddress)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xd74a2a50.
//
// Solidity: function setHttpAddress(uint256 nodeId, string httpAddress) returns()
func (_Nodes *NodesTransactorSession) SetHttpAddress(nodeId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.SetHttpAddress(&_Nodes.TransactOpts, nodeId, httpAddress)
}

// SetIsApiEnabled is a paid mutator transaction binding the contract method 0x895620b7.
//
// Solidity: function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) returns()
func (_Nodes *NodesTransactor) SetIsApiEnabled(opts *bind.TransactOpts, nodeId *big.Int, isApiEnabled bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setIsApiEnabled", nodeId, isApiEnabled)
}

// SetIsApiEnabled is a paid mutator transaction binding the contract method 0x895620b7.
//
// Solidity: function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) returns()
func (_Nodes *NodesSession) SetIsApiEnabled(nodeId *big.Int, isApiEnabled bool) (*types.Transaction, error) {
	return _Nodes.Contract.SetIsApiEnabled(&_Nodes.TransactOpts, nodeId, isApiEnabled)
}

// SetIsApiEnabled is a paid mutator transaction binding the contract method 0x895620b7.
//
// Solidity: function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) returns()
func (_Nodes *NodesTransactorSession) SetIsApiEnabled(nodeId *big.Int, isApiEnabled bool) (*types.Transaction, error) {
	return _Nodes.Contract.SetIsApiEnabled(&_Nodes.TransactOpts, nodeId, isApiEnabled)
}

// SetIsReplicationEnabled is a paid mutator transaction binding the contract method 0xe18cb254.
//
// Solidity: function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_Nodes *NodesTransactor) SetIsReplicationEnabled(opts *bind.TransactOpts, nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setIsReplicationEnabled", nodeId, isReplicationEnabled)
}

// SetIsReplicationEnabled is a paid mutator transaction binding the contract method 0xe18cb254.
//
// Solidity: function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_Nodes *NodesSession) SetIsReplicationEnabled(nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _Nodes.Contract.SetIsReplicationEnabled(&_Nodes.TransactOpts, nodeId, isReplicationEnabled)
}

// SetIsReplicationEnabled is a paid mutator transaction binding the contract method 0xe18cb254.
//
// Solidity: function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) returns()
func (_Nodes *NodesTransactorSession) SetIsReplicationEnabled(nodeId *big.Int, isReplicationEnabled bool) (*types.Transaction, error) {
	return _Nodes.Contract.SetIsReplicationEnabled(&_Nodes.TransactOpts, nodeId, isReplicationEnabled)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_Nodes *NodesTransactor) SetMaxActiveNodes(opts *bind.TransactOpts, newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setMaxActiveNodes", newMaxActiveNodes)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_Nodes *NodesSession) SetMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _Nodes.Contract.SetMaxActiveNodes(&_Nodes.TransactOpts, newMaxActiveNodes)
}

// SetMaxActiveNodes is a paid mutator transaction binding the contract method 0x8ed9ea34.
//
// Solidity: function setMaxActiveNodes(uint8 newMaxActiveNodes) returns()
func (_Nodes *NodesTransactorSession) SetMaxActiveNodes(newMaxActiveNodes uint8) (*types.Transaction, error) {
	return _Nodes.Contract.SetMaxActiveNodes(&_Nodes.TransactOpts, newMaxActiveNodes)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) returns()
func (_Nodes *NodesTransactor) SetMinMonthlyFee(opts *bind.TransactOpts, nodeId *big.Int, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setMinMonthlyFee", nodeId, minMonthlyFeeMicroDollars)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) returns()
func (_Nodes *NodesSession) SetMinMonthlyFee(nodeId *big.Int, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.SetMinMonthlyFee(&_Nodes.TransactOpts, nodeId, minMonthlyFeeMicroDollars)
}

// SetMinMonthlyFee is a paid mutator transaction binding the contract method 0xce999489.
//
// Solidity: function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) returns()
func (_Nodes *NodesTransactorSession) SetMinMonthlyFee(nodeId *big.Int, minMonthlyFeeMicroDollars *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.SetMinMonthlyFee(&_Nodes.TransactOpts, nodeId, minMonthlyFeeMicroDollars)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_Nodes *NodesTransactor) SetNodeOperatorCommissionPercent(opts *bind.TransactOpts, newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "setNodeOperatorCommissionPercent", newCommissionPercent)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_Nodes *NodesSession) SetNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.SetNodeOperatorCommissionPercent(&_Nodes.TransactOpts, newCommissionPercent)
}

// SetNodeOperatorCommissionPercent is a paid mutator transaction binding the contract method 0xc4741f31.
//
// Solidity: function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) returns()
func (_Nodes *NodesTransactorSession) SetNodeOperatorCommissionPercent(newCommissionPercent *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.SetNodeOperatorCommissionPercent(&_Nodes.TransactOpts, newCommissionPercent)
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

// NodesApiDisabledIterator is returned from FilterApiDisabled and is used to iterate over the raw logs and unpacked data for ApiDisabled events raised by the Nodes contract.
type NodesApiDisabledIterator struct {
	Event *NodesApiDisabled // Event containing the contract specifics and raw log

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
func (it *NodesApiDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesApiDisabled)
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
		it.Event = new(NodesApiDisabled)
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
func (it *NodesApiDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesApiDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesApiDisabled represents a ApiDisabled event raised by the Nodes contract.
type NodesApiDisabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterApiDisabled is a free log retrieval operation binding the contract event 0x9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f3.
//
// Solidity: event ApiDisabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) FilterApiDisabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesApiDisabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ApiDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesApiDisabledIterator{contract: _Nodes.contract, event: "ApiDisabled", logs: logs, sub: sub}, nil
}

// WatchApiDisabled is a free log subscription operation binding the contract event 0x9e61606f143f77c9ef06d68819ce159e3b98756d8eb558b91db82c8ca42357f3.
//
// Solidity: event ApiDisabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) WatchApiDisabled(opts *bind.WatchOpts, sink chan<- *NodesApiDisabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ApiDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesApiDisabled)
				if err := _Nodes.contract.UnpackLog(event, "ApiDisabled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseApiDisabled(log types.Log) (*NodesApiDisabled, error) {
	event := new(NodesApiDisabled)
	if err := _Nodes.contract.UnpackLog(event, "ApiDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesApiEnabledIterator is returned from FilterApiEnabled and is used to iterate over the raw logs and unpacked data for ApiEnabled events raised by the Nodes contract.
type NodesApiEnabledIterator struct {
	Event *NodesApiEnabled // Event containing the contract specifics and raw log

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
func (it *NodesApiEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesApiEnabled)
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
		it.Event = new(NodesApiEnabled)
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
func (it *NodesApiEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesApiEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesApiEnabled represents a ApiEnabled event raised by the Nodes contract.
type NodesApiEnabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterApiEnabled is a free log retrieval operation binding the contract event 0x3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da173.
//
// Solidity: event ApiEnabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) FilterApiEnabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesApiEnabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ApiEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesApiEnabledIterator{contract: _Nodes.contract, event: "ApiEnabled", logs: logs, sub: sub}, nil
}

// WatchApiEnabled is a free log subscription operation binding the contract event 0x3cd0abf09f2e68bb82445a4d250b6a3e30b2b2c711b322e3f3817927a07da173.
//
// Solidity: event ApiEnabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) WatchApiEnabled(opts *bind.WatchOpts, sink chan<- *NodesApiEnabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ApiEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesApiEnabled)
				if err := _Nodes.contract.UnpackLog(event, "ApiEnabled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseApiEnabled(log types.Log) (*NodesApiEnabled, error) {
	event := new(NodesApiEnabled)
	if err := _Nodes.contract.UnpackLog(event, "ApiEnabled", log); err != nil {
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
	NodeId                    *big.Int
	MinMonthlyFeeMicroDollars *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterMinMonthlyFeeUpdated is a free log retrieval operation binding the contract event 0x27a815a14bf8281048d2768dcd6b695fbd4e98af4e3fb52d92c8c65384320d4a.
//
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFeeMicroDollars)
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
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFeeMicroDollars)
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
// Solidity: event MinMonthlyFeeUpdated(uint256 indexed nodeId, uint256 minMonthlyFeeMicroDollars)
func (_Nodes *NodesFilterer) ParseMinMonthlyFeeUpdated(log types.Log) (*NodesMinMonthlyFeeUpdated, error) {
	event := new(NodesMinMonthlyFeeUpdated)
	if err := _Nodes.contract.UnpackLog(event, "MinMonthlyFeeUpdated", log); err != nil {
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
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars)
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
// Solidity: event NodeAdded(uint256 indexed nodeId, address indexed owner, bytes signingKeyPub, string httpAddress, uint256 minMonthlyFeeMicroDollars)
func (_Nodes *NodesFilterer) ParseNodeAdded(log types.Log) (*NodesNodeAdded, error) {
	event := new(NodesNodeAdded)
	if err := _Nodes.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesNodeDisabledIterator is returned from FilterNodeDisabled and is used to iterate over the raw logs and unpacked data for NodeDisabled events raised by the Nodes contract.
type NodesNodeDisabledIterator struct {
	Event *NodesNodeDisabled // Event containing the contract specifics and raw log

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
func (it *NodesNodeDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeDisabled)
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
		it.Event = new(NodesNodeDisabled)
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
func (it *NodesNodeDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeDisabled represents a NodeDisabled event raised by the Nodes contract.
type NodesNodeDisabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeDisabled is a free log retrieval operation binding the contract event 0xa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8.
//
// Solidity: event NodeDisabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) FilterNodeDisabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesNodeDisabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesNodeDisabledIterator{contract: _Nodes.contract, event: "NodeDisabled", logs: logs, sub: sub}, nil
}

// WatchNodeDisabled is a free log subscription operation binding the contract event 0xa6c942fbe3ded4df132dc2c4adbb95359afebc3c361393a3d7217e3c310923e8.
//
// Solidity: event NodeDisabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) WatchNodeDisabled(opts *bind.WatchOpts, sink chan<- *NodesNodeDisabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeDisabled)
				if err := _Nodes.contract.UnpackLog(event, "NodeDisabled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseNodeDisabled(log types.Log) (*NodesNodeDisabled, error) {
	event := new(NodesNodeDisabled)
	if err := _Nodes.contract.UnpackLog(event, "NodeDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesNodeEnabledIterator is returned from FilterNodeEnabled and is used to iterate over the raw logs and unpacked data for NodeEnabled events raised by the Nodes contract.
type NodesNodeEnabledIterator struct {
	Event *NodesNodeEnabled // Event containing the contract specifics and raw log

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
func (it *NodesNodeEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeEnabled)
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
		it.Event = new(NodesNodeEnabled)
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
func (it *NodesNodeEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeEnabled represents a NodeEnabled event raised by the Nodes contract.
type NodesNodeEnabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeEnabled is a free log retrieval operation binding the contract event 0xf044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a.
//
// Solidity: event NodeEnabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) FilterNodeEnabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesNodeEnabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesNodeEnabledIterator{contract: _Nodes.contract, event: "NodeEnabled", logs: logs, sub: sub}, nil
}

// WatchNodeEnabled is a free log subscription operation binding the contract event 0xf044a2d72ef98b7636ca9d9f8c0fc60e24309bbb8d472fdecbbaca55fe166d0a.
//
// Solidity: event NodeEnabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) WatchNodeEnabled(opts *bind.WatchOpts, sink chan<- *NodesNodeEnabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeEnabled)
				if err := _Nodes.contract.UnpackLog(event, "NodeEnabled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseNodeEnabled(log types.Log) (*NodesNodeEnabled, error) {
	event := new(NodesNodeEnabled)
	if err := _Nodes.contract.UnpackLog(event, "NodeEnabled", log); err != nil {
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

// NodesReplicationDisabledIterator is returned from FilterReplicationDisabled and is used to iterate over the raw logs and unpacked data for ReplicationDisabled events raised by the Nodes contract.
type NodesReplicationDisabledIterator struct {
	Event *NodesReplicationDisabled // Event containing the contract specifics and raw log

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
func (it *NodesReplicationDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesReplicationDisabled)
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
		it.Event = new(NodesReplicationDisabled)
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
func (it *NodesReplicationDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesReplicationDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesReplicationDisabled represents a ReplicationDisabled event raised by the Nodes contract.
type NodesReplicationDisabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReplicationDisabled is a free log retrieval operation binding the contract event 0xa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce280761808.
//
// Solidity: event ReplicationDisabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) FilterReplicationDisabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesReplicationDisabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ReplicationDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesReplicationDisabledIterator{contract: _Nodes.contract, event: "ReplicationDisabled", logs: logs, sub: sub}, nil
}

// WatchReplicationDisabled is a free log subscription operation binding the contract event 0xa9837328431beea294d22d476aeafca23f85f320de41750a9b9c3ce280761808.
//
// Solidity: event ReplicationDisabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) WatchReplicationDisabled(opts *bind.WatchOpts, sink chan<- *NodesReplicationDisabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ReplicationDisabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesReplicationDisabled)
				if err := _Nodes.contract.UnpackLog(event, "ReplicationDisabled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseReplicationDisabled(log types.Log) (*NodesReplicationDisabled, error) {
	event := new(NodesReplicationDisabled)
	if err := _Nodes.contract.UnpackLog(event, "ReplicationDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesReplicationEnabledIterator is returned from FilterReplicationEnabled and is used to iterate over the raw logs and unpacked data for ReplicationEnabled events raised by the Nodes contract.
type NodesReplicationEnabledIterator struct {
	Event *NodesReplicationEnabled // Event containing the contract specifics and raw log

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
func (it *NodesReplicationEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesReplicationEnabled)
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
		it.Event = new(NodesReplicationEnabled)
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
func (it *NodesReplicationEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesReplicationEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesReplicationEnabled represents a ReplicationEnabled event raised by the Nodes contract.
type NodesReplicationEnabled struct {
	NodeId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReplicationEnabled is a free log retrieval operation binding the contract event 0xd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d0.
//
// Solidity: event ReplicationEnabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) FilterReplicationEnabled(opts *bind.FilterOpts, nodeId []*big.Int) (*NodesReplicationEnabledIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "ReplicationEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodesReplicationEnabledIterator{contract: _Nodes.contract, event: "ReplicationEnabled", logs: logs, sub: sub}, nil
}

// WatchReplicationEnabled is a free log subscription operation binding the contract event 0xd9199c75487673396ebe8093e82e5cf7902ccfb90befe22763a8bc4c36b976d0.
//
// Solidity: event ReplicationEnabled(uint256 indexed nodeId)
func (_Nodes *NodesFilterer) WatchReplicationEnabled(opts *bind.WatchOpts, sink chan<- *NodesReplicationEnabled, nodeId []*big.Int) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "ReplicationEnabled", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesReplicationEnabled)
				if err := _Nodes.contract.UnpackLog(event, "ReplicationEnabled", log); err != nil {
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
func (_Nodes *NodesFilterer) ParseReplicationEnabled(log types.Log) (*NodesReplicationEnabled, error) {
	event := new(NodesReplicationEnabled)
	if err := _Nodes.contract.UnpackLog(event, "ReplicationEnabled", log); err != nil {
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
