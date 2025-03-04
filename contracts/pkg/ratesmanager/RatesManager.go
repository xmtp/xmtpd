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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RATES_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRates\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRates\",\"inputs\":[{\"name\":\"fromIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rates\",\"type\":\"tuple[]\",\"internalType\":\"structRatesManager.Rates[]\",\"components\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"hasMore\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRatesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RatesAdded\",\"inputs\":[{\"name\":\"messageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"storageFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"congestionFee\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x608060405234801561000f575f5ffd5b5062015180338061003957604051636116401160e11b81525f600482015260240160405180910390fd5b600180546001600160d01b0316600160d01b65ffffffffffff8516021790556100625f8261009b565b5061007e91505f516020611b385f395f51905f5290505f61010a565b6100955f516020611b385f395f51905f523361009b565b50610227565b5f826100f7575f6100b46002546001600160a01b031690565b6001600160a01b0316146100db57604051631fe1e13d60e11b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0384161790555b6101018383610136565b90505b92915050565b8161012857604051631fe1e13d60e11b815260040160405180910390fd5b61013282826101dd565b5050565b5f828152602081815260408083206001600160a01b038516845290915281205460ff166101d6575f838152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905561018e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610104565b505f610104565b5f82815260208190526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b611904806102345f395ff3fe608060405234801561000f575f5ffd5b5060043610610179575f3560e01c806384ef8ffc116100d2578063a217fddf11610088578063cf6eefb711610063578063cf6eefb714610372578063d547741f146103be578063d602b9fd146103d1575f5ffd5b8063a217fddf1461035b578063cc8463c814610362578063cefc14291461036a575f5ffd5b806391d14854116100b857806391d14854146102ca578063970b6d871461030d578063a1eda53c14610334575f5ffd5b806384ef8ffc146102835780638da5cb5b146102c2575f5ffd5b80632da7229111610132578063444121071161010d578063444121071461024a578063634e93da1461025d578063649a5ec714610270575f5ffd5b80632da722911461021c5780632f2ff15d1461022457806336568abe14610237575f5ffd5b8063081802b111610162578063081802b1146101c15780630aa6220b146101e2578063248a9ca3146101ec575f5ffd5b806301ffc9a71461017d578063022d63fb146101a5575b5f5ffd5b61019061018b3660046115fd565b6103d9565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff909116815260200161019c565b6101d46101cf36600461163c565b610434565b60405161019c929190611653565b6101ea6106c0565b005b61020e6101fa36600461163c565b5f9081526020819052604090206001015490565b60405190815260200161019c565b60035461020e565b6101ea610232366004611715565b6106d5565b6101ea610245366004611715565b61071a565b6101ea610258366004611756565b61081f565b6101ea61026b3660046117a7565b610aba565b6101ea61027e3660046117c0565b610acd565b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161019c565b61029d610ae0565b6101906102d8366004611715565b5f9182526020828152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b61020e7fe1f251c81ea8f25043f9fc7ec35c0a2a10dd591e11346051677e0580f5887a3381565b61033c610b05565b6040805165ffffffffffff93841681529290911660208301520161019c565b61020e5f81565b6101aa610b7f565b6101ea610c1c565b6001546040805173ffffffffffffffffffffffffffffffffffffffff831681527401000000000000000000000000000000000000000090920465ffffffffffff1660208301520161019c565b6101ea6103cc366004611715565b610c78565b6101ea610cb9565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f3149878600000000000000000000000000000000000000000000000000000000148061042e575061042e82610ccb565b92915050565b6003546060905f90158015610447575082155b156104bc57604080515f80825260208201909252906104b2565b604080516080810182525f8082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816104615790505b50935f9350915050565b600354831061052c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f66726f6d496e646578206f7574206f662072616e67650000000000000000000060448201526064015b60405180910390fd5b60325f610539828661183f565b60035490915081111561054b57506003545b5f6105568683611852565b90505f8167ffffffffffffffff811115610572576105726117e5565b6040519080825280602002602001820160405280156105e157816020015b604080516080810182525f8082526020808301829052928201819052606082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816105905790505b5090505f5b828110156106ae5760036105fa828a61183f565b8154811061060a5761060a611865565b5f91825260209182902060408051608081018252929091015467ffffffffffffffff80821684526801000000000000000082048116948401949094527001000000000000000000000000000000008104841691830191909152780100000000000000000000000000000000000000000000000090049091166060820152825183908390811061069b5761069b611865565b60209081029190910101526001016105e6565b50600354909792109550909350505050565b5f6106ca81610d61565b6106d2610d6b565b50565b8161070c576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107168282610d77565b5050565b81158015610742575060025473ffffffffffffffffffffffffffffffffffffffff8281169116145b156108155760015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff1681151580610796575065ffffffffffff8116155b806107a957504265ffffffffffff821610155b156107ea576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610523565b5050600180547fffffffffffff000000000000ffffffffffffffffffffffffffffffffffffffff1690555b6107168282610da1565b7fe1f251c81ea8f25043f9fc7ec35c0a2a10dd591e11346051677e0580f5887a3361084981610d61565b60035415806108ab57506003805461086390600190611852565b8154811061087357610873611865565b5f9182526020909120015467ffffffffffffffff78010000000000000000000000000000000000000000000000009091048116908316115b610937576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f737461727454696d65206d7573742062652067726561746572207468616e207460448201527f6865206c61737420737461727454696d650000000000000000000000000000006064820152608401610523565b6040805160808101825267ffffffffffffffff80881682528681166020830190815286821683850190815286831660608501908152600380546001810182555f9190915294517fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b9095018054935192519151851678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff92861670010000000000000000000000000000000002929092166fffffffffffffffffffffffffffffffff93861668010000000000000000027fffffffffffffffffffffffffffffffff000000000000000000000000000000009095169690951695909517929092171691909117179055517f3bff7b1e021b47f5dfd21d1d3fe2daaf3b9cbaca733c6e15b3a0da546657f19a90610aab90879087908790879067ffffffffffffffff948516815292841660208401529083166040830152909116606082015260800190565b60405180910390a15050505050565b5f610ac481610d61565b61071682610dff565b5f610ad781610d61565b61071682610e7e565b5f610b0060025473ffffffffffffffffffffffffffffffffffffffff1690565b905090565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015610b4757504265ffffffffffff821610155b610b52575f5f610b77565b60025474010000000000000000000000000000000000000000900465ffffffffffff16815b915091509091565b6002545f907a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015158015610bc057504265ffffffffffff8216105b610bf2576001547a010000000000000000000000000000000000000000000000000000900465ffffffffffff16610c16565b60025474010000000000000000000000000000000000000000900465ffffffffffff165b91505090565b60015473ffffffffffffffffffffffffffffffffffffffff16338114610c70576040517fc22c8022000000000000000000000000000000000000000000000000000000008152336004820152602401610523565b6106d2610eed565b81610caf576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107168282610fde565b5f610cc381610d61565b6106d2611002565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061042e57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000083161461042e565b6106d2813361100c565b610d755f5f611091565b565b5f82815260208190526040902060010154610d9181610d61565b610d9b83836111ea565b50505050565b73ffffffffffffffffffffffffffffffffffffffff81163314610df0576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610dfa82826112af565b505050565b5f610e08610b7f565b610e1142611310565b610e1b9190611892565b9050610e27828261135f565b60405165ffffffffffff8216815273ffffffffffffffffffffffffffffffffffffffff8316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b5f610e88826113fa565b610e9142611310565b610e9b9190611892565b9050610ea78282611091565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b910160405180910390a15050565b60015473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900465ffffffffffff16801580610f3d57504265ffffffffffff821610155b15610f7e576040517f19ca5ebb00000000000000000000000000000000000000000000000000000000815265ffffffffffff82166004820152602401610523565b610fa65f610fa160025473ffffffffffffffffffffffffffffffffffffffff1690565b6112af565b50610fb15f836111ea565b5050600180547fffffffffffff000000000000000000000000000000000000000000000000000016905550565b5f82815260208190526040902060010154610ff881610d61565b610d9b83836112af565b610d755f5f61135f565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610716576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260248101839052604401610523565b6002547a010000000000000000000000000000000000000000000000000000900465ffffffffffff168015611165574265ffffffffffff8216101561113c576002546001805479ffffffffffffffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000090920465ffffffffffff167a01000000000000000000000000000000000000000000000000000002919091179055611165565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5905f90a15b506002805473ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000065ffffffffffff9485160279ffffffffffffffffffffffffffffffffffffffffffffffffffff16177a0100000000000000000000000000000000000000000000000000009290931691909102919091179055565b5f8261129e575f61121060025473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff161461125d576040517f3fc3c27a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84161790555b6112a8838361144b565b9392505050565b5f821580156112d8575060025473ffffffffffffffffffffffffffffffffffffffff8381169116145b1561130657600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b6112a88383611544565b5f65ffffffffffff82111561135b576040517f6dfcc6500000000000000000000000000000000000000000000000000000000081526030600482015260248101839052604401610523565b5090565b600180547401000000000000000000000000000000000000000065ffffffffffff84811682027fffffffffffff0000000000000000000000000000000000000000000000000000841673ffffffffffffffffffffffffffffffffffffffff881617179093559004168015610dfa576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109905f90a1505050565b5f5f611404610b7f565b90508065ffffffffffff168365ffffffffffff161161142c5761142783826118b0565b6112a8565b6112a865ffffffffffff8416620697805f8282188284100282186112a8565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff1661153d575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556114db3390565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161042e565b505f61042e565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff161561153d575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161042e565b5f6020828403121561160d575f5ffd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146112a8575f5ffd5b5f6020828403121561164c575f5ffd5b5035919050565b604080825283519082018190525f9060208501906060840190835b818110156116da57835167ffffffffffffffff815116845267ffffffffffffffff602082015116602085015267ffffffffffffffff604082015116604085015267ffffffffffffffff60608201511660608501525060808301925060208401935060018101905061166e565b5050841515602085015291506112a89050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611710575f5ffd5b919050565b5f5f60408385031215611726575f5ffd5b82359150611736602084016116ed565b90509250929050565b803567ffffffffffffffff81168114611710575f5ffd5b5f5f5f5f60808587031215611769575f5ffd5b6117728561173f565b93506117806020860161173f565b925061178e6040860161173f565b915061179c6060860161173f565b905092959194509250565b5f602082840312156117b7575f5ffd5b6112a8826116ed565b5f602082840312156117d0575f5ffd5b813565ffffffffffff811681146112a8575f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b8082018082111561042e5761042e611812565b8181038181111561042e5761042e611812565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b65ffffffffffff818116838216019081111561042e5761042e611812565b65ffffffffffff828116828216039081111561042e5761042e61181256fea26469706673582212204564120b3bee6ddc2c134598c7aa9396c886c4b6edd56f444b5c5e78a43af61564736f6c634300081c0033e1f251c81ea8f25043f9fc7ec35c0a2a10dd591e11346051677e0580f5887a33",
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

// RATESADMINROLE is a free data retrieval call binding the contract method 0x970b6d87.
//
// Solidity: function RATES_ADMIN_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerCaller) RATESADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "RATES_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RATESADMINROLE is a free data retrieval call binding the contract method 0x970b6d87.
//
// Solidity: function RATES_ADMIN_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerSession) RATESADMINROLE() ([32]byte, error) {
	return _RatesManager.Contract.RATESADMINROLE(&_RatesManager.CallOpts)
}

// RATESADMINROLE is a free data retrieval call binding the contract method 0x970b6d87.
//
// Solidity: function RATES_ADMIN_ROLE() view returns(bytes32)
func (_RatesManager *RatesManagerCallerSession) RATESADMINROLE() ([32]byte, error) {
	return _RatesManager.Contract.RATESADMINROLE(&_RatesManager.CallOpts)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_RatesManager *RatesManagerCaller) DefaultAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "defaultAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_RatesManager *RatesManagerSession) DefaultAdmin() (common.Address, error) {
	return _RatesManager.Contract.DefaultAdmin(&_RatesManager.CallOpts)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_RatesManager *RatesManagerCallerSession) DefaultAdmin() (common.Address, error) {
	return _RatesManager.Contract.DefaultAdmin(&_RatesManager.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_RatesManager *RatesManagerCaller) DefaultAdminDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "defaultAdminDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_RatesManager *RatesManagerSession) DefaultAdminDelay() (*big.Int, error) {
	return _RatesManager.Contract.DefaultAdminDelay(&_RatesManager.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_RatesManager *RatesManagerCallerSession) DefaultAdminDelay() (*big.Int, error) {
	return _RatesManager.Contract.DefaultAdminDelay(&_RatesManager.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_RatesManager *RatesManagerCaller) DefaultAdminDelayIncreaseWait(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "defaultAdminDelayIncreaseWait")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_RatesManager *RatesManagerSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _RatesManager.Contract.DefaultAdminDelayIncreaseWait(&_RatesManager.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_RatesManager *RatesManagerCallerSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _RatesManager.Contract.DefaultAdminDelayIncreaseWait(&_RatesManager.CallOpts)
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
// Solidity: function getRatesCount() view returns(uint256)
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
// Solidity: function getRatesCount() view returns(uint256)
func (_RatesManager *RatesManagerSession) GetRatesCount() (*big.Int, error) {
	return _RatesManager.Contract.GetRatesCount(&_RatesManager.CallOpts)
}

// GetRatesCount is a free data retrieval call binding the contract method 0x2da72291.
//
// Solidity: function getRatesCount() view returns(uint256)
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RatesManager *RatesManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RatesManager *RatesManagerSession) Owner() (common.Address, error) {
	return _RatesManager.Contract.Owner(&_RatesManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RatesManager *RatesManagerCallerSession) Owner() (common.Address, error) {
	return _RatesManager.Contract.Owner(&_RatesManager.CallOpts)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_RatesManager *RatesManagerCaller) PendingDefaultAdmin(opts *bind.CallOpts) (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "pendingDefaultAdmin")

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
func (_RatesManager *RatesManagerSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _RatesManager.Contract.PendingDefaultAdmin(&_RatesManager.CallOpts)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_RatesManager *RatesManagerCallerSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _RatesManager.Contract.PendingDefaultAdmin(&_RatesManager.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_RatesManager *RatesManagerCaller) PendingDefaultAdminDelay(opts *bind.CallOpts) (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _RatesManager.contract.Call(opts, &out, "pendingDefaultAdminDelay")

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
func (_RatesManager *RatesManagerSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _RatesManager.Contract.PendingDefaultAdminDelay(&_RatesManager.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_RatesManager *RatesManagerCallerSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _RatesManager.Contract.PendingDefaultAdminDelay(&_RatesManager.CallOpts)
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

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_RatesManager *RatesManagerTransactor) AcceptDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "acceptDefaultAdminTransfer")
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_RatesManager *RatesManagerSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _RatesManager.Contract.AcceptDefaultAdminTransfer(&_RatesManager.TransactOpts)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_RatesManager *RatesManagerTransactorSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _RatesManager.Contract.AcceptDefaultAdminTransfer(&_RatesManager.TransactOpts)
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

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_RatesManager *RatesManagerTransactor) BeginDefaultAdminTransfer(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "beginDefaultAdminTransfer", newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_RatesManager *RatesManagerSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.BeginDefaultAdminTransfer(&_RatesManager.TransactOpts, newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_RatesManager *RatesManagerTransactorSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.BeginDefaultAdminTransfer(&_RatesManager.TransactOpts, newAdmin)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_RatesManager *RatesManagerTransactor) CancelDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "cancelDefaultAdminTransfer")
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_RatesManager *RatesManagerSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _RatesManager.Contract.CancelDefaultAdminTransfer(&_RatesManager.TransactOpts)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_RatesManager *RatesManagerTransactorSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _RatesManager.Contract.CancelDefaultAdminTransfer(&_RatesManager.TransactOpts)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_RatesManager *RatesManagerTransactor) ChangeDefaultAdminDelay(opts *bind.TransactOpts, newDelay *big.Int) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "changeDefaultAdminDelay", newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_RatesManager *RatesManagerSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _RatesManager.Contract.ChangeDefaultAdminDelay(&_RatesManager.TransactOpts, newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_RatesManager *RatesManagerTransactorSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _RatesManager.Contract.ChangeDefaultAdminDelay(&_RatesManager.TransactOpts, newDelay)
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

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.RenounceRole(&_RatesManager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_RatesManager *RatesManagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RatesManager.Contract.RenounceRole(&_RatesManager.TransactOpts, role, account)
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

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_RatesManager *RatesManagerTransactor) RollbackDefaultAdminDelay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RatesManager.contract.Transact(opts, "rollbackDefaultAdminDelay")
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_RatesManager *RatesManagerSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _RatesManager.Contract.RollbackDefaultAdminDelay(&_RatesManager.TransactOpts)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_RatesManager *RatesManagerTransactorSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _RatesManager.Contract.RollbackDefaultAdminDelay(&_RatesManager.TransactOpts)
}

// RatesManagerDefaultAdminDelayChangeCanceledIterator is returned from FilterDefaultAdminDelayChangeCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeCanceled events raised by the RatesManager contract.
type RatesManagerDefaultAdminDelayChangeCanceledIterator struct {
	Event *RatesManagerDefaultAdminDelayChangeCanceled // Event containing the contract specifics and raw log

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
func (it *RatesManagerDefaultAdminDelayChangeCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerDefaultAdminDelayChangeCanceled)
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
		it.Event = new(RatesManagerDefaultAdminDelayChangeCanceled)
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
func (it *RatesManagerDefaultAdminDelayChangeCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerDefaultAdminDelayChangeCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the RatesManager contract.
type RatesManagerDefaultAdminDelayChangeCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeCanceled is a free log retrieval operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_RatesManager *RatesManagerFilterer) FilterDefaultAdminDelayChangeCanceled(opts *bind.FilterOpts) (*RatesManagerDefaultAdminDelayChangeCanceledIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return &RatesManagerDefaultAdminDelayChangeCanceledIterator{contract: _RatesManager.contract, event: "DefaultAdminDelayChangeCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeCanceled is a free log subscription operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_RatesManager *RatesManagerFilterer) WatchDefaultAdminDelayChangeCanceled(opts *bind.WatchOpts, sink chan<- *RatesManagerDefaultAdminDelayChangeCanceled) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerDefaultAdminDelayChangeCanceled)
				if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseDefaultAdminDelayChangeCanceled(log types.Log) (*RatesManagerDefaultAdminDelayChangeCanceled, error) {
	event := new(RatesManagerDefaultAdminDelayChangeCanceled)
	if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerDefaultAdminDelayChangeScheduledIterator is returned from FilterDefaultAdminDelayChangeScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeScheduled events raised by the RatesManager contract.
type RatesManagerDefaultAdminDelayChangeScheduledIterator struct {
	Event *RatesManagerDefaultAdminDelayChangeScheduled // Event containing the contract specifics and raw log

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
func (it *RatesManagerDefaultAdminDelayChangeScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerDefaultAdminDelayChangeScheduled)
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
		it.Event = new(RatesManagerDefaultAdminDelayChangeScheduled)
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
func (it *RatesManagerDefaultAdminDelayChangeScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerDefaultAdminDelayChangeScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the RatesManager contract.
type RatesManagerDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeScheduled is a free log retrieval operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_RatesManager *RatesManagerFilterer) FilterDefaultAdminDelayChangeScheduled(opts *bind.FilterOpts) (*RatesManagerDefaultAdminDelayChangeScheduledIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return &RatesManagerDefaultAdminDelayChangeScheduledIterator{contract: _RatesManager.contract, event: "DefaultAdminDelayChangeScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeScheduled is a free log subscription operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_RatesManager *RatesManagerFilterer) WatchDefaultAdminDelayChangeScheduled(opts *bind.WatchOpts, sink chan<- *RatesManagerDefaultAdminDelayChangeScheduled) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerDefaultAdminDelayChangeScheduled)
				if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseDefaultAdminDelayChangeScheduled(log types.Log) (*RatesManagerDefaultAdminDelayChangeScheduled, error) {
	event := new(RatesManagerDefaultAdminDelayChangeScheduled)
	if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerDefaultAdminTransferCanceledIterator is returned from FilterDefaultAdminTransferCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferCanceled events raised by the RatesManager contract.
type RatesManagerDefaultAdminTransferCanceledIterator struct {
	Event *RatesManagerDefaultAdminTransferCanceled // Event containing the contract specifics and raw log

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
func (it *RatesManagerDefaultAdminTransferCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerDefaultAdminTransferCanceled)
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
		it.Event = new(RatesManagerDefaultAdminTransferCanceled)
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
func (it *RatesManagerDefaultAdminTransferCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerDefaultAdminTransferCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the RatesManager contract.
type RatesManagerDefaultAdminTransferCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferCanceled is a free log retrieval operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_RatesManager *RatesManagerFilterer) FilterDefaultAdminTransferCanceled(opts *bind.FilterOpts) (*RatesManagerDefaultAdminTransferCanceledIterator, error) {

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return &RatesManagerDefaultAdminTransferCanceledIterator{contract: _RatesManager.contract, event: "DefaultAdminTransferCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferCanceled is a free log subscription operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_RatesManager *RatesManagerFilterer) WatchDefaultAdminTransferCanceled(opts *bind.WatchOpts, sink chan<- *RatesManagerDefaultAdminTransferCanceled) (event.Subscription, error) {

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerDefaultAdminTransferCanceled)
				if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseDefaultAdminTransferCanceled(log types.Log) (*RatesManagerDefaultAdminTransferCanceled, error) {
	event := new(RatesManagerDefaultAdminTransferCanceled)
	if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RatesManagerDefaultAdminTransferScheduledIterator is returned from FilterDefaultAdminTransferScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferScheduled events raised by the RatesManager contract.
type RatesManagerDefaultAdminTransferScheduledIterator struct {
	Event *RatesManagerDefaultAdminTransferScheduled // Event containing the contract specifics and raw log

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
func (it *RatesManagerDefaultAdminTransferScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RatesManagerDefaultAdminTransferScheduled)
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
		it.Event = new(RatesManagerDefaultAdminTransferScheduled)
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
func (it *RatesManagerDefaultAdminTransferScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RatesManagerDefaultAdminTransferScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RatesManagerDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the RatesManager contract.
type RatesManagerDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferScheduled is a free log retrieval operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_RatesManager *RatesManagerFilterer) FilterDefaultAdminTransferScheduled(opts *bind.FilterOpts, newAdmin []common.Address) (*RatesManagerDefaultAdminTransferScheduledIterator, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _RatesManager.contract.FilterLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return &RatesManagerDefaultAdminTransferScheduledIterator{contract: _RatesManager.contract, event: "DefaultAdminTransferScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferScheduled is a free log subscription operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_RatesManager *RatesManagerFilterer) WatchDefaultAdminTransferScheduled(opts *bind.WatchOpts, sink chan<- *RatesManagerDefaultAdminTransferScheduled, newAdmin []common.Address) (event.Subscription, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _RatesManager.contract.WatchLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RatesManagerDefaultAdminTransferScheduled)
				if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
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
func (_RatesManager *RatesManagerFilterer) ParseDefaultAdminTransferScheduled(log types.Log) (*RatesManagerDefaultAdminTransferScheduled, error) {
	event := new(RatesManagerDefaultAdminTransferScheduled)
	if err := _RatesManager.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
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
