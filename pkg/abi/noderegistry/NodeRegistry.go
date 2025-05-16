// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package noderegistry

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

// INodeRegistryNode is an auto generated low-level Go binding around an user-defined struct.
type INodeRegistryNode struct {
	Signer           common.Address
	IsCanonical      bool
	SigningPublicKey []byte
	HttpAddress      string
}

// INodeRegistryNodeWithId is an auto generated low-level Go binding around an user-defined struct.
type INodeRegistryNodeWithId struct {
	NodeId uint32
	Node   INodeRegistryNode
}

// NodeRegistryMetaData contains all meta data concerning the NodeRegistry contract.
var NodeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"parameterRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"NODE_INCREMENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addNode\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signingPublicKey_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"signer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addToNetwork\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"admin_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"adminParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"canonicalNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"canonicalNodesCount_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"allNodes_\",\"type\":\"tuple[]\",\"internalType\":\"structINodeRegistry.NodeWithId[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"node\",\"type\":\"tuple\",\"internalType\":\"structINodeRegistry.Node\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isCanonical\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signingPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeCount_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsCanonicalNode\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"isCanonicalNode_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"node_\",\"type\":\"tuple\",\"internalType\":\"structINodeRegistry.Node\",\"components\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isCanonical\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signingPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSigner\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"signer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"implementation_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxCanonicalNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"maxCanonicalNodes_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxCanonicalNodesParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migratorParameterKey\",\"inputs\":[],\"outputs\":[{\"name\":\"key_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"parameterRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeFromNetwork\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"baseURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setHttpAddress\",\"inputs\":[{\"name\":\"nodeId_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"httpAddress_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAdmin\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxCanonicalNodes\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminUpdated\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BaseURIUpdated\",\"inputs\":[{\"name\":\"baseURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HttpAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxCanonicalNodesUpdated\",\"inputs\":[{\"name\":\"maxCanonicalNodes\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"migrator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAdded\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"signingPublicKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"httpAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeAddedToCanonicalNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeRemovedFromCanonicalNetwork\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EmptyCode\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FailedToAddNodeToCanonicalNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToRemoveNodeFromCanonicalNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHttpAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigningPublicKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCanonicalNodesBelowCurrentCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCanonicalNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxNodesReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationFailed\",\"inputs\":[{\"name\":\"migrator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertData_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNodeOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterOutOfTypeBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroMigrator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroParameterRegistry\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051613d55380380613d5583398101604081905261002e9161011a565b6001600160a01b038116608081905261005a5760405163d973fd8d60e01b815260040160405180910390fd5b610062610068565b50610147565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100b85760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146101175780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b5f6020828403121561012a575f5ffd5b81516001600160a01b0381168114610140575f5ffd5b9392505050565b608051613be16101745f395f81816102a101528181610fe4015281816112b00152611aaa0152613be15ff3fe608060405234801561000f575f5ffd5b506004361061024f575f3560e01c80638cf20c681161013d578063c18e273d116100b8578063e06f876f11610088578063f84ce8b91161006e578063f84ce8b91461068e578063f851a440146106a1578063fd667d1e146106de575f5ffd5b8063e06f876f14610607578063e985e9c514610627575f5ffd5b8063c18e273d14610559578063c87b56dd146105a9578063c9c02a02146105bc578063d3b2f598146105ff575f5ffd5b8063a0eae81d1161010d578063a22cb465116100f3578063a22cb46514610520578063ad03d0a514610533578063b88d4fde14610546575f5ffd5b8063a0eae81d146104c7578063a1174e7d1461050b575f5ffd5b80638cf20c681461046b5780638fd3ab801461047e57806395d89b41146104865780639f40b6251461048e575f5ffd5b806350d0215f116101cd57806368501a3e1161019d5780638129fc1c116101835780638129fc1c1461042257806382a5cfc31461042a5780638aab82ba14610432575f5ffd5b806368501a3e146103ee57806370a0823114610401575f5ffd5b806350d0215f1461034957806355f804b3146103a15780635c60da1b146103b45780636352211e146103db575f5ffd5b8063081812fc11610222578063236b6eb811610208578063236b6eb81461031057806323b872dd1461032357806342842e0e14610336575f5ffd5b8063081812fc146102e8578063095ea7b3146102fb575f5ffd5b80630124b8821461025357806301ffc9a71461027157806306fdde03146102945780630723499e1461029c575b5f5ffd5b61025b6106e6565b604051610268919061302e565b60405180910390f35b61028461027f36600461306d565b610706565b6040519015158152602001610268565b61025b6107ea565b6102c37f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610268565b6102c36102f6366004613088565b61089e565b61030e6103093660046130c7565b6108f1565b005b61030e61031e366004613102565b610900565b61030e61033136600461311b565b610aaa565b61030e61034436600461311b565b610b9e565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200054760100000000000000000000000000000000000000000000900463ffffffff165b60405163ffffffff9091168152602001610268565b61030e6103af36600461319a565b610bbd565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546102c3565b6102c36103e9366004613088565b610d05565b6102c36103fc366004613102565b610d0f565b61041461040f3660046131d9565b610d6d565b604051908152602001610268565b61030e610e0b565b61030e610fde565b60408051808201909152601a81527f786d74702e6e6f646552656769737472792e6d69677261746f72000000000000602082015261025b565b61030e610479366004613102565b611160565b61030e6112a8565b61025b611315565b60408051808201909152601781527f786d74702e6e6f646552656769737472792e61646d696e000000000000000000602082015261025b565b6104da6104d53660046131f2565b611366565b6040805163ffffffff909316835273ffffffffffffffffffffffffffffffffffffffff909116602083015201610268565b61051361170c565b60405161026891906132d4565b61030e61052e366004613372565b6119b9565b610284610541366004613102565b6119c4565b61030e6105543660046133d8565b611a27565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a20005474010000000000000000000000000000000000000000900460ff165b60405160ff9091168152602001610268565b61025b6105b7366004613088565b611a3f565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2000547501000000000000000000000000000000000000000000900460ff16610597565b61030e611aa4565b61061a610615366004613102565b611be8565b60405161026891906134f2565b610284610635366004613504565b73ffffffffffffffffffffffffffffffffffffffff9182165f9081527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab00793056020908152604080832093909416825291909152205460ff1690565b61030e61069c366004613535565b611dc6565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a20005473ffffffffffffffffffffffffffffffffffffffff166102c3565b61038c606481565b6060604051806060016040528060238152602001613b8960239139905090565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd00000000000000000000000000000000000000000000000000000000148061079857507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b5e139f00000000000000000000000000000000000000000000000000000000145b806107e457507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b7f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab0079300805460609190819061081c90613584565b80601f016020809104026020016040519081016040528092919081815260200182805461084890613584565b80156108935780601f1061086a57610100808354040283529160200191610893565b820191905f5260205f20905b81548152906001019060200180831161087657829003601f168201915b505050505091505090565b5f6108a882611e86565b505f8281527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab0079304602052604090205473ffffffffffffffffffffffffffffffffffffffff166107e4565b6108fc828233611f03565b5050565b610908611f10565b6109178163ffffffff16611e86565b5063ffffffff81165f9081527fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200160205260409020547fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a20009074010000000000000000000000000000000000000000900460ff1615610992575050565b805460ff74010000000000000000000000000000000000000000820481169183916015916109da91750100000000000000000000000000000000000000000090910416613602565b91906101000a81548160ff021916908360ff160217905560ff161115610a2c576040517f5811df3000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff82165f81815260018301602052604080822080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055517f13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f89190a2505b50565b73ffffffffffffffffffffffffffffffffffffffff8216610afe576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f60048201526024015b60405180910390fd5b5f610b0a838333611f80565b90508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610b98576040517f64283d7b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff80861660048301526024820184905282166044820152606401610af5565b50505050565b610bb883838360405180602001604052805f815250611a27565b505050565b610bc5611f10565b5f819003610bff576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f2f000000000000000000000000000000000000000000000000000000000000008282610c2d600182613620565b818110610c3c57610c3c613633565b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610c9f576040517f3ba0191100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2002610cec8385836136a4565b604051610cf991906137ba565b60405180910390a15050565b5f6107e482611e86565b5f610d1f8263ffffffff16611e86565b505063ffffffff165f9081527fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2001602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b5f7f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab007930073ffffffffffffffffffffffffffffffffffffffff8316610dde576040517f89c62b640000000000000000000000000000000000000000000000000000000081525f6004820152602401610af5565b73ffffffffffffffffffffffffffffffffffffffff9092165f908152600390920160205250604090205490565b5f610e1461213f565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f81158015610e405750825b90505f8267ffffffffffffffff166001148015610e5c5750303b155b905081158015610e6a575080155b15610ea1576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001660011785558315610f025784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b610f766040518060400160405280600a81526020017f584d5450204e6f646573000000000000000000000000000000000000000000008152506040518060400160405280600581526020017f6e584d5450000000000000000000000000000000000000000000000000000000815250612167565b8315610fd75784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b5f6110107f000000000000000000000000000000000000000000000000000000000000000061100b6106e6565b612179565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200080549192509074010000000000000000000000000000000000000000900460ff9081169083160361108e576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805460ff7501000000000000000000000000000000000000000000909104811690831610156110e9576040517f472b0bf100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000060ff84169081029190911782556040519081527f581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff7690602001610cf9565b611168611f10565b6111778163ffffffff16611e86565b5063ffffffff81165f9081527fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200160205260409020547fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a20009074010000000000000000000000000000000000000000900460ff166111f1575050565b63ffffffff82165f908152600182016020526040902080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16905580548190601590611259907501000000000000000000000000000000000000000000900460ff1661385f565b91906101000a81548160ff021916908360ff1602179055508163ffffffff167f7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec60405160405180910390a25050565b61131361130e7f000000000000000000000000000000000000000000000000000000000000000061130960408051808201909152601a81527f786d74702e6e6f646552656769737472792e6d69677261746f72000000000000602082015290565b6121c2565b6121d5565b565b7f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab007930180546060917f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab00793009161081c90613584565b5f5f611370611f10565b73ffffffffffffffffffffffffffffffffffffffff87166113bd576040517f49e27cff00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8590036113f7576040517fbf51f54700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f839003611431576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2000805463ffffffff9060649061148490760100000000000000000000000000000000000000000000900483166001613899565b61148e91906138ac565b11156114c6576040517f957d208000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80547fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff81167601000000000000000000000000000000000000000000009182900463ffffffff9081166001019081169092021782556040516064909102935061153290889088906138c3565b6040805191829003822060808301825273ffffffffffffffffffffffffffffffffffffffff811683525f6020808501919091528251601f8b0182900482028101820184528a8152919550918301918a908a90819084018382808284375f92019190915250505090825250604080516020601f8901819004810282018101909252878152918101919088908890819084018382808284375f92018290525093909452505063ffffffff86168152600180850160209081526040928390208551815492870151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090931673ffffffffffffffffffffffffffffffffffffffff909116179190911781559184015191925082019061166690826138d2565b506060820151600282019061167b90826138d2565b5090505061168f888463ffffffff166123d7565b8173ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff168463ffffffff167f9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a48a8a8a8a6040516116f99493929190613a30565b60405180910390a4509550959350505050565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2000805460609190760100000000000000000000000000000000000000000000900463ffffffff1667ffffffffffffffff81111561176b5761176b6133ab565b6040519080825280602002602001820160405280156117a457816020015b611791612f85565b8152602001906001900390816117895790505b5091505f5b815463ffffffff760100000000000000000000000000000000000000000000909104811690821610156119b4575f6117e2826001613a61565b6117ed906064613a7d565b60408051808201825263ffffffff83168082525f90815260018781016020908152918490208451608081018652815473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000900460ff161515818501529181018054969750939592860194919390929184019161187590613584565b80601f01602080910402602001604051908101604052809291908181526020018280546118a190613584565b80156118ec5780601f106118c3576101008083540402835291602001916118ec565b820191905f5260205f20905b8154815290600101906020018083116118cf57829003601f168201915b5050505050815260200160028201805461190590613584565b80601f016020809104026020016040519081016040528092919081815260200182805461193190613584565b801561197c5780601f106119535761010080835404028352916020019161197c565b820191905f5260205f20905b81548152906001019060200180831161195f57829003601f168201915b505050505081525050815250848363ffffffff16815181106119a0576119a0613633565b6020908102919091010152506001016117a9565b505090565b6108fc338383612484565b5f6119d48263ffffffff16611e86565b505063ffffffff165f9081527fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2001602052604090205474010000000000000000000000000000000000000000900460ff1690565b611a32848484610aaa565b610b9833858585856125a4565b6060611a4a82611e86565b505f611a54612799565b90505f815111611a725760405180602001604052805f815250611a9d565b80611a7c8461284b565b604051602001611a8d929190613aba565b6040516020818303038152906040525b9392505050565b5f611b037f000000000000000000000000000000000000000000000000000000000000000061130960408051808201909152601781527f786d74702e6e6f646552656769737472792e61646d696e000000000000000000602082015290565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200080549192509073ffffffffffffffffffffffffffffffffffffffff90811690831603611b7c576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff831690811782556040517f54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d905f90a25050565b604080516080810182525f8082526020820152606091810182905281810191909152611c198263ffffffff16611e86565b5063ffffffff82165f9081527fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200160209081526040918290208251608081018452815473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000900460ff161515928101929092526001810180549293919291840191611caf90613584565b80601f0160208091040260200160405190810160405280929190818152602001828054611cdb90613584565b8015611d265780601f10611cfd57610100808354040283529160200191611d26565b820191905f5260205f20905b815481529060010190602001808311611d0957829003601f168201915b50505050508152602001600282018054611d3f90613584565b80601f0160208091040260200160405190810160405280929190818152602001828054611d6b90613584565b8015611db65780601f10611d8d57610100808354040283529160200191611db6565b820191905f5260205f20905b815481529060010190602001808311611d9957829003601f168201915b5050505050815250509050919050565b611dcf83612907565b5f819003611e09576040517fcbd6898900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff83165f8181527fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a2001602052604090207f5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e52471290600201611e6c8486836136a4565b604051611e7991906137ba565b60405180910390a2505050565b5f8181527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab0079302602052604081205473ffffffffffffffffffffffffffffffffffffffff16806107e4576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101849052602401610af5565b610bb88383836001612964565b7fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a20005473ffffffffffffffffffffffffffffffffffffffff163314611313576040517f7bfa4b9f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8281527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab007930260205260408120547f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab00793009073ffffffffffffffffffffffffffffffffffffffff90811690841615611ffa57611ffa818587612b6d565b73ffffffffffffffffffffffffffffffffffffffff81161561206f576120225f865f5f612964565b73ffffffffffffffffffffffffffffffffffffffff81165f908152600383016020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190555b73ffffffffffffffffffffffffffffffffffffffff8616156120b95773ffffffffffffffffffffffffffffffffffffffff86165f9081526003830160205260409020805460010190555b5f85815260028301602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8a811691821790925591518893918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a495945050505050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006107e4565b61216f612c1d565b6108fc8282612c5b565b5f5f6121858484612c9e565b905060ff811115611a9d576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f611a9d6121d08484612c9e565b612d31565b73ffffffffffffffffffffffffffffffffffffffff8116612222576040517f0d626a3200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff8216907fa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098905f90a25f5f8273ffffffffffffffffffffffffffffffffffffffff166040515f60405180830381855af49150503d805f81146122b6576040519150601f19603f3d011682016040523d82523d5f602084013e6122bb565b606091505b5091509150816122fb5782816040517f68b0b16b000000000000000000000000000000000000000000000000000000008152600401610af5929190613ace565b805115801561231f575073ffffffffffffffffffffffffffffffffffffffff83163b155b1561236e576040517f626c416100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610af5565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff167fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b60405160405180910390a2505050565b73ffffffffffffffffffffffffffffffffffffffff8216612426576040517f64a0ae920000000000000000000000000000000000000000000000000000000081525f6004820152602401610af5565b5f61243283835f611f80565b905073ffffffffffffffffffffffffffffffffffffffff811615610bb8576040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610af5565b7f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab007930073ffffffffffffffffffffffffffffffffffffffff831661250a576040517f5b08ba1800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610af5565b73ffffffffffffffffffffffffffffffffffffffff8481165f81815260058401602090815260408083209488168084529482529182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001687151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a350505050565b73ffffffffffffffffffffffffffffffffffffffff83163b15610fd7576040517f150b7a0200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84169063150b7a0290612619908890889087908790600401613afc565b6020604051808303815f875af1925050508015612671575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261266e91810190613b56565b60015b6126fe573d80801561269e576040519150601f19603f3d011682016040523d82523d5f602084013e6126a3565b606091505b5080515f036126f6576040517f64a0ae9200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610af5565b805181602001fd5b7fffffffff0000000000000000000000000000000000000000000000000000000081167f150b7a020000000000000000000000000000000000000000000000000000000014612791576040517f64a0ae9200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610af5565b505050505050565b60607fd48713bc7b5e2644bcb4e26ace7d67dc9027725a9a1ee11596536cc6096a200060020180546127ca90613584565b80601f01602080910402602001604051908101604052809291908181526020018280546127f690613584565b80156128415780601f1061281857610100808354040283529160200191612841565b820191905f5260205f20905b81548152906001019060200180831161282457829003601f168201915b5050505050905090565b60605f61285783612d84565b60010190505f8167ffffffffffffffff811115612876576128766133ab565b6040519080825280601f01601f1916602001820160405280156128a0576020820181803683370190505b5090508181016020015b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85049450846128aa57509392505050565b3361291763ffffffff8316611e86565b73ffffffffffffffffffffffffffffffffffffffff1614610aa7576040517fd08a05d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab007930081806129a6575073ffffffffffffffffffffffffffffffffffffffff831615155b15612b18575f6129b585611e86565b905073ffffffffffffffffffffffffffffffffffffffff841615801590612a0857508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b8015612a65575073ffffffffffffffffffffffffffffffffffffffff8082165f9081527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab0079305602090815260408083209388168352929052205460ff16155b15612ab4576040517fa9fbf51f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610af5565b8215612b1657848673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b5f93845260040160205250506040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b612b78838383612e65565b610bb85773ffffffffffffffffffffffffffffffffffffffff8316612bcc576040517f7e27328900000000000000000000000000000000000000000000000000000000815260048101829052602401610af5565b6040517f177e802f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015260248101829052604401610af5565b612c25612f67565b611313576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612c63612c1d565b7f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab007930080612c8f84826138d2565b5060018101610b9883826138d2565b6040517fd6d7d5250000000000000000000000000000000000000000000000000000000081525f9073ffffffffffffffffffffffffffffffffffffffff84169063d6d7d52590612cf290859060040161302e565b602060405180830381865afa158015612d0d573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611a9d9190613b71565b5f73ffffffffffffffffffffffffffffffffffffffff821115612d80576040517f37f4f14800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090565b5f807a184f03e93ff9f4daa797ed6e38ed64bf6a1f0100000000000000008310612dcc577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef81000000008310612df8576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310612e1657662386f26fc10000830492506010015b6305f5e1008310612e2e576305f5e100830492506008015b6127108310612e4257612710830492506004015b60648310612e54576064830492506002015b600a83106107e45760010192915050565b5f73ffffffffffffffffffffffffffffffffffffffff831615801590612f5f57508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161480612f10575073ffffffffffffffffffffffffffffffffffffffff8085165f9081527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab0079305602090815260408083209387168352929052205460ff165b80612f5f57505f8281527f80bb2b638cc20bc4d0a60d66940f3ab4a00c1d7b313497ca82fb0b4ab0079304602052604090205473ffffffffffffffffffffffffffffffffffffffff8481169116145b949350505050565b5f612f7061213f565b5468010000000000000000900460ff16919050565b60405180604001604052805f63ffffffff168152602001612fdd60405180608001604052805f73ffffffffffffffffffffffffffffffffffffffff1681526020015f1515815260200160608152602001606081525090565b905290565b5f81518084528060208401602086015e5f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081525f611a9d6020830184612fe2565b7fffffffff0000000000000000000000000000000000000000000000000000000081168114610aa7575f5ffd5b5f6020828403121561307d575f5ffd5b8135611a9d81613040565b5f60208284031215613098575f5ffd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146130c2575f5ffd5b919050565b5f5f604083850312156130d8575f5ffd5b6130e18361309f565b946020939093013593505050565b803563ffffffff811681146130c2575f5ffd5b5f60208284031215613112575f5ffd5b611a9d826130ef565b5f5f5f6060848603121561312d575f5ffd5b6131368461309f565b92506131446020850161309f565b929592945050506040919091013590565b5f5f83601f840112613165575f5ffd5b50813567ffffffffffffffff81111561317c575f5ffd5b602083019150836020828501011115613193575f5ffd5b9250929050565b5f5f602083850312156131ab575f5ffd5b823567ffffffffffffffff8111156131c1575f5ffd5b6131cd85828601613155565b90969095509350505050565b5f602082840312156131e9575f5ffd5b611a9d8261309f565b5f5f5f5f5f60608688031215613206575f5ffd5b61320f8661309f565b9450602086013567ffffffffffffffff81111561322a575f5ffd5b61323688828901613155565b909550935050604086013567ffffffffffffffff811115613255575f5ffd5b61326188828901613155565b969995985093965092949392505050565b73ffffffffffffffffffffffffffffffffffffffff81511682526020810151151560208301525f6040820151608060408501526132b26080850182612fe2565b9050606083015184820360608601526132cb8282612fe2565b95945050505050565b5f602082016020835280845180835260408501915060408160051b8601019250602086015f5b82811015613366577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0878603018452815163ffffffff815116865260208101519050604060208701526133506040870182613272565b95505060209384019391909101906001016132fa565b50929695505050505050565b5f5f60408385031215613383575f5ffd5b61338c8361309f565b9150602083013580151581146133a0575f5ffd5b809150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f5f5f608085870312156133eb575f5ffd5b6133f48561309f565b93506134026020860161309f565b925060408501359150606085013567ffffffffffffffff811115613424575f5ffd5b8501601f81018713613434575f5ffd5b803567ffffffffffffffff81111561344e5761344e6133ab565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156134ba576134ba6133ab565b6040528181528282016020018910156134d1575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b602081525f611a9d6020830184613272565b5f5f60408385031215613515575f5ffd5b61351e8361309f565b915061352c6020840161309f565b90509250929050565b5f5f5f60408486031215613547575f5ffd5b613550846130ef565b9250602084013567ffffffffffffffff81111561356b575f5ffd5b61357786828701613155565b9497909650939450505050565b600181811c9082168061359857607f821691505b6020821081036135cf577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f60ff821660ff8103613617576136176135d5565b60010192915050565b818103818111156107e4576107e46135d5565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b601f821115610bb857805f5260205f20601f840160051c810160208510156136855750805b601f840160051c820191505b81811015610fd7575f8155600101613691565b67ffffffffffffffff8311156136bc576136bc6133ab565b6136d0836136ca8354613584565b83613660565b5f601f841160018114613720575f85156136ea5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355610fd7565b5f838152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08716915b8281101561376d578685013582556020948501946001909201910161374d565b50868210156137a8577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b602081525f5f83546137cb81613584565b806020860152600182165f81146137e9576001811461382357613854565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083166040870152604082151560051b8701019350613854565b865f5260205f205f5b8381101561384b5781548882016040015260019091019060200161382c565b87016040019450505b509195945050505050565b5f60ff821680613871576138716135d5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b808201808211156107e4576107e46135d5565b80820281158282048414176107e4576107e46135d5565b818382375f9101908152919050565b815167ffffffffffffffff8111156138ec576138ec6133ab565b613900816138fa8454613584565b84613660565b6020601f821160018114613951575f831561391b5750848201515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600385901b1c1916600184901b178455610fd7565b5f848152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08516915b8281101561399e578785015182556020948501946001909201910161397e565b50848210156139da57868401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b60f8161c191681555b50505050600190811b01905550565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b604081525f613a436040830186886139e9565b8281036020840152613a568185876139e9565b979650505050505050565b63ffffffff81811683821601908111156107e4576107e46135d5565b63ffffffff8181168382160290811690818114613a9c57613a9c6135d5565b5092915050565b5f81518060208401855e5f93019283525090919050565b5f612f5f613ac88386613aa3565b84613aa3565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201525f612f5f6040830184612fe2565b73ffffffffffffffffffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff84166020820152826040820152608060608201525f613b4c6080830184612fe2565b9695505050505050565b5f60208284031215613b66575f5ffd5b8151611a9d81613040565b5f60208284031215613b81575f5ffd5b505191905056fe786d74702e6e6f646552656769737472792e6d617843616e6f6e6963616c4e6f646573a2646970667358221220b49872ac43291720edf1deb60c7f40a22d6c8001cb2a8580b05ea2fb325aaf5964736f6c634300081c0033",
}

// NodeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeRegistryMetaData.ABI instead.
var NodeRegistryABI = NodeRegistryMetaData.ABI

// NodeRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodeRegistryMetaData.Bin instead.
var NodeRegistryBin = NodeRegistryMetaData.Bin

// DeployNodeRegistry deploys a new Ethereum contract, binding an instance of NodeRegistry to it.
func DeployNodeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, parameterRegistry_ common.Address) (common.Address, *types.Transaction, *NodeRegistry, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodeRegistryBin), backend, parameterRegistry_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NodeRegistry is an auto generated Go binding around an Ethereum contract.
type NodeRegistry struct {
	NodeRegistryCaller     // Read-only binding to the contract
	NodeRegistryTransactor // Write-only binding to the contract
	NodeRegistryFilterer   // Log filterer for contract events
}

// NodeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeRegistrySession struct {
	Contract     *NodeRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeRegistryCallerSession struct {
	Contract *NodeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// NodeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeRegistryTransactorSession struct {
	Contract     *NodeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NodeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeRegistryRaw struct {
	Contract *NodeRegistry // Generic contract binding to access the raw methods on
}

// NodeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeRegistryCallerRaw struct {
	Contract *NodeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// NodeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeRegistryTransactorRaw struct {
	Contract *NodeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeRegistry creates a new instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistry(address common.Address, backend bind.ContractBackend) (*NodeRegistry, error) {
	contract, err := bindNodeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeRegistry{NodeRegistryCaller: NodeRegistryCaller{contract: contract}, NodeRegistryTransactor: NodeRegistryTransactor{contract: contract}, NodeRegistryFilterer: NodeRegistryFilterer{contract: contract}}, nil
}

// NewNodeRegistryCaller creates a new read-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryCaller(address common.Address, caller bind.ContractCaller) (*NodeRegistryCaller, error) {
	contract, err := bindNodeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryCaller{contract: contract}, nil
}

// NewNodeRegistryTransactor creates a new write-only instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeRegistryTransactor, error) {
	contract, err := bindNodeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransactor{contract: contract}, nil
}

// NewNodeRegistryFilterer creates a new log filterer instance of NodeRegistry, bound to a specific deployed contract.
func NewNodeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeRegistryFilterer, error) {
	contract, err := bindNodeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryFilterer{contract: contract}, nil
}

// bindNodeRegistry binds a generic wrapper to an already deployed contract.
func bindNodeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodeRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.NodeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.NodeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeRegistry *NodeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeRegistry *NodeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeRegistry.Contract.contract.Transact(opts, method, params...)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistryCaller) NODEINCREMENT(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "NODE_INCREMENT")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistrySession) NODEINCREMENT() (uint32, error) {
	return _NodeRegistry.Contract.NODEINCREMENT(&_NodeRegistry.CallOpts)
}

// NODEINCREMENT is a free data retrieval call binding the contract method 0xfd667d1e.
//
// Solidity: function NODE_INCREMENT() view returns(uint32)
func (_NodeRegistry *NodeRegistryCallerSession) NODEINCREMENT() (uint32, error) {
	return _NodeRegistry.Contract.NODEINCREMENT(&_NodeRegistry.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address admin_)
func (_NodeRegistry *NodeRegistryCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address admin_)
func (_NodeRegistry *NodeRegistrySession) Admin() (common.Address, error) {
	return _NodeRegistry.Contract.Admin(&_NodeRegistry.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address admin_)
func (_NodeRegistry *NodeRegistryCallerSession) Admin() (common.Address, error) {
	return _NodeRegistry.Contract.Admin(&_NodeRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCaller) AdminParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "adminParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistrySession) AdminParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.AdminParameterKey(&_NodeRegistry.CallOpts)
}

// AdminParameterKey is a free data retrieval call binding the contract method 0x9f40b625.
//
// Solidity: function adminParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCallerSession) AdminParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.AdminParameterKey(&_NodeRegistry.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistryCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistrySession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeRegistry.Contract.BalanceOf(&_NodeRegistry.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NodeRegistry *NodeRegistryCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NodeRegistry.Contract.BalanceOf(&_NodeRegistry.CallOpts, owner)
}

// CanonicalNodesCount is a free data retrieval call binding the contract method 0xc9c02a02.
//
// Solidity: function canonicalNodesCount() view returns(uint8 canonicalNodesCount_)
func (_NodeRegistry *NodeRegistryCaller) CanonicalNodesCount(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "canonicalNodesCount")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CanonicalNodesCount is a free data retrieval call binding the contract method 0xc9c02a02.
//
// Solidity: function canonicalNodesCount() view returns(uint8 canonicalNodesCount_)
func (_NodeRegistry *NodeRegistrySession) CanonicalNodesCount() (uint8, error) {
	return _NodeRegistry.Contract.CanonicalNodesCount(&_NodeRegistry.CallOpts)
}

// CanonicalNodesCount is a free data retrieval call binding the contract method 0xc9c02a02.
//
// Solidity: function canonicalNodesCount() view returns(uint8 canonicalNodesCount_)
func (_NodeRegistry *NodeRegistryCallerSession) CanonicalNodesCount() (uint8, error) {
	return _NodeRegistry.Contract.CanonicalNodesCount(&_NodeRegistry.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint32,(address,bool,bytes,string))[] allNodes_)
func (_NodeRegistry *NodeRegistryCaller) GetAllNodes(opts *bind.CallOpts) ([]INodeRegistryNodeWithId, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getAllNodes")

	if err != nil {
		return *new([]INodeRegistryNodeWithId), err
	}

	out0 := *abi.ConvertType(out[0], new([]INodeRegistryNodeWithId)).(*[]INodeRegistryNodeWithId)

	return out0, err

}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint32,(address,bool,bytes,string))[] allNodes_)
func (_NodeRegistry *NodeRegistrySession) GetAllNodes() ([]INodeRegistryNodeWithId, error) {
	return _NodeRegistry.Contract.GetAllNodes(&_NodeRegistry.CallOpts)
}

// GetAllNodes is a free data retrieval call binding the contract method 0xa1174e7d.
//
// Solidity: function getAllNodes() view returns((uint32,(address,bool,bytes,string))[] allNodes_)
func (_NodeRegistry *NodeRegistryCallerSession) GetAllNodes() ([]INodeRegistryNodeWithId, error) {
	return _NodeRegistry.Contract.GetAllNodes(&_NodeRegistry.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint32 nodeCount_)
func (_NodeRegistry *NodeRegistryCaller) GetAllNodesCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getAllNodesCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint32 nodeCount_)
func (_NodeRegistry *NodeRegistrySession) GetAllNodesCount() (uint32, error) {
	return _NodeRegistry.Contract.GetAllNodesCount(&_NodeRegistry.CallOpts)
}

// GetAllNodesCount is a free data retrieval call binding the contract method 0x50d0215f.
//
// Solidity: function getAllNodesCount() view returns(uint32 nodeCount_)
func (_NodeRegistry *NodeRegistryCallerSession) GetAllNodesCount() (uint32, error) {
	return _NodeRegistry.Contract.GetAllNodesCount(&_NodeRegistry.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistrySession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.GetApproved(&_NodeRegistry.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.GetApproved(&_NodeRegistry.CallOpts, tokenId)
}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xad03d0a5.
//
// Solidity: function getIsCanonicalNode(uint32 nodeId_) view returns(bool isCanonicalNode_)
func (_NodeRegistry *NodeRegistryCaller) GetIsCanonicalNode(opts *bind.CallOpts, nodeId_ uint32) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getIsCanonicalNode", nodeId_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xad03d0a5.
//
// Solidity: function getIsCanonicalNode(uint32 nodeId_) view returns(bool isCanonicalNode_)
func (_NodeRegistry *NodeRegistrySession) GetIsCanonicalNode(nodeId_ uint32) (bool, error) {
	return _NodeRegistry.Contract.GetIsCanonicalNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetIsCanonicalNode is a free data retrieval call binding the contract method 0xad03d0a5.
//
// Solidity: function getIsCanonicalNode(uint32 nodeId_) view returns(bool isCanonicalNode_)
func (_NodeRegistry *NodeRegistryCallerSession) GetIsCanonicalNode(nodeId_ uint32) (bool, error) {
	return _NodeRegistry.Contract.GetIsCanonicalNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetNode is a free data retrieval call binding the contract method 0xe06f876f.
//
// Solidity: function getNode(uint32 nodeId_) view returns((address,bool,bytes,string) node_)
func (_NodeRegistry *NodeRegistryCaller) GetNode(opts *bind.CallOpts, nodeId_ uint32) (INodeRegistryNode, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getNode", nodeId_)

	if err != nil {
		return *new(INodeRegistryNode), err
	}

	out0 := *abi.ConvertType(out[0], new(INodeRegistryNode)).(*INodeRegistryNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0xe06f876f.
//
// Solidity: function getNode(uint32 nodeId_) view returns((address,bool,bytes,string) node_)
func (_NodeRegistry *NodeRegistrySession) GetNode(nodeId_ uint32) (INodeRegistryNode, error) {
	return _NodeRegistry.Contract.GetNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetNode is a free data retrieval call binding the contract method 0xe06f876f.
//
// Solidity: function getNode(uint32 nodeId_) view returns((address,bool,bytes,string) node_)
func (_NodeRegistry *NodeRegistryCallerSession) GetNode(nodeId_ uint32) (INodeRegistryNode, error) {
	return _NodeRegistry.Contract.GetNode(&_NodeRegistry.CallOpts, nodeId_)
}

// GetSigner is a free data retrieval call binding the contract method 0x68501a3e.
//
// Solidity: function getSigner(uint32 nodeId_) view returns(address signer_)
func (_NodeRegistry *NodeRegistryCaller) GetSigner(opts *bind.CallOpts, nodeId_ uint32) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "getSigner", nodeId_)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSigner is a free data retrieval call binding the contract method 0x68501a3e.
//
// Solidity: function getSigner(uint32 nodeId_) view returns(address signer_)
func (_NodeRegistry *NodeRegistrySession) GetSigner(nodeId_ uint32) (common.Address, error) {
	return _NodeRegistry.Contract.GetSigner(&_NodeRegistry.CallOpts, nodeId_)
}

// GetSigner is a free data retrieval call binding the contract method 0x68501a3e.
//
// Solidity: function getSigner(uint32 nodeId_) view returns(address signer_)
func (_NodeRegistry *NodeRegistryCallerSession) GetSigner(nodeId_ uint32) (common.Address, error) {
	return _NodeRegistry.Contract.GetSigner(&_NodeRegistry.CallOpts, nodeId_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_NodeRegistry *NodeRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_NodeRegistry *NodeRegistrySession) Implementation() (common.Address, error) {
	return _NodeRegistry.Contract.Implementation(&_NodeRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address implementation_)
func (_NodeRegistry *NodeRegistryCallerSession) Implementation() (common.Address, error) {
	return _NodeRegistry.Contract.Implementation(&_NodeRegistry.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistryCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistrySession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodeRegistry.Contract.IsApprovedForAll(&_NodeRegistry.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NodeRegistry *NodeRegistryCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NodeRegistry.Contract.IsApprovedForAll(&_NodeRegistry.CallOpts, owner, operator)
}

// MaxCanonicalNodes is a free data retrieval call binding the contract method 0xc18e273d.
//
// Solidity: function maxCanonicalNodes() view returns(uint8 maxCanonicalNodes_)
func (_NodeRegistry *NodeRegistryCaller) MaxCanonicalNodes(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "maxCanonicalNodes")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MaxCanonicalNodes is a free data retrieval call binding the contract method 0xc18e273d.
//
// Solidity: function maxCanonicalNodes() view returns(uint8 maxCanonicalNodes_)
func (_NodeRegistry *NodeRegistrySession) MaxCanonicalNodes() (uint8, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodes(&_NodeRegistry.CallOpts)
}

// MaxCanonicalNodes is a free data retrieval call binding the contract method 0xc18e273d.
//
// Solidity: function maxCanonicalNodes() view returns(uint8 maxCanonicalNodes_)
func (_NodeRegistry *NodeRegistryCallerSession) MaxCanonicalNodes() (uint8, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodes(&_NodeRegistry.CallOpts)
}

// MaxCanonicalNodesParameterKey is a free data retrieval call binding the contract method 0x0124b882.
//
// Solidity: function maxCanonicalNodesParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCaller) MaxCanonicalNodesParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "maxCanonicalNodesParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MaxCanonicalNodesParameterKey is a free data retrieval call binding the contract method 0x0124b882.
//
// Solidity: function maxCanonicalNodesParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistrySession) MaxCanonicalNodesParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodesParameterKey(&_NodeRegistry.CallOpts)
}

// MaxCanonicalNodesParameterKey is a free data retrieval call binding the contract method 0x0124b882.
//
// Solidity: function maxCanonicalNodesParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCallerSession) MaxCanonicalNodesParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MaxCanonicalNodesParameterKey(&_NodeRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCaller) MigratorParameterKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "migratorParameterKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistrySession) MigratorParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MigratorParameterKey(&_NodeRegistry.CallOpts)
}

// MigratorParameterKey is a free data retrieval call binding the contract method 0x8aab82ba.
//
// Solidity: function migratorParameterKey() pure returns(bytes key_)
func (_NodeRegistry *NodeRegistryCallerSession) MigratorParameterKey() ([]byte, error) {
	return _NodeRegistry.Contract.MigratorParameterKey(&_NodeRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistryCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistrySession) Name() (string, error) {
	return _NodeRegistry.Contract.Name(&_NodeRegistry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) Name() (string, error) {
	return _NodeRegistry.Contract.Name(&_NodeRegistry.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistrySession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.OwnerOf(&_NodeRegistry.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NodeRegistry.Contract.OwnerOf(&_NodeRegistry.CallOpts, tokenId)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_NodeRegistry *NodeRegistryCaller) ParameterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "parameterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_NodeRegistry *NodeRegistrySession) ParameterRegistry() (common.Address, error) {
	return _NodeRegistry.Contract.ParameterRegistry(&_NodeRegistry.CallOpts)
}

// ParameterRegistry is a free data retrieval call binding the contract method 0x0723499e.
//
// Solidity: function parameterRegistry() view returns(address)
func (_NodeRegistry *NodeRegistryCallerSession) ParameterRegistry() (common.Address, error) {
	return _NodeRegistry.Contract.ParameterRegistry(&_NodeRegistry.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodeRegistry *NodeRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodeRegistry *NodeRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodeRegistry.Contract.SupportsInterface(&_NodeRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NodeRegistry *NodeRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NodeRegistry.Contract.SupportsInterface(&_NodeRegistry.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistryCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistrySession) Symbol() (string, error) {
	return _NodeRegistry.Contract.Symbol(&_NodeRegistry.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) Symbol() (string, error) {
	return _NodeRegistry.Contract.Symbol(&_NodeRegistry.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistryCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NodeRegistry.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistrySession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodeRegistry.Contract.TokenURI(&_NodeRegistry.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_NodeRegistry *NodeRegistryCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _NodeRegistry.Contract.TokenURI(&_NodeRegistry.CallOpts, tokenId)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address owner_, bytes signingPublicKey_, string httpAddress_) returns(uint32 nodeId_, address signer_)
func (_NodeRegistry *NodeRegistryTransactor) AddNode(opts *bind.TransactOpts, owner_ common.Address, signingPublicKey_ []byte, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "addNode", owner_, signingPublicKey_, httpAddress_)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address owner_, bytes signingPublicKey_, string httpAddress_) returns(uint32 nodeId_, address signer_)
func (_NodeRegistry *NodeRegistrySession) AddNode(owner_ common.Address, signingPublicKey_ []byte, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddNode(&_NodeRegistry.TransactOpts, owner_, signingPublicKey_, httpAddress_)
}

// AddNode is a paid mutator transaction binding the contract method 0xa0eae81d.
//
// Solidity: function addNode(address owner_, bytes signingPublicKey_, string httpAddress_) returns(uint32 nodeId_, address signer_)
func (_NodeRegistry *NodeRegistryTransactorSession) AddNode(owner_ common.Address, signingPublicKey_ []byte, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddNode(&_NodeRegistry.TransactOpts, owner_, signingPublicKey_, httpAddress_)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0x236b6eb8.
//
// Solidity: function addToNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactor) AddToNetwork(opts *bind.TransactOpts, nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "addToNetwork", nodeId_)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0x236b6eb8.
//
// Solidity: function addToNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistrySession) AddToNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddToNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// AddToNetwork is a paid mutator transaction binding the contract method 0x236b6eb8.
//
// Solidity: function addToNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) AddToNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.AddToNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.Approve(&_NodeRegistry.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.Approve(&_NodeRegistry.TransactOpts, to, tokenId)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NodeRegistry *NodeRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NodeRegistry *NodeRegistrySession) Initialize() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Initialize(&_NodeRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Initialize(&_NodeRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_NodeRegistry *NodeRegistryTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_NodeRegistry *NodeRegistrySession) Migrate() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Migrate(&_NodeRegistry.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) Migrate() (*types.Transaction, error) {
	return _NodeRegistry.Contract.Migrate(&_NodeRegistry.TransactOpts)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0x8cf20c68.
//
// Solidity: function removeFromNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactor) RemoveFromNetwork(opts *bind.TransactOpts, nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "removeFromNetwork", nodeId_)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0x8cf20c68.
//
// Solidity: function removeFromNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistrySession) RemoveFromNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RemoveFromNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// RemoveFromNetwork is a paid mutator transaction binding the contract method 0x8cf20c68.
//
// Solidity: function removeFromNetwork(uint32 nodeId_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) RemoveFromNetwork(nodeId_ uint32) (*types.Transaction, error) {
	return _NodeRegistry.Contract.RemoveFromNetwork(&_NodeRegistry.TransactOpts, nodeId_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistryTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistrySession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom0(&_NodeRegistry.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SafeTransferFrom0(&_NodeRegistry.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistrySession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetApprovalForAll(&_NodeRegistry.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetApprovalForAll(&_NodeRegistry.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetBaseURI(opts *bind.TransactOpts, baseURI_ string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setBaseURI", baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_NodeRegistry *NodeRegistrySession) SetBaseURI(baseURI_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetBaseURI(&_NodeRegistry.TransactOpts, baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetBaseURI(baseURI_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetBaseURI(&_NodeRegistry.TransactOpts, baseURI_)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xf84ce8b9.
//
// Solidity: function setHttpAddress(uint32 nodeId_, string httpAddress_) returns()
func (_NodeRegistry *NodeRegistryTransactor) SetHttpAddress(opts *bind.TransactOpts, nodeId_ uint32, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "setHttpAddress", nodeId_, httpAddress_)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xf84ce8b9.
//
// Solidity: function setHttpAddress(uint32 nodeId_, string httpAddress_) returns()
func (_NodeRegistry *NodeRegistrySession) SetHttpAddress(nodeId_ uint32, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetHttpAddress(&_NodeRegistry.TransactOpts, nodeId_, httpAddress_)
}

// SetHttpAddress is a paid mutator transaction binding the contract method 0xf84ce8b9.
//
// Solidity: function setHttpAddress(uint32 nodeId_, string httpAddress_) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) SetHttpAddress(nodeId_ uint32, httpAddress_ string) (*types.Transaction, error) {
	return _NodeRegistry.Contract.SetHttpAddress(&_NodeRegistry.TransactOpts, nodeId_, httpAddress_)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistrySession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.TransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NodeRegistry *NodeRegistryTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NodeRegistry.Contract.TransferFrom(&_NodeRegistry.TransactOpts, from, to, tokenId)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xd3b2f598.
//
// Solidity: function updateAdmin() returns()
func (_NodeRegistry *NodeRegistryTransactor) UpdateAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "updateAdmin")
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xd3b2f598.
//
// Solidity: function updateAdmin() returns()
func (_NodeRegistry *NodeRegistrySession) UpdateAdmin() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateAdmin(&_NodeRegistry.TransactOpts)
}

// UpdateAdmin is a paid mutator transaction binding the contract method 0xd3b2f598.
//
// Solidity: function updateAdmin() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) UpdateAdmin() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateAdmin(&_NodeRegistry.TransactOpts)
}

// UpdateMaxCanonicalNodes is a paid mutator transaction binding the contract method 0x82a5cfc3.
//
// Solidity: function updateMaxCanonicalNodes() returns()
func (_NodeRegistry *NodeRegistryTransactor) UpdateMaxCanonicalNodes(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeRegistry.contract.Transact(opts, "updateMaxCanonicalNodes")
}

// UpdateMaxCanonicalNodes is a paid mutator transaction binding the contract method 0x82a5cfc3.
//
// Solidity: function updateMaxCanonicalNodes() returns()
func (_NodeRegistry *NodeRegistrySession) UpdateMaxCanonicalNodes() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateMaxCanonicalNodes(&_NodeRegistry.TransactOpts)
}

// UpdateMaxCanonicalNodes is a paid mutator transaction binding the contract method 0x82a5cfc3.
//
// Solidity: function updateMaxCanonicalNodes() returns()
func (_NodeRegistry *NodeRegistryTransactorSession) UpdateMaxCanonicalNodes() (*types.Transaction, error) {
	return _NodeRegistry.Contract.UpdateMaxCanonicalNodes(&_NodeRegistry.TransactOpts)
}

// NodeRegistryAdminUpdatedIterator is returned from FilterAdminUpdated and is used to iterate over the raw logs and unpacked data for AdminUpdated events raised by the NodeRegistry contract.
type NodeRegistryAdminUpdatedIterator struct {
	Event *NodeRegistryAdminUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryAdminUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryAdminUpdated)
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
		it.Event = new(NodeRegistryAdminUpdated)
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
func (it *NodeRegistryAdminUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryAdminUpdated represents a AdminUpdated event raised by the NodeRegistry contract.
type NodeRegistryAdminUpdated struct {
	Admin common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAdminUpdated is a free log retrieval operation binding the contract event 0x54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d.
//
// Solidity: event AdminUpdated(address indexed admin)
func (_NodeRegistry *NodeRegistryFilterer) FilterAdminUpdated(opts *bind.FilterOpts, admin []common.Address) (*NodeRegistryAdminUpdatedIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "AdminUpdated", adminRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryAdminUpdatedIterator{contract: _NodeRegistry.contract, event: "AdminUpdated", logs: logs, sub: sub}, nil
}

// WatchAdminUpdated is a free log subscription operation binding the contract event 0x54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d.
//
// Solidity: event AdminUpdated(address indexed admin)
func (_NodeRegistry *NodeRegistryFilterer) WatchAdminUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryAdminUpdated, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "AdminUpdated", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryAdminUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "AdminUpdated", log); err != nil {
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

// ParseAdminUpdated is a log parse operation binding the contract event 0x54e4612788f90384e6843298d7854436f3a585b2c3831ab66abf1de63bfa6c2d.
//
// Solidity: event AdminUpdated(address indexed admin)
func (_NodeRegistry *NodeRegistryFilterer) ParseAdminUpdated(log types.Log) (*NodeRegistryAdminUpdated, error) {
	event := new(NodeRegistryAdminUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "AdminUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the NodeRegistry contract.
type NodeRegistryApprovalIterator struct {
	Event *NodeRegistryApproval // Event containing the contract specifics and raw log

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
func (it *NodeRegistryApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryApproval)
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
		it.Event = new(NodeRegistryApproval)
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
func (it *NodeRegistryApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryApproval represents a Approval event raised by the NodeRegistry contract.
type NodeRegistryApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NodeRegistryApprovalIterator, error) {

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

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryApprovalIterator{contract: _NodeRegistry.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NodeRegistryApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryApproval)
				if err := _NodeRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseApproval(log types.Log) (*NodeRegistryApproval, error) {
	event := new(NodeRegistryApproval)
	if err := _NodeRegistry.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the NodeRegistry contract.
type NodeRegistryApprovalForAllIterator struct {
	Event *NodeRegistryApprovalForAll // Event containing the contract specifics and raw log

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
func (it *NodeRegistryApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryApprovalForAll)
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
		it.Event = new(NodeRegistryApprovalForAll)
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
func (it *NodeRegistryApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryApprovalForAll represents a ApprovalForAll event raised by the NodeRegistry contract.
type NodeRegistryApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*NodeRegistryApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryApprovalForAllIterator{contract: _NodeRegistry.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NodeRegistry *NodeRegistryFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NodeRegistryApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryApprovalForAll)
				if err := _NodeRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseApprovalForAll(log types.Log) (*NodeRegistryApprovalForAll, error) {
	event := new(NodeRegistryApprovalForAll)
	if err := _NodeRegistry.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryBaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the NodeRegistry contract.
type NodeRegistryBaseURIUpdatedIterator struct {
	Event *NodeRegistryBaseURIUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryBaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryBaseURIUpdated)
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
		it.Event = new(NodeRegistryBaseURIUpdated)
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
func (it *NodeRegistryBaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryBaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryBaseURIUpdated represents a BaseURIUpdated event raised by the NodeRegistry contract.
type NodeRegistryBaseURIUpdated struct {
	BaseURI string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string baseURI)
func (_NodeRegistry *NodeRegistryFilterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*NodeRegistryBaseURIUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryBaseURIUpdatedIterator{contract: _NodeRegistry.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string baseURI)
func (_NodeRegistry *NodeRegistryFilterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryBaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryBaseURIUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
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
// Solidity: event BaseURIUpdated(string baseURI)
func (_NodeRegistry *NodeRegistryFilterer) ParseBaseURIUpdated(log types.Log) (*NodeRegistryBaseURIUpdated, error) {
	event := new(NodeRegistryBaseURIUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryHttpAddressUpdatedIterator is returned from FilterHttpAddressUpdated and is used to iterate over the raw logs and unpacked data for HttpAddressUpdated events raised by the NodeRegistry contract.
type NodeRegistryHttpAddressUpdatedIterator struct {
	Event *NodeRegistryHttpAddressUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryHttpAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryHttpAddressUpdated)
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
		it.Event = new(NodeRegistryHttpAddressUpdated)
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
func (it *NodeRegistryHttpAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryHttpAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryHttpAddressUpdated represents a HttpAddressUpdated event raised by the NodeRegistry contract.
type NodeRegistryHttpAddressUpdated struct {
	NodeId      uint32
	HttpAddress string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterHttpAddressUpdated is a free log retrieval operation binding the contract event 0x5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e524712.
//
// Solidity: event HttpAddressUpdated(uint32 indexed nodeId, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) FilterHttpAddressUpdated(opts *bind.FilterOpts, nodeId []uint32) (*NodeRegistryHttpAddressUpdatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryHttpAddressUpdatedIterator{contract: _NodeRegistry.contract, event: "HttpAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchHttpAddressUpdated is a free log subscription operation binding the contract event 0x5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e524712.
//
// Solidity: event HttpAddressUpdated(uint32 indexed nodeId, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) WatchHttpAddressUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryHttpAddressUpdated, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "HttpAddressUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryHttpAddressUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
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

// ParseHttpAddressUpdated is a log parse operation binding the contract event 0x5698a22512088407e91d125d2eb43d829d9694a71f664ab0dc2aea3a8e524712.
//
// Solidity: event HttpAddressUpdated(uint32 indexed nodeId, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) ParseHttpAddressUpdated(log types.Log) (*NodeRegistryHttpAddressUpdated, error) {
	event := new(NodeRegistryHttpAddressUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "HttpAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NodeRegistry contract.
type NodeRegistryInitializedIterator struct {
	Event *NodeRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *NodeRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryInitialized)
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
		it.Event = new(NodeRegistryInitialized)
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
func (it *NodeRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryInitialized represents a Initialized event raised by the NodeRegistry contract.
type NodeRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NodeRegistry *NodeRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*NodeRegistryInitializedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryInitializedIterator{contract: _NodeRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NodeRegistry *NodeRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NodeRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryInitialized)
				if err := _NodeRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseInitialized(log types.Log) (*NodeRegistryInitialized, error) {
	event := new(NodeRegistryInitialized)
	if err := _NodeRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryMaxCanonicalNodesUpdatedIterator is returned from FilterMaxCanonicalNodesUpdated and is used to iterate over the raw logs and unpacked data for MaxCanonicalNodesUpdated events raised by the NodeRegistry contract.
type NodeRegistryMaxCanonicalNodesUpdatedIterator struct {
	Event *NodeRegistryMaxCanonicalNodesUpdated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryMaxCanonicalNodesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryMaxCanonicalNodesUpdated)
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
		it.Event = new(NodeRegistryMaxCanonicalNodesUpdated)
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
func (it *NodeRegistryMaxCanonicalNodesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryMaxCanonicalNodesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryMaxCanonicalNodesUpdated represents a MaxCanonicalNodesUpdated event raised by the NodeRegistry contract.
type NodeRegistryMaxCanonicalNodesUpdated struct {
	MaxCanonicalNodes uint8
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterMaxCanonicalNodesUpdated is a free log retrieval operation binding the contract event 0x581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff76.
//
// Solidity: event MaxCanonicalNodesUpdated(uint8 maxCanonicalNodes)
func (_NodeRegistry *NodeRegistryFilterer) FilterMaxCanonicalNodesUpdated(opts *bind.FilterOpts) (*NodeRegistryMaxCanonicalNodesUpdatedIterator, error) {

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "MaxCanonicalNodesUpdated")
	if err != nil {
		return nil, err
	}
	return &NodeRegistryMaxCanonicalNodesUpdatedIterator{contract: _NodeRegistry.contract, event: "MaxCanonicalNodesUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxCanonicalNodesUpdated is a free log subscription operation binding the contract event 0x581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff76.
//
// Solidity: event MaxCanonicalNodesUpdated(uint8 maxCanonicalNodes)
func (_NodeRegistry *NodeRegistryFilterer) WatchMaxCanonicalNodesUpdated(opts *bind.WatchOpts, sink chan<- *NodeRegistryMaxCanonicalNodesUpdated) (event.Subscription, error) {

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "MaxCanonicalNodesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryMaxCanonicalNodesUpdated)
				if err := _NodeRegistry.contract.UnpackLog(event, "MaxCanonicalNodesUpdated", log); err != nil {
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

// ParseMaxCanonicalNodesUpdated is a log parse operation binding the contract event 0x581c4d2fc386422e99f02a47a9735e8936050b0c2a384b98c8a6740786d9ff76.
//
// Solidity: event MaxCanonicalNodesUpdated(uint8 maxCanonicalNodes)
func (_NodeRegistry *NodeRegistryFilterer) ParseMaxCanonicalNodesUpdated(log types.Log) (*NodeRegistryMaxCanonicalNodesUpdated, error) {
	event := new(NodeRegistryMaxCanonicalNodesUpdated)
	if err := _NodeRegistry.contract.UnpackLog(event, "MaxCanonicalNodesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the NodeRegistry contract.
type NodeRegistryMigratedIterator struct {
	Event *NodeRegistryMigrated // Event containing the contract specifics and raw log

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
func (it *NodeRegistryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryMigrated)
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
		it.Event = new(NodeRegistryMigrated)
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
func (it *NodeRegistryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryMigrated represents a Migrated event raised by the NodeRegistry contract.
type NodeRegistryMigrated struct {
	Migrator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_NodeRegistry *NodeRegistryFilterer) FilterMigrated(opts *bind.FilterOpts, migrator []common.Address) (*NodeRegistryMigratedIterator, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryMigratedIterator{contract: _NodeRegistry.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_NodeRegistry *NodeRegistryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *NodeRegistryMigrated, migrator []common.Address) (event.Subscription, error) {

	var migratorRule []interface{}
	for _, migratorItem := range migrator {
		migratorRule = append(migratorRule, migratorItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Migrated", migratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryMigrated)
				if err := _NodeRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
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

// ParseMigrated is a log parse operation binding the contract event 0xa2e7361c23d7820040603b83c0cd3f494d377bac69736377d75bb56c651a5098.
//
// Solidity: event Migrated(address indexed migrator)
func (_NodeRegistry *NodeRegistryFilterer) ParseMigrated(log types.Log) (*NodeRegistryMigrated, error) {
	event := new(NodeRegistryMigrated)
	if err := _NodeRegistry.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeAddedIterator is returned from FilterNodeAdded and is used to iterate over the raw logs and unpacked data for NodeAdded events raised by the NodeRegistry contract.
type NodeRegistryNodeAddedIterator struct {
	Event *NodeRegistryNodeAdded // Event containing the contract specifics and raw log

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
func (it *NodeRegistryNodeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeAdded)
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
		it.Event = new(NodeRegistryNodeAdded)
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
func (it *NodeRegistryNodeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeAdded represents a NodeAdded event raised by the NodeRegistry contract.
type NodeRegistryNodeAdded struct {
	NodeId           uint32
	Owner            common.Address
	Signer           common.Address
	SigningPublicKey []byte
	HttpAddress      string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNodeAdded is a free log retrieval operation binding the contract event 0x9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a4.
//
// Solidity: event NodeAdded(uint32 indexed nodeId, address indexed owner, address indexed signer, bytes signingPublicKey, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeAdded(opts *bind.FilterOpts, nodeId []uint32, owner []common.Address, signer []common.Address) (*NodeRegistryNodeAddedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeAdded", nodeIdRule, ownerRule, signerRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeAddedIterator{contract: _NodeRegistry.contract, event: "NodeAdded", logs: logs, sub: sub}, nil
}

// WatchNodeAdded is a free log subscription operation binding the contract event 0x9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a4.
//
// Solidity: event NodeAdded(uint32 indexed nodeId, address indexed owner, address indexed signer, bytes signingPublicKey, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeAdded(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeAdded, nodeId []uint32, owner []common.Address, signer []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeAdded", nodeIdRule, ownerRule, signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeAdded)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
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

// ParseNodeAdded is a log parse operation binding the contract event 0x9b385c30e390e1e15ab8a2e34c4caa40b3c59882c17185fcbc3f87b2bf6658a4.
//
// Solidity: event NodeAdded(uint32 indexed nodeId, address indexed owner, address indexed signer, bytes signingPublicKey, string httpAddress)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeAdded(log types.Log) (*NodeRegistryNodeAdded, error) {
	event := new(NodeRegistryNodeAdded)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeAddedToCanonicalNetworkIterator is returned from FilterNodeAddedToCanonicalNetwork and is used to iterate over the raw logs and unpacked data for NodeAddedToCanonicalNetwork events raised by the NodeRegistry contract.
type NodeRegistryNodeAddedToCanonicalNetworkIterator struct {
	Event *NodeRegistryNodeAddedToCanonicalNetwork // Event containing the contract specifics and raw log

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
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeAddedToCanonicalNetwork)
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
		it.Event = new(NodeRegistryNodeAddedToCanonicalNetwork)
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
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeAddedToCanonicalNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeAddedToCanonicalNetwork represents a NodeAddedToCanonicalNetwork event raised by the NodeRegistry contract.
type NodeRegistryNodeAddedToCanonicalNetwork struct {
	NodeId uint32
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeAddedToCanonicalNetwork is a free log retrieval operation binding the contract event 0x13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f8.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeAddedToCanonicalNetwork(opts *bind.FilterOpts, nodeId []uint32) (*NodeRegistryNodeAddedToCanonicalNetworkIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeAddedToCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeAddedToCanonicalNetworkIterator{contract: _NodeRegistry.contract, event: "NodeAddedToCanonicalNetwork", logs: logs, sub: sub}, nil
}

// WatchNodeAddedToCanonicalNetwork is a free log subscription operation binding the contract event 0x13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f8.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeAddedToCanonicalNetwork(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeAddedToCanonicalNetwork, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeAddedToCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeAddedToCanonicalNetwork)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeAddedToCanonicalNetwork", log); err != nil {
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

// ParseNodeAddedToCanonicalNetwork is a log parse operation binding the contract event 0x13695734a48552c5f7d826df6e02f4094ed655e28bcedb3ccc3645997f6b47f8.
//
// Solidity: event NodeAddedToCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeAddedToCanonicalNetwork(log types.Log) (*NodeRegistryNodeAddedToCanonicalNetwork, error) {
	event := new(NodeRegistryNodeAddedToCanonicalNetwork)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeAddedToCanonicalNetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryNodeRemovedFromCanonicalNetworkIterator is returned from FilterNodeRemovedFromCanonicalNetwork and is used to iterate over the raw logs and unpacked data for NodeRemovedFromCanonicalNetwork events raised by the NodeRegistry contract.
type NodeRegistryNodeRemovedFromCanonicalNetworkIterator struct {
	Event *NodeRegistryNodeRemovedFromCanonicalNetwork // Event containing the contract specifics and raw log

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
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryNodeRemovedFromCanonicalNetwork)
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
		it.Event = new(NodeRegistryNodeRemovedFromCanonicalNetwork)
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
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryNodeRemovedFromCanonicalNetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryNodeRemovedFromCanonicalNetwork represents a NodeRemovedFromCanonicalNetwork event raised by the NodeRegistry contract.
type NodeRegistryNodeRemovedFromCanonicalNetwork struct {
	NodeId uint32
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeRemovedFromCanonicalNetwork is a free log retrieval operation binding the contract event 0x7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) FilterNodeRemovedFromCanonicalNetwork(opts *bind.FilterOpts, nodeId []uint32) (*NodeRegistryNodeRemovedFromCanonicalNetworkIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "NodeRemovedFromCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryNodeRemovedFromCanonicalNetworkIterator{contract: _NodeRegistry.contract, event: "NodeRemovedFromCanonicalNetwork", logs: logs, sub: sub}, nil
}

// WatchNodeRemovedFromCanonicalNetwork is a free log subscription operation binding the contract event 0x7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) WatchNodeRemovedFromCanonicalNetwork(opts *bind.WatchOpts, sink chan<- *NodeRegistryNodeRemovedFromCanonicalNetwork, nodeId []uint32) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "NodeRemovedFromCanonicalNetwork", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryNodeRemovedFromCanonicalNetwork)
				if err := _NodeRegistry.contract.UnpackLog(event, "NodeRemovedFromCanonicalNetwork", log); err != nil {
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

// ParseNodeRemovedFromCanonicalNetwork is a log parse operation binding the contract event 0x7cf9bcdd519495a485911496098851db2c18ee9a708b453dd48f2822098e16ec.
//
// Solidity: event NodeRemovedFromCanonicalNetwork(uint32 indexed nodeId)
func (_NodeRegistry *NodeRegistryFilterer) ParseNodeRemovedFromCanonicalNetwork(log types.Log) (*NodeRegistryNodeRemovedFromCanonicalNetwork, error) {
	event := new(NodeRegistryNodeRemovedFromCanonicalNetwork)
	if err := _NodeRegistry.contract.UnpackLog(event, "NodeRemovedFromCanonicalNetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the NodeRegistry contract.
type NodeRegistryTransferIterator struct {
	Event *NodeRegistryTransfer // Event containing the contract specifics and raw log

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
func (it *NodeRegistryTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryTransfer)
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
		it.Event = new(NodeRegistryTransfer)
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
func (it *NodeRegistryTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryTransfer represents a Transfer event raised by the NodeRegistry contract.
type NodeRegistryTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NodeRegistryTransferIterator, error) {

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

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryTransferIterator{contract: _NodeRegistry.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NodeRegistry *NodeRegistryFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NodeRegistryTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryTransfer)
				if err := _NodeRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseTransfer(log types.Log) (*NodeRegistryTransfer, error) {
	event := new(NodeRegistryTransfer)
	if err := _NodeRegistry.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NodeRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the NodeRegistry contract.
type NodeRegistryUpgradedIterator struct {
	Event *NodeRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *NodeRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NodeRegistryUpgraded)
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
		it.Event = new(NodeRegistryUpgraded)
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
func (it *NodeRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NodeRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NodeRegistryUpgraded represents a Upgraded event raised by the NodeRegistry contract.
type NodeRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NodeRegistry *NodeRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*NodeRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NodeRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &NodeRegistryUpgradedIterator{contract: _NodeRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NodeRegistry *NodeRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *NodeRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NodeRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NodeRegistryUpgraded)
				if err := _NodeRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_NodeRegistry *NodeRegistryFilterer) ParseUpgraded(log types.Log) (*NodeRegistryUpgraded, error) {
	event := new(NodeRegistryUpgraded)
	if err := _NodeRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
