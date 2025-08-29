package blockchain

import (
	"fmt"
	"regexp"
)

var (
	protocolErrorsDictionary = map[string]string{
		"0x3c355a89": "AlreadyClaimed(uint32 originatorNodeId, uint256 payerReportIndex)",
		"0xa24a13a6": "ArrayLengthMismatch()",
		"0xb4f54111": "DeployFailed()",
		"0xfcc36c5f": "EmptyAdmins()",
		"0x521299a9": "EmptyArray()",
		"0x21744a59": "EmptyBytecode()",
		"0x626c4161": "EmptyCode(address implementation)",
		"0xb6cc7531": "EndIndexOutOfRange()",
		"0x64283d7b": "ERC721IncorrectOwner(address sender, uint256 tokenId, address owner)",
		"0x177e802f": "ERC721InsufficientApproval(address operator, uint256 tokenId)",
		"0xa9fbf51f": "ERC721InvalidApprover(address approver)",
		"0x5b08ba18": "ERC721InvalidOperator(address operator)",
		"0x89c62b64": "ERC721InvalidOwner(address owner)",
		"0x7e273289": "ERC721InvalidOwner(address owner)",
		"0x64a0ae92": "ERC721InvalidReceiver(address receiver)",
		"0x73c6ac6e": "ERC721InvalidSender(address sender)",
		"0x4992486d": "FailedToAddNodeToCanonicalNetwork()",
		"0xe31ff236": "FailedToRemoveNodeFromCanonicalNetwork()",
		"0xea61fe70": "FromIndexOutOfRange()",
		"0x23615171": "InitializationFailed(bytes)",
		"0xf4d678b8": "InsufficientBalance()",
		"0x3cb16751": "InsufficientDeposit(uint96 amount, uint96 minimumDeposit)",
		"0x31f1a313": "InsufficientSignatures(uint8 validSignatureCount, uint8 requiredSignatureCount)",
		"0x7de51b2c": "InvalidBitCount32Input()",
		"0xcbd68989": "InvalidHttpAddress()",
		"0x68155f9a": "InvalidImplementation()",
		"0x3438704d": "InvalidLeafCount()",
		"0x1d8e7a4a": "InvalidMaxPayloadSize()",
		"0xe219e4f0": "InvalidMinPayloadSize()",
		"0x49e27cff": "InvalidOwner()",
		"0x93b7abe6": "InvalidPayloadSize(uint256 actualSize_, uint256 minSize_, uint256 maxSize_)",
		"0x09bde339": "InvalidProof()",
		"0x82eeb3b2": "InvalidProtocolFeeRate()",
		"0xa7ee0517": "InvalidSequenceIds()",
		"0xbf51f547": "InvalidSigningPublicKey()",
		"0x84e23433": "InvalidStartSequenceId(uint64 startSequenceId, uint64 lastSequenceId)",
		"0x3ba01911": "InvalidURI()",
		"0x472b0bf1": "MaxCanonicalNodesBelowCurrentCount()",
		"0x5811df30": "MaxCanonicalNodesReached()",
		"0x957d2080": "MaxNodesReached()",
		"0x68b0b16b": "MigrationFailed(address migrator_, bytes revertData_)",
		"0xdf917c5a": "NoChainIds()",
		"0xa88ee577": "NoChange()",
		"0x6c163b7e": "NoExcess()",
		"0x988ebee0": "NoFeesOwed()",
		"0x2c566e6f": "NoKeyComponents()",
		"0x01909835": "NoKeys()",
		"0xbaec3d9a": "NoLeaves()",
		"0x9121b84f": "NoPendingWithdrawal()",
		"0x320499c4": "NoProofElements()",
		"0x7bfa4b9f": "NotAdmin()",
		"0x4125cde8": "NotInPayerReport(uint32 originatorNodeId, uint256 payerReportIndex)",
		"0xd08a05d5": "NotNodeOwner()",
		"0x6cd60201": "NotPaused()",
		"0xc525f923": "NotPayloadBootstrapper()",
		"0x2f55f067": "NotSettlementChainGateway()",
		"0x05b94333": "NotSettler()",
		"0x37f4f148": "ParameterOutOfTypeBounds()",
		"0x9e87fac8": "Paused()",
		"0xac135e07": "PayerFeesLengthTooLong()",
		"0x906dc378": "PayerInDebt()",
		"0xe7d70a4f": "PayerReportEntirelySettled()",
		"0xc3e0931f": "PayerReportIndexOutOfBounds()",
		"0x7fddb8df": "PayerReportNotSettled(uint32 originatorNodeId, uint256 payerReportIndex)",
		"0x75c41473": "PendingWithdrawalExists()",
		"0xcaa2acb9": "SettleUsageFailed(bytes)",
		"0x90b8ec18": "TransferFailed()",
		"0x7939f424": "TransferFromFailed()",
		"0x99a67242": "UnorderedNodeIds()",
		"0xa5dab5fe": "UnsupportedChainId(uint256 chainId)",
		"0x8f8db830": "WithdrawalNotReady(uint32 timestamp, uint32 withdrawableTimestamp)",
		"0x7289db0e": "ZeroAdmin()",
		"0x1f2a2005": "ZeroAmount()",
		"0x253e0003": "ZeroAppChainId()",
		"0x99b81da3": "ZeroAppChainGateway()",
		"0x2d8dc664": "ZeroAvailableBalance()",
		"0x669567ea": "ZeroBalance()",
		"0x047b9cec": "ZeroCount()",
		"0xa5febaf3": "ZeroFeeDistributor()",
		"0xf22ca9da": "ZeroFeeToken()",
		"0x4208d2eb": "ZeroImplementation()",
		"0x0d626a32": "ZeroMigrator()",
		"0x5bc1c4a0": "ZeroMinimumDeposit()",
		"0xaa63199b": "ZeroNodeRegistry()",
		"0xa26bef69": "ZeroPayer()",
		"0x8695d35a": "ZeroPayerRegistry()",
		"0xe1ac5586": "ZeroPayerReportManager()",
		"0xd973fd8d": "ZeroParameterRegistry()",
		"0xd27b4443": "ZeroRecipient()",
		"0x2d23123e": "ZeroSettlementChainGateway()",
		"0x0450b01d": "ZeroSettler()",
		"0x180d3393": "ZeroTotalAmount()",
		"0xcf986a84": "ZeroUnderlying()",
		"0xbd512241": "ZeroWithdrawalAmount()",
	}

	ErrCodeNotFound = fmt.Errorf("error message does not contain a valid error code")
	ErrCodeNotInDic = fmt.Errorf("code not found in protocol errors dictionary")
	ErrCompileRegex = fmt.Errorf("error compiling regex")
)

type ProtocolError interface {
	error
	IsNoChange() bool
}

type BlockchainError struct {
	msg string
	err error
}

func NewBlockchainError(e error) *BlockchainError {
	if e == nil {
		return nil
	}

	message, err := tryExtractProtocolError(e)
	if err != nil {
		message = err.Error()
	}

	return &BlockchainError{msg: message, err: e}
}

func (e BlockchainError) Error() string {
	switch e.msg {
	case ErrCodeNotFound.Error(), ErrCodeNotInDic.Error(), ErrCompileRegex.Error():
		return e.err.Error()
	default:
		return fmt.Sprintf("%s, original error: %s", e.msg, e.err.Error())
	}
}

func (e BlockchainError) IsNoChange() bool {
	return e.msg == "NoChange()"
}

// tryExtractProtocolError tries to extract the protocol error from the error message.
// Error codes are 4 bytes hex strings, in example: 0x31f1a313.
func tryExtractProtocolError(e error) (message string, err error) {
	re, err := regexp.Compile(
		`(0x[0-9a-fA-F]{8})`,
	)
	if err != nil {
		return "", ErrCompileRegex
	}

	matches := re.FindStringSubmatch(e.Error())
	if len(matches) != 2 {
		return "", ErrCodeNotFound
	}

	message, exists := protocolErrorsDictionary[matches[1]]
	if !exists {
		return "", ErrCodeNotInDic
	}

	return message, nil
}
