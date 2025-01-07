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

// NodesNode is an auto generated low-level Go binding around an user-defined struct.
type NodesNode struct {
	SigningKeyPub []byte
	HttpAddress   string
	IsHealthy     bool
}

// NodesNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type NodesNodeWithId struct {
	NodeId uint32
	Node   NodesNode
}

// NodesMetaData contains all meta data concerning the Nodes contract.
var NodesMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structNodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structNodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"internalType\":\"bool\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"healthyNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structNodes.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structNodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"internalType\":\"bool\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateHealth\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateHttpAddress\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"node\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structNodes.Node\",\"components\":[{\"name\":\"signingKeyPub\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isHealthy\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60806040526006805463ffffffff60a01b1916905534801561001f575f5ffd5b5033604051806040016040528060128152602001712c26aa28102737b2329027b832b930ba37b960711b815250604051806040016040528060048152602001630584d54560e41b815250815f908161007791906101ac565b50600161008482826101ac565b5050506001600160a01b0381166100b457604051631e4fbdf760e01b81525f600482015260240160405180910390fd5b6100bd816100c3565b50610266565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061013c57607f821691505b60208210810361015a57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156101a757805f5260205f20601f840160051c810160208510156101855750805b601f840160051c820191505b818110156101a4575f8155600101610191565b50505b505050565b81516001600160401b038111156101c5576101c5610114565b6101d9816101d38454610128565b84610160565b6020601f82116001811461020b575f83156101f45750848201515b5f19600385901b1c1916600184901b1784556101a4565b5f84815260208120601f198516915b8281101561023a578785015182556020948501946001909201910161021a565b508482101561025757868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b612539806102735f395ff3fe608060405234801561000f575f5ffd5b5060043610610179575f3560e01c806370a08231116100d2578063a0eae81d11610088578063c87b56dd11610063578063c87b56dd14610326578063e985e9c514610339578063f2fde38b14610374575f5ffd5b8063a0eae81d146102d8578063a22cb46514610300578063b88d4fde14610313575f5ffd5b80638354196f116100b85780638354196f146102ac5780638da5cb5b146102bf57806395d89b41146102d0575f5ffd5b806370a0823114610283578063715018a6146102a4575f5ffd5b806317b5b840116101325780634f0f4aa91161010d5780634f0f4aa91461023d57806350d809931461025d5780636352211e14610270575f5ffd5b806317b5b8401461020f57806323b872dd1461021757806342842e0e1461022a575f5ffd5b806306fdde031161016257806306fdde03146101ba578063081812fc146101cf578063095ea7b3146101fa575f5ffd5b806301ffc9a71461017d57806302f6f1ba146101a5575b5f5ffd5b61019061018b366004611bd3565b610387565b60405190151581526020015b60405180910390f35b6101ad61046b565b60405161019c9190611c62565b6101c26106f7565b60405161019c9190611d00565b6101e26101dd366004611d12565b610786565b6040516001600160a01b03909116815260200161019c565b61020d610208366004611d44565b6107ad565b005b6101ad6107bc565b61020d610225366004611d6c565b610aec565b61020d610238366004611d6c565b610b9b565b61025061024b366004611d12565b610bb5565b60405161019c9190611da6565b61020d61026b366004611dfd565b610d28565b6101e261027e366004611d12565b610df6565b610296610291366004611e45565b610e00565b60405190815260200161019c565b61020d610e5e565b61020d6102ba366004611e6d565b610e71565b6006546001600160a01b03166101e2565b6101c2610ec6565b6102eb6102e6366004611e97565b610ed5565b60405163ffffffff909116815260200161019c565b61020d61030e366004611f17565b611069565b61020d610321366004611f6c565b611074565b6101c2610334366004611d12565b611092565b61019061034736600461204a565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b61020d610382366004611e45565b611103565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd00000000000000000000000000000000000000000000000000000000148061041957507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b5e139f00000000000000000000000000000000000000000000000000000000145b8061046557507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6006546060905f9074010000000000000000000000000000000000000000900463ffffffff1667ffffffffffffffff8111156104a9576104a9611f3f565b6040519080825280602002602001820160405280156104e257816020015b6104cf611b65565b8152602001906001900390816104c75790505b5090505f5b60065463ffffffff74010000000000000000000000000000000000000000909104811690821610156106f1575f61051f82600161209f565b61052a9060646120bb565b90506105528163ffffffff165f908152600260205260409020546001600160a01b0316151590565b156106e85760405180604001604052808263ffffffff16815260200160075f8463ffffffff1681526020019081526020015f206040518060600160405290815f8201805461059f906120e1565b80601f01602080910402602001604051908101604052809291908181526020018280546105cb906120e1565b80156106165780601f106105ed57610100808354040283529160200191610616565b820191905f5260205f20905b8154815290600101906020018083116105f957829003601f168201915b5050505050815260200160018201805461062f906120e1565b80601f016020809104026020016040519081016040528092919081815260200182805461065b906120e1565b80156106a65780601f1061067d576101008083540402835291602001916106a6565b820191905f5260205f20905b81548152906001019060200180831161068957829003601f168201915b50505091835250506002919091015460ff16151560209091015290528351849063ffffffff85169081106106dc576106dc61212c565b60200260200101819052505b506001016104e7565b50919050565b60605f8054610705906120e1565b80601f0160208091040260200160405190810160405280929190818152602001828054610731906120e1565b801561077c5780601f106107535761010080835404028352916020019161077c565b820191905f5260205f20905b81548152906001019060200180831161075f57829003601f168201915b5050505050905090565b5f61079082611159565b505f828152600460205260409020546001600160a01b0316610465565b6107b88282336111aa565b5050565b60605f805b60065474010000000000000000000000000000000000000000900463ffffffff16811015610854575f6107f5826001612159565b61080090606461216c565b5f818152600260205260409020549091506001600160a01b03161515801561083857505f8181526007602052604090206002015460ff165b1561084b578261084781612183565b9350505b506001016107c1565b505f8167ffffffffffffffff81111561086f5761086f611f3f565b6040519080825280602002602001820160405280156108a857816020015b610895611b65565b81526020019060019003908161088d5790505b5090505f805b60065463ffffffff7401000000000000000000000000000000000000000090910481169082161015610ae3575f6108e682600161209f565b6108f19060646120bb565b90506109198163ffffffff165f908152600260205260409020546001600160a01b0316151590565b801561093c575063ffffffff81165f9081526007602052604090206002015460ff165b15610ada5760405180604001604052808263ffffffff16815260200160075f8463ffffffff1681526020019081526020015f206040518060600160405290815f82018054610989906120e1565b80601f01602080910402602001604051908101604052809291908181526020018280546109b5906120e1565b8015610a005780601f106109d757610100808354040283529160200191610a00565b820191905f5260205f20905b8154815290600101906020018083116109e357829003601f168201915b50505050508152602001600182018054610a19906120e1565b80601f0160208091040260200160405190810160405280929190818152602001828054610a45906120e1565b8015610a905780601f10610a6757610100808354040283529160200191610a90565b820191905f5260205f20905b815481529060010190602001808311610a7357829003601f168201915b50505091835250506002919091015460ff16151560209091015290528451859085908110610ac057610ac061212c565b60200260200101819052508280610ad690612183565b9350505b506001016108ae565b50909392505050565b6006546001600160a01b03163314610b8b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603360248201527f4f6e6c792074686520636f6e7472616374206f776e65722063616e207472616e60448201527f73666572204e6f6465206f776e6572736869700000000000000000000000000060648201526084015b60405180910390fd5b610b968383836111b7565b505050565b610b9683838360405180602001604052805f815250611074565b604080516060808201835280825260208201525f91810191909152610bd982611159565b505f8281526007602052604090819020815160608101909252805482908290610c01906120e1565b80601f0160208091040260200160405190810160405280929190818152602001828054610c2d906120e1565b8015610c785780601f10610c4f57610100808354040283529160200191610c78565b820191905f5260205f20905b815481529060010190602001808311610c5b57829003601f168201915b50505050508152602001600182018054610c91906120e1565b80601f0160208091040260200160405190810160405280929190818152602001828054610cbd906120e1565b8015610d085780601f10610cdf57610100808354040283529160200191610d08565b820191905f5260205f20905b815481529060010190602001808311610ceb57829003601f168201915b50505091835250506002919091015460ff16151560209091015292915050565b610d3183610df6565b6001600160a01b0316336001600160a01b031614610dd1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603a60248201527f4f6e6c7920746865206f776e6572206f6620746865204e6f6465204e4654206360448201527f616e2075706461746520697473206874747020616464726573730000000000006064820152608401610b82565b5f838152600760205260409020600101610dec8284836121df565b50610b968361126c565b5f61046582611159565b5f6001600160a01b038216610e43576040517f89c62b640000000000000000000000000000000000000000000000000000000081525f6004820152602401610b82565b506001600160a01b03165f9081526003602052604090205490565b610e666112b6565b610e6f5f6112fc565b565b610e796112b6565b610e8282611159565b505f82815260076020526040902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168215151790556107b88261126c565b606060018054610705906120e1565b5f610ede6112b6565b6006805474010000000000000000000000000000000000000000900463ffffffff16906014610f0c83612299565b91906101000a81548163ffffffff021916908363ffffffff160217905550505f6064600660149054906101000a900463ffffffff16610f4b91906120bb565b9050610f5d878263ffffffff16611365565b604051806060016040528087878080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250505090825250604080516020601f8801819004810282018101909252868152918101919087908790819084018382808284375f920182905250938552505060016020938401525063ffffffff841681526007909152604090208151819061100190826122bd565b506020820151600182019061101690826122bd565b5060409190910151600290910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905561105f63ffffffff821661126c565b9695505050505050565b6107b83383836113f8565b61107f848484610aec565b61108c33858585856114cd565b50505050565b606061109d82611159565b505f6110b360408051602081019091525f815290565b90505f8151116110d15760405180602001604052805f8152506110fc565b806110db84611671565b6040516020016110ec92919061238f565b6040516020818303038152906040525b9392505050565b61110b6112b6565b6001600160a01b03811661114d576040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f6004820152602401610b82565b611156816112fc565b50565b5f818152600260205260408120546001600160a01b031680610465576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101849052602401610b82565b610b96838383600161170e565b6001600160a01b0382166111f9576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610b82565b5f611205838333611861565b9050836001600160a01b0316816001600160a01b03161461108c576040517f64283d7b0000000000000000000000000000000000000000000000000000000081526001600160a01b0380861660048301526024820184905282166044820152606401610b82565b5f818152600760205260409081902090517f925a803604a97550dcfaba2eb298fb6e6f828340bdf9d467d2a34c50cfbba96a916112ab91849190612440565b60405180910390a150565b6006546001600160a01b03163314610e6f576040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152602401610b82565b600680546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b6001600160a01b0382166113a7576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610b82565b5f6113b383835f611861565b90506001600160a01b03811615610b96576040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610b82565b6001600160a01b038216611443576040517f5b08ba180000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602401610b82565b6001600160a01b038381165f8181526005602090815260408083209487168084529482529182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b1561166a576040517f150b7a020000000000000000000000000000000000000000000000000000000081526001600160a01b0384169063150b7a02906115289088908890879087906004016124b2565b6020604051808303815f875af1925050508015611562575060408051601f3d908101601f1916820190925261155f918101906124e8565b60015b6115e2573d80801561158f576040519150601f19603f3d011682016040523d82523d5f602084013e611594565b606091505b5080515f036115da576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610b82565b805181602001fd5b7fffffffff0000000000000000000000000000000000000000000000000000000081167f150b7a020000000000000000000000000000000000000000000000000000000014611668576040517f64a0ae920000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610b82565b505b5050505050565b60605f61167d8361196b565b60010190505f8167ffffffffffffffff81111561169c5761169c611f3f565b6040519080825280601f01601f1916602001820160405280156116c6576020820181803683370190505b5090508181016020015b5f19017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85049450846116d057509392505050565b808061172257506001600160a01b03821615155b1561181a575f61173184611159565b90506001600160a01b0383161580159061175d5750826001600160a01b0316816001600160a01b031614155b801561178e57506001600160a01b038082165f9081526005602090815260408083209387168352929052205460ff16155b156117d0576040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526001600160a01b0384166004820152602401610b82565b81156118185783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260046020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b5f828152600260205260408120546001600160a01b039081169083161561188d5761188d818486611a4c565b6001600160a01b038116156118c7576118a85f855f5f61170e565b6001600160a01b0381165f90815260036020526040902080545f190190555b6001600160a01b038516156118f5576001600160a01b0385165f908152600360205260409020805460010190555b5f8481526002602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b5f807a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083106119b3577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef810000000083106119df576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc1000083106119fd57662386f26fc10000830492506010015b6305f5e1008310611a15576305f5e100830492506008015b6127108310611a2957612710830492506004015b60648310611a3b576064830492506002015b600a83106104655760010192915050565b611a57838383611ae2565b610b96576001600160a01b038316611a9e576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101829052602401610b82565b6040517f177e802f0000000000000000000000000000000000000000000000000000000081526001600160a01b038316600482015260248101829052604401610b82565b5f6001600160a01b03831615801590611b5d5750826001600160a01b0316846001600160a01b03161480611b3a57506001600160a01b038085165f9081526005602090815260408083209387168352929052205460ff165b80611b5d57505f828152600460205260409020546001600160a01b038481169116145b949350505050565b60405180604001604052805f63ffffffff168152602001611ba1604051806060016040528060608152602001606081526020015f151581525090565b905290565b7fffffffff0000000000000000000000000000000000000000000000000000000081168114611156575f5ffd5b5f60208284031215611be3575f5ffd5b81356110fc81611ba6565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b5f815160608452611c306060850182611bee565b905060208301518482036020860152611c498282611bee565b9150506040830151151560408501528091505092915050565b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b82811015611cf4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0878603018452815163ffffffff81511686526020810151905060406020870152611cde6040870182611c1c565b9550506020938401939190910190600101611c88565b50929695505050505050565b602081525f6110fc6020830184611bee565b5f60208284031215611d22575f5ffd5b5035919050565b80356001600160a01b0381168114611d3f575f5ffd5b919050565b5f5f60408385031215611d55575f5ffd5b611d5e83611d29565b946020939093013593505050565b5f5f5f60608486031215611d7e575f5ffd5b611d8784611d29565b9250611d9560208501611d29565b929592945050506040919091013590565b602081525f6110fc6020830184611c1c565b5f5f83601f840112611dc8575f5ffd5b50813567ffffffffffffffff811115611ddf575f5ffd5b602083019150836020828501011115611df6575f5ffd5b9250929050565b5f5f5f60408486031215611e0f575f5ffd5b83359250602084013567ffffffffffffffff811115611e2c575f5ffd5b611e3886828701611db8565b9497909650939450505050565b5f60208284031215611e55575f5ffd5b6110fc82611d29565b80358015158114611d3f575f5ffd5b5f5f60408385031215611e7e575f5ffd5b82359150611e8e60208401611e5e565b90509250929050565b5f5f5f5f5f60608688031215611eab575f5ffd5b611eb486611d29565b9450602086013567ffffffffffffffff811115611ecf575f5ffd5b611edb88828901611db8565b909550935050604086013567ffffffffffffffff811115611efa575f5ffd5b611f0688828901611db8565b969995985093965092949392505050565b5f5f60408385031215611f28575f5ffd5b611f3183611d29565b9150611e8e60208401611e5e565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f5f5f60808587031215611f7f575f5ffd5b611f8885611d29565b9350611f9660208601611d29565b925060408501359150606085013567ffffffffffffffff811115611fb8575f5ffd5b8501601f81018713611fc8575f5ffd5b803567ffffffffffffffff811115611fe257611fe2611f3f565b604051601f19603f601f19601f8501160116810181811067ffffffffffffffff8211171561201257612012611f3f565b604052818152828201602001891015612029575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b5f5f6040838503121561205b575f5ffd5b61206483611d29565b9150611e8e60208401611d29565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b63ffffffff818116838216019081111561046557610465612072565b63ffffffff81811683821602908116908181146120da576120da612072565b5092915050565b600181811c908216806120f557607f821691505b6020821081036106f1577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b8082018082111561046557610465612072565b808202811582820484141761046557610465612072565b5f5f19820361219457612194612072565b5060010190565b601f821115610b9657805f5260205f20601f840160051c810160208510156121c05750805b601f840160051c820191505b8181101561166a575f81556001016121cc565b67ffffffffffffffff8311156121f7576121f7611f3f565b61220b8361220583546120e1565b8361219b565b5f601f84116001811461223c575f85156122255750838201355b5f19600387901b1c1916600186901b17835561166a565b5f83815260208120601f198716915b8281101561226b578685013582556020948501946001909201910161224b565b5086821015612287575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b5f63ffffffff821663ffffffff81036122b4576122b4612072565b60010192915050565b815167ffffffffffffffff8111156122d7576122d7611f3f565b6122eb816122e584546120e1565b8461219b565b6020601f82116001811461231d575f83156123065750848201515b5f19600385901b1c1916600184901b17845561166a565b5f84815260208120601f198516915b8281101561234c578785015182556020948501946001909201910161232c565b508482101561236957868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b5f81518060208401855e5f93019283525090919050565b5f611b5d61239d8386612378565b84612378565b5f81546123af816120e1565b8085526001821680156123c9576001811461240357612437565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083166020870152602082151560051b8701019350612437565b845f5260205f205f5b8381101561242e5781546020828a01015260018201915060208101905061240c565b87016020019450505b50505092915050565b82815260406020820152606060408201525f61245f60a08301846123a3565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc083820301606084015261249681600186016123a3565b905060ff60028501541615156080840152809150509392505050565b6001600160a01b03851681526001600160a01b0384166020820152826040820152608060608201525f61105f6080830184611bee565b5f602082840312156124f8575f5ffd5b81516110fc81611ba656fea2646970667358221220666875a4755217bd116c1476c7d2d9e0e8228592ebb579176663769d172daabc64736f6c634300081c0033",
}

// NodesABI is the input ABI used to generate the binding from.
// Deprecated: Use NodesMetaData.ABI instead.
var NodesABI = NodesMetaData.ABI

// NodesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodesMetaData.Bin instead.
var NodesBin = NodesMetaData.Bin

// DeployNodes deploys a new Ethereum contract, binding an instance of Nodes to it.
func DeployNodes(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Nodes, error) {
	parsed, err := NodesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodesBin), backend)
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

// AllNodes is a free data retrieval call binding the contract method 0x02f6f1ba.
//
// Solidity: function allNodes() view returns((uint32,(bytes,string,bool))[])
func (_Nodes *NodesCaller) AllNodes(opts *bind.CallOpts) ([]NodesNodeWithId, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "allNodes")

	if err != nil {
		return *new([]NodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]NodesNodeWithId)).(*[]NodesNodeWithId)

	return out0, err

}

// AllNodes is a free data retrieval call binding the contract method 0x02f6f1ba.
//
// Solidity: function allNodes() view returns((uint32,(bytes,string,bool))[])
func (_Nodes *NodesSession) AllNodes() ([]NodesNodeWithId, error) {
	return _Nodes.Contract.AllNodes(&_Nodes.CallOpts)
}

// AllNodes is a free data retrieval call binding the contract method 0x02f6f1ba.
//
// Solidity: function allNodes() view returns((uint32,(bytes,string,bool))[])
func (_Nodes *NodesCallerSession) AllNodes() ([]NodesNodeWithId, error) {
	return _Nodes.Contract.AllNodes(&_Nodes.CallOpts)
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
// Solidity: function getNode(uint256 tokenId) view returns((bytes,string,bool))
func (_Nodes *NodesCaller) GetNode(opts *bind.CallOpts, tokenId *big.Int) (NodesNode, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "getNode", tokenId)

	if err != nil {
		return *new(NodesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(NodesNode)).(*NodesNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 tokenId) view returns((bytes,string,bool))
func (_Nodes *NodesSession) GetNode(tokenId *big.Int) (NodesNode, error) {
	return _Nodes.Contract.GetNode(&_Nodes.CallOpts, tokenId)
}

// GetNode is a free data retrieval call binding the contract method 0x4f0f4aa9.
//
// Solidity: function getNode(uint256 tokenId) view returns((bytes,string,bool))
func (_Nodes *NodesCallerSession) GetNode(tokenId *big.Int) (NodesNode, error) {
	return _Nodes.Contract.GetNode(&_Nodes.CallOpts, tokenId)
}

// HealthyNodes is a free data retrieval call binding the contract method 0x17b5b840.
//
// Solidity: function healthyNodes() view returns((uint32,(bytes,string,bool))[])
func (_Nodes *NodesCaller) HealthyNodes(opts *bind.CallOpts) ([]NodesNodeWithId, error) {
	var out []interface{}
	err := _Nodes.contract.Call(opts, &out, "healthyNodes")

	if err != nil {
		return *new([]NodesNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]NodesNodeWithId)).(*[]NodesNodeWithId)

	return out0, err

}

// HealthyNodes is a free data retrieval call binding the contract method 0x17b5b840.
//
// Solidity: function healthyNodes() view returns((uint32,(bytes,string,bool))[])
func (_Nodes *NodesSession) HealthyNodes() ([]NodesNodeWithId, error) {
	return _Nodes.Contract.HealthyNodes(&_Nodes.CallOpts)
}

// HealthyNodes is a free data retrieval call binding the contract method 0x17b5b840.
//
// Solidity: function healthyNodes() view returns((uint32,(bytes,string,bool))[])
func (_Nodes *NodesCallerSession) HealthyNodes() ([]NodesNodeWithId, error) {
	return _Nodes.Contract.HealthyNodes(&_Nodes.CallOpts)
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

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress) returns(uint32)
func (_Nodes *NodesTransactor) AddNode(opts *bind.TransactOpts, to common.Address, signingKeyPub []byte, httpAddress string) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "addNode", to, signingKeyPub, httpAddress)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress) returns(uint32)
func (_Nodes *NodesSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, to, signingKeyPub, httpAddress)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address to, bytes signingKeyPub, string httpAddress) returns(uint32)
func (_Nodes *NodesTransactorSession) AddNode(to common.Address, signingKeyPub []byte, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.AddNode(&_Nodes.TransactOpts, to, signingKeyPub, httpAddress)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nodes *NodesTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nodes *NodesSession) RenounceOwnership() (*types.Transaction, error) {
	return _Nodes.Contract.RenounceOwnership(&_Nodes.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nodes *NodesTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Nodes.Contract.RenounceOwnership(&_Nodes.TransactOpts)
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

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Nodes *NodesTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Nodes *NodesSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.TransferFrom(&_Nodes.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Nodes *NodesTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nodes.Contract.TransferFrom(&_Nodes.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nodes *NodesTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nodes *NodesSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.TransferOwnership(&_Nodes.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nodes *NodesTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Nodes.Contract.TransferOwnership(&_Nodes.TransactOpts, newOwner)
}

// UpdateHealth is a paid mutator transaction binding the contract method 0x8354196f.
//
// Solidity: function updateHealth(uint256 tokenId, bool isHealthy) returns()
func (_Nodes *NodesTransactor) UpdateHealth(opts *bind.TransactOpts, tokenId *big.Int, isHealthy bool) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateHealth", tokenId, isHealthy)
}

// UpdateHealth is a paid mutator transaction binding the contract method 0x8354196f.
//
// Solidity: function updateHealth(uint256 tokenId, bool isHealthy) returns()
func (_Nodes *NodesSession) UpdateHealth(tokenId *big.Int, isHealthy bool) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateHealth(&_Nodes.TransactOpts, tokenId, isHealthy)
}

// UpdateHealth is a paid mutator transaction binding the contract method 0x8354196f.
//
// Solidity: function updateHealth(uint256 tokenId, bool isHealthy) returns()
func (_Nodes *NodesTransactorSession) UpdateHealth(tokenId *big.Int, isHealthy bool) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateHealth(&_Nodes.TransactOpts, tokenId, isHealthy)
}

// UpdateHttpAddress is a paid mutator transaction binding the contract method 0x50d80993.
//
// Solidity: function updateHttpAddress(uint256 tokenId, string httpAddress) returns()
func (_Nodes *NodesTransactor) UpdateHttpAddress(opts *bind.TransactOpts, tokenId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.contract.Transact(opts, "updateHttpAddress", tokenId, httpAddress)
}

// UpdateHttpAddress is a paid mutator transaction binding the contract method 0x50d80993.
//
// Solidity: function updateHttpAddress(uint256 tokenId, string httpAddress) returns()
func (_Nodes *NodesSession) UpdateHttpAddress(tokenId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateHttpAddress(&_Nodes.TransactOpts, tokenId, httpAddress)
}

// UpdateHttpAddress is a paid mutator transaction binding the contract method 0x50d80993.
//
// Solidity: function updateHttpAddress(uint256 tokenId, string httpAddress) returns()
func (_Nodes *NodesTransactorSession) UpdateHttpAddress(tokenId *big.Int, httpAddress string) (*types.Transaction, error) {
	return _Nodes.Contract.UpdateHttpAddress(&_Nodes.TransactOpts, tokenId, httpAddress)
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

// NodesNodeUpdatedIterator is returned from FilterNodeUpdated and is used to iterate over the raw logs and unpacked data for NodeUpdated events raised by the Nodes contract.
type NodesNodeUpdatedIterator struct {
	Event *NodesNodeUpdated // Event containing the contract specifics and raw log

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
func (it *NodesNodeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesNodeUpdated)
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
		it.Event = new(NodesNodeUpdated)
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
func (it *NodesNodeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesNodeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesNodeUpdated represents a NodeUpdated event raised by the Nodes contract.
type NodesNodeUpdated struct {
	NodeId *big.Int
	Node   NodesNode
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeUpdated is a free log retrieval operation binding the contract event 0x925a803604a97550dcfaba2eb298fb6e6f828340bdf9d467d2a34c50cfbba96a.
//
// Solidity: event NodeUpdated(uint256 nodeId, (bytes,string,bool) node)
func (_Nodes *NodesFilterer) FilterNodeUpdated(opts *bind.FilterOpts) (*NodesNodeUpdatedIterator, error) {

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "NodeUpdated")
	if err != nil {
		return nil, err
	}
	return &NodesNodeUpdatedIterator{contract: _Nodes.contract, event: "NodeUpdated", logs: logs, sub: sub}, nil
}

// WatchNodeUpdated is a free log subscription operation binding the contract event 0x925a803604a97550dcfaba2eb298fb6e6f828340bdf9d467d2a34c50cfbba96a.
//
// Solidity: event NodeUpdated(uint256 nodeId, (bytes,string,bool) node)
func (_Nodes *NodesFilterer) WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *NodesNodeUpdated) (event.Subscription, error) {

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "NodeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesNodeUpdated)
				if err := _Nodes.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
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

// ParseNodeUpdated is a log parse operation binding the contract event 0x925a803604a97550dcfaba2eb298fb6e6f828340bdf9d467d2a34c50cfbba96a.
//
// Solidity: event NodeUpdated(uint256 nodeId, (bytes,string,bool) node)
func (_Nodes *NodesFilterer) ParseNodeUpdated(log types.Log) (*NodesNodeUpdated, error) {
	event := new(NodesNodeUpdated)
	if err := _Nodes.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodesOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Nodes contract.
type NodesOwnershipTransferredIterator struct {
	Event *NodesOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NodesOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodesOwnershipTransferred)
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
		it.Event = new(NodesOwnershipTransferred)
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
func (it *NodesOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodesOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodesOwnershipTransferred represents a OwnershipTransferred event raised by the Nodes contract.
type NodesOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nodes *NodesFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NodesOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Nodes.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NodesOwnershipTransferredIterator{contract: _Nodes.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nodes *NodesFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NodesOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Nodes.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodesOwnershipTransferred)
				if err := _Nodes.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nodes *NodesFilterer) ParseOwnershipTransferred(log types.Log) (*NodesOwnershipTransferred, error) {
	event := new(NodesOwnershipTransferred)
	if err := _Nodes.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
