// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abis

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

// IdentityUpdatesMetaData contains all meta data concerning the IdentityUpdates contract.
var IdentityUpdatesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_PAYLOAD_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addIdentityUpdate\",\"inputs\":[{\"name\":\"inboxId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"update\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"IdentityUpdateCreated\",\"inputs\":[{\"name\":\"inboxId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"update\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sequenceId\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeAuthorized\",\"inputs\":[{\"name\":\"upgrader\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newImplementation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayloadSize\",\"inputs\":[{\"name\":\"actualSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAdminAddress\",\"inputs\":[]}]",
	Bin: "0x60a0604052306080523480156012575f5ffd5b506080516117be6100395f395f8181610b4601528181610b6f0152610e8401526117be5ff3fe608060405260043610610109575f3560e01c80638456cb59116100a1578063ad3cb1cc11610071578063c4d66de811610057578063c4d66de81461036e578063d547741f1461038d578063f5368506146103ac575f5ffd5b8063ad3cb1cc146102fa578063ba74fc7c1461034f575f5ffd5b80638456cb591461024d57806391d1485414610261578063a217fddf146102d1578063a2ba1934146102e4575f5ffd5b80633f4ba83a116100dc5780633f4ba83a146101dc5780634f1ef286146101f057806352d1902d146102035780635c975abb14610217575f5ffd5b806301ffc9a71461010d578063248a9ca3146101415780632f2ff15d1461019c57806336568abe146101bd575b5f5ffd5b348015610118575f5ffd5b5061012c61012736600461143c565b6103c0565b60405190151581526020015b60405180910390f35b34801561014c575f5ffd5b5061018e61015b36600461147b565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b604051908152602001610138565b3480156101a7575f5ffd5b506101bb6101b63660046114ba565b610458565b005b3480156101c8575f5ffd5b506101bb6101d73660046114ba565b6104a1565b3480156101e7575f5ffd5b506101bb6104ff565b6101bb6101fe366004611511565b610514565b34801561020e575f5ffd5b5061018e610533565b348015610222575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1661012c565b348015610258575f5ffd5b506101bb610561565b34801561026c575f5ffd5b5061012c61027b3660046114ba565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b3480156102dc575f5ffd5b5061018e5f81565b3480156102ef575f5ffd5b5061018e6240000081565b348015610305575f5ffd5b506103426040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516101389190611612565b34801561035a575f5ffd5b506101bb610369366004611665565b610573565b348015610379575f5ffd5b506101bb6103883660046116dc565b610665565b348015610398575f5ffd5b506101bb6103a73660046114ba565b61084a565b3480156103b7575f5ffd5b5061018e606881565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061045257507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546104918161088d565b61049b8383610897565b50505050565b73ffffffffffffffffffffffffffffffffffffffff811633146104f0576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104fa82826109b5565b505050565b5f6105098161088d565b610511610a91565b50565b61051c610b2e565b61052582610c34565b61052f8282610d33565b5050565b5f61053c610e6c565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b5f61056b8161088d565b610511610edb565b61057b610f54565b6068811080159061058f5750624000008111155b819060689062400000906105e5576040517f93b7abe60000000000000000000000000000000000000000000000000000000081526004810193909352602483019190915260448201526064015b60405180910390fd5b50505f805467ffffffffffffffff808216600101167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090911681179091556040517fc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd59250610658918691869186916116f5565b60405180910390a1505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156106af5750825b90505f8267ffffffffffffffff1660011480156106cb5750303b155b9050811580156106d9575080155b15610710576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016600117855583156107715784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b73ffffffffffffffffffffffffffffffffffffffff86166107be576040517f3ef39b8100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107c6610fb0565b6107ce610fb0565b6107d6610fb8565b6107e05f87610897565b5083156108425784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546108838161088d565b61049b83836109b5565b6105118133610fc8565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff166109ac575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556109483390565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050610452565b5f915050610452565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff16156109ac575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a46001915050610452565b610a9961106e565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a150565b3073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161480610bfb57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610be27f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b15610c32576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f610c3e8161088d565b73ffffffffffffffffffffffffffffffffffffffff8216610ce1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f4e657720696d706c656d656e746174696f6e2063616e6e6f74206265207a657260448201527f6f2061646472657373000000000000000000000000000000000000000000000060648201526084016105dc565b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201527fd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa910160405180910390a15050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610db8575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201909252610db59181019061175b565b60015b610e06576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016105dc565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8114610e62576040517faa1d49a4000000000000000000000000000000000000000000000000000000008152600481018290526024016105dc565b6104fa83836110c9565b3073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610c32576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ee3610f54565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833610b03565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610c32576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c3261112b565b610fc061112b565b610c32611192565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1661052f576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602481018390526044016105dc565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610c32576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110d2826111e3565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115611123576104fa82826112b1565b61052f611330565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610c32576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61119a61112b565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b8073ffffffffffffffffffffffffffffffffffffffff163b5f0361124b576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016105dc565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff16846040516112da9190611772565b5f60405180830381855af49150503d805f8114611312576040519150601f19603f3d011682016040523d82523d5f602084013e611317565b606091505b5091509150611327858383611368565b95945050505050565b3415610c32576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608261137d57611378826113fa565b6113f3565b81511580156113a1575073ffffffffffffffffffffffffffffffffffffffff84163b155b156113f0576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024016105dc565b50805b9392505050565b80511561140a5780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6020828403121561144c575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146113f3575f5ffd5b5f6020828403121561148b575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146114b5575f5ffd5b919050565b5f5f604083850312156114cb575f5ffd5b823591506114db60208401611492565b90509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611522575f5ffd5b61152b83611492565b9150602083013567ffffffffffffffff811115611546575f5ffd5b8301601f81018513611556575f5ffd5b803567ffffffffffffffff811115611570576115706114e4565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156115dc576115dc6114e4565b6040528181528282016020018710156115f3575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f5f5f60408486031215611677575f5ffd5b83359250602084013567ffffffffffffffff811115611694575f5ffd5b8401601f810186136116a4575f5ffd5b803567ffffffffffffffff8111156116ba575f5ffd5b8660208284010111156116cb575f5ffd5b939660209190910195509293505050565b5f602082840312156116ec575f5ffd5b6113f382611492565b84815260606020820152826060820152828460808301375f608084830101525f60807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116830101905067ffffffffffffffff8316604083015295945050505050565b5f6020828403121561176b575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fea26469706673582212202d65f062d705cb761a34070256e7c4f8eab20ae6669f6158830cb02fcd42c16164736f6c634300081c0033",
}

// IdentityUpdatesABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityUpdatesMetaData.ABI instead.
var IdentityUpdatesABI = IdentityUpdatesMetaData.ABI

// IdentityUpdatesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IdentityUpdatesMetaData.Bin instead.
var IdentityUpdatesBin = IdentityUpdatesMetaData.Bin

// DeployIdentityUpdates deploys a new Ethereum contract, binding an instance of IdentityUpdates to it.
func DeployIdentityUpdates(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IdentityUpdates, error) {
	parsed, err := IdentityUpdatesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IdentityUpdatesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IdentityUpdates{IdentityUpdatesCaller: IdentityUpdatesCaller{contract: contract}, IdentityUpdatesTransactor: IdentityUpdatesTransactor{contract: contract}, IdentityUpdatesFilterer: IdentityUpdatesFilterer{contract: contract}}, nil
}

// IdentityUpdates is an auto generated Go binding around an Ethereum contract.
type IdentityUpdates struct {
	IdentityUpdatesCaller     // Read-only binding to the contract
	IdentityUpdatesTransactor // Write-only binding to the contract
	IdentityUpdatesFilterer   // Log filterer for contract events
}

// IdentityUpdatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityUpdatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityUpdatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityUpdatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityUpdatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityUpdatesSession struct {
	Contract     *IdentityUpdates  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IdentityUpdatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityUpdatesCallerSession struct {
	Contract *IdentityUpdatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IdentityUpdatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityUpdatesTransactorSession struct {
	Contract     *IdentityUpdatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IdentityUpdatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityUpdatesRaw struct {
	Contract *IdentityUpdates // Generic contract binding to access the raw methods on
}

// IdentityUpdatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityUpdatesCallerRaw struct {
	Contract *IdentityUpdatesCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityUpdatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityUpdatesTransactorRaw struct {
	Contract *IdentityUpdatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityUpdates creates a new instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdates(address common.Address, backend bind.ContractBackend) (*IdentityUpdates, error) {
	contract, err := bindIdentityUpdates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdates{IdentityUpdatesCaller: IdentityUpdatesCaller{contract: contract}, IdentityUpdatesTransactor: IdentityUpdatesTransactor{contract: contract}, IdentityUpdatesFilterer: IdentityUpdatesFilterer{contract: contract}}, nil
}

// NewIdentityUpdatesCaller creates a new read-only instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdatesCaller(address common.Address, caller bind.ContractCaller) (*IdentityUpdatesCaller, error) {
	contract, err := bindIdentityUpdates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesCaller{contract: contract}, nil
}

// NewIdentityUpdatesTransactor creates a new write-only instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdatesTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityUpdatesTransactor, error) {
	contract, err := bindIdentityUpdates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesTransactor{contract: contract}, nil
}

// NewIdentityUpdatesFilterer creates a new log filterer instance of IdentityUpdates, bound to a specific deployed contract.
func NewIdentityUpdatesFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityUpdatesFilterer, error) {
	contract, err := bindIdentityUpdates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesFilterer{contract: contract}, nil
}

// bindIdentityUpdates binds a generic wrapper to an already deployed contract.
func bindIdentityUpdates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityUpdatesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityUpdates *IdentityUpdatesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityUpdates.Contract.IdentityUpdatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityUpdates *IdentityUpdatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.IdentityUpdatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityUpdates *IdentityUpdatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.IdentityUpdatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityUpdates *IdentityUpdatesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityUpdates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityUpdates *IdentityUpdatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityUpdates *IdentityUpdatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _IdentityUpdates.Contract.DEFAULTADMINROLE(&_IdentityUpdates.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _IdentityUpdates.Contract.DEFAULTADMINROLE(&_IdentityUpdates.CallOpts)
}

// MAXPAYLOADSIZE is a free data retrieval call binding the contract method 0xa2ba1934.
//
// Solidity: function MAX_PAYLOAD_SIZE() view returns(uint256)
func (_IdentityUpdates *IdentityUpdatesCaller) MAXPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "MAX_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXPAYLOADSIZE is a free data retrieval call binding the contract method 0xa2ba1934.
//
// Solidity: function MAX_PAYLOAD_SIZE() view returns(uint256)
func (_IdentityUpdates *IdentityUpdatesSession) MAXPAYLOADSIZE() (*big.Int, error) {
	return _IdentityUpdates.Contract.MAXPAYLOADSIZE(&_IdentityUpdates.CallOpts)
}

// MAXPAYLOADSIZE is a free data retrieval call binding the contract method 0xa2ba1934.
//
// Solidity: function MAX_PAYLOAD_SIZE() view returns(uint256)
func (_IdentityUpdates *IdentityUpdatesCallerSession) MAXPAYLOADSIZE() (*big.Int, error) {
	return _IdentityUpdates.Contract.MAXPAYLOADSIZE(&_IdentityUpdates.CallOpts)
}

// MINPAYLOADSIZE is a free data retrieval call binding the contract method 0xf5368506.
//
// Solidity: function MIN_PAYLOAD_SIZE() view returns(uint256)
func (_IdentityUpdates *IdentityUpdatesCaller) MINPAYLOADSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "MIN_PAYLOAD_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINPAYLOADSIZE is a free data retrieval call binding the contract method 0xf5368506.
//
// Solidity: function MIN_PAYLOAD_SIZE() view returns(uint256)
func (_IdentityUpdates *IdentityUpdatesSession) MINPAYLOADSIZE() (*big.Int, error) {
	return _IdentityUpdates.Contract.MINPAYLOADSIZE(&_IdentityUpdates.CallOpts)
}

// MINPAYLOADSIZE is a free data retrieval call binding the contract method 0xf5368506.
//
// Solidity: function MIN_PAYLOAD_SIZE() view returns(uint256)
func (_IdentityUpdates *IdentityUpdatesCallerSession) MINPAYLOADSIZE() (*big.Int, error) {
	return _IdentityUpdates.Contract.MINPAYLOADSIZE(&_IdentityUpdates.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_IdentityUpdates *IdentityUpdatesCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_IdentityUpdates *IdentityUpdatesSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _IdentityUpdates.Contract.UPGRADEINTERFACEVERSION(&_IdentityUpdates.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_IdentityUpdates *IdentityUpdatesCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _IdentityUpdates.Contract.UPGRADEINTERFACEVERSION(&_IdentityUpdates.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _IdentityUpdates.Contract.GetRoleAdmin(&_IdentityUpdates.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _IdentityUpdates.Contract.GetRoleAdmin(&_IdentityUpdates.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_IdentityUpdates *IdentityUpdatesCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_IdentityUpdates *IdentityUpdatesSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _IdentityUpdates.Contract.HasRole(&_IdentityUpdates.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_IdentityUpdates *IdentityUpdatesCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _IdentityUpdates.Contract.HasRole(&_IdentityUpdates.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IdentityUpdates *IdentityUpdatesCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IdentityUpdates *IdentityUpdatesSession) Paused() (bool, error) {
	return _IdentityUpdates.Contract.Paused(&_IdentityUpdates.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IdentityUpdates *IdentityUpdatesCallerSession) Paused() (bool, error) {
	return _IdentityUpdates.Contract.Paused(&_IdentityUpdates.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesSession) ProxiableUUID() ([32]byte, error) {
	return _IdentityUpdates.Contract.ProxiableUUID(&_IdentityUpdates.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IdentityUpdates *IdentityUpdatesCallerSession) ProxiableUUID() ([32]byte, error) {
	return _IdentityUpdates.Contract.ProxiableUUID(&_IdentityUpdates.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_IdentityUpdates *IdentityUpdatesCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _IdentityUpdates.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_IdentityUpdates *IdentityUpdatesSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IdentityUpdates.Contract.SupportsInterface(&_IdentityUpdates.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_IdentityUpdates *IdentityUpdatesCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IdentityUpdates.Contract.SupportsInterface(&_IdentityUpdates.CallOpts, interfaceId)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId, bytes update) returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) AddIdentityUpdate(opts *bind.TransactOpts, inboxId [32]byte, update []byte) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "addIdentityUpdate", inboxId, update)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId, bytes update) returns()
func (_IdentityUpdates *IdentityUpdatesSession) AddIdentityUpdate(inboxId [32]byte, update []byte) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.AddIdentityUpdate(&_IdentityUpdates.TransactOpts, inboxId, update)
}

// AddIdentityUpdate is a paid mutator transaction binding the contract method 0xba74fc7c.
//
// Solidity: function addIdentityUpdate(bytes32 inboxId, bytes update) returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) AddIdentityUpdate(inboxId [32]byte, update []byte) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.AddIdentityUpdate(&_IdentityUpdates.TransactOpts, inboxId, update)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_IdentityUpdates *IdentityUpdatesSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.GrantRole(&_IdentityUpdates.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.GrantRole(&_IdentityUpdates.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) Initialize(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "initialize", _admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_IdentityUpdates *IdentityUpdatesSession) Initialize(_admin common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.Initialize(&_IdentityUpdates.TransactOpts, _admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _admin) returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) Initialize(_admin common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.Initialize(&_IdentityUpdates.TransactOpts, _admin)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IdentityUpdates *IdentityUpdatesSession) Pause() (*types.Transaction, error) {
	return _IdentityUpdates.Contract.Pause(&_IdentityUpdates.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) Pause() (*types.Transaction, error) {
	return _IdentityUpdates.Contract.Pause(&_IdentityUpdates.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_IdentityUpdates *IdentityUpdatesSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.RenounceRole(&_IdentityUpdates.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.RenounceRole(&_IdentityUpdates.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_IdentityUpdates *IdentityUpdatesSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.RevokeRole(&_IdentityUpdates.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.RevokeRole(&_IdentityUpdates.TransactOpts, role, account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IdentityUpdates *IdentityUpdatesSession) Unpause() (*types.Transaction, error) {
	return _IdentityUpdates.Contract.Unpause(&_IdentityUpdates.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) Unpause() (*types.Transaction, error) {
	return _IdentityUpdates.Contract.Unpause(&_IdentityUpdates.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_IdentityUpdates *IdentityUpdatesTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _IdentityUpdates.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_IdentityUpdates *IdentityUpdatesSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.UpgradeToAndCall(&_IdentityUpdates.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_IdentityUpdates *IdentityUpdatesTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _IdentityUpdates.Contract.UpgradeToAndCall(&_IdentityUpdates.TransactOpts, newImplementation, data)
}

// IdentityUpdatesIdentityUpdateCreatedIterator is returned from FilterIdentityUpdateCreated and is used to iterate over the raw logs and unpacked data for IdentityUpdateCreated events raised by the IdentityUpdates contract.
type IdentityUpdatesIdentityUpdateCreatedIterator struct {
	Event *IdentityUpdatesIdentityUpdateCreated // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesIdentityUpdateCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesIdentityUpdateCreated)
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
		it.Event = new(IdentityUpdatesIdentityUpdateCreated)
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
func (it *IdentityUpdatesIdentityUpdateCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesIdentityUpdateCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesIdentityUpdateCreated represents a IdentityUpdateCreated event raised by the IdentityUpdates contract.
type IdentityUpdatesIdentityUpdateCreated struct {
	InboxId    [32]byte
	Update     []byte
	SequenceId uint64
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterIdentityUpdateCreated is a free log retrieval operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterIdentityUpdateCreated(opts *bind.FilterOpts) (*IdentityUpdatesIdentityUpdateCreatedIterator, error) {

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "IdentityUpdateCreated")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesIdentityUpdateCreatedIterator{contract: _IdentityUpdates.contract, event: "IdentityUpdateCreated", logs: logs, sub: sub}, nil
}

// WatchIdentityUpdateCreated is a free log subscription operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchIdentityUpdateCreated(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesIdentityUpdateCreated) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "IdentityUpdateCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesIdentityUpdateCreated)
				if err := _IdentityUpdates.contract.UnpackLog(event, "IdentityUpdateCreated", log); err != nil {
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

// ParseIdentityUpdateCreated is a log parse operation binding the contract event 0xc1a40f292090ec0435e939cdfe248e0322a88566679a90a50c4e9e5ef762dbd5.
//
// Solidity: event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId)
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseIdentityUpdateCreated(log types.Log) (*IdentityUpdatesIdentityUpdateCreated, error) {
	event := new(IdentityUpdatesIdentityUpdateCreated)
	if err := _IdentityUpdates.contract.UnpackLog(event, "IdentityUpdateCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IdentityUpdates contract.
type IdentityUpdatesInitializedIterator struct {
	Event *IdentityUpdatesInitialized // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesInitialized)
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
		it.Event = new(IdentityUpdatesInitialized)
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
func (it *IdentityUpdatesInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesInitialized represents a Initialized event raised by the IdentityUpdates contract.
type IdentityUpdatesInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterInitialized(opts *bind.FilterOpts) (*IdentityUpdatesInitializedIterator, error) {

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesInitializedIterator{contract: _IdentityUpdates.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesInitialized) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesInitialized)
				if err := _IdentityUpdates.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseInitialized(log types.Log) (*IdentityUpdatesInitialized, error) {
	event := new(IdentityUpdatesInitialized)
	if err := _IdentityUpdates.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the IdentityUpdates contract.
type IdentityUpdatesPausedIterator struct {
	Event *IdentityUpdatesPaused // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesPaused)
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
		it.Event = new(IdentityUpdatesPaused)
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
func (it *IdentityUpdatesPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesPaused represents a Paused event raised by the IdentityUpdates contract.
type IdentityUpdatesPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterPaused(opts *bind.FilterOpts) (*IdentityUpdatesPausedIterator, error) {

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesPausedIterator{contract: _IdentityUpdates.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesPaused) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesPaused)
				if err := _IdentityUpdates.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParsePaused(log types.Log) (*IdentityUpdatesPaused, error) {
	event := new(IdentityUpdatesPaused)
	if err := _IdentityUpdates.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the IdentityUpdates contract.
type IdentityUpdatesRoleAdminChangedIterator struct {
	Event *IdentityUpdatesRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesRoleAdminChanged)
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
		it.Event = new(IdentityUpdatesRoleAdminChanged)
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
func (it *IdentityUpdatesRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesRoleAdminChanged represents a RoleAdminChanged event raised by the IdentityUpdates contract.
type IdentityUpdatesRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*IdentityUpdatesRoleAdminChangedIterator, error) {

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

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesRoleAdminChangedIterator{contract: _IdentityUpdates.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesRoleAdminChanged)
				if err := _IdentityUpdates.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseRoleAdminChanged(log types.Log) (*IdentityUpdatesRoleAdminChanged, error) {
	event := new(IdentityUpdatesRoleAdminChanged)
	if err := _IdentityUpdates.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the IdentityUpdates contract.
type IdentityUpdatesRoleGrantedIterator struct {
	Event *IdentityUpdatesRoleGranted // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesRoleGranted)
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
		it.Event = new(IdentityUpdatesRoleGranted)
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
func (it *IdentityUpdatesRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesRoleGranted represents a RoleGranted event raised by the IdentityUpdates contract.
type IdentityUpdatesRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*IdentityUpdatesRoleGrantedIterator, error) {

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

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesRoleGrantedIterator{contract: _IdentityUpdates.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesRoleGranted)
				if err := _IdentityUpdates.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseRoleGranted(log types.Log) (*IdentityUpdatesRoleGranted, error) {
	event := new(IdentityUpdatesRoleGranted)
	if err := _IdentityUpdates.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the IdentityUpdates contract.
type IdentityUpdatesRoleRevokedIterator struct {
	Event *IdentityUpdatesRoleRevoked // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesRoleRevoked)
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
		it.Event = new(IdentityUpdatesRoleRevoked)
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
func (it *IdentityUpdatesRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesRoleRevoked represents a RoleRevoked event raised by the IdentityUpdates contract.
type IdentityUpdatesRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*IdentityUpdatesRoleRevokedIterator, error) {

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

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesRoleRevokedIterator{contract: _IdentityUpdates.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesRoleRevoked)
				if err := _IdentityUpdates.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseRoleRevoked(log types.Log) (*IdentityUpdatesRoleRevoked, error) {
	event := new(IdentityUpdatesRoleRevoked)
	if err := _IdentityUpdates.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the IdentityUpdates contract.
type IdentityUpdatesUnpausedIterator struct {
	Event *IdentityUpdatesUnpaused // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesUnpaused)
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
		it.Event = new(IdentityUpdatesUnpaused)
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
func (it *IdentityUpdatesUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesUnpaused represents a Unpaused event raised by the IdentityUpdates contract.
type IdentityUpdatesUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IdentityUpdatesUnpausedIterator, error) {

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesUnpausedIterator{contract: _IdentityUpdates.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesUnpaused) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesUnpaused)
				if err := _IdentityUpdates.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseUnpaused(log types.Log) (*IdentityUpdatesUnpaused, error) {
	event := new(IdentityUpdatesUnpaused)
	if err := _IdentityUpdates.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesUpgradeAuthorizedIterator is returned from FilterUpgradeAuthorized and is used to iterate over the raw logs and unpacked data for UpgradeAuthorized events raised by the IdentityUpdates contract.
type IdentityUpdatesUpgradeAuthorizedIterator struct {
	Event *IdentityUpdatesUpgradeAuthorized // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesUpgradeAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesUpgradeAuthorized)
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
		it.Event = new(IdentityUpdatesUpgradeAuthorized)
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
func (it *IdentityUpdatesUpgradeAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesUpgradeAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesUpgradeAuthorized represents a UpgradeAuthorized event raised by the IdentityUpdates contract.
type IdentityUpdatesUpgradeAuthorized struct {
	Upgrader          common.Address
	NewImplementation common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeAuthorized is a free log retrieval operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterUpgradeAuthorized(opts *bind.FilterOpts) (*IdentityUpdatesUpgradeAuthorizedIterator, error) {

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesUpgradeAuthorizedIterator{contract: _IdentityUpdates.contract, event: "UpgradeAuthorized", logs: logs, sub: sub}, nil
}

// WatchUpgradeAuthorized is a free log subscription operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchUpgradeAuthorized(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesUpgradeAuthorized) (event.Subscription, error) {

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesUpgradeAuthorized)
				if err := _IdentityUpdates.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseUpgradeAuthorized(log types.Log) (*IdentityUpdatesUpgradeAuthorized, error) {
	event := new(IdentityUpdatesUpgradeAuthorized)
	if err := _IdentityUpdates.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityUpdatesUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IdentityUpdates contract.
type IdentityUpdatesUpgradedIterator struct {
	Event *IdentityUpdatesUpgraded // Event containing the contract specifics and raw log

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
func (it *IdentityUpdatesUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityUpdatesUpgraded)
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
		it.Event = new(IdentityUpdatesUpgraded)
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
func (it *IdentityUpdatesUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityUpdatesUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityUpdatesUpgraded represents a Upgraded event raised by the IdentityUpdates contract.
type IdentityUpdatesUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IdentityUpdates *IdentityUpdatesFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IdentityUpdatesUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IdentityUpdates.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IdentityUpdatesUpgradedIterator{contract: _IdentityUpdates.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IdentityUpdates *IdentityUpdatesFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IdentityUpdatesUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IdentityUpdates.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityUpdatesUpgraded)
				if err := _IdentityUpdates.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_IdentityUpdates *IdentityUpdatesFilterer) ParseUpgraded(log types.Log) (*IdentityUpdatesUpgraded, error) {
	event := new(IdentityUpdatesUpgraded)
	if err := _IdentityUpdates.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
