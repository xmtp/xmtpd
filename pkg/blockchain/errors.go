package blockchain

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrAlreadyClaimed = fmt.Errorf(
		"AlreadyClaimed(uint32,uint256)",
	)
	ErrArrayLengthMismatch  = fmt.Errorf("ArrayLengthMismatch()")
	ErrDeployFailed         = fmt.Errorf("DeployFailed()")
	ErrEmptyAdmins          = fmt.Errorf("EmptyAdmins()")
	ErrEmptyArray           = fmt.Errorf("EmptyArray()")
	ErrEmptyBytecode        = fmt.Errorf("EmptyBytecode()")
	ErrEmptyCode            = fmt.Errorf("EmptyCode(address)")
	ErrEndIndexOutOfRange   = fmt.Errorf("EndIndexOutOfRange()")
	ErrERC721IncorrectOwner = fmt.Errorf(
		"ERC721IncorrectOwner(address,uint256,address)",
	)
	ErrERC721InsufficientApproval = fmt.Errorf(
		"ERC721InsufficientApproval(address,uint256)",
	)
	ErrERC721InvalidApprover             = fmt.Errorf("ERC721InvalidApprover(address)")
	ErrERC721InvalidOperator             = fmt.Errorf("ERC721InvalidOperator(address)")
	ErrERC721InvalidOwner                = fmt.Errorf("ERC721InvalidOwner(address)")
	ErrERC721InvalidReceiver             = fmt.Errorf("ERC721InvalidReceiver(address)")
	ErrERC721InvalidSender               = fmt.Errorf("ERC721InvalidSender(address)")
	ErrFailedToAddNodeToCanonicalNetwork = fmt.Errorf(
		"FailedToAddNodeToCanonicalNetwork()",
	)
	ErrFailedToRemoveNodeFromCanonicalNetwork = fmt.Errorf(
		"FailedToRemoveNodeFromCanonicalNetwork()",
	)
	ErrFromIndexOutOfRange  = fmt.Errorf("FromIndexOutOfRange()")
	ErrInitializationFailed = fmt.Errorf("InitializationFailed(bytes)")
	ErrInsufficientBalance  = fmt.Errorf("InsufficientBalance()")
	ErrInsufficientDeposit  = fmt.Errorf(
		"InsufficientDeposit(uint96,uint96)",
	)
	ErrInsufficientSignatures = fmt.Errorf(
		"InsufficientSignatures(uint8,uint8)",
	)
	ErrInvalidBitCount32Input = fmt.Errorf("InvalidBitCount32Input()")
	ErrInvalidHTTPAddress     = fmt.Errorf("InvalidHttpAddress()")
	ErrInvalidImplementation  = fmt.Errorf("InvalidImplementation()")
	ErrInvalidLeafCount       = fmt.Errorf("InvalidLeafCount()")
	ErrInvalidMaxPayloadSize  = fmt.Errorf("InvalidMaxPayloadSize()")
	ErrInvalidMinPayloadSize  = fmt.Errorf("InvalidMinPayloadSize()")
	ErrInvalidOwner           = fmt.Errorf("InvalidOwner()")
	ErrInvalidPayloadSize     = fmt.Errorf(
		"InvalidPayloadSize(uint256,uint256,uint256)",
	)
	ErrInvalidProof            = fmt.Errorf("InvalidProof()")
	ErrInvalidProtocolFeeRate  = fmt.Errorf("InvalidProtocolFeeRate()")
	ErrInvalidSequenceIDs      = fmt.Errorf("InvalidSequenceIds()")
	ErrInvalidSigningPublicKey = fmt.Errorf("InvalidSigningPublicKey()")
	ErrInvalidStartSequenceID  = fmt.Errorf(
		"InvalidStartSequenceId(uint64,uint64)",
	)
	ErrInvalidURI                         = fmt.Errorf("InvalidURI()")
	ErrMaxCanonicalNodesBelowCurrentCount = fmt.Errorf(
		"MaxCanonicalNodesBelowCurrentCount()",
	)
	ErrMaxCanonicalNodesReached = fmt.Errorf("MaxCanonicalNodesReached()")
	ErrMaxNodesReached          = fmt.Errorf("MaxNodesReached()")
	ErrMigrationFailed          = fmt.Errorf(
		"MigrationFailed(address,bytes)",
	)
	ErrNoChainIds          = fmt.Errorf("NoChainIds()")
	ErrNoChange            = fmt.Errorf("NoChange()")
	ErrNoExcess            = fmt.Errorf("NoExcess()")
	ErrNoFeesOwed          = fmt.Errorf("NoFeesOwed()")
	ErrNoKeyComponents     = fmt.Errorf("NoKeyComponents()")
	ErrNoKeys              = fmt.Errorf("NoKeys()")
	ErrNoLeaves            = fmt.Errorf("NoLeaves()")
	ErrNoPendingWithdrawal = fmt.Errorf("NoPendingWithdrawal()")
	ErrNoProofElements     = fmt.Errorf("NoProofElements()")
	ErrNotAdmin            = fmt.Errorf("NotAdmin()")
	ErrNotInPayerReport    = fmt.Errorf(
		"NotInPayerReport(uint32,uint256)",
	)
	ErrNotNodeOwner                = fmt.Errorf("NotNodeOwner()")
	ErrNotPaused                   = fmt.Errorf("NotPaused()")
	ErrNotPayloadBootstrapper      = fmt.Errorf("NotPayloadBootstrapper()")
	ErrNotSettlementChainGateway   = fmt.Errorf("NotSettlementChainGateway()")
	ErrNotSettler                  = fmt.Errorf("NotSettler()")
	ErrParameterOutOfTypeBounds    = fmt.Errorf("ParameterOutOfTypeBounds()")
	ErrPaused                      = fmt.Errorf("Paused()")
	ErrPayerFeesLengthTooLong      = fmt.Errorf("PayerFeesLengthTooLong()")
	ErrPayerInDebt                 = fmt.Errorf("PayerInDebt()")
	ErrPayerReportEntirelySettled  = fmt.Errorf("PayerReportEntirelySettled()")
	ErrPayerReportIndexOutOfBounds = fmt.Errorf("PayerReportIndexOutOfBounds()")
	ErrPayerReportNotSettled       = fmt.Errorf(
		"PayerReportNotSettled(uint32,uint256)",
	)
	ErrPendingWithdrawalExists = fmt.Errorf("PendingWithdrawalExists()")
	ErrSettleUsageFailed       = fmt.Errorf("SettleUsageFailed(bytes)")
	ErrTransferFailed          = fmt.Errorf("TransferFailed()")
	ErrTransferFromFailed      = fmt.Errorf("TransferFromFailed()")
	ErrUnorderedNodeIDs        = fmt.Errorf("UnorderedNodeIds()")
	ErrUnsupportedChainID      = fmt.Errorf("UnsupportedChainId(uint256)")
	ErrWithdrawalNotReady      = fmt.Errorf(
		"WithdrawalNotReady(uint32,uint32)",
	)
	ErrZeroAdmin                  = fmt.Errorf("ZeroAdmin()")
	ErrZeroAmount                 = fmt.Errorf("ZeroAmount()")
	ErrZeroAppChainID             = fmt.Errorf("ZeroAppChainId()")
	ErrZeroAppChainGateway        = fmt.Errorf("ZeroAppChainGateway()")
	ErrZeroAvailableBalance       = fmt.Errorf("ZeroAvailableBalance()")
	ErrZeroBalance                = fmt.Errorf("ZeroBalance()")
	ErrZeroCount                  = fmt.Errorf("ZeroCount()")
	ErrZeroFeeDistributor         = fmt.Errorf("ZeroFeeDistributor()")
	ErrZeroFeeToken               = fmt.Errorf("ZeroFeeToken()")
	ErrZeroImplementation         = fmt.Errorf("ZeroImplementation()")
	ErrZeroMigrator               = fmt.Errorf("ZeroMigrator()")
	ErrZeroMinimumDeposit         = fmt.Errorf("ZeroMinimumDeposit()")
	ErrZeroNodeRegistry           = fmt.Errorf("ZeroNodeRegistry()")
	ErrZeroPayer                  = fmt.Errorf("ZeroPayer()")
	ErrZeroPayerRegistry          = fmt.Errorf("ZeroPayerRegistry()")
	ErrZeroPayerReportManager     = fmt.Errorf("ZeroPayerReportManager()")
	ErrZeroParameterRegistry      = fmt.Errorf("ZeroParameterRegistry()")
	ErrZeroRecipient              = fmt.Errorf("ZeroRecipient()")
	ErrZeroSettlementChainGateway = fmt.Errorf("ZeroSettlementChainGateway()")
	ErrZeroSettler                = fmt.Errorf("ZeroSettler()")
	ErrZeroTotalAmount            = fmt.Errorf("ZeroTotalAmount()")
	ErrZeroUnderlying             = fmt.Errorf("ZeroUnderlying()")
	ErrZeroWithdrawalAmount       = fmt.Errorf("ZeroWithdrawalAmount()")

	protocolErrorsDictionary = map[string]error{
		"0x3c355a89": ErrAlreadyClaimed,
		"0xa24a13a6": ErrArrayLengthMismatch,
		"0xb4f54111": ErrDeployFailed,
		"0xfcc36c5f": ErrEmptyAdmins,
		"0x521299a9": ErrEmptyArray,
		"0x21744a59": ErrEmptyBytecode,
		"0x626c4161": ErrEmptyCode,
		"0xb6cc7531": ErrEndIndexOutOfRange,
		"0x64283d7b": ErrERC721IncorrectOwner,
		"0x177e802f": ErrERC721InsufficientApproval,
		"0xa9fbf51f": ErrERC721InvalidApprover,
		"0x5b08ba18": ErrERC721InvalidOperator,
		"0x89c62b64": ErrERC721InvalidOwner,
		"0x64a0ae92": ErrERC721InvalidReceiver,
		"0x73c6ac6e": ErrERC721InvalidSender,
		"0x4992486d": ErrFailedToAddNodeToCanonicalNetwork,
		"0xe31ff236": ErrFailedToRemoveNodeFromCanonicalNetwork,
		"0xea61fe70": ErrFromIndexOutOfRange,
		"0x23615171": ErrInitializationFailed,
		"0xf4d678b8": ErrInsufficientBalance,
		"0x3cb16751": ErrInsufficientDeposit,
		"0x31f1a313": ErrInsufficientSignatures,
		"0x7de51b2c": ErrInvalidBitCount32Input,
		"0xcbd68989": ErrInvalidHTTPAddress,
		"0x68155f9a": ErrInvalidImplementation,
		"0x3438704d": ErrInvalidLeafCount,
		"0x1d8e7a4a": ErrInvalidMaxPayloadSize,
		"0xe219e4f0": ErrInvalidMinPayloadSize,
		"0x49e27cff": ErrInvalidOwner,
		"0x93b7abe6": ErrInvalidPayloadSize,
		"0x09bde339": ErrInvalidProof,
		"0x82eeb3b2": ErrInvalidProtocolFeeRate,
		"0xa7ee0517": ErrInvalidSequenceIDs,
		"0xbf51f547": ErrInvalidSigningPublicKey,
		"0x84e23433": ErrInvalidStartSequenceID,
		"0x3ba01911": ErrInvalidURI,
		"0x472b0bf1": ErrMaxCanonicalNodesBelowCurrentCount,
		"0x5811df30": ErrMaxCanonicalNodesReached,
		"0x957d2080": ErrMaxNodesReached,
		"0x68b0b16b": ErrMigrationFailed,
		"0xdf917c5a": ErrNoChainIds,
		"0xa88ee577": ErrNoChange,
		"0x6c163b7e": ErrNoExcess,
		"0x988ebee0": ErrNoFeesOwed,
		"0x2c566e6f": ErrNoKeyComponents,
		"0x01909835": ErrNoKeys,
		"0xbaec3d9a": ErrNoLeaves,
		"0x9121b84f": ErrNoPendingWithdrawal,
		"0x320499c4": ErrNoProofElements,
		"0x7bfa4b9f": ErrNotAdmin,
		"0x4125cde8": ErrNotInPayerReport,
		"0xd08a05d5": ErrNotNodeOwner,
		"0x6cd60201": ErrNotPaused,
		"0xc525f923": ErrNotPayloadBootstrapper,
		"0x2f55f067": ErrNotSettlementChainGateway,
		"0x05b94333": ErrNotSettler,
		"0x37f4f148": ErrParameterOutOfTypeBounds,
		"0x9e87fac8": ErrPaused,
		"0xac135e07": ErrPayerFeesLengthTooLong,
		"0x906dc378": ErrPayerInDebt,
		"0xe7d70a4f": ErrPayerReportEntirelySettled,
		"0xc3e0931f": ErrPayerReportIndexOutOfBounds,
		"0x7fddb8df": ErrPayerReportNotSettled,
		"0x75c41473": ErrPendingWithdrawalExists,
		"0xcaa2acb9": ErrSettleUsageFailed,
		"0x90b8ec18": ErrTransferFailed,
		"0x7939f424": ErrTransferFromFailed,
		"0x99a67242": ErrUnorderedNodeIDs,
		"0xa5dab5fe": ErrUnsupportedChainID,
		"0x8f8db830": ErrWithdrawalNotReady,
		"0x7289db0e": ErrZeroAdmin,
		"0x1f2a2005": ErrZeroAmount,
		"0x253e0003": ErrZeroAppChainID,
		"0x99b81da3": ErrZeroAppChainGateway,
		"0x2d8dc664": ErrZeroAvailableBalance,
		"0x669567ea": ErrZeroBalance,
		"0x047b9cec": ErrZeroCount,
		"0xa5febaf3": ErrZeroFeeDistributor,
		"0xf22ca9da": ErrZeroFeeToken,
		"0x4208d2eb": ErrZeroImplementation,
		"0x0d626a32": ErrZeroMigrator,
		"0x5bc1c4a0": ErrZeroMinimumDeposit,
		"0xaa63199b": ErrZeroNodeRegistry,
		"0xa26bef69": ErrZeroPayer,
		"0x8695d35a": ErrZeroPayerRegistry,
		"0xe1ac5586": ErrZeroPayerReportManager,
		"0xd973fd8d": ErrZeroParameterRegistry,
		"0xd27b4443": ErrZeroRecipient,
		"0x2d23123e": ErrZeroSettlementChainGateway,
		"0x0450b01d": ErrZeroSettler,
		"0x180d3393": ErrZeroTotalAmount,
		"0xcf986a84": ErrZeroUnderlying,
		"0xbd512241": ErrZeroWithdrawalAmount,
	}

	ErrCodeNotFound = fmt.Errorf("error message does not contain a valid error code")
	ErrCodeNotInDic = fmt.Errorf("code not found in protocol errors dictionary")
	ErrCompileRegex = fmt.Errorf("error compiling regex")
)

type ProtocolError interface {
	error
	Unwrap() error
	IsNoChange() bool
	IsErrInvalidSequenceIDs() bool
}

type BlockchainError struct {
	protocolErr error
	originalErr error
}

func NewBlockchainError(originalErr error) *BlockchainError {
	if originalErr == nil {
		return nil
	}

	protocolErr, err := tryExtractProtocolError(originalErr)
	if err != nil {
		return &BlockchainError{
			protocolErr: err,
			originalErr: originalErr,
		}
	}

	return &BlockchainError{protocolErr: protocolErr, originalErr: originalErr}
}

func (e BlockchainError) Error() string {
	if e.protocolErr == nil || e.protocolErr == ErrCodeNotFound ||
		e.protocolErr == ErrCodeNotInDic ||
		e.protocolErr == ErrCompileRegex {
		return e.originalErr.Error()
	}

	return e.protocolErr.Error()
}

// Unwrap returns the protocol error for errors.Is() checks
func (e *BlockchainError) Unwrap() error {
	if e.protocolErr != nil {
		return e.protocolErr
	}

	return e.originalErr
}

func (e BlockchainError) IsNoChange() bool {
	return e.protocolErr != nil && errors.Is(e.protocolErr, ErrNoChange)
}

// IsErrInvalidSequenceIDs returns true if the error is an invalid sequence ID error.
// That can happen because the report:
// - Was submitted with a wrong start sequence ID.
// - The end sequence ID is smaller than the start sequence ID.
func (e BlockchainError) IsErrInvalidSequenceIDs() bool {
	return e.protocolErr != nil && (errors.Is(e.protocolErr, ErrInvalidSequenceIDs) ||
		errors.Is(e.protocolErr, ErrInvalidStartSequenceID))
}

// tryExtractProtocolError tries to extract the protocol error from the error message.
// Error codes are 4 bytes hex strings, in example: 0x31f1a313.
func tryExtractProtocolError(e error) (message, err error) {
	re, err := regexp.Compile(
		`(0x[0-9a-fA-F]{8})`,
	)
	if err != nil {
		return nil, ErrCompileRegex
	}

	matches := re.FindStringSubmatch(e.Error())
	if len(matches) != 2 {
		return nil, ErrCodeNotFound
	}

	protocolError, exists := protocolErrorsDictionary[matches[1]]
	if !exists {
		return nil, ErrCodeNotInDic
	}

	return protocolError, nil
}
