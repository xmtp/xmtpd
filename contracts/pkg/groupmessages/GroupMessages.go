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
	ABI: "[{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addMessage\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"groupId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeAuthorized\",\"inputs\":[{\"name\":\"upgrader\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newImplementation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAdminAddress\",\"inputs\":[]}]",
	Bin: "0x60a0604052306080523480156012575f5ffd5b506080516117be6100395f395f8181610ba401528181610bcd0152610ee001526117be5ff3fe608060405260043610610109575f3560e01c80635c975abb116100a1578063a2ba193411610071578063c4d66de811610057578063c4d66de81461036e578063d547741f1461038d578063f5368506146103ac575f5ffd5b8063a2ba193414610303578063ad3cb1cc14610319575f5ffd5b80635c975abb146102365780638456cb591461026c57806391d1485414610280578063a217fddf146102f0575f5ffd5b80633f4ba83a116100dc5780633f4ba83a146101dc5780634dff26b5146101f05780634f1ef2861461020f57806352d1902d14610222575f5ffd5b806301ffc9a71461010d578063248a9ca3146101415780632f2ff15d1461019c57806336568abe146101bd575b5f5ffd5b348015610118575f5ffd5b5061012c61012736600461143c565b6103c0565b60405190151581526020015b60405180910390f35b34801561014c575f5ffd5b5061018e61015b36600461147b565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b604051908152602001610138565b3480156101a7575f5ffd5b506101bb6101b63660046114ba565b610458565b005b3480156101c8575f5ffd5b506101bb6101d73660046114ba565b6104a1565b3480156101e7575f5ffd5b506101bb6104ff565b3480156101fb575f5ffd5b506101bb61020a3660046114e4565b610514565b6101bb61021d366004611588565b610606565b34801561022d575f5ffd5b5061018e610625565b348015610241575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1661012c565b348015610277575f5ffd5b506101bb610653565b34801561028b575f5ffd5b5061012c61029a3660046114ba565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b3480156102fb575f5ffd5b5061018e5f81565b34801561030e575f5ffd5b5061018e6240000081565b348015610324575f5ffd5b506103616040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516101389190611689565b348015610379575f5ffd5b506101bb6103883660046116dc565b610665565b348015610398575f5ffd5b506101bb6103a73660046114ba565b61084a565b3480156103b7575f5ffd5b5061018e604e81565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061045257507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546104918161088d565b61049b8383610897565b50505050565b73ffffffffffffffffffffffffffffffffffffffff811633146104f0576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104fa82826109b5565b505050565b5f6105098161088d565b610511610a91565b50565b61051c610b2e565b604e81108015906105305750624000008111155b8190604e906240000090610586576040517f93b7abe60000000000000000000000000000000000000000000000000000000081526004810193909352602483019190915260448201526064015b60405180910390fd5b50505f805467ffffffffffffffff808216600101167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090911681179091556040517f91f47151424884a46811ed593aa8a02ee5012e9332a4dcf1e9236a8ed4443c3e92506105f9918691869186916116f5565b60405180910390a1505050565b61060e610b8c565b61061782610c90565b6106218282610d8f565b5050565b5f61062e610ec8565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b5f61065d8161088d565b610511610f37565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156106af5750825b90505f8267ffffffffffffffff1660011480156106cb5750303b155b9050811580156106d9575080155b15610710576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156107715784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b73ffffffffffffffffffffffffffffffffffffffff86166107be576040517f3ef39b8100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107c6610fb0565b6107ce610fb0565b6107d6610fb8565b6107e05f87610897565b5083156108425784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546108838161088d565b61049b83836109b5565b6105118133610fc8565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff166109ac575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556109483390565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050610452565b5f915050610452565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff16156109ac575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a46001915050610452565b610a9961106e565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a150565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610b8a576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b3073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161480610c5957507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610c407f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b15610b8a576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610c9a8161088d565b73ffffffffffffffffffffffffffffffffffffffff8216610d3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f4e657720696d706c656d656e746174696f6e2063616e6e6f74206265207a657260448201527f6f20616464726573730000000000000000000000000000000000000000000000606482015260840161057d565b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201527fd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa910160405180910390a15050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610e14575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201909252610e119181019061175b565b60015b610e62576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260240161057d565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8114610ebe576040517faa1d49a40000000000000000000000000000000000000000000000000000000081526004810182905260240161057d565b6104fa83836110c9565b3073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610b8a576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610f3f610b2e565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833610b03565b610b8a61112b565b610fc061112b565b610b8a611192565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610621576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024810183905260440161057d565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610b8a576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110d2826111e3565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115611123576104fa82826112b1565b610621611330565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610b8a576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61119a61112b565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b8073ffffffffffffffffffffffffffffffffffffffff163b5f0361124b576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161057d565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff16846040516112da9190611772565b5f60405180830381855af49150503d805f8114611312576040519150601f19603f3d011682016040523d82523d5f602084013e611317565b606091505b5091509150611327858383611368565b95945050505050565b3415610b8a576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608261137d57611378826113fa565b6113f3565b81511580156113a1575073ffffffffffffffffffffffffffffffffffffffff84163b155b156113f0576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260240161057d565b50805b9392505050565b80511561140a5780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6020828403121561144c575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146113f3575f5ffd5b5f6020828403121561148b575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146114b5575f5ffd5b919050565b5f5f604083850312156114cb575f5ffd5b823591506114db60208401611492565b90509250929050565b5f5f5f604084860312156114f6575f5ffd5b83359250602084013567ffffffffffffffff811115611513575f5ffd5b8401601f81018613611523575f5ffd5b803567ffffffffffffffff811115611539575f5ffd5b86602082840101111561154a575f5ffd5b939660209190910195509293505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611599575f5ffd5b6115a283611492565b9150602083013567ffffffffffffffff8111156115bd575f5ffd5b8301601f810185136115cd575f5ffd5b803567ffffffffffffffff8111156115e7576115e761155b565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156116535761165361155b565b60405281815282820160200187101561166a575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f602082840312156116ec575f5ffd5b6113f382611492565b84815260606020820152826060820152828460808301375f608084830101525f60807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116830101905067ffffffffffffffff8316604083015295945050505050565b5f6020828403121561176b575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fea264697066735822122072d7566503253678adc66fd27920aaabe0ac779b70ca5ff95eb4b05309ab3a9064736f6c634300081c0033",
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

// MAXPAYLOADSIZE is a free data retrieval call binding the contract method 0xa2ba1934.
//
// Solidity: function MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCaller) MAXPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "MAX_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXPAYLOADSIZE is a free data retrieval call binding the contract method 0xa2ba1934.
//
// Solidity: function MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesSession) MAXPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.MAXPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// MAXPAYLOADSIZE is a free data retrieval call binding the contract method 0xa2ba1934.
//
// Solidity: function MAX_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCallerSession) MAXPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.MAXPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// MINPAYLOADSIZE is a free data retrieval call binding the contract method 0xf5368506.
//
// Solidity: function MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCaller) MINPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GroupMessages.contract.Call(opts, &out, "MIN_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINPAYLOADSIZE is a free data retrieval call binding the contract method 0xf5368506.
//
// Solidity: function MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesSession) MINPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.MINPAYLOADSIZE(&_GroupMessages.CallOpts)
}

// MINPAYLOADSIZE is a free data retrieval call binding the contract method 0xf5368506.
//
// Solidity: function MIN_PAYLOAD_SIZE() view returns(uint256)
func (_GroupMessages *GroupMessagesCallerSession) MINPAYLOADSIZE() (*big.Int, error) {
	return _GroupMessages.Contract.MINPAYLOADSIZE(&_GroupMessages.CallOpts)
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
// Solidity: function initialize(address _admin) returns()
func (_GroupMessages *GroupMessagesTransactor) Initialize(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _GroupMessages.contract.Transact(opts, "initialize", _admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_GroupMessages *GroupMessagesSession) Initialize(_admin common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.Initialize(&_GroupMessages.TransactOpts, _admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_GroupMessages *GroupMessagesTransactorSession) Initialize(_admin common.Address) (*types.Transaction, error) {
	return _GroupMessages.Contract.Initialize(&_GroupMessages.TransactOpts, _admin)
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
