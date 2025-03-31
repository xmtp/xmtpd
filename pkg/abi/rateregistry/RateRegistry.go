// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rateregistry

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

// RateRegistryRates is an auto generated low-level Go binding around an user-defined struct.
type RateRegistryRates struct {
	MessageFee          uint64
	StorageFee          uint64
	CongestionFee       uint64
	TargetRatePerMinute uint64
	StartTime           uint64
}

// RateRegistryMetaData contains all meta data concerning the RateRegistry contract.
var RateRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAGE_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RATES_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addRates\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRates\",\"inputs\":[{\"name\":\"fromIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rates\",\"type\":\"tuple[]\",\"internalType\":\"structRateRegistry.Rates[]\",\"components\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"hasMore\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRatesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RatesAdded\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"targetRatePerMinute\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeAuthorized\",\"inputs\":[{\"name\":\"upgrader\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newImplementation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FromIndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidStartTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAdminAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroImplementationAddress\",\"inputs\":[]}]",
	Bin: "0x60a0604052306080523480156012575f5ffd5b50608051611e0e6100395f395f8181611080015281816110a9015261136d0152611e0e5ff3fe60806040526004361061013d575f3560e01c80634f1ef286116100bb578063a217fddf11610071578063c4d66de811610057578063c4d66de814610400578063d547741f1461041f578063e33967351461043e575f5ffd5b8063a217fddf14610398578063ad3cb1cc146103ab575f5ffd5b80635c975abb116100a15780635c975abb146102de5780638456cb591461031457806391d1485414610328575f5ffd5b80634f1ef286146102b757806352d1902d146102ca575f5ffd5b80632f2ff15d1161011057806336568abe116100f657806336568abe146102655780633f4ba83a146102845780634cc8c51514610298575f5ffd5b80632f2ff15d1461023057806331b2866414610251575f5ffd5b806301ffc9a714610141578063081802b114610175578063248a9ca3146101a25780632da72291146101fd575b5f5ffd5b34801561014c575f5ffd5b5061016061015b3660046119c3565b610471565b60405190151581526020015b60405180910390f35b348015610180575f5ffd5b5061019461018f366004611a02565b610509565b60405161016c929190611a19565b3480156101ad575f5ffd5b506101ef6101bc366004611a02565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b60405190815260200161016c565b348015610208575f5ffd5b507f6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d000546101ef565b34801561023b575f5ffd5b5061024f61024a366004611aef565b61078d565b005b34801561025c575f5ffd5b506101ef603281565b348015610270575f5ffd5b5061024f61027f366004611aef565b6107d6565b34801561028f575f5ffd5b5061024f610834565b3480156102a3575f5ffd5b5061024f6102b2366004611b30565b610849565b61024f6102c5366004611bbe565b610ad4565b3480156102d5575f5ffd5b506101ef610af3565b3480156102e9575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610160565b34801561031f575f5ffd5b5061024f610b21565b348015610333575f5ffd5b50610160610342366004611aef565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b3480156103a3575f5ffd5b506101ef5f81565b3480156103b6575f5ffd5b506103f36040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b60405161016c9190611cbf565b34801561040b575f5ffd5b5061024f61041a366004611d12565b610b33565b34801561042a575f5ffd5b5061024f610439366004611aef565b610d6d565b348015610449575f5ffd5b506101ef7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e81565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061050357507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b7f6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d00080546060915f9115801561053c575083155b156105b957604080515f80825260208201909252906105ae565b6040805160a0810182525f808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816105565790505b50945f945092505050565b805484106105f3576040517fea61fe7000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610609610602603287611d58565b8354610db0565b90506106158582611d6b565b67ffffffffffffffff81111561062d5761062d611b91565b6040519080825280602002602001820160405280156106a357816020015b6040805160a0810182525f808252602080830182905292820181905260608201819052608082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90920191018161064b5790505b5093505f5b845181101561078157826106bc8288611d58565b815481106106cc576106cc611d7e565b5f9182526020918290206040805160a0810182526002909302909101805467ffffffffffffffff80821685526801000000000000000082048116958501959095527001000000000000000000000000000000008104851692840192909252780100000000000000000000000000000000000000000000000090910483166060830152600101549091166080820152855186908390811061076e5761076e611d7e565b60209081029190910101526001016106a8565b50905492949211925050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546107c681610dc7565b6107d08383610dd1565b50505050565b73ffffffffffffffffffffffffffffffffffffffff81163314610825576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61082f8282610eef565b505050565b5f61083e81610dc7565b610846610fcb565b50565b7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e61087381610dc7565b7f6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d0008054158015906108e45750805481906108af90600190611d6b565b815481106108bf576108bf611d7e565b5f91825260209091206001600290920201015467ffffffffffffffff90811690841611155b1561091b576040517fb290253c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160a08101825267ffffffffffffffff808a16825288811660208084019182528983168486019081528984166060860190815289851660808701908152885460018181018b555f8b815295909520975160029091029097018054955193519251871678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff93881670010000000000000000000000000000000002939093166fffffffffffffffffffffffffffffffff94881668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000000000000000000090971698881698909817959095179290921695909517949094178255925192018054929091167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000092909216919091179055517f6892d26f6f7f805fab93d55f64f67960ab38bae8d4134cbdd02e023d21885cbc90610ac3908990899089908990899067ffffffffffffffff95861681529385166020850152918416604084015283166060830152909116608082015260a00190565b60405180910390a150505050505050565b610adc611068565b610ae58261116e565b610aef8282611217565b5050565b5f610afc611355565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b5f610b2b81610dc7565b6108466113c4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f81158015610b7d5750825b90505f8267ffffffffffffffff166001148015610b995750303b155b905081158015610ba7575080155b15610bde576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001660011785558315610c3f5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b73ffffffffffffffffffffffffffffffffffffffff8616610c8c576040517f3ef39b8100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c9461143d565b610c9c61143d565b610ca4611445565b610cce7f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e5f611455565b610cd85f87610dd1565b50610d037f64b4740f54156feb06b7a9f424e5bce966a52344cf27635887cf63c0ebf2a61e87610dd1565b508315610d655784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610da681610dc7565b6107d08383610eef565b5f818310610dbe5781610dc0565b825b9392505050565b61084681336114f6565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff16610ee6575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610e823390565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050610503565b5f915050610503565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020818152604080842073ffffffffffffffffffffffffffffffffffffffff8616855290915282205460ff1615610ee6575f8481526020828152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a46001915050610503565b610fd361159c565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a150565b3073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016148061113557507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661111c7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b1561116c576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f61117881610dc7565b73ffffffffffffffffffffffffffffffffffffffff82166111c5576040517fd02c623d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805133815273ffffffffffffffffffffffffffffffffffffffff841660208201527fd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa910160405180910390a15050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561129c575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261129991810190611dab565b60015b6112ef576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc811461134b576040517faa1d49a4000000000000000000000000000000000000000000000000000000008152600481018290526024016112e6565b61082f83836115f7565b3073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461116c576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113cc611659565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2583361103d565b61116c6116b5565b61144d6116b5565b61116c61171c565b7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268005f6114ae845f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b5f85815260208490526040808220600101869055519192508491839187917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a450505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610aef576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602481018390526044016112e6565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1661116c576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6116008261176d565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156116515761082f828261183b565b610aef6118ba565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff161561116c576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1661116c576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117246116b5565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b8073ffffffffffffffffffffffffffffffffffffffff163b5f036117d5576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016112e6565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff16846040516118649190611dc2565b5f60405180830381855af49150503d805f811461189c576040519150601f19603f3d011682016040523d82523d5f602084013e6118a1565b606091505b50915091506118b18583836118f2565b95945050505050565b341561116c576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060826119075761190282611981565b610dc0565b815115801561192b575073ffffffffffffffffffffffffffffffffffffffff84163b155b1561197a576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024016112e6565b5092915050565b8051156119915780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f602082840312156119d3575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610dc0575f5ffd5b5f60208284031215611a12575f5ffd5b5035919050565b604080825283519082018190525f9060208501906060840190835b81811015611ab457835167ffffffffffffffff815116845267ffffffffffffffff602082015116602085015267ffffffffffffffff604082015116604085015267ffffffffffffffff606082015116606085015267ffffffffffffffff60808201511660808501525060a083019250602084019350600181019050611a34565b505084151560208501529150610dc09050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611aea575f5ffd5b919050565b5f5f60408385031215611b00575f5ffd5b82359150611b1060208401611ac7565b90509250929050565b803567ffffffffffffffff81168114611aea575f5ffd5b5f5f5f5f5f60a08688031215611b44575f5ffd5b611b4d86611b19565b9450611b5b60208701611b19565b9350611b6960408701611b19565b9250611b7760608701611b19565b9150611b8560808701611b19565b90509295509295909350565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611bcf575f5ffd5b611bd883611ac7565b9150602083013567ffffffffffffffff811115611bf3575f5ffd5b8301601f81018513611c03575f5ffd5b803567ffffffffffffffff811115611c1d57611c1d611b91565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff82111715611c8957611c89611b91565b604052818152828201602001871015611ca0575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f60208284031215611d22575f5ffd5b610dc082611ac7565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8082018082111561050357610503611d2b565b8181038181111561050357610503611d2b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f60208284031215611dbb575f5ffd5b5051919050565b5f82518060208501845e5f92019182525091905056fea2646970667358221220c5e3a0762af955cd1f577e87492bec171b8113e9c9cf7f5415adca799508f9f964736f6c634300081c0033",
}

// RateRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RateRegistryMetaData.ABI instead.
var RateRegistryABI = RateRegistryMetaData.ABI

// RateRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RateRegistryMetaData.Bin instead.
var RateRegistryBin = RateRegistryMetaData.Bin

// DeployRateRegistry deploys a new Ethereum contract, binding an instance of RateRegistry to it.
func DeployRateRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RateRegistry, error) {
	parsed, err := RateRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RateRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RateRegistry{RateRegistryCaller: RateRegistryCaller{contract: contract}, RateRegistryTransactor: RateRegistryTransactor{contract: contract}, RateRegistryFilterer: RateRegistryFilterer{contract: contract}}, nil
}

// RateRegistry is an auto generated Go binding around an Ethereum contract.
type RateRegistry struct {
	RateRegistryCaller     // Read-only binding to the contract
	RateRegistryTransactor // Write-only binding to the contract
	RateRegistryFilterer   // Log filterer for contract events
}

// RateRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RateRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RateRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RateRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RateRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RateRegistrySession struct {
	Contract     *RateRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RateRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RateRegistryCallerSession struct {
	Contract *RateRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RateRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RateRegistryTransactorSession struct {
	Contract     *RateRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RateRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RateRegistryRaw struct {
	Contract *RateRegistry // Generic contract binding to access the raw methods on
}

// RateRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RateRegistryCallerRaw struct {
	Contract *RateRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RateRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RateRegistryTransactorRaw struct {
	Contract *RateRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRateRegistry creates a new instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistry(address common.Address, backend bind.ContractBackend) (*RateRegistry, error) {
	contract, err := bindRateRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RateRegistry{RateRegistryCaller: RateRegistryCaller{contract: contract}, RateRegistryTransactor: RateRegistryTransactor{contract: contract}, RateRegistryFilterer: RateRegistryFilterer{contract: contract}}, nil
}

// NewRateRegistryCaller creates a new read-only instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryCaller(address common.Address, caller bind.ContractCaller) (*RateRegistryCaller, error) {
	contract, err := bindRateRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RateRegistryCaller{contract: contract}, nil
}

// NewRateRegistryTransactor creates a new write-only instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RateRegistryTransactor, error) {
	contract, err := bindRateRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RateRegistryTransactor{contract: contract}, nil
}

// NewRateRegistryFilterer creates a new log filterer instance of RateRegistry, bound to a specific deployed contract.
func NewRateRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RateRegistryFilterer, error) {
	contract, err := bindRateRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RateRegistryFilterer{contract: contract}, nil
}

// bindRateRegistry binds a generic wrapper to an already deployed contract.
func bindRateRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RateRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateRegistry *RateRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RateRegistry.Contract.RateRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateRegistry *RateRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.Contract.RateRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateRegistry *RateRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateRegistry.Contract.RateRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RateRegistry *RateRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RateRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RateRegistry *RateRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RateRegistry *RateRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RateRegistry.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RateRegistry *RateRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RateRegistry *RateRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RateRegistry.Contract.DEFAULTADMINROLE(&_RateRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RateRegistry *RateRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RateRegistry.Contract.DEFAULTADMINROLE(&_RateRegistry.CallOpts)
}

// PAGESIZE is a free data retrieval call binding the contract method 0x31b28664.
//
// Solidity: function PAGE_SIZE() view returns(uint256)
func (_RateRegistry *RateRegistryCaller) PAGESIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "PAGE_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PAGESIZE is a free data retrieval call binding the contract method 0x31b28664.
//
// Solidity: function PAGE_SIZE() view returns(uint256)
func (_RateRegistry *RateRegistrySession) PAGESIZE() (*big.Int, error) {
	return _RateRegistry.Contract.PAGESIZE(&_RateRegistry.CallOpts)
}

// PAGESIZE is a free data retrieval call binding the contract method 0x31b28664.
//
// Solidity: function PAGE_SIZE() view returns(uint256)
func (_RateRegistry *RateRegistryCallerSession) PAGESIZE() (*big.Int, error) {
	return _RateRegistry.Contract.PAGESIZE(&_RateRegistry.CallOpts)
}

// RATESMANAGERROLE is a free data retrieval call binding the contract method 0xe3396735.
//
// Solidity: function RATES_MANAGER_ROLE() view returns(bytes32)
func (_RateRegistry *RateRegistryCaller) RATESMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "RATES_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RATESMANAGERROLE is a free data retrieval call binding the contract method 0xe3396735.
//
// Solidity: function RATES_MANAGER_ROLE() view returns(bytes32)
func (_RateRegistry *RateRegistrySession) RATESMANAGERROLE() ([32]byte, error) {
	return _RateRegistry.Contract.RATESMANAGERROLE(&_RateRegistry.CallOpts)
}

// RATESMANAGERROLE is a free data retrieval call binding the contract method 0xe3396735.
//
// Solidity: function RATES_MANAGER_ROLE() view returns(bytes32)
func (_RateRegistry *RateRegistryCallerSession) RATESMANAGERROLE() ([32]byte, error) {
	return _RateRegistry.Contract.RATESMANAGERROLE(&_RateRegistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_RateRegistry *RateRegistryCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_RateRegistry *RateRegistrySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _RateRegistry.Contract.UPGRADEINTERFACEVERSION(&_RateRegistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_RateRegistry *RateRegistryCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _RateRegistry.Contract.UPGRADEINTERFACEVERSION(&_RateRegistry.CallOpts)
}

// GetRates is a free data retrieval call binding the contract method 0x081802b1.
//
// Solidity: function getRates(uint256 fromIndex) view returns((uint64,uint64,uint64,uint64,uint64)[] rates, bool hasMore)
func (_RateRegistry *RateRegistryCaller) GetRates(opts *bind.CallOpts, fromIndex *big.Int) (struct {
	Rates   []RateRegistryRates
	HasMore bool
}, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRates", fromIndex)

	outstruct := new(struct {
		Rates   []RateRegistryRates
		HasMore bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Rates = *abi.ConvertType(out[0], new([]RateRegistryRates)).(*[]RateRegistryRates)
	outstruct.HasMore = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// GetRates is a free data retrieval call binding the contract method 0x081802b1.
//
// Solidity: function getRates(uint256 fromIndex) view returns((uint64,uint64,uint64,uint64,uint64)[] rates, bool hasMore)
func (_RateRegistry *RateRegistrySession) GetRates(fromIndex *big.Int) (struct {
	Rates   []RateRegistryRates
	HasMore bool
}, error) {
	return _RateRegistry.Contract.GetRates(&_RateRegistry.CallOpts, fromIndex)
}

// GetRates is a free data retrieval call binding the contract method 0x081802b1.
//
// Solidity: function getRates(uint256 fromIndex) view returns((uint64,uint64,uint64,uint64,uint64)[] rates, bool hasMore)
func (_RateRegistry *RateRegistryCallerSession) GetRates(fromIndex *big.Int) (struct {
	Rates   []RateRegistryRates
	HasMore bool
}, error) {
	return _RateRegistry.Contract.GetRates(&_RateRegistry.CallOpts, fromIndex)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count)
func (_RateRegistry *RateRegistryCaller) GetRatesCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRatesCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count)
func (_RateRegistry *RateRegistrySession) GetRatesCount() (*big.Int, error) {
	return _RateRegistry.Contract.GetRatesCount(&_RateRegistry.CallOpts)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256 count)
func (_RateRegistry *RateRegistryCallerSession) GetRatesCount() (*big.Int, error) {
	return _RateRegistry.Contract.GetRatesCount(&_RateRegistry.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RateRegistry *RateRegistryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RateRegistry *RateRegistrySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RateRegistry.Contract.GetRoleAdmin(&_RateRegistry.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RateRegistry *RateRegistryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RateRegistry.Contract.GetRoleAdmin(&_RateRegistry.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RateRegistry *RateRegistryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RateRegistry *RateRegistrySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RateRegistry.Contract.HasRole(&_RateRegistry.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RateRegistry *RateRegistryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RateRegistry.Contract.HasRole(&_RateRegistry.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RateRegistry *RateRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RateRegistry *RateRegistrySession) Paused() (bool, error) {
	return _RateRegistry.Contract.Paused(&_RateRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RateRegistry *RateRegistryCallerSession) Paused() (bool, error) {
	return _RateRegistry.Contract.Paused(&_RateRegistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RateRegistry *RateRegistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RateRegistry *RateRegistrySession) ProxiableUUID() ([32]byte, error) {
	return _RateRegistry.Contract.ProxiableUUID(&_RateRegistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RateRegistry *RateRegistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _RateRegistry.Contract.ProxiableUUID(&_RateRegistry.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RateRegistry *RateRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RateRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RateRegistry *RateRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RateRegistry.Contract.SupportsInterface(&_RateRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RateRegistry *RateRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RateRegistry.Contract.SupportsInterface(&_RateRegistry.CallOpts, interfaceId)
}

// AddRates is a paid mutator transaction binding the contract method 0x4cc8c515.
//
// Solidity: function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime) returns()
func (_RateRegistry *RateRegistryTransactor) AddRates(opts *bind.TransactOpts, messageFee uint64, storageFee uint64, congestionFee uint64, targetRatePerMinute uint64, startTime uint64) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "addRates", messageFee, storageFee, congestionFee, targetRatePerMinute, startTime)
}

// AddRates is a paid mutator transaction binding the contract method 0x4cc8c515.
//
// Solidity: function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime) returns()
func (_RateRegistry *RateRegistrySession) AddRates(messageFee uint64, storageFee uint64, congestionFee uint64, targetRatePerMinute uint64, startTime uint64) (*types.Transaction, error) {
	return _RateRegistry.Contract.AddRates(&_RateRegistry.TransactOpts, messageFee, storageFee, congestionFee, targetRatePerMinute, startTime)
}

// AddRates is a paid mutator transaction binding the contract method 0x4cc8c515.
//
// Solidity: function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime) returns()
func (_RateRegistry *RateRegistryTransactorSession) AddRates(messageFee uint64, storageFee uint64, congestionFee uint64, targetRatePerMinute uint64, startTime uint64) (*types.Transaction, error) {
	return _RateRegistry.Contract.AddRates(&_RateRegistry.TransactOpts, messageFee, storageFee, congestionFee, targetRatePerMinute, startTime)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RateRegistry *RateRegistryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RateRegistry *RateRegistrySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.GrantRole(&_RateRegistry.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RateRegistry *RateRegistryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.GrantRole(&_RateRegistry.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_RateRegistry *RateRegistryTransactor) Initialize(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "initialize", admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_RateRegistry *RateRegistrySession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.Initialize(&_RateRegistry.TransactOpts, admin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address admin) returns()
func (_RateRegistry *RateRegistryTransactorSession) Initialize(admin common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.Initialize(&_RateRegistry.TransactOpts, admin)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RateRegistry *RateRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RateRegistry *RateRegistrySession) Pause() (*types.Transaction, error) {
	return _RateRegistry.Contract.Pause(&_RateRegistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RateRegistry *RateRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _RateRegistry.Contract.Pause(&_RateRegistry.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RateRegistry *RateRegistryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RateRegistry *RateRegistrySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.RenounceRole(&_RateRegistry.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RateRegistry *RateRegistryTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.RenounceRole(&_RateRegistry.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RateRegistry *RateRegistryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RateRegistry *RateRegistrySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.RevokeRole(&_RateRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RateRegistry *RateRegistryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RateRegistry.Contract.RevokeRole(&_RateRegistry.TransactOpts, role, account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RateRegistry *RateRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RateRegistry *RateRegistrySession) Unpause() (*types.Transaction, error) {
	return _RateRegistry.Contract.Unpause(&_RateRegistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RateRegistry *RateRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _RateRegistry.Contract.Unpause(&_RateRegistry.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RateRegistry *RateRegistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RateRegistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RateRegistry *RateRegistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RateRegistry.Contract.UpgradeToAndCall(&_RateRegistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RateRegistry *RateRegistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RateRegistry.Contract.UpgradeToAndCall(&_RateRegistry.TransactOpts, newImplementation, data)
}

// RateRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RateRegistry contract.
type RateRegistryInitializedIterator struct {
	Event *RateRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *RateRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryInitialized)
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
		it.Event = new(RateRegistryInitialized)
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
func (it *RateRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryInitialized represents a Initialized event raised by the RateRegistry contract.
type RateRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*RateRegistryInitializedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RateRegistryInitializedIterator{contract: _RateRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RateRegistry *RateRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RateRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryInitialized)
				if err := _RateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseInitialized(log types.Log) (*RateRegistryInitialized, error) {
	event := new(RateRegistryInitialized)
	if err := _RateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RateRegistry contract.
type RateRegistryPausedIterator struct {
	Event *RateRegistryPaused // Event containing the contract specifics and raw log

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
func (it *RateRegistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryPaused)
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
		it.Event = new(RateRegistryPaused)
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
func (it *RateRegistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryPaused represents a Paused event raised by the RateRegistry contract.
type RateRegistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RateRegistry *RateRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*RateRegistryPausedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RateRegistryPausedIterator{contract: _RateRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RateRegistry *RateRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RateRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryPaused)
				if err := _RateRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParsePaused(log types.Log) (*RateRegistryPaused, error) {
	event := new(RateRegistryPaused)
	if err := _RateRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryRatesAddedIterator is returned from FilterRatesAdded and is used to iterate over the raw logs and unpacked data for RatesAdded events raised by the RateRegistry contract.
type RateRegistryRatesAddedIterator struct {
	Event *RateRegistryRatesAdded // Event containing the contract specifics and raw log

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
func (it *RateRegistryRatesAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryRatesAdded)
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
		it.Event = new(RateRegistryRatesAdded)
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
func (it *RateRegistryRatesAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryRatesAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryRatesAdded represents a RatesAdded event raised by the RateRegistry contract.
type RateRegistryRatesAdded struct {
	MessageFee          uint64
	StorageFee          uint64
	CongestionFee       uint64
	TargetRatePerMinute uint64
	StartTime           uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterRatesAdded is a free log retrieval operation binding the contract event 0x6892d26f6f7f805fab93d55f64f67960ab38bae8d4134cbdd02e023d21885cbc.
//
// Solidity: event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
func (_RateRegistry *RateRegistryFilterer) FilterRatesAdded(opts *bind.FilterOpts) (*RateRegistryRatesAddedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "RatesAdded")
	if err != nil {
		return nil, err
	}
	return &RateRegistryRatesAddedIterator{contract: _RateRegistry.contract, event: "RatesAdded", logs: logs, sub: sub}, nil
}

// WatchRatesAdded is a free log subscription operation binding the contract event 0x6892d26f6f7f805fab93d55f64f67960ab38bae8d4134cbdd02e023d21885cbc.
//
// Solidity: event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
func (_RateRegistry *RateRegistryFilterer) WatchRatesAdded(opts *bind.WatchOpts, sink chan<- *RateRegistryRatesAdded) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "RatesAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryRatesAdded)
				if err := _RateRegistry.contract.UnpackLog(event, "RatesAdded", log); err != nil {
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

// ParseRatesAdded is a log parse operation binding the contract event 0x6892d26f6f7f805fab93d55f64f67960ab38bae8d4134cbdd02e023d21885cbc.
//
// Solidity: event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
func (_RateRegistry *RateRegistryFilterer) ParseRatesAdded(log types.Log) (*RateRegistryRatesAdded, error) {
	event := new(RateRegistryRatesAdded)
	if err := _RateRegistry.contract.UnpackLog(event, "RatesAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RateRegistry contract.
type RateRegistryRoleAdminChangedIterator struct {
	Event *RateRegistryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RateRegistryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryRoleAdminChanged)
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
		it.Event = new(RateRegistryRoleAdminChanged)
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
func (it *RateRegistryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the RateRegistry contract.
type RateRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RateRegistry *RateRegistryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RateRegistryRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryRoleAdminChangedIterator{contract: _RateRegistry.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RateRegistry *RateRegistryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RateRegistryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryRoleAdminChanged)
				if err := _RateRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseRoleAdminChanged(log types.Log) (*RateRegistryRoleAdminChanged, error) {
	event := new(RateRegistryRoleAdminChanged)
	if err := _RateRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RateRegistry contract.
type RateRegistryRoleGrantedIterator struct {
	Event *RateRegistryRoleGranted // Event containing the contract specifics and raw log

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
func (it *RateRegistryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryRoleGranted)
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
		it.Event = new(RateRegistryRoleGranted)
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
func (it *RateRegistryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryRoleGranted represents a RoleGranted event raised by the RateRegistry contract.
type RateRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RateRegistry *RateRegistryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RateRegistryRoleGrantedIterator, error) {

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

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryRoleGrantedIterator{contract: _RateRegistry.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RateRegistry *RateRegistryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RateRegistryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryRoleGranted)
				if err := _RateRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseRoleGranted(log types.Log) (*RateRegistryRoleGranted, error) {
	event := new(RateRegistryRoleGranted)
	if err := _RateRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RateRegistry contract.
type RateRegistryRoleRevokedIterator struct {
	Event *RateRegistryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RateRegistryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryRoleRevoked)
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
		it.Event = new(RateRegistryRoleRevoked)
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
func (it *RateRegistryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryRoleRevoked represents a RoleRevoked event raised by the RateRegistry contract.
type RateRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RateRegistry *RateRegistryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RateRegistryRoleRevokedIterator, error) {

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

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryRoleRevokedIterator{contract: _RateRegistry.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RateRegistry *RateRegistryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RateRegistryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryRoleRevoked)
				if err := _RateRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseRoleRevoked(log types.Log) (*RateRegistryRoleRevoked, error) {
	event := new(RateRegistryRoleRevoked)
	if err := _RateRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RateRegistry contract.
type RateRegistryUnpausedIterator struct {
	Event *RateRegistryUnpaused // Event containing the contract specifics and raw log

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
func (it *RateRegistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryUnpaused)
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
		it.Event = new(RateRegistryUnpaused)
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
func (it *RateRegistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryUnpaused represents a Unpaused event raised by the RateRegistry contract.
type RateRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RateRegistry *RateRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RateRegistryUnpausedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RateRegistryUnpausedIterator{contract: _RateRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RateRegistry *RateRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RateRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryUnpaused)
				if err := _RateRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseUnpaused(log types.Log) (*RateRegistryUnpaused, error) {
	event := new(RateRegistryUnpaused)
	if err := _RateRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryUpgradeAuthorizedIterator is returned from FilterUpgradeAuthorized and is used to iterate over the raw logs and unpacked data for UpgradeAuthorized events raised by the RateRegistry contract.
type RateRegistryUpgradeAuthorizedIterator struct {
	Event *RateRegistryUpgradeAuthorized // Event containing the contract specifics and raw log

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
func (it *RateRegistryUpgradeAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryUpgradeAuthorized)
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
		it.Event = new(RateRegistryUpgradeAuthorized)
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
func (it *RateRegistryUpgradeAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryUpgradeAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryUpgradeAuthorized represents a UpgradeAuthorized event raised by the RateRegistry contract.
type RateRegistryUpgradeAuthorized struct {
	Upgrader          common.Address
	NewImplementation common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpgradeAuthorized is a free log retrieval operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_RateRegistry *RateRegistryFilterer) FilterUpgradeAuthorized(opts *bind.FilterOpts) (*RateRegistryUpgradeAuthorizedIterator, error) {

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return &RateRegistryUpgradeAuthorizedIterator{contract: _RateRegistry.contract, event: "UpgradeAuthorized", logs: logs, sub: sub}, nil
}

// WatchUpgradeAuthorized is a free log subscription operation binding the contract event 0xd30e1d298bf814ea43d22b4ce8298062b08609cd67496483769d836157dd52fa.
//
// Solidity: event UpgradeAuthorized(address upgrader, address newImplementation)
func (_RateRegistry *RateRegistryFilterer) WatchUpgradeAuthorized(opts *bind.WatchOpts, sink chan<- *RateRegistryUpgradeAuthorized) (event.Subscription, error) {

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "UpgradeAuthorized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryUpgradeAuthorized)
				if err := _RateRegistry.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseUpgradeAuthorized(log types.Log) (*RateRegistryUpgradeAuthorized, error) {
	event := new(RateRegistryUpgradeAuthorized)
	if err := _RateRegistry.contract.UnpackLog(event, "UpgradeAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RateRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the RateRegistry contract.
type RateRegistryUpgradedIterator struct {
	Event *RateRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *RateRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RateRegistryUpgraded)
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
		it.Event = new(RateRegistryUpgraded)
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
func (it *RateRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RateRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RateRegistryUpgraded represents a Upgraded event raised by the RateRegistry contract.
type RateRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RateRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RateRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RateRegistryUpgradedIterator{contract: _RateRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RateRegistry *RateRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RateRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RateRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RateRegistryUpgraded)
				if err := _RateRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_RateRegistry *RateRegistryFilterer) ParseUpgraded(log types.Log) (*RateRegistryUpgraded, error) {
	event := new(RateRegistryUpgraded)
	if err := _RateRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
