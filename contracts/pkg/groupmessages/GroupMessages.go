// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package groupmessages

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

// GroupMessagesMetaData contains all meta data concerning the GroupMessages contract.
var GroupMessagesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"ABSOLUTE_MAX_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ABSOLUTE_MIN_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addMessage\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minPayloadSize\",\"inputs\":[],\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxPayloadSize\",\"inputs\":[{\"name\":\"maxPayloadSizeRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinPayloadSize\",\"inputs\":[{\"name\":\"minPayloadSizeRequest\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxPayloadSizeUpdated\",\"inputs\":[{\"name\":\"oldSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinPayloadSizeUpdated\",\"inputs\":[{\"name\":\"oldSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newSize\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeAuthorized\",\"inputs\":[{\"name\":\"upgrader\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newImplementation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMaxPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMinPayloadSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAdminAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroImplementationAddress\",\"inputs\":[]}]",
	Bin: "0x60a0604052306080523480156012575f5ffd5b50608051611a4b6100395f395f8181610e8701528181610eb0015261116d0152611a4b5ff3fe608060405260043610610161575f3560e01c806352d1902d116100c6578063a217fddf1161007c578063d547741f11610057578063d547741f1461044b578063f96927ac1461046a578063fe8e37a31461049d575f5ffd5b8063a217fddf146103c4578063ad3cb1cc146103d7578063c4d66de81461042c575f5ffd5b80635c975abb116100ac5780635c975abb1461030a5780638456cb591461034057806391d1485414610354575f5ffd5b806352d1902d146102c357806358e3e94c146102d7575f5ffd5b8063314a100e1161011b5780633f4ba83a116101015780633f4ba83a1461027d5780634dff26b5146102915780634f1ef286146102b0575f5ffd5b8063314a100e1461023f57806336568abe1461025e575f5ffd5b80631de015991161014b5780631de01599146101bd578063248a9ca3146101d15780632f2ff15d1461021e575f5ffd5b806209e1271461016557806301ffc9a71461018e575b5f5ffd5b348015610170575f5ffd5b5061017b6240000081565b6040519081526020015b60405180910390f35b348015610199575f5ffd5b506101ad6101a83660046116c9565b6104bc565b6040519015158152602001610185565b3480156101c8575f5ffd5b5061017b604e81565b3480156101dc575f5ffd5b5061017b6101eb366004611708565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b348015610229575f5ffd5b5061023d610238366004611747565b610554565b005b34801561024a575f5ffd5b5061023d610259366004611708565b61059d565b348015610269575f5ffd5b5061023d610278366004611747565b610674565b348015610288575f5ffd5b5061023d6106d2565b34801561029c575f5ffd5b5061023d6102ab366004611771565b6106e7565b61023d6102be366004611815565b6107e9565b3480156102ce575f5ffd5b5061017b610808565b3480156102e2575f5ffd5b507f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a610015461017b565b348015610315575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff166101ad565b34801561034b575f5ffd5b5061023d610836565b34801561035f575f5ffd5b506101ad61036e366004611747565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b3480156103cf575f5ffd5b5061017b5f81565b3480156103e2575f5ffd5b5061041f6040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516101859190611916565b348015610437575f5ffd5b5061023d610446366004611969565b610848565b348015610456575f5ffd5b5061023d610465366004611747565b610a7a565b348015610475575f5ffd5b507f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a610005461017b565b3480156104a8575f5ffd5b5061023d6104b7366004611708565b610abd565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061054e57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015461058d81610b70565b6105978383610b7a565b50505050565b5f6105a781610b70565b7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61001547f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61000908311806105f85750604e83105b1561062f576040517fe219e4f000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805483825560408051828152602081018690527f1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd91015b60405180910390a150505050565b73ffffffffffffffffffffffffffffffffffffffff811633146106c3576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106cd8282610c98565b505050565b5f6106dc81610b70565b6106e4610d74565b50565b6106ef610e11565b7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a6100080548210806107225750600181015482115b1561077557805460018201546040517f93b7abe600000000000000000000000000000000000000000000000000000000815260048101859052602481019290925260448201526064015b60405180910390fd5b6002810180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008116600167ffffffffffffffff92831601918216179091556040517f91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e916106669187918791879190611982565b6107f1610e6f565b6107fa82610f73565b610804828261101c565b5050565b5f610811611155565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b5f61084081610b70565b6106e46111c4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156108925750825b90505f8267ffffffffffffffff1660011480156108ae5750303b155b9050811580156108bc575080155b156108f3576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156109545784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b73ffffffffffffffffffffffffffffffffffffffff86166109a1576040517f3ef39b8100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109a961123d565b6109b161123d565b6109b9611245565b604e7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61000908155624000007f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a6100155610a0f5f88610b7a565b50508315610a725784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610ab381610b70565b6105978383610c98565b5f610ac781610b70565b7f5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a610008054831080610af957506240000083115b15610b30576040517f1d8e7a4a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001810180549084905560408051828152602081018690527ff59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe909101610666565b6106e48133611255565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff16610c8f575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610c2b3390565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4600191505061054e565b5f91505061054e565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff1615610c8f575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4600191505061054e565b610d7c6112fb565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a150565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610e6d576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b3073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161480610f3c57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610f237f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b15610e6d576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610f7d81610b70565b73ffffffffffffffffffffffffffffffffffffffff8216610fca576040517fd02c623d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201527fd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa910160405180910390a15050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156110a1575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261109e918101906119e8565b60015b6110ef576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260240161076c565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc811461114b576040517faa1d49a40000000000000000000000000000000000000000000000000000000081526004810182905260240161076c565b6106cd8383611356565b3073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610e6d576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6111cc610e11565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833610de6565b610e6d6113b8565b61124d6113b8565b610e6d61141f565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610804576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024810183905260440161076c565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610e6d576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61135f82611470565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156113b0576106cd828261153e565b6108046115bd565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610e6d576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6114276113b8565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b8073ffffffffffffffffffffffffffffffffffffffff163b5f036114d8576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161076c565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff168460405161156791906119ff565b5f60405180830381855af49150503d805f811461159f576040519150601f19603f3d011682016040523d82523d5f602084013e6115a4565b606091505b50915091506115b48583836115f5565b95945050505050565b3415610e6d576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608261160a5761160582611687565b611680565b815115801561162e575073ffffffffffffffffffffffffffffffffffffffff84163b155b1561167d576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260240161076c565b50805b9392505050565b8051156116975780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f602082840312156116d9575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611680575f5ffd5b5f60208284031215611718575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611742575f5ffd5b919050565b5f5f60408385031215611758575f5ffd5b823591506117686020840161171f565b90509250929050565b5f5f5f60408486031215611783575f5ffd5b83359250602084013567ffffffffffffffff8111156117a0575f5ffd5b8401601f810186136117b0575f5ffd5b803567ffffffffffffffff8111156117c6575f5ffd5b8660208284010111156117d7575f5ffd5b939660209190910195509293505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611826575f5ffd5b61182f8361171f565b9150602083013567ffffffffffffffff81111561184a575f5ffd5b8301601f8101851361185a575f5ffd5b803567ffffffffffffffff811115611874576118746117e8565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156118e0576118e06117e8565b6040528181528282016020018710156118f7575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f60208284031215611979575f5ffd5b6116808261171f565b84815260606020820152826060820152828460808301375f608084830101525f60807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116830101905067ffffffffffffffff8316604083015295945050505050565b5f602082840312156119f8575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fea2646970667358221220a9d8450861b726977a199e2e96685ffa25c46d0cc14014015e81c1518113c59464736f6c634300081c0033",
}

// GroupMessagesABI is the input ABI used to generate the binding from.
// Deprecated: Use GroupMessagesMetaData.ABI instead.
var GroupMessagesABI = GroupMessagesMetaData.ABI

// GroupMessagesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GroupMessagesMetaData.Bin instead.
var GroupMessagesBin = GroupMessagesMetaData.Bin

// DeployGroupMessages deploys a new Ethereum contract, binding an instance of GroupMessages to it.
func DeployGroupMessages(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GroupMessages, error) {
	parsed, err := GroupMessagesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GroupMessagesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GroupMessages{GroupMessagesCaller: GroupMessagesCaller{contract: contract}, GroupMessagesTransactor: GroupMessagesTransactor{contract: contract}, GroupMessagesFilterer: GroupMessagesFilterer{contract: contract}}, nil
}

// GroupMessages is an auto generated Go binding around an Ethereum contract.
type GroupMessages struct {
	GroupMessagesCaller     // Read-only binding to the contract
	GroupMessagesTransactor // Write-only binding to the contract
	GroupMessagesFilterer   // Log filterer for contract events
}

// GroupMessagesCaller is an auto generated read-only Go binding around an Ethereum contract.
type GroupMessagesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessagesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GroupMessagesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessagesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GroupMessagesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GroupMessagesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GroupMessagesSession struct {
	Contract     *GroupMessages    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GroupMessagesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GroupMessagesCallerSession struct {
	Contract *GroupMessagesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// GroupMessagesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GroupMessagesTransactorSession struct {
	Contract     *GroupMessagesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GroupMessagesRaw is an auto generated low-level Go binding around an Ethereum contract.
type GroupMessagesRaw struct {
	Contract *GroupMessages // Generic contract binding to access the raw methods on
}

// GroupMessagesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GroupMessagesCallerRaw struct {
	Contract *GroupMessagesCaller // Generic read-only contract binding to access the raw methods on
}

// GroupMessagesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GroupMessagesTransactorRaw struct {
	Contract *GroupMessagesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGroupMessages creates a new instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessages(address common.Address, backend bind.ContractBackend) (*GroupMessages, error) {
	contract, err := bindGroupMessages(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GroupMessages{GroupMessagesCaller: GroupMessagesCaller{contract: contract}, GroupMessagesTransactor: GroupMessagesTransactor{contract: contract}, GroupMessagesFilterer: GroupMessagesFilterer{contract: contract}}, nil
}

// NewGroupMessagesCaller creates a new read-only instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessagesCaller(address common.Address, caller bind.ContractCaller) (*GroupMessagesCaller, error) {
	contract, err := bindGroupMessages(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesCaller{contract: contract}, nil
}

// NewGroupMessagesTransactor creates a new write-only instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessagesTransactor(address common.Address, transactor bind.ContractTransactor) (*GroupMessagesTransactor, error) {
	contract, err := bindGroupMessages(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesTransactor{contract: contract}, nil
}

// NewGroupMessagesFilterer creates a new log filterer instance of GroupMessages, bound to a specific deployed contract.
func NewGroupMessagesFilterer(address common.Address, filterer bind.ContractFilterer) (*GroupMessagesFilterer, error) {
	contract, err := bindGroupMessages(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesFilterer{contract: contract}, nil
}

// bindGroupMessages binds a generic wrapper to an already deployed contract.
func bindGroupMessages(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GroupMessagesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessages *GroupMessagesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessages.Contract.GroupMessagesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessages *GroupMessagesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessages.Contract.GroupMessagesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessages *GroupMessagesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessages.Contract.GroupMessagesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GroupMessages *GroupMessagesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GroupMessages.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GroupMessages *GroupMessagesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessages.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GroupMessages *GroupMessagesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GroupMessages.Contract.contract.Transact(opts, method, params...)
}

// ABSOLUTEMAXPAYLOADSIZE is a free data retrieval call binding the contract method 0x0009e127.
//
// Solidity: function ABSOLUTE_MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCaller) ABSOLUTEMAXPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "ABSOLUTE_MAX_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ABSOLUTEMAXPAYLOADSIZE is a free data retrieval call binding the contract method 0x0009e127.
//
// Solidity: function ABSOLUTE_MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesSession) ABSOLUTEMAXPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.ABSOLUTEMAXPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// ABSOLUTEMAXPAYLOADSIZE is a free data retrieval call binding the contract method 0x0009e127.
//
// Solidity: function ABSOLUTE_MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCallerSession) ABSOLUTEMAXPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.ABSOLUTEMAXPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// ABSOLUTEMINPAYLOADSIZE is a free data retrieval call binding the contract method 0x1de01599.
//
// Solidity: function ABSOLUTE_MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCaller) ABSOLUTEMINPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "ABSOLUTE_MIN_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ABSOLUTEMINPAYLOADSIZE is a free data retrieval call binding the contract method 0x1de01599.
//
// Solidity: function ABSOLUTE_MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesSession) ABSOLUTEMINPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.ABSOLUTEMINPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// ABSOLUTEMINPAYLOADSIZE is a free data retrieval call binding the contract method 0x1de01599.
//
// Solidity: function ABSOLUTE_MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCallerSession) ABSOLUTEMINPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.ABSOLUTEMINPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GroupMessages *GroupMessagesCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GroupMessages *GroupMessagesSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GroupMessages.Contract.DEFAULTADMINROLE(&_GroupMessages.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GroupMessages *GroupMessagesCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GroupMessages.Contract.DEFAULTADMINROLE(&_GroupMessages.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GroupMessages *GroupMessagesCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GroupMessages *GroupMessagesSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _GroupMessages.Contract.UPGRADEINTERFACEVERSION(&_GroupMessages.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GroupMessages *GroupMessagesCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _GroupMessages.Contract.UPGRADEINTERFACEVERSION(&_GroupMessages.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GroupMessages *GroupMessagesCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GroupMessages *GroupMessagesSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GroupMessages.Contract.GetRoleAdmin(&_GroupMessages.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GroupMessages *GroupMessagesCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GroupMessages.Contract.GetRoleAdmin(&_GroupMessages.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GroupMessages *GroupMessagesCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GroupMessages *GroupMessagesSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GroupMessages.Contract.HasRole(&_GroupMessages.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GroupMessages *GroupMessagesCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GroupMessages.Contract.HasRole(&_GroupMessages.CallOpts, role, account)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint256 size)
func (_GroupMessages *GroupMessagesCaller) MaxPayloadSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "maxPayloadSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint256 size)
func (_GroupMessages *GroupMessagesSession) MaxPayloadSize() (*big.Int, error) {
	return _GroupMessages.Contract.MaxPayloadSize(&_GroupMessages.CallOpts)
}

// MaxPayloadSize is a free data retrieval call binding the contract method 0x58e3e94c.
//
// Solidity: function maxPayloadSize() view returns(uint256 size)
func (_GroupMessages *GroupMessagesCallerSession) MaxPayloadSize() (*big.Int, error) {
	return _GroupMessages.Contract.MaxPayloadSize(&_GroupMessages.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint256 size)
func (_GroupMessages *GroupMessagesCaller) MinPayloadSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "minPayloadSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint256 size)
func (_GroupMessages *GroupMessagesSession) MinPayloadSize() (*big.Int, error) {
	return _GroupMessages.Contract.MinPayloadSize(&_GroupMessages.CallOpts)
}

// MinPayloadSize is a free data retrieval call binding the contract method 0xf96927ac.
//
// Solidity: function minPayloadSize() view returns(uint256 size)
func (_GroupMessages *GroupMessagesCallerSession) MinPayloadSize() (*big.Int, error) {
	return _GroupMessages.Contract.MinPayloadSize(&_GroupMessages.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GroupMessages *GroupMessagesCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GroupMessages *GroupMessagesSession) Paused() (bool, error) {
	return _GroupMessages.Contract.Paused(&_GroupMessages.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GroupMessages *GroupMessagesCallerSession) Paused() (bool, error) {
	return _GroupMessages.Contract.Paused(&_GroupMessages.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GroupMessages *GroupMessagesCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GroupMessages *GroupMessagesSession) ProxiableUUID() ([32]byte, error) {
	return _GroupMessages.Contract.ProxiableUUID(&_GroupMessages.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GroupMessages *GroupMessagesCallerSession) ProxiableUUID() ([32]byte, error) {
	return _GroupMessages.Contract.ProxiableUUID(&_GroupMessages.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GroupMessages *GroupMessagesCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GroupMessages *GroupMessagesSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GroupMessages.Contract.SupportsInterface(&_GroupMessages.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GroupMessages *GroupMessagesCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GroupMessages.Contract.SupportsInterface(&_GroupMessages.CallOpts, interfaceId)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessages *GroupMessagesTransactor) AddMessage(opts *bind.TransactOpts, groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "addMessage", groupId, message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessages *GroupMessagesSession) AddMessage(groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessages.Contract.AddMessage(&_GroupMessages.TransactOpts, groupId, message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x4dff26b5.
//
// Solidity: function addMessage(bytes32 groupId, bytes message) returns()
func (_GroupMessages *GroupMessagesTransactorSession) AddMessage(groupId [32]byte, message []byte) (*types.Transaction, error) {
	return _GroupMessages.Contract.AddMessage(&_GroupMessages.TransactOpts, groupId, message)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GroupMessages *GroupMessagesTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GroupMessages *GroupMessagesSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.GrantRole(&_GroupMessages.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GroupMessages *GroupMessagesTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.GrantRole(&_GroupMessages.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_GroupMessages *GroupMessagesTransactor) Initialize(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "initialize", admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_GroupMessages *GroupMessagesSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.Initialize(&_GroupMessages.TransactOpts, admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_GroupMessages *GroupMessagesTransactorSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.Initialize(&_GroupMessages.TransactOpts, admin)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GroupMessages *GroupMessagesTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GroupMessages *GroupMessagesSession) Pause() (*types.Transaction, error) {
	return _GroupMessages.Contract.Pause(&_GroupMessages.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GroupMessages *GroupMessagesTransactorSession) Pause() (*types.Transaction, error) {
	return _GroupMessages.Contract.Pause(&_GroupMessages.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GroupMessages *GroupMessagesTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GroupMessages *GroupMessagesSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.RenounceRole(&_GroupMessages.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GroupMessages *GroupMessagesTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.RenounceRole(&_GroupMessages.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GroupMessages *GroupMessagesTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GroupMessages *GroupMessagesSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.RevokeRole(&_GroupMessages.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GroupMessages *GroupMessagesTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.RevokeRole(&_GroupMessages.TransactOpts, role, account)
}

// SetMaxPayloadSize is a paid mutator transaction binding the contract method 0xfe8e37a3.
//
// Solidity: function setMaxPayloadSize(uint256 maxPayloadSizeRequest) returns()
func (_GroupMessages *GroupMessagesTransactor) SetMaxPayloadSize(opts *bind.TransactOpts, maxPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "setMaxPayloadSize", maxPayloadSizeRequest)
}

// SetMaxPayloadSize is a paid mutator transaction binding the contract method 0xfe8e37a3.
//
// Solidity: function setMaxPayloadSize(uint256 maxPayloadSizeRequest) returns()
func (_GroupMessages *GroupMessagesSession) SetMaxPayloadSize(maxPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessages.Contract.SetMaxPayloadSize(&_GroupMessages.TransactOpts, maxPayloadSizeRequest)
}

// SetMaxPayloadSize is a paid mutator transaction binding the contract method 0xfe8e37a3.
//
// Solidity: function setMaxPayloadSize(uint256 maxPayloadSizeRequest) returns()
func (_GroupMessages *GroupMessagesTransactorSession) SetMaxPayloadSize(maxPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessages.Contract.SetMaxPayloadSize(&_GroupMessages.TransactOpts, maxPayloadSizeRequest)
}

// SetMinPayloadSize is a paid mutator transaction binding the contract method 0x314a100e.
//
// Solidity: function setMinPayloadSize(uint256 minPayloadSizeRequest) returns()
func (_GroupMessages *GroupMessagesTransactor) SetMinPayloadSize(opts *bind.TransactOpts, minPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "setMinPayloadSize", minPayloadSizeRequest)
}

// SetMinPayloadSize is a paid mutator transaction binding the contract method 0x314a100e.
//
// Solidity: function setMinPayloadSize(uint256 minPayloadSizeRequest) returns()
func (_GroupMessages *GroupMessagesSession) SetMinPayloadSize(minPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessages.Contract.SetMinPayloadSize(&_GroupMessages.TransactOpts, minPayloadSizeRequest)
}

// SetMinPayloadSize is a paid mutator transaction binding the contract method 0x314a100e.
//
// Solidity: function setMinPayloadSize(uint256 minPayloadSizeRequest) returns()
func (_GroupMessages *GroupMessagesTransactorSession) SetMinPayloadSize(minPayloadSizeRequest *big.Int) (*types.Transaction, error) {
	return _GroupMessages.Contract.SetMinPayloadSize(&_GroupMessages.TransactOpts, minPayloadSizeRequest)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GroupMessages *GroupMessagesTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GroupMessages *GroupMessagesSession) Unpause() (*types.Transaction, error) {
	return _GroupMessages.Contract.Unpause(&_GroupMessages.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GroupMessages *GroupMessagesTransactorSession) Unpause() (*types.Transaction, error) {
	return _GroupMessages.Contract.Unpause(&_GroupMessages.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GroupMessages *GroupMessagesTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GroupMessages *GroupMessagesSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GroupMessages.Contract.UpgradeToAndCall(&_GroupMessages.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GroupMessages *GroupMessagesTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GroupMessages.Contract.UpgradeToAndCall(&_GroupMessages.TransactOpts, newImplementation, data)
}

// GroupMessagesInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GroupMessages contract.
type GroupMessagesInitializedIterator struct {
	Event *GroupMessagesInitialized // Event containing the contract specifics and raw log

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
func (it *GroupMessagesInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesInitialized)
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
		it.Event = new(GroupMessagesInitialized)
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
func (it *GroupMessagesInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesInitialized represents a Initialized event raised by the GroupMessages contract.
type GroupMessagesInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessages *GroupMessagesFilterer) FilterInitialized(opts *bind.FilterOpts) (*GroupMessagesInitializedIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesInitializedIterator{contract: _GroupMessages.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GroupMessages *GroupMessagesFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GroupMessagesInitialized) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesInitialized)
				if err := _GroupMessages.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseInitialized(log types.Log) (*GroupMessagesInitialized, error) {
	event := new(GroupMessagesInitialized)
	if err := _GroupMessages.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesMaxPayloadSizeUpdatedIterator is returned from FilterMaxPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MaxPayloadSizeUpdated events raised by the GroupMessages contract.
type GroupMessagesMaxPayloadSizeUpdatedIterator struct {
	Event *GroupMessagesMaxPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessagesMaxPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesMaxPayloadSizeUpdated)
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
		it.Event = new(GroupMessagesMaxPayloadSizeUpdated)
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
func (it *GroupMessagesMaxPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesMaxPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesMaxPayloadSizeUpdated represents a MaxPayloadSizeUpdated event raised by the GroupMessages contract.
type GroupMessagesMaxPayloadSizeUpdated struct {
	OldSize *big.Int
	NewSize *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMaxPayloadSizeUpdated is a free log retrieval operation binding the contract event 0xf59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe90.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessages *GroupMessagesFilterer) FilterMaxPayloadSizeUpdated(opts *bind.FilterOpts) (*GroupMessagesMaxPayloadSizeUpdatedIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "MaxPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesMaxPayloadSizeUpdatedIterator{contract: _GroupMessages.contract, event: "MaxPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxPayloadSizeUpdated is a free log subscription operation binding the contract event 0xf59e99f8f54d2696b7cf184949ab2b4bbd6336ec1816b36f58ae9948d868fe90.
//
// Solidity: event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessages *GroupMessagesFilterer) WatchMaxPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessagesMaxPayloadSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "MaxPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesMaxPayloadSizeUpdated)
				if err := _GroupMessages.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseMaxPayloadSizeUpdated(log types.Log) (*GroupMessagesMaxPayloadSizeUpdated, error) {
	event := new(GroupMessagesMaxPayloadSizeUpdated)
	if err := _GroupMessages.contract.UnpackLog(event, "MaxPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesMessageSentIterator is returned from FilterMessageSent and is used to iterate over the raw logs and unpacked data for MessageSent events raised by the GroupMessages contract.
type GroupMessagesMessageSentIterator struct {
	Event *GroupMessagesMessageSent // Event containing the contract specifics and raw log

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
func (it *GroupMessagesMessageSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesMessageSent)
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
		it.Event = new(GroupMessagesMessageSent)
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
func (it *GroupMessagesMessageSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesMessageSent represents a MessageSent event raised by the GroupMessages contract.
type GroupMessagesMessageSent struct {
	GroupId    [32]byte
	Message    []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMessageSent is a free log retrieval operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessages *GroupMessagesFilterer) FilterMessageSent(opts *bind.FilterOpts) (*GroupMessagesMessageSentIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesMessageSentIterator{contract: _GroupMessages.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

// WatchMessageSent is a free log subscription operation binding the contract event 0x91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e.
//
// Solidity: event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId)
func (_GroupMessages *GroupMessagesFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *GroupMessagesMessageSent) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesMessageSent)
				if err := _GroupMessages.contract.UnpackLog(event, "MessageSent", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseMessageSent(log types.Log) (*GroupMessagesMessageSent, error) {
	event := new(GroupMessagesMessageSent)
	if err := _GroupMessages.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesMinPayloadSizeUpdatedIterator is returned from FilterMinPayloadSizeUpdated and is used to iterate over the raw logs and unpacked data for MinPayloadSizeUpdated events raised by the GroupMessages contract.
type GroupMessagesMinPayloadSizeUpdatedIterator struct {
	Event *GroupMessagesMinPayloadSizeUpdated // Event containing the contract specifics and raw log

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
func (it *GroupMessagesMinPayloadSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesMinPayloadSizeUpdated)
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
		it.Event = new(GroupMessagesMinPayloadSizeUpdated)
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
func (it *GroupMessagesMinPayloadSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesMinPayloadSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesMinPayloadSizeUpdated represents a MinPayloadSizeUpdated event raised by the GroupMessages contract.
type GroupMessagesMinPayloadSizeUpdated struct {
	OldSize *big.Int
	NewSize *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinPayloadSizeUpdated is a free log retrieval operation binding the contract event 0x1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd.
//
// Solidity: event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessages *GroupMessagesFilterer) FilterMinPayloadSizeUpdated(opts *bind.FilterOpts) (*GroupMessagesMinPayloadSizeUpdatedIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "MinPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesMinPayloadSizeUpdatedIterator{contract: _GroupMessages.contract, event: "MinPayloadSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinPayloadSizeUpdated is a free log subscription operation binding the contract event 0x1ee836faee0e7c61d20a079d0b5b4e1ee9c536e18268ef6f7c620dcec82f72cd.
//
// Solidity: event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize)
func (_GroupMessages *GroupMessagesFilterer) WatchMinPayloadSizeUpdated(opts *bind.WatchOpts, sink chan<- *GroupMessagesMinPayloadSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "MinPayloadSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesMinPayloadSizeUpdated)
				if err := _GroupMessages.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseMinPayloadSizeUpdated(log types.Log) (*GroupMessagesMinPayloadSizeUpdated, error) {
	event := new(GroupMessagesMinPayloadSizeUpdated)
	if err := _GroupMessages.contract.UnpackLog(event, "MinPayloadSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GroupMessages contract.
type GroupMessagesPausedIterator struct {
	Event *GroupMessagesPaused // Event containing the contract specifics and raw log

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
func (it *GroupMessagesPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesPaused)
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
		it.Event = new(GroupMessagesPaused)
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
func (it *GroupMessagesPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesPaused represents a Paused event raised by the GroupMessages contract.
type GroupMessagesPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GroupMessages *GroupMessagesFilterer) FilterPaused(opts *bind.FilterOpts) (*GroupMessagesPausedIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesPausedIterator{contract: _GroupMessages.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GroupMessages *GroupMessagesFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GroupMessagesPaused) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesPaused)
				if err := _GroupMessages.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParsePaused(log types.Log) (*GroupMessagesPaused, error) {
	event := new(GroupMessagesPaused)
	if err := _GroupMessages.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the GroupMessages contract.
type GroupMessagesRoleAdminChangedIterator struct {
	Event *GroupMessagesRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *GroupMessagesRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesRoleAdminChanged)
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
		it.Event = new(GroupMessagesRoleAdminChanged)
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
func (it *GroupMessagesRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesRoleAdminChanged represents a RoleAdminChanged event raised by the GroupMessages contract.
type GroupMessagesRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GroupMessages *GroupMessagesFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*GroupMessagesRoleAdminChangedIterator, error) {

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

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesRoleAdminChangedIterator{contract: _GroupMessages.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GroupMessages *GroupMessagesFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *GroupMessagesRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesRoleAdminChanged)
				if err := _GroupMessages.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseRoleAdminChanged(log types.Log) (*GroupMessagesRoleAdminChanged, error) {
	event := new(GroupMessagesRoleAdminChanged)
	if err := _GroupMessages.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the GroupMessages contract.
type GroupMessagesRoleGrantedIterator struct {
	Event *GroupMessagesRoleGranted // Event containing the contract specifics and raw log

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
func (it *GroupMessagesRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesRoleGranted)
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
		it.Event = new(GroupMessagesRoleGranted)
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
func (it *GroupMessagesRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesRoleGranted represents a RoleGranted event raised by the GroupMessages contract.
type GroupMessagesRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessages *GroupMessagesFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GroupMessagesRoleGrantedIterator, error) {

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

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesRoleGrantedIterator{contract: _GroupMessages.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessages *GroupMessagesFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *GroupMessagesRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesRoleGranted)
				if err := _GroupMessages.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseRoleGranted(log types.Log) (*GroupMessagesRoleGranted, error) {
	event := new(GroupMessagesRoleGranted)
	if err := _GroupMessages.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the GroupMessages contract.
type GroupMessagesRoleRevokedIterator struct {
	Event *GroupMessagesRoleRevoked // Event containing the contract specifics and raw log

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
func (it *GroupMessagesRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesRoleRevoked)
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
		it.Event = new(GroupMessagesRoleRevoked)
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
func (it *GroupMessagesRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesRoleRevoked represents a RoleRevoked event raised by the GroupMessages contract.
type GroupMessagesRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessages *GroupMessagesFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GroupMessagesRoleRevokedIterator, error) {

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

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesRoleRevokedIterator{contract: _GroupMessages.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GroupMessages *GroupMessagesFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *GroupMessagesRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesRoleRevoked)
				if err := _GroupMessages.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseRoleRevoked(log types.Log) (*GroupMessagesRoleRevoked, error) {
	event := new(GroupMessagesRoleRevoked)
	if err := _GroupMessages.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GroupMessages contract.
type GroupMessagesUnpausedIterator struct {
	Event *GroupMessagesUnpaused // Event containing the contract specifics and raw log

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
func (it *GroupMessagesUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesUnpaused)
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
		it.Event = new(GroupMessagesUnpaused)
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
func (it *GroupMessagesUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesUnpaused represents a Unpaused event raised by the GroupMessages contract.
type GroupMessagesUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GroupMessages *GroupMessagesFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GroupMessagesUnpausedIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesUnpausedIterator{contract: _GroupMessages.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GroupMessages *GroupMessagesFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GroupMessagesUnpaused) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesUnpaused)
				if err := _GroupMessages.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseUnpaused(log types.Log) (*GroupMessagesUnpaused, error) {
	event := new(GroupMessagesUnpaused)
	if err := _GroupMessages.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesUpgradeAuthorizedIterator is returned from FilterUpgradeAuthorized and is used to iterate over the raw logs and unpacked data for UpgradeAuthorized events raised by the GroupMessages contract.
type GroupMessagesUpgradeAuthorizedIterator struct {
	Event *GroupMessagesUpgradeAuthorized // Event containing the contract specifics and raw log

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
func (it *GroupMessagesUpgradeAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesUpgradeAuthorized)
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
		it.Event = new(GroupMessagesUpgradeAuthorized)
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
func (it *GroupMessagesUpgradeAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesUpgradeAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesUpgradeAuthorized represents a UpgradeAuthorized event raised by the GroupMessages contract.
type GroupMessagesUpgradeAuthorized struct {
	Upgrader          common.Address
	NewImplementation common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeAuthorized is a free log retrieval operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_GroupMessages *GroupMessagesFilterer) FilterUpgradeAuthorized(opts *bind.FilterOpts) (*GroupMessagesUpgradeAuthorizedIterator, error) {

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return &GroupMessagesUpgradeAuthorizedIterator{contract: _GroupMessages.contract, event: "UpgradeAuthorized", logs: logs, sub: sub}, nil
}

// WatchUpgradeAuthorized is a free log subscription operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_GroupMessages *GroupMessagesFilterer) WatchUpgradeAuthorized(opts *bind.WatchOpts, sink chan<- *GroupMessagesUpgradeAuthorized) (event.Subscription, error) {

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesUpgradeAuthorized)
				if err := _GroupMessages.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseUpgradeAuthorized(log types.Log) (*GroupMessagesUpgradeAuthorized, error) {
	event := new(GroupMessagesUpgradeAuthorized)
	if err := _GroupMessages.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GroupMessagesUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the GroupMessages contract.
type GroupMessagesUpgradedIterator struct {
	Event *GroupMessagesUpgraded // Event containing the contract specifics and raw log

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
func (it *GroupMessagesUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GroupMessagesUpgraded)
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
		it.Event = new(GroupMessagesUpgraded)
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
func (it *GroupMessagesUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GroupMessagesUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GroupMessagesUpgraded represents a Upgraded event raised by the GroupMessages contract.
type GroupMessagesUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessages *GroupMessagesFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*GroupMessagesUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GroupMessages.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &GroupMessagesUpgradedIterator{contract: _GroupMessages.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GroupMessages *GroupMessagesFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *GroupMessagesUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GroupMessages.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GroupMessagesUpgraded)
				if err := _GroupMessages.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_GroupMessages *GroupMessagesFilterer) ParseUpgraded(log types.Log) (*GroupMessagesUpgraded, error) {
	event := new(GroupMessagesUpgraded)
	if err := _GroupMessages.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
