package payerreport

import "github.com/ethereum/go-ethereum/accounts/abi"

var payerReportMessageHash = abi.Arguments{
	{
		Name: "originatorNodeID",
		Type: abi.Type{T: abi.UintTy, Size: 32},
	},
	{
		Name: "startSequenceID",
		Type: abi.Type{T: abi.UintTy, Size: 64},
	},
	{
		Name: "endSequenceID",
		Type: abi.Type{T: abi.UintTy, Size: 64},
	},
	{
		Name: "payersMerkleRoot",
		Type: abi.Type{T: abi.FixedBytesTy, Size: 32},
	},
	{
		Name: "activeNodeIds",
		Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.UintTy, Size: 32}},
	},
}
