// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ratesmanager

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

// RatesManagerRates is an auto generated low-level Go binding around an user-defined struct.
type RatesManagerRates struct {
	MessageFee    uint64
	StorageFee    uint64
	CongestionFee uint64
	StartTime     uint64
}

// RatesManagerMetaData contains all meta data concerning the RatesManager contract.
var RatesManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAGE_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RATES_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addRates\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRates\",\"inputs\":[{\"name\":\"fromIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rates\",\"type\":\"tuple[]\",\"internalType\":\"structRatesManager.Rates[]\",\"components\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"hasMore\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRatesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RatesAdded\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeAuthorized\",\"inputs\":[{\"name\":\"upgrader\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newImplementation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FromIndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidStartTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAdminAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroImplementationAddress\",\"inputs\":[]}]",
	Bin: "0x60a0604052306080523480156012575f5ffd5b50608051611d726100395f395f81816110080152818161103101526112f50152611d725ff3fe60806040526004361061013d575f3560e01c80634f1ef286116100bb578063a217fddf11610071578063c4d66de811610057578063c4d66de814610400578063d547741f1461041f578063e33967351461043e575f5ffd5b8063a217fddf14610398578063ad3cb1cc146103ab575f5ffd5b80635c975abb116100a15780635c975abb146102de5780638456cb591461031457806391d1485414610328575f5ffd5b80634f1ef286146102b757806352d1902d146102ca575f5ffd5b80632f2ff15d1161011057806336568abe116100f657806336568abe146102655780633f4ba83a146102845780634441210714610298575f5ffd5b80632f2ff15d1461023057806331b2866414610251575f5ffd5b806301ffc9a714610141578063081802b114610175578063248a9ca3146101a25780632da72291146101fd575b5f5ffd5b34801561014c575f5ffd5b5061016061015b36600461194b565b610471565b60405190151581526020015b60405180910390f35b348015610180575f5ffd5b5061019461018f36600461198a565b610509565b60405161016c9291906119a1565b3480156101ad575f5ffd5b506101ef6101bc36600461198a565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b60405190815260200161016c565b348015610208575f5ffd5b507f6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d000546101ef565b34801561023b575f5ffd5b5061024f61024a366004611a63565b61076e565b005b34801561025c575f5ffd5b506101ef603281565b348015610270575f5ffd5b5061024f61027f366004611a63565b6107b7565b34801561028f575f5ffd5b5061024f610815565b3480156102a3575f5ffd5b5061024f6102b2366004611aa4565b61082a565b61024f6102c5366004611b22565b610a60565b3480156102d5575f5ffd5b506101ef610a7f565b3480156102e9575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610160565b34801561031f575f5ffd5b5061024f610aad565b348015610333575f5ffd5b50610160610342366004611a63565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b3480156103a3575f5ffd5b506101ef5f81565b3480156103b6575f5ffd5b506103f36040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b60405161016c9190611c23565b34801561040b575f5ffd5b5061024f61041a366004611c76565b610abf565b34801561042a575f5ffd5b5061024f610439366004611a63565b610cf5565b348015610449575f5ffd5b506101ef7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e81565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061050357507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b7f6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d00080546060915f9115801561053c575083155b156105b257604080515f80825260208201909252906105a7565b604080516080810182525f8082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816105565790505b50945f945092505050565b805484106105ec576040517fea61fe7000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6106026105fb603287611cbc565b8354610d38565b905061060e8582611ccf565b67ffffffffffffffff81111561062657610626611af5565b60405190808252806020026020018201604052801561069557816020015b604080516080810182525f8082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816106445790505b5093505f5b845181101561076257826106ae8288611cbc565b815481106106be576106be611ce2565b5f91825260209182902060408051608081018252929091015467ffffffffffffffff80821684526801000000000000000082048116948401949094527001000000000000000000000000000000008104841691830191909152780100000000000000000000000000000000000000000000000090049091166060820152855186908390811061074f5761074f611ce2565b602090810291909101015260010161069a565b50905492949211925050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546107a781610d4f565b6107b18383610d59565b50505050565b73ffffffffffffffffffffffffffffffffffffffff81163314610806576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108108282610e77565b505050565b5f61081f81610d4f565b610827610f53565b50565b7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e61085481610d4f565b7f6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d0008054158015906108d957508054819061089090600190611ccf565b815481106108a0576108a0611ce2565b5f9182526020909120015467ffffffffffffffff7801000000000000000000000000000000000000000000000000909104811690841611155b15610910576040517fb290253c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516080808201835267ffffffffffffffff89811680845289821660208086018281528b85168789018181528c87166060808b018281528d54600181018f555f8f81528890209c519c0180549651945191518b1678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff928c1670010000000000000000000000000000000002929092166fffffffffffffffffffffffffffffffff958c1668010000000000000000027fffffffffffffffffffffffffffffffff000000000000000000000000000000009098169d909b169c909c179590951792909216979097179290921790975587519384529083019190915294810193909352928201929092527f3bff7b1e021b47f5dfd21d1d3fe2daaf3b9cbaca733c6e15b3a0da546657f19a91015b60405180910390a1505050505050565b610a68610ff0565b610a71826110f6565b610a7b828261119f565b5050565b5f610a886112dd565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b5f610ab781610d4f565b61082761134c565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f81158015610b095750825b90505f8267ffffffffffffffff166001148015610b255750303b155b905081158015610b33575080155b15610b6a576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001660011785558315610bcb5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b73ffffffffffffffffffffffffffffffffffffffff8616610c18576040517f3ef39b8100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c206113c5565b610c286113c5565b610c306113cd565b610c5a7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e5f6113dd565b610c645f87610d59565b50610c8f7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e87610d59565b508315610ced5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602001610a50565b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610d2e81610d4f565b6107b18383610e77565b5f818310610d465781610d48565b825b9392505050565b610827813361147e565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff16610e6e575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610e0a3390565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050610503565b5f915050610503565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff1615610e6e575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a46001915050610503565b610f5b611524565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a150565b3073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614806110bd57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166110a47f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b156110f4576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f61110081610d4f565b73ffffffffffffffffffffffffffffffffffffffff821661114d576040517fd02c623d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201527fd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa910160405180910390a15050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611224575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261122191810190611d0f565b60015b611277576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81146112d3576040517faa1d49a40000000000000000000000000000000000000000000000000000000081526004810182905260240161126e565b610810838361157f565b3073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146110f4576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113546115e1565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833610fc5565b6110f461163d565b6113d561163d565b6110f46116a4565b7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268005f611436845f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b5f85815260208490526040808220600101869055519192508491839187917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a450505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610a7b576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024810183905260440161126e565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff166110f4576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611588826116f5565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156115d95761081082826117c3565b610a7b611842565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16156110f4576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166110f4576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6116ac61163d565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b8073ffffffffffffffffffffffffffffffffffffffff163b5f0361175d576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161126e565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff16846040516117ec9190611d26565b5f60405180830381855af49150503d805f8114611824576040519150601f19603f3d011682016040523d82523d5f602084013e611829565b606091505b509150915061183985838361187a565b95945050505050565b34156110f4576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608261188f5761188a82611909565b610d48565b81511580156118b3575073ffffffffffffffffffffffffffffffffffffffff84163b155b15611902576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8516600482015260240161126e565b5092915050565b8051156119195780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6020828403121561195b575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610d48575f5ffd5b5f6020828403121561199a575f5ffd5b5035919050565b604080825283519082018190525f9060208501906060840190835b81811015611a2857835167ffffffffffffffff815116845267ffffffffffffffff602082015116602085015267ffffffffffffffff604082015116604085015267ffffffffffffffff6060820151166060850152506080830192506020840193506001810190506119bc565b505084151560208501529150610d489050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611a5e575f5ffd5b919050565b5f5f60408385031215611a74575f5ffd5b82359150611a8460208401611a3b565b90509250929050565b803567ffffffffffffffff81168114611a5e575f5ffd5b5f5f5f5f60808587031215611ab7575f5ffd5b611ac085611a8d565b9350611ace60208601611a8d565b9250611adc60408601611a8d565b9150611aea60608601611a8d565b905092959194509250565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611b33575f5ffd5b611b3c83611a3b565b9150602083013567ffffffffffffffff811115611b57575f5ffd5b8301601f81018513611b67575f5ffd5b803567ffffffffffffffff811115611b8157611b81611af5565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff82111715611bed57611bed611af5565b604052818152828201602001871015611c04575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f60208284031215611c86575f5ffd5b610d4882611a3b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8082018082111561050357610503611c8f565b8181038181111561050357610503611c8f565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f60208284031215611d1f575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fea2646970667358221220d3aa48445f0d099a9bca4f5102cc9670e7586b2ad3a3367766854a893fcf5f6164736f6c634300081c0033",
}

// RatesManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use RatesManagerMetaData.ABI instead.
var RatesManagerABI = RatesManagerMetaData.ABI

// RatesManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RatesManagerMetaData.Bin instead.
var RatesManagerBin = RatesManagerMetaData.Bin

// DeployRatesManager deploys a new Ethereum contract, binding an instance of RatesManager to it.
func DeployRatesManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RatesManager, error) {
	parsed, err := RatesManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RatesManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RatesManager{RatesManagerCaller: RatesManagerCaller{contract: contract}, RatesManagerTransactor: RatesManagerTransactor{contract: contract}, RatesManagerFilterer: RatesManagerFilterer{contract: contract}}, nil
}

// RatesManager is an auto generated Go binding around an Ethereum contract.
type RatesManager struct {
	RatesManagerCaller     // Read-only binding to the contract
	RatesManagerTransactor // Write-only binding to the contract
	RatesManagerFilterer   // Log filterer for contract events
}

// RatesManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RatesManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RatesManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RatesManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RatesManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RatesManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RatesManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RatesManagerSession struct {
	Contract     *RatesManager     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RatesManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RatesManagerCallerSession struct {
	Contract *RatesManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RatesManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RatesManagerTransactorSession struct {
	Contract     *RatesManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RatesManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RatesManagerRaw struct {
	Contract *RatesManager // Generic contract binding to access the raw methods on
}

// RatesManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RatesManagerCallerRaw struct {
	Contract *RatesManagerCaller // Generic read-only contract binding to access the raw methods on
}

// RatesManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RatesManagerTransactorRaw struct {
	Contract *RatesManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRatesManager creates a new instance of RatesManager, bound to a specific deployed contract.
func NewRatesManager(address common.Address, backend bind.ContractBackend) (*RatesManager, error) {
	contract, err := bindRatesManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RatesManager{RatesManagerCaller: RatesManagerCaller{contract: contract}, RatesManagerTransactor: RatesManagerTransactor{contract: contract}, RatesManagerFilterer: RatesManagerFilterer{contract: contract}}, nil
}

// NewRatesManagerCaller creates a new read-only instance of RatesManager, bound to a specific deployed contract.
func NewRatesManagerCaller(address common.Address, caller bind.ContractCaller) (*RatesManagerCaller, error) {
	contract, err := bindRatesManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RatesManagerCaller{contract: contract}, nil
}

// NewRatesManagerTransactor creates a new write-only instance of RatesManager, bound to a specific deployed contract.
func NewRatesManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RatesManagerTransactor, error) {
	contract, err := bindRatesManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RatesManagerTransactor{contract: contract}, nil
}

// NewRatesManagerFilterer creates a new log filterer instance of RatesManager, bound to a specific deployed contract.
func NewRatesManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RatesManagerFilterer, error) {
	contract, err := bindRatesManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RatesManagerFilterer{contract: contract}, nil
}

// bindRatesManager binds a generic wrapper to an already deployed contract.
func bindRatesManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RatesManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RatesManager *RatesManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RatesManager.Contract.RatesManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RatesManager *RatesManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.Contract.RatesManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RatesManager *RatesManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RatesManager.Contract.RatesManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RatesManager *RatesManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RatesManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RatesManager *RatesManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RatesManager *RatesManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RatesManager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RatesManager.Contract.DEFAULTADMINROLE(&_RatesManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RatesManager.Contract.DEFAULTADMINROLE(&_RatesManager.CallOpts)
}

// PAGESIZE is a free data retrieval call binding the contract method 0x31b28664.
//
// Solidity: function PAGE_SIZE() view returns(uint256)
func (_RatesManager *RatesManagerCaller) PAGESIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "PAGE_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PAGESIZE is a free data retrieval call binding the contract method 0x31b28664.
//
// Solidity: function PAGE_SIZE() view returns(uint256)
func (_RatesManager *RatesManagerSession) PAGESIZE() (*big.Int, error) {
	return _RatesManager.Contract.PAGESIZE(&_RatesManager.CallOpts)
}

// PAGESIZE is a free data retrieval call binding the contract method 0x31b28664.
//
// Solidity: function PAGE_SIZE() view returns(uint256)
func (_RatesManager *RatesManagerCallerSession) PAGESIZE() (*big.Int, error) {
	return _RatesManager.Contract.PAGESIZE(&_RatesManager.CallOpts)
}

// RATESMANAGERROLE is a free data retrieval call binding the contract method 0xe3396735.
//
// Solidity: function RATES_MANAGER_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerCaller) RATESMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "RATES_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RATESMANAGERROLE is a free data retrieval call binding the contract method 0xe3396735.
//
// Solidity: function RATES_MANAGER_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerSession) RATESMANAGERROLE() ([32]byte, error) {
	return _RatesManager.Contract.RATESMANAGERROLE(&_RatesManager.CallOpts)
}

// RATESMANAGERROLE is a free data retrieval call binding the contract method 0xe3396735.
//
// Solidity: function RATES_MANAGER_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerCallerSession) RATESMANAGERROLE() ([32]byte, error) {
	return _RatesManager.Contract.RATESMANAGERROLE(&_RatesManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_RatesManager *RatesManagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_RatesManager *RatesManagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _RatesManager.Contract.UPGRADEINTERFACEVERSION(&_RatesManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_RatesManager *RatesManagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _RatesManager.Contract.UPGRADEINTERFACEVERSION(&_RatesManager.CallOpts)
}

// GetRates is a free data retrieval call binding the contract method 0x081802b1.
//
// Solidity: function getRates(uint256 fromIndex) view returns((uint64,uint64,uint64,uint64)[] rates, bool hasMore)
func (_RatesManager *RatesManagerCaller) GetRates(opts *bind.CallOpts, fromIndex *big.Int) (struct {
	Rates   []RatesManagerRates
	HasMore bool
}, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "getRates", fromIndex)

	outstruct := new(struct {
		Rates   []RatesManagerRates
		HasMore bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Rates = *abi.ConvertType(out[0], new([]RatesManagerRates)).(*[]RatesManagerRates)
	outstruct.HasMore = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// GetRates is a free data retrieval call binding the contract method 0x081802b1.
//
// Solidity: function getRates(uint256 fromIndex) view returns((uint64,uint64,uint64,uint64)[] rates, bool hasMore)
func (_RatesManager *RatesManagerSession) GetRates(fromIndex *big.Int) (struct {
	Rates   []RatesManagerRates
	HasMore bool
}, error) {
	return _RatesManager.Contract.GetRates(&_RatesManager.CallOpts, fromIndex)
}

// GetRates is a free data retrieval call binding the contract method 0x081802b1.
//
// Solidity: function getRates(uint256 fromIndex) view returns((uint64,uint64,uint64,uint64)[] rates, bool hasMore)
func (_RatesManager *RatesManagerCallerSession) GetRates(fromIndex *big.Int) (struct {
	Rates   []RatesManagerRates
	HasMore bool
}, error) {
	return _RatesManager.Contract.GetRates(&_RatesManager.CallOpts, fromIndex)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count)
func (_RatesManager *RatesManagerCaller) GetRatesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "getRatesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count)
func (_RatesManager *RatesManagerSession) GetRatesCount() (*big.Int, error) {
	return _RatesManager.Contract.GetRatesCount(&_RatesManager.CallOpts)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count)
func (_RatesManager *RatesManagerCallerSession) GetRatesCount() (*big.Int, error) {
	return _RatesManager.Contract.GetRatesCount(&_RatesManager.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RatesManager *RatesManagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RatesManager *RatesManagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RatesManager.Contract.GetRoleAdmin(&_RatesManager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RatesManager *RatesManagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RatesManager.Contract.GetRoleAdmin(&_RatesManager.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RatesManager *RatesManagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RatesManager *RatesManagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RatesManager.Contract.HasRole(&_RatesManager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RatesManager *RatesManagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RatesManager.Contract.HasRole(&_RatesManager.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RatesManager *RatesManagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RatesManager *RatesManagerSession) Paused() (bool, error) {
	return _RatesManager.Contract.Paused(&_RatesManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RatesManager *RatesManagerCallerSession) Paused() (bool, error) {
	return _RatesManager.Contract.Paused(&_RatesManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RatesManager *RatesManagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RatesManager *RatesManagerSession) ProxiableUUID() ([32]byte, error) {
	return _RatesManager.Contract.ProxiableUUID(&_RatesManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RatesManager *RatesManagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _RatesManager.Contract.ProxiableUUID(&_RatesManager.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RatesManager *RatesManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RatesManager *RatesManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RatesManager.Contract.SupportsInterface(&_RatesManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RatesManager *RatesManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RatesManager.Contract.SupportsInterface(&_RatesManager.CallOpts, interfaceId)
}

// AddRates is a paid mutator transaction binding the contract method 0x44412107.
//
// Solidity: function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime) returns()
func (_RatesManager *RatesManagerTransactor) AddRates(opts *bind.TransactOpts, messageFee uint64, storageFee uint64, congestionFee uint64, startTime uint64) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "addRates", messageFee, storageFee, congestionFee, startTime)
}

// AddRates is a paid mutator transaction binding the contract method 0x44412107.
//
// Solidity: function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime) returns()
func (_RatesManager *RatesManagerSession) AddRates(messageFee uint64, storageFee uint64, congestionFee uint64, startTime uint64) (*types.Transaction, error) {
	return _RatesManager.Contract.AddRates(&_RatesManager.TransactOpts, messageFee, storageFee, congestionFee, startTime)
}

// AddRates is a paid mutator transaction binding the contract method 0x44412107.
//
// Solidity: function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime) returns()
func (_RatesManager *RatesManagerTransactorSession) AddRates(messageFee uint64, storageFee uint64, congestionFee uint64, startTime uint64) (*types.Transaction, error) {
	return _RatesManager.Contract.AddRates(&_RatesManager.TransactOpts, messageFee, storageFee, congestionFee, startTime)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.GrantRole(&_RatesManager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.GrantRole(&_RatesManager.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_RatesManager *RatesManagerTransactor) Initialize(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "initialize", admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_RatesManager *RatesManagerSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.Initialize(&_RatesManager.TransactOpts, admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_RatesManager *RatesManagerTransactorSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.Initialize(&_RatesManager.TransactOpts, admin)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RatesManager *RatesManagerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RatesManager *RatesManagerSession) Pause() (*types.Transaction, error) {
	return _RatesManager.Contract.Pause(&_RatesManager.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RatesManager *RatesManagerTransactorSession) Pause() (*types.Transaction, error) {
	return _RatesManager.Contract.Pause(&_RatesManager.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RatesManager *RatesManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RatesManager *RatesManagerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.RenounceRole(&_RatesManager.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RatesManager *RatesManagerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.RenounceRole(&_RatesManager.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.RevokeRole(&_RatesManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.RevokeRole(&_RatesManager.TransactOpts, role, account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RatesManager *RatesManagerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RatesManager *RatesManagerSession) Unpause() (*types.Transaction, error) {
	return _RatesManager.Contract.Unpause(&_RatesManager.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RatesManager *RatesManagerTransactorSession) Unpause() (*types.Transaction, error) {
	return _RatesManager.Contract.Unpause(&_RatesManager.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RatesManager *RatesManagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RatesManager *RatesManagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RatesManager.Contract.UpgradeToAndCall(&_RatesManager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RatesManager *RatesManagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RatesManager.Contract.UpgradeToAndCall(&_RatesManager.TransactOpts, newImplementation, data)
}

// RatesManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RatesManager contract.
type RatesManagerInitializedIterator struct {
	Event *RatesManagerInitialized // Event containing the contract specifics and raw log

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
func (it *RatesManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerInitialized)
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
		it.Event = new(RatesManagerInitialized)
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
func (it *RatesManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerInitialized represents a Initialized event raised by the RatesManager contract.
type RatesManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RatesManager *RatesManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*RatesManagerInitializedIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RatesManagerInitializedIterator{contract: _RatesManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RatesManager *RatesManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RatesManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerInitialized)
				if err := _RatesManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseInitialized(log types.Log) (*RatesManagerInitialized, error) {
	event := new(RatesManagerInitialized)
	if err := _RatesManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RatesManager contract.
type RatesManagerPausedIterator struct {
	Event *RatesManagerPaused // Event containing the contract specifics and raw log

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
func (it *RatesManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerPaused)
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
		it.Event = new(RatesManagerPaused)
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
func (it *RatesManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerPaused represents a Paused event raised by the RatesManager contract.
type RatesManagerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RatesManager *RatesManagerFilterer) FilterPaused(opts *bind.FilterOpts) (*RatesManagerPausedIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RatesManagerPausedIterator{contract: _RatesManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RatesManager *RatesManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RatesManagerPaused) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerPaused)
				if err := _RatesManager.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParsePaused(log types.Log) (*RatesManagerPaused, error) {
	event := new(RatesManagerPaused)
	if err := _RatesManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerRatesAddedIterator is returned from FilterRatesAdded and is used to iterate over the raw logs and unpacked data for RatesAdded events raised by the RatesManager contract.
type RatesManagerRatesAddedIterator struct {
	Event *RatesManagerRatesAdded // Event containing the contract specifics and raw log

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
func (it *RatesManagerRatesAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerRatesAdded)
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
		it.Event = new(RatesManagerRatesAdded)
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
func (it *RatesManagerRatesAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerRatesAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerRatesAdded represents a RatesAdded event raised by the RatesManager contract.
type RatesManagerRatesAdded struct {
	MessageFee    uint64
	StorageFee    uint64
	CongestionFee uint64
	StartTime     uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRatesAdded is a free log retrieval operation binding the contract event 0x3bff7b1e021b47f5dfd21d1d3fe2daaf3b9cbaca733c6e15b3a0da546657f19a.
//
// Solidity: event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime)
func (_RatesManager *RatesManagerFilterer) FilterRatesAdded(opts *bind.FilterOpts) (*RatesManagerRatesAddedIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "RatesAdded")
	if err != nil {
		return nil, err
	}
	return &RatesManagerRatesAddedIterator{contract: _RatesManager.contract, event: "RatesAdded", logs: logs, sub: sub}, nil
}

// WatchRatesAdded is a free log subscription operation binding the contract event 0x3bff7b1e021b47f5dfd21d1d3fe2daaf3b9cbaca733c6e15b3a0da546657f19a.
//
// Solidity: event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime)
func (_RatesManager *RatesManagerFilterer) WatchRatesAdded(opts *bind.WatchOpts, sink chan<- *RatesManagerRatesAdded) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "RatesAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerRatesAdded)
				if err := _RatesManager.contract.UnpackLog(event, "RatesAdded", log); err != nil {
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

// ParseRatesAdded is a log parse operation binding the contract event 0x3bff7b1e021b47f5dfd21d1d3fe2daaf3b9cbaca733c6e15b3a0da546657f19a.
//
// Solidity: event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime)
func (_RatesManager *RatesManagerFilterer) ParseRatesAdded(log types.Log) (*RatesManagerRatesAdded, error) {
	event := new(RatesManagerRatesAdded)
	if err := _RatesManager.contract.UnpackLog(event, "RatesAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RatesManager contract.
type RatesManagerRoleAdminChangedIterator struct {
	Event *RatesManagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RatesManagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerRoleAdminChanged)
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
		it.Event = new(RatesManagerRoleAdminChanged)
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
func (it *RatesManagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerRoleAdminChanged represents a RoleAdminChanged event raised by the RatesManager contract.
type RatesManagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RatesManager *RatesManagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RatesManagerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RatesManagerRoleAdminChangedIterator{contract: _RatesManager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RatesManager *RatesManagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RatesManagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerRoleAdminChanged)
				if err := _RatesManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseRoleAdminChanged(log types.Log) (*RatesManagerRoleAdminChanged, error) {
	event := new(RatesManagerRoleAdminChanged)
	if err := _RatesManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RatesManager contract.
type RatesManagerRoleGrantedIterator struct {
	Event *RatesManagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *RatesManagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerRoleGranted)
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
		it.Event = new(RatesManagerRoleGranted)
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
func (it *RatesManagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerRoleGranted represents a RoleGranted event raised by the RatesManager contract.
type RatesManagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RatesManager *RatesManagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RatesManagerRoleGrantedIterator, error) {

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

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RatesManagerRoleGrantedIterator{contract: _RatesManager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RatesManager *RatesManagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RatesManagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerRoleGranted)
				if err := _RatesManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseRoleGranted(log types.Log) (*RatesManagerRoleGranted, error) {
	event := new(RatesManagerRoleGranted)
	if err := _RatesManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RatesManager contract.
type RatesManagerRoleRevokedIterator struct {
	Event *RatesManagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RatesManagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerRoleRevoked)
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
		it.Event = new(RatesManagerRoleRevoked)
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
func (it *RatesManagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerRoleRevoked represents a RoleRevoked event raised by the RatesManager contract.
type RatesManagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RatesManager *RatesManagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RatesManagerRoleRevokedIterator, error) {

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

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RatesManagerRoleRevokedIterator{contract: _RatesManager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RatesManager *RatesManagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RatesManagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerRoleRevoked)
				if err := _RatesManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseRoleRevoked(log types.Log) (*RatesManagerRoleRevoked, error) {
	event := new(RatesManagerRoleRevoked)
	if err := _RatesManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RatesManager contract.
type RatesManagerUnpausedIterator struct {
	Event *RatesManagerUnpaused // Event containing the contract specifics and raw log

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
func (it *RatesManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerUnpaused)
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
		it.Event = new(RatesManagerUnpaused)
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
func (it *RatesManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerUnpaused represents a Unpaused event raised by the RatesManager contract.
type RatesManagerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RatesManager *RatesManagerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RatesManagerUnpausedIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RatesManagerUnpausedIterator{contract: _RatesManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RatesManager *RatesManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RatesManagerUnpaused) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerUnpaused)
				if err := _RatesManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseUnpaused(log types.Log) (*RatesManagerUnpaused, error) {
	event := new(RatesManagerUnpaused)
	if err := _RatesManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerUpgradeAuthorizedIterator is returned from FilterUpgradeAuthorized and is used to iterate over the raw logs and unpacked data for UpgradeAuthorized events raised by the RatesManager contract.
type RatesManagerUpgradeAuthorizedIterator struct {
	Event *RatesManagerUpgradeAuthorized // Event containing the contract specifics and raw log

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
func (it *RatesManagerUpgradeAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerUpgradeAuthorized)
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
		it.Event = new(RatesManagerUpgradeAuthorized)
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
func (it *RatesManagerUpgradeAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerUpgradeAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerUpgradeAuthorized represents a UpgradeAuthorized event raised by the RatesManager contract.
type RatesManagerUpgradeAuthorized struct {
	Upgrader          common.Address
	NewImplementation common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeAuthorized is a free log retrieval operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_RatesManager *RatesManagerFilterer) FilterUpgradeAuthorized(opts *bind.FilterOpts) (*RatesManagerUpgradeAuthorizedIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return &RatesManagerUpgradeAuthorizedIterator{contract: _RatesManager.contract, event: "UpgradeAuthorized", logs: logs, sub: sub}, nil
}

// WatchUpgradeAuthorized is a free log subscription operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_RatesManager *RatesManagerFilterer) WatchUpgradeAuthorized(opts *bind.WatchOpts, sink chan<- *RatesManagerUpgradeAuthorized) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerUpgradeAuthorized)
				if err := _RatesManager.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseUpgradeAuthorized(log types.Log) (*RatesManagerUpgradeAuthorized, error) {
	event := new(RatesManagerUpgradeAuthorized)
	if err := _RatesManager.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the RatesManager contract.
type RatesManagerUpgradedIterator struct {
	Event *RatesManagerUpgraded // Event containing the contract specifics and raw log

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
func (it *RatesManagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerUpgraded)
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
		it.Event = new(RatesManagerUpgraded)
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
func (it *RatesManagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerUpgraded represents a Upgraded event raised by the RatesManager contract.
type RatesManagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RatesManager *RatesManagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RatesManagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RatesManagerUpgradedIterator{contract: _RatesManager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RatesManager *RatesManagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RatesManagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerUpgraded)
				if err := _RatesManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseUpgraded(log types.Log) (*RatesManagerUpgraded, error) {
	event := new(RatesManagerUpgraded)
	if err := _RatesManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
