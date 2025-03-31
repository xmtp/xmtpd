// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package groupmessagebroadcaster

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

// GroupMessageBroadcasterMetaData contains all meta data concerning the GroupMessageBroadcaster contract.
var GroupMessageBroadcasterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"ABSOLUTE_MAX_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ABSOLUTE_MIN_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addMessage\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxPayloadSize\",\"inputs\":[{\"name\":\"maxPayloadSizeRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinPayloadSize\",\"inputs\":[{\"name\":\"minPayloadSizeRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxPayloadSizeUpdated\",\"inputs\":[{\"name\":\"oldSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinPayloadSizeUpdated\",\"inputs\":[{\"name\":\"oldSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeAuthorized\",\"inputs\":[{\"name\":\"upgrader\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newImplementation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMaxPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMinPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAdminAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroImplementationAddress\",\"inputs\":[]}]",
	Bin: "0x60a0604052306080523480156012575f5ffd5b50608051611a4b6100395f395f8181610e8701528181610eb0015261116d0152611a4b5ff3fe608060405260043610610161575f3560e01c806352d1902d116100c6578063a217fddf1161007c578063d547741f11610057578063d547741f1461044b578063f96927ac1461046a578063fe8e37a31461049d575f5ffd5b8063a217fddf146103c4578063ad3cb1cc146103d7578063c4d66de81461042c575f5ffd5b80635c975abb116100ac5780635c975abb1461030a5780638456cb591461034057806391d1485414610354575f5ffd5b806352d1902d146102c357806358e3e94c146102d7575f5ffd5b8063314a100e1161011b5780633f4ba83a116101015780633f4ba83a1461027d5780634dff26b5146102915780634f1ef286146102b0575f5ffd5b8063314a100e1461023f57806336568abe1461025e575f5ffd5b80631de015991161014b5780631de01599146101bd578063248a9ca3146101d15780632f2ff15d1461021e575f5ffd5b806209e1271461016557806301ffc9a71461018e575b5f5ffd5b348015610170575f5ffd5b5061017b6240000081565b6040519081526020015b60405180910390f35b348015610199575f5ffd5b506101ad6101a83660046116c9565b6104bc565b6040519015158152602001610185565b3480156101c8575f5ffd5b5061017b604e81565b3480156101dc575f5ffd5b5061017b6101eb366004611708565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b348015610229575f5ffd5b5061023d610238366004611747565b610554565b005b34801561024a575f5ffd5b5061023d610259366004611708565b61059d565b348015610269575f5ffd5b5061023d610278366004611747565b610674565b348015610288575f5ffd5b5061023d6106d2565b34801561029c575f5ffd5b5061023d6102ab366004611771565b6106e7565b61023d6102be366004611815565b6107e9565b3480156102ce575f5ffd5b5061017b610808565b3480156102e2575f5ffd5b507f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a610015461017b565b348015610315575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff166101ad565b34801561034b575f5ffd5b5061023d610836565b34801561035f575f5ffd5b506101ad61036e366004611747565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b3480156103cf575f5ffd5b5061017b5f81565b3480156103e2575f5ffd5b5061041f6040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516101859190611916565b348015610437575f5ffd5b5061023d610446366004611969565b610848565b348015610456575f5ffd5b5061023d610465366004611747565b610a7a565b348015610475575f5ffd5b507f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a610005461017b565b3480156104a8575f5ffd5b5061023d6104b7366004611708565b610abd565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061054e57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015461058d81610b70565b6105978383610b7a565b50505050565b5f6105a781610b70565b7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61001547f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61000908311806105f85750604e83105b1561062f576040517fe219e4f000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805483825560408051828152602081018690527f1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd91015b60405180910390a150505050565b73ffffffffffffffffffffffffffffffffffffffff811633146106c3576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106cd8282610c98565b505050565b5f6106dc81610b70565b6106e4610d74565b50565b6106ef610e11565b7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a6100080548210806107225750600181015482115b1561077557805460018201546040517f93b7abe600000000000000000000000000000000000000000000000000000000815260048101859052602481019290925260448201526064015b60405180910390fd5b6002810180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008116600167ffffffffffffffff92831601918216179091556040517f91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e916106669187918791879190611982565b6107f1610e6f565b6107fa82610f73565b610804828261101c565b5050565b5f610811611155565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b5f61084081610b70565b6106e46111c4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156108925750825b90505f8267ffffffffffffffff1660011480156108ae5750303b155b9050811580156108bc575080155b156108f3576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156109545784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b73ffffffffffffffffffffffffffffffffffffffff86166109a1576040517f3ef39b8100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109a961123d565b6109b161123d565b6109b9611245565b604e7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61000908155624000007f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a6100155610a0f5f88610b7a565b50508315610a725784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610ab381610b70565b6105978383610c98565b5f610ac781610b70565b7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a610008054831080610af957506240000083115b15610b30576040517f1d8e7a4a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001810180549084905560408051828152602081018690527ff59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe909101610666565b6106e48133611255565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff16610c8f575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610c2b3390565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4600191505061054e565b5f91505061054e565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff1615610c8f575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4600191505061054e565b610d7c6112fb565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a150565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610e6d576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b3073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161480610f3c57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610f237f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b15610e6d576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610f7d81610b70565b73ffffffffffffffffffffffffffffffffffffffff8216610fca576040517fd02c623d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201527fd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa910160405180910390a15050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156110a1575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261109e918101906119e8565b60015b6110ef576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260240161076c565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc811461114b576040517faa1d49a40000000000000000000000000000000000000000000000000000000081526004810182905260240161076c565b6106cd8383611356565b3073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610e6d576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6111cc610e11565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833610de6565b610e6d6113b8565b61124d6113b8565b610e6d61141f565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610804576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024810183905260440161076c565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610e6d576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61135f82611470565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156113b0576106cd828261153e565b6108046115bd565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610e6d576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6114276113b8565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b8073ffffffffffffffffffffffffffffffffffffffff163b5f036114d8576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161076c565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff168460405161156791906119ff565b5f60405180830381855af49150503d805f811461159f576040519150601f19603f3d011682016040523d82523d5f602084013e6115a4565b606091505b50915091506115b48583836115f5565b95945050505050565b3415610e6d576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608261160a5761160582611687565b611680565b815115801561162e575073ffffffffffffffffffffffffffffffffffffffff84163b155b1561167d576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260240161076c565b50805b9392505050565b8051156116975780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f602082840312156116d9575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611680575f5ffd5b5f60208284031215611718575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611742575f5ffd5b919050565b5f5f60408385031215611758575f5ffd5b823591506117686020840161171f565b90509250929050565b5f5f5f60408486031215611783575f5ffd5b83359250602084013567ffffffffffffffff8111156117a0575f5ffd5b8401601f810186136117b0575f5ffd5b803567ffffffffffffffff8111156117c6575f5ffd5b8660208284010111156117d7575f5ffd5b939660209190910195509293505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611826575f5ffd5b61182f8361171f565b9150602083013567ffffffffffffffff81111561184a575f5ffd5b8301601f8101851361185a575f5ffd5b803567ffffffffffffffff811115611874576118746117e8565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156118e0576118e06117e8565b6040528181528282016020018710156118f7575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f60208284031215611979575f5ffd5b6116808261171f565b84815260606020820152826060820152828460808301375f608084830101525f60807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116830101905067ffffffffffffffff8316604083015295945050505050565b5f602082840312156119f8575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fea26469706673582212208e05f0c93fc8bfda298139941e4198d4f130901ae1e4000a9824ab4ec5802f2c64736f6c634300081c0033",
}

// GroupMessageBroadcasterABI is the input ABI used to generate the binding from.
// Deprecated: Use GroupMessageBroadcasterMetaData.ABI instead.
var GroupMessageBroadcasterABI = GroupMessageBroadcasterMetaData.ABI

// GroupMessageBroadcasterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GroupMessageBroadcasterMetaData.Bin instead.
var GroupMessageBroadcasterBin = GroupMessageBroadcasterMetaData.Bin

// DeployGroupMessageBroadcaster deploys a new Ethereum contract, binding an instance of GroupMessageBroadcaster to it.
func DeployGroupMessageBroadcaster(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GroupMessageBroadcaster, error) {
	parsed, err := GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GroupMessageBroadcasterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GroupMessageBroadcaster{GroupMessageBroadcasterCaller: GroupMessageBroadcasterCaller{contract: contract}, GroupMessageBroadcasterTransactor: GroupMessageBroadcasterTransactor{contract: contract}, GroupMessageBroadcasterFilterer: GroupMessageBroadcasterFilterer{contract: contract}}, nil
}

// GroupMessageBroadcaster is an auto generated Go binding around an Ethereum contract.
type GroupMessageBroadcaster struct {
	GroupMessageBroadcasterCaller     // Read-only binding to the contract
	GroupMessageBroadcasterTransactor // Write-only binding to the contract
	GroupMessageBroadcasterFilterer   // Log filterer for contract events
}

// GroupMessageBroadcasterCaller is an auto generated read-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessageBroadcasterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessageBroadcasterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GroupMessageBroadcasterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessageBroadcasterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GroupMessageBroadcasterSession struct {
	Contract     *GroupMessageBroadcaster // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GroupMessageBroadcasterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GroupMessageBroadcasterCallerSession struct {
	Contract *GroupMessageBroadcasterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// GroupMessageBroadcasterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GroupMessageBroadcasterTransactorSession struct {
	Contract     *GroupMessageBroadcasterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// GroupMessageBroadcasterRaw is an auto generated low-level Go binding around an Ethereum contract.
type GroupMessageBroadcasterRaw struct {
	Contract *GroupMessageBroadcaster // Generic contract binding to access the raw methods on
}

// GroupMessageBroadcasterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterCallerRaw struct {
	Contract *GroupMessageBroadcasterCaller // Generic read-only contract binding to access the raw methods on
}

// GroupMessageBroadcasterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GroupMessageBroadcasterTransactorRaw struct {
	Contract *GroupMessageBroadcasterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGroupMessageBroadcaster creates a new instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcaster(address common.Address, backend bind.ContractBackend) (*GroupMessageBroadcaster, error) {
	contract, err := bindGroupMessageBroadcaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcaster{GroupMessageBroadcasterCaller: GroupMessageBroadcasterCaller{contract: contract}, GroupMessageBroadcasterTransactor: GroupMessageBroadcasterTransactor{contract: contract}, GroupMessageBroadcasterFilterer: GroupMessageBroadcasterFilterer{contract: contract}}, nil
}

// NewGroupMessageBroadcasterCaller creates a new read-only instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcasterCaller(address common.Address, caller bind.ContractCaller) (*GroupMessageBroadcasterCaller, error) {
	contract, err := bindGroupMessageBroadcaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterCaller{contract: contract}, nil
}

// NewGroupMessageBroadcasterTransactor creates a new write-only instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcasterTransactor(address common.Address, transactor bind.ContractTransactor) (*GroupMessageBroadcasterTransactor, error) {
	contract, err := bindGroupMessageBroadcaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterTransactor{contract: contract}, nil
}

// NewGroupMessageBroadcasterFilterer creates a new log filterer instance of GroupMessageBroadcaster, bound to a specific deployed contract.
func NewGroupMessageBroadcasterFilterer(address common.Address, filterer bind.ContractFilterer) (*GroupMessageBroadcasterFilterer, error) {
	contract, err := bindGroupMessageBroadcaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterFilterer{contract: contract}, nil
}

// bindGroupMessageBroadcaster binds a generic wrapper to an already deployed contract.
func bindGroupMessageBroadcaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessageBroadcaster.Contract.GroupMessageBroadcasterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.GroupMessageBroadcasterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.GroupMessageBroadcasterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessageBroadcaster.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.contract.Transact(opts, method, params...)
}

// ABSOLUTEMAXPAYLOADSIZE is a free data retrieval call binding the contract method 0x0009e127.
//
// Solidity: function ABSOLUTE_MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) ABSOLUTEMAXPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "ABSOLUTE_MAX_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ABSOLUTEMAXPAYLOADSIZE is a free data retrieval call binding the contract method 0x0009e127.
//
// Solidity: function ABSOLUTE_MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) ABSOLUTEMAXPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.ABSOLUTEMAXPAYLOADSIZE(&_GroupMessageBroadcaster.CallOpts)
}

// ABSOLUTEMAXPAYLOADSIZE is a free data retrieval call binding the contract method 0x0009e127.
//
// Solidity: function ABSOLUTE_MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) ABSOLUTEMAXPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.ABSOLUTEMAXPAYLOADSIZE(&_GroupMessageBroadcaster.CallOpts)
}

// ABSOLUTEMINPAYLOADSIZE is a free data retrieval call binding the contract method 0x1de01599.
//
// Solidity: function ABSOLUTE_MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) ABSOLUTEMINPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "ABSOLUTE_MIN_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ABSOLUTEMINPAYLOADSIZE is a free data retrieval call binding the contract method 0x1de01599.
//
// Solidity: function ABSOLUTE_MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) ABSOLUTEMINPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.ABSOLUTEMINPAYLOADSIZE(&_GroupMessageBroadcaster.CallOpts)
}

// ABSOLUTEMINPAYLOADSIZE is a free data retrieval call binding the contract method 0x1de01599.
//
// Solidity: function ABSOLUTE_MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) ABSOLUTEMINPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.ABSOLUTEMINPAYLOADSIZE(&_GroupMessageBroadcaster.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GroupMessageBroadcaster.Contract.DEFAULTADMINROLE(&_GroupMessageBroadcaster.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GroupMessageBroadcaster.Contract.DEFAULTADMINROLE(&_GroupMessageBroadcaster.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _GroupMessageBroadcaster.Contract.UPGRADEINTERFACEVERSION(&_GroupMessageBroadcaster.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _GroupMessageBroadcaster.Contract.UPGRADEINTERFACEVERSION(&_GroupMessageBroadcaster.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GroupMessageBroadcaster.Contract.GetRoleAdmin(&_GroupMessageBroadcaster.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GroupMessageBroadcaster.Contract.GetRoleAdmin(&_GroupMessageBroadcaster.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GroupMessageBroadcaster.Contract.HasRole(&_GroupMessageBroadcaster.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GroupMessageBroadcaster.Contract.HasRole(&_GroupMessageBroadcaster.CallOpts, role, account)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint256 size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MaxPayloadSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "maxPayloadSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint256 size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MaxPayloadSize() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint256 size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MaxPayloadSize() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.MaxPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint256 size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) MinPayloadSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "minPayloadSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint256 size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) MinPayloadSize() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint256 size)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) MinPayloadSize() (*big.Int, error) {
	return _GroupMessageBroadcaster.Contract.MinPayloadSize(&_GroupMessageBroadcaster.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Paused() (bool, error) {
	return _GroupMessageBroadcaster.Contract.Paused(&_GroupMessageBroadcaster.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) Paused() (bool, error) {
	return _GroupMessageBroadcaster.Contract.Paused(&_GroupMessageBroadcaster.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) ProxiableUUID() ([32]byte, error) {
	return _GroupMessageBroadcaster.Contract.ProxiableUUID(&_GroupMessageBroadcaster.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) ProxiableUUID() ([32]byte, error) {
	return _GroupMessageBroadcaster.Contract.ProxiableUUID(&_GroupMessageBroadcaster.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GroupMessageBroadcaster.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GroupMessageBroadcaster.Contract.SupportsInterface(&_GroupMessageBroadcaster.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GroupMessageBroadcaster.Contract.SupportsInterface(&_GroupMessageBroadcaster.CallOpts, interfaceId)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) AddMessage(opts *bind.TransactOpts, groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "addMessage", groupId, message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) AddMessage(groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.AddMessage(&_GroupMessageBroadcaster.TransactOpts, groupId, message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) AddMessage(groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.AddMessage(&_GroupMessageBroadcaster.TransactOpts, groupId, message)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.GrantRole(&_GroupMessageBroadcaster.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.GrantRole(&_GroupMessageBroadcaster.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) Initialize(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "initialize", admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Initialize(&_GroupMessageBroadcaster.TransactOpts, admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Initialize(&_GroupMessageBroadcaster.TransactOpts, admin)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Pause() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Pause(&_GroupMessageBroadcaster.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) Pause() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Pause(&_GroupMessageBroadcaster.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.RenounceRole(&_GroupMessageBroadcaster.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.RenounceRole(&_GroupMessageBroadcaster.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.RevokeRole(&_GroupMessageBroadcaster.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.RevokeRole(&_GroupMessageBroadcaster.TransactOpts, role, account)
}

// SetMaxPayloadSize is a paid mutator transaction binding the contract method 0xfe8e37a3.
//
// Solidity: function setMaxPayloadSize(uint256 maxPayloadSizeRequest) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) SetMaxPayloadSize(opts *bind.TransactOpts, maxPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "setMaxPayloadSize", maxPayloadSizeRequest)
}

// SetMaxPayloadSize is a paid mutator transaction binding the contract method 0xfe8e37a3.
//
// Solidity: function setMaxPayloadSize(uint256 maxPayloadSizeRequest) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) SetMaxPayloadSize(maxPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.SetMaxPayloadSize(&_GroupMessageBroadcaster.TransactOpts, maxPayloadSizeRequest)
}

// SetMaxPayloadSize is a paid mutator transaction binding the contract method 0xfe8e37a3.
//
// Solidity: function setMaxPayloadSize(uint256 maxPayloadSizeRequest) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) SetMaxPayloadSize(maxPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.SetMaxPayloadSize(&_GroupMessageBroadcaster.TransactOpts, maxPayloadSizeRequest)
}

// SetMinPayloadSize is a paid mutator transaction binding the contract method 0x314a100e.
//
// Solidity: function setMinPayloadSize(uint256 minPayloadSizeRequest) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) SetMinPayloadSize(opts *bind.TransactOpts, minPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "setMinPayloadSize", minPayloadSizeRequest)
}

// SetMinPayloadSize is a paid mutator transaction binding the contract method 0x314a100e.
//
// Solidity: function setMinPayloadSize(uint256 minPayloadSizeRequest) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) SetMinPayloadSize(minPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.SetMinPayloadSize(&_GroupMessageBroadcaster.TransactOpts, minPayloadSizeRequest)
}

// SetMinPayloadSize is a paid mutator transaction binding the contract method 0x314a100e.
//
// Solidity: function setMinPayloadSize(uint256 minPayloadSizeRequest) returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) SetMinPayloadSize(minPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.SetMinPayloadSize(&_GroupMessageBroadcaster.TransactOpts, minPayloadSizeRequest)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) Unpause() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Unpause(&_GroupMessageBroadcaster.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) Unpause() (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.Unpause(&_GroupMessageBroadcaster.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpgradeToAndCall(&_GroupMessageBroadcaster.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GroupMessageBroadcaster *GroupMessageBroadcasterTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GroupMessageBroadcaster.Contract.UpgradeToAndCall(&_GroupMessageBroadcaster.TransactOpts, newImplementation, data)
}

// GroupMessageBroadcasterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterInitializedIterator struct {
	Event *GroupMessageBroadcasterInitialized // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterInitialized)
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
		it.Event = new(GroupMessageBroadcasterInitialized)
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
func (it *GroupMessageBroadcasterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterInitialized represents a Initialized event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterInitialized(opts *bind.FilterOpts) (*GroupMessageBroadcasterInitializedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterInitializedIterator{contract: _GroupMessageBroadcaster.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterInitialized) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterInitialized)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseInitialized(log types.Log) (*GroupMessageBroadcasterInitialized, error) {
	event := new(GroupMessageBroadcasterInitialized)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator is returned from FilterMaxPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MaxPayloadSizeUpdated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator struct {
	Event *GroupMessageBroadcasterMaxPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
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
		it.Event = new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
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
func (it *GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMaxPayloadSizeUpdated represents a MaxPayloadSizeUpdated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMaxPayloadSizeUpdated struct {
	OldSize *big.Int
	NewSize *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMaxPayloadSizeUpdated is a free log retrieval operation binding the contract event 0xf59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe90.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMaxPayloadSizeUpdated(opts *bind.FilterOpts) (*GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "MaxPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMaxPayloadSizeUpdatedIterator{contract: _GroupMessageBroadcaster.contract, event: "MaxPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxPayloadSizeUpdated is a free log subscription operation binding the contract event 0xf59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe90.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMaxPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMaxPayloadSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "MaxPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
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

// ParseMaxPayloadSizeUpdated is a log parse operation binding the contract event 0xf59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe90.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMaxPayloadSizeUpdated(log types.Log) (*GroupMessageBroadcasterMaxPayloadSizeUpdated, error) {
	event := new(GroupMessageBroadcasterMaxPayloadSizeUpdated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMessageSentIterator is returned from FilterMessageSent and is used to iterate over the raw logs and unpacked data for MessageSent events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMessageSentIterator struct {
	Event *GroupMessageBroadcasterMessageSent // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMessageSent)
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
		it.Event = new(GroupMessageBroadcasterMessageSent)
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
func (it *GroupMessageBroadcasterMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMessageSent represents a MessageSent event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMessageSent struct {
	GroupId    [32]byte
	Message    []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMessageSent(opts *bind.FilterOpts) (*GroupMessageBroadcasterMessageSentIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMessageSentIterator{contract: _GroupMessageBroadcaster.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

// WatchMessageSent is a free log subscription operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMessageSent) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMessageSent)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

// ParseMessageSent is a log parse operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMessageSent(log types.Log) (*GroupMessageBroadcasterMessageSent, error) {
	event := new(GroupMessageBroadcasterMessageSent)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterMinPayloadSizeUpdatedIterator is returned from FilterMinPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MinPayloadSizeUpdated events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMinPayloadSizeUpdatedIterator struct {
	Event *GroupMessageBroadcasterMinPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterMinPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterMinPayloadSizeUpdated)
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
		it.Event = new(GroupMessageBroadcasterMinPayloadSizeUpdated)
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
func (it *GroupMessageBroadcasterMinPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterMinPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterMinPayloadSizeUpdated represents a MinPayloadSizeUpdated event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterMinPayloadSizeUpdated struct {
	OldSize *big.Int
	NewSize *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinPayloadSizeUpdated is a free log retrieval operation binding the contract event 0x1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd.
//
// Solidity: event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterMinPayloadSizeUpdated(opts *bind.FilterOpts) (*GroupMessageBroadcasterMinPayloadSizeUpdatedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "MinPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterMinPayloadSizeUpdatedIterator{contract: _GroupMessageBroadcaster.contract, event: "MinPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinPayloadSizeUpdated is a free log subscription operation binding the contract event 0x1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd.
//
// Solidity: event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchMinPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterMinPayloadSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "MinPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterMinPayloadSizeUpdated)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
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

// ParseMinPayloadSizeUpdated is a log parse operation binding the contract event 0x1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd.
//
// Solidity: event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseMinPayloadSizeUpdated(log types.Log) (*GroupMessageBroadcasterMinPayloadSizeUpdated, error) {
	event := new(GroupMessageBroadcasterMinPayloadSizeUpdated)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterPausedIterator struct {
	Event *GroupMessageBroadcasterPaused // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterPaused)
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
		it.Event = new(GroupMessageBroadcasterPaused)
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
func (it *GroupMessageBroadcasterPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterPaused represents a Paused event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterPaused(opts *bind.FilterOpts) (*GroupMessageBroadcasterPausedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterPausedIterator{contract: _GroupMessageBroadcaster.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterPaused) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterPaused)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParsePaused(log types.Log) (*GroupMessageBroadcasterPaused, error) {
	event := new(GroupMessageBroadcasterPaused)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterRoleAdminChangedIterator struct {
	Event *GroupMessageBroadcasterRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterRoleAdminChanged)
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
		it.Event = new(GroupMessageBroadcasterRoleAdminChanged)
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
func (it *GroupMessageBroadcasterRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterRoleAdminChanged represents a RoleAdminChanged event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*GroupMessageBroadcasterRoleAdminChangedIterator, error) {

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

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterRoleAdminChangedIterator{contract: _GroupMessageBroadcaster.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterRoleAdminChanged)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseRoleAdminChanged(log types.Log) (*GroupMessageBroadcasterRoleAdminChanged, error) {
	event := new(GroupMessageBroadcasterRoleAdminChanged)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterRoleGrantedIterator struct {
	Event *GroupMessageBroadcasterRoleGranted // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterRoleGranted)
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
		it.Event = new(GroupMessageBroadcasterRoleGranted)
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
func (it *GroupMessageBroadcasterRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterRoleGranted represents a RoleGranted event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GroupMessageBroadcasterRoleGrantedIterator, error) {

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

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterRoleGrantedIterator{contract: _GroupMessageBroadcaster.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterRoleGranted)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseRoleGranted(log types.Log) (*GroupMessageBroadcasterRoleGranted, error) {
	event := new(GroupMessageBroadcasterRoleGranted)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterRoleRevokedIterator struct {
	Event *GroupMessageBroadcasterRoleRevoked // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterRoleRevoked)
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
		it.Event = new(GroupMessageBroadcasterRoleRevoked)
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
func (it *GroupMessageBroadcasterRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterRoleRevoked represents a RoleRevoked event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GroupMessageBroadcasterRoleRevokedIterator, error) {

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

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterRoleRevokedIterator{contract: _GroupMessageBroadcaster.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterRoleRevoked)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseRoleRevoked(log types.Log) (*GroupMessageBroadcasterRoleRevoked, error) {
	event := new(GroupMessageBroadcasterRoleRevoked)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUnpausedIterator struct {
	Event *GroupMessageBroadcasterUnpaused // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterUnpaused)
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
		it.Event = new(GroupMessageBroadcasterUnpaused)
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
func (it *GroupMessageBroadcasterUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterUnpaused represents a Unpaused event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GroupMessageBroadcasterUnpausedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterUnpausedIterator{contract: _GroupMessageBroadcaster.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterUnpaused) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterUnpaused)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseUnpaused(log types.Log) (*GroupMessageBroadcasterUnpaused, error) {
	event := new(GroupMessageBroadcasterUnpaused)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterUpgradeAuthorizedIterator is returned from FilterUpgradeAuthorized and is used to iterate over the raw logs and unpacked data for UpgradeAuthorized events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUpgradeAuthorizedIterator struct {
	Event *GroupMessageBroadcasterUpgradeAuthorized // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterUpgradeAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterUpgradeAuthorized)
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
		it.Event = new(GroupMessageBroadcasterUpgradeAuthorized)
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
func (it *GroupMessageBroadcasterUpgradeAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterUpgradeAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterUpgradeAuthorized represents a UpgradeAuthorized event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUpgradeAuthorized struct {
	Upgrader          common.Address
	NewImplementation common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeAuthorized is a free log retrieval operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterUpgradeAuthorized(opts *bind.FilterOpts) (*GroupMessageBroadcasterUpgradeAuthorizedIterator, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterUpgradeAuthorizedIterator{contract: _GroupMessageBroadcaster.contract, event: "UpgradeAuthorized", logs: logs, sub: sub}, nil
}

// WatchUpgradeAuthorized is a free log subscription operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchUpgradeAuthorized(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterUpgradeAuthorized) (event.Subscription, error) {

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterUpgradeAuthorized)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
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

// ParseUpgradeAuthorized is a log parse operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseUpgradeAuthorized(log types.Log) (*GroupMessageBroadcasterUpgradeAuthorized, error) {
	event := new(GroupMessageBroadcasterUpgradeAuthorized)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessageBroadcasterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUpgradedIterator struct {
	Event *GroupMessageBroadcasterUpgraded // Event containing the contract specifics and raw log

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
func (it *GroupMessageBroadcasterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessageBroadcasterUpgraded)
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
		it.Event = new(GroupMessageBroadcasterUpgraded)
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
func (it *GroupMessageBroadcasterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessageBroadcasterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessageBroadcasterUpgraded represents a Upgraded event raised by the GroupMessageBroadcaster contract.
type GroupMessageBroadcasterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*GroupMessageBroadcasterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessageBroadcasterUpgradedIterator{contract: _GroupMessageBroadcaster.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *GroupMessageBroadcasterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GroupMessageBroadcaster.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessageBroadcasterUpgraded)
				if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessageBroadcaster *GroupMessageBroadcasterFilterer) ParseUpgraded(log types.Log) (*GroupMessageBroadcasterUpgraded, error) {
	event := new(GroupMessageBroadcasterUpgraded)
	if err := _GroupMessageBroadcaster.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
