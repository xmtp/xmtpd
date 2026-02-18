package blockchain

import (
	"errors"
	"regexp"
)

var (
	ErrAlreadyClaimed              = errors.New("AlreadyClaimed(uint32,uint256)")
	ErrPayerReportAlreadySubmitted = errors.New(
		"PayerReportAlreadySubmitted(uint32,uint64,uint64)",
	)
	ErrArrayLengthMismatch  = errors.New("ArrayLengthMismatch()")
	ErrDeployFailed         = errors.New("DeployFailed()")
	ErrEmptyAdmins          = errors.New("EmptyAdmins()")
	ErrEmptyArray           = errors.New("EmptyArray()")
	ErrEmptyBytecode        = errors.New("EmptyBytecode()")
	ErrEmptyCode            = errors.New("EmptyCode(address)")
	ErrEndIndexOutOfRange   = errors.New("EndIndexOutOfRange()")
	ErrERC721IncorrectOwner = errors.New(
		"ERC721IncorrectOwner(address,uint256,address)",
	)
	ErrERC721InsufficientApproval = errors.New(
		"ERC721InsufficientApproval(address,uint256)",
	)
	ErrERC721InvalidApprover                  = errors.New("ERC721InvalidApprover(address)")
	ErrERC721InvalidOperator                  = errors.New("ERC721InvalidOperator(address)")
	ErrERC721InvalidOwner                     = errors.New("ERC721InvalidOwner(address)")
	ErrERC721InvalidReceiver                  = errors.New("ERC721InvalidReceiver(address)")
	ErrERC721InvalidSender                    = errors.New("ERC721InvalidSender(address)")
	ErrFailedToAddNodeToCanonicalNetwork      = errors.New("FailedToAddNodeToCanonicalNetwork()")
	ErrFailedToRemoveNodeFromCanonicalNetwork = errors.New(
		"FailedToRemoveNodeFromCanonicalNetwork()",
	)
	ErrFromIndexOutOfRange    = errors.New("FromIndexOutOfRange()")
	ErrInitializationFailed   = errors.New("InitializationFailed(bytes)")
	ErrInsufficientBalance    = errors.New("InsufficientBalance()")
	ErrInsufficientDeposit    = errors.New("InsufficientDeposit(uint96,uint96)")
	ErrInsufficientSignatures = errors.New("InsufficientSignatures(uint8,uint8)")
	ErrInvalidBitCount32Input = errors.New("InvalidBitCount32Input()")
	ErrInvalidHTTPAddress     = errors.New("InvalidHttpAddress()")
	ErrInvalidImplementation  = errors.New("InvalidImplementation()")
	ErrInvalidLeafCount       = errors.New("InvalidLeafCount()")
	ErrInvalidMaxPayloadSize  = errors.New("InvalidMaxPayloadSize()")
	ErrInvalidMinPayloadSize  = errors.New("InvalidMinPayloadSize()")
	ErrInvalidOwner           = errors.New("InvalidOwner()")
	ErrInvalidPayloadSize     = errors.New(
		"InvalidPayloadSize(uint256,uint256,uint256)",
	)
	ErrInvalidProof                       = errors.New("InvalidProof()")
	ErrInvalidProtocolFeeRate             = errors.New("InvalidProtocolFeeRate()")
	ErrInvalidSequenceIDs                 = errors.New("InvalidSequenceIds()")
	ErrInvalidSigningPublicKey            = errors.New("InvalidSigningPublicKey()")
	ErrInvalidStartSequenceID             = errors.New("InvalidStartSequenceId(uint64,uint64)")
	ErrInvalidURI                         = errors.New("InvalidURI()")
	ErrMaxCanonicalNodesBelowCurrentCount = errors.New("MaxCanonicalNodesBelowCurrentCount()")
	ErrMaxCanonicalNodesReached           = errors.New("MaxCanonicalNodesReached()")
	ErrMaxNodesReached                    = errors.New("MaxNodesReached()")
	ErrMigrationFailed                    = errors.New("MigrationFailed(address,bytes)")
	ErrNoChainIds                         = errors.New("NoChainIds()")
	ErrNoChange                           = errors.New("NoChange()")
	ErrNoExcess                           = errors.New("NoExcess()")
	ErrNoFeesOwed                         = errors.New("NoFeesOwed()")
	ErrNoKeyComponents                    = errors.New("NoKeyComponents()")
	ErrNoKeys                             = errors.New("NoKeys()")
	ErrNoLeaves                           = errors.New("NoLeaves()")
	ErrNoPendingWithdrawal                = errors.New("NoPendingWithdrawal()")
	ErrNoProofElements                    = errors.New("NoProofElements()")
	ErrNotAdmin                           = errors.New("NotAdmin()")
	ErrNotInPayerReport                   = errors.New("NotInPayerReport(uint32,uint256)")
	ErrNotNodeOwner                       = errors.New("NotNodeOwner()")
	ErrNotPaused                          = errors.New("NotPaused()")
	ErrNotPayloadBootstrapper             = errors.New("NotPayloadBootstrapper()")
	ErrNotSettlementChainGateway          = errors.New("NotSettlementChainGateway()")
	ErrNotSettler                         = errors.New("NotSettler()")
	ErrParameterOutOfTypeBounds           = errors.New("ParameterOutOfTypeBounds()")
	ErrPaused                             = errors.New("Paused()")
	ErrPayerFeesLengthTooLong             = errors.New("PayerFeesLengthTooLong()")
	ErrPayerInDebt                        = errors.New("PayerInDebt()")
	ErrPayerReportEntirelySettled         = errors.New("PayerReportEntirelySettled()")
	ErrPayerReportIndexOutOfBounds        = errors.New("PayerReportIndexOutOfBounds()")
	ErrPayerReportNotSettled              = errors.New("PayerReportNotSettled(uint32,uint256)")
	ErrPendingWithdrawalExists            = errors.New("PendingWithdrawalExists()")
	ErrTransferFailed                     = errors.New("TransferFailed()")
	ErrTransferFromFailed                 = errors.New("TransferFromFailed()")
	ErrUnorderedNodeIDs                   = errors.New("UnorderedNodeIds()")
	ErrUnsupportedChainID                 = errors.New("UnsupportedChainId(uint256)")
	ErrWithdrawalNotReady                 = errors.New("WithdrawalNotReady(uint32,uint32)")
	ErrZeroAdmin                          = errors.New("ZeroAdmin()")
	ErrZeroAmount                         = errors.New("ZeroAmount()")
	ErrZeroAppChainID                     = errors.New("ZeroAppChainId()")
	ErrZeroAppChainGateway                = errors.New("ZeroAppChainGateway()")
	ErrZeroAvailableBalance               = errors.New("ZeroAvailableBalance()")
	ErrZeroBalance                        = errors.New("ZeroBalance()")
	ErrZeroCount                          = errors.New("ZeroCount()")
	ErrZeroFeeDistributor                 = errors.New("ZeroFeeDistributor()")
	ErrZeroFeeToken                       = errors.New("ZeroFeeToken()")
	ErrZeroImplementation                 = errors.New("ZeroImplementation()")
	ErrZeroMigrator                       = errors.New("ZeroMigrator()")
	ErrZeroMinimumDeposit                 = errors.New("ZeroMinimumDeposit()")
	ErrZeroNodeRegistry                   = errors.New("ZeroNodeRegistry()")
	ErrZeroPayer                          = errors.New("ZeroPayer()")
	ErrZeroPayerRegistry                  = errors.New("ZeroPayerRegistry()")
	ErrZeroPayerReportManager             = errors.New("ZeroPayerReportManager()")
	ErrZeroParameterRegistry              = errors.New("ZeroParameterRegistry()")
	ErrZeroRecipient                      = errors.New("ZeroRecipient()")
	ErrZeroSettlementChainGateway         = errors.New("ZeroSettlementChainGateway()")
	ErrZeroSettler                        = errors.New("ZeroSettler()")
	ErrZeroTotalAmount                    = errors.New("ZeroTotalAmount()")
	ErrZeroUnderlying                     = errors.New("ZeroUnderlying()")
	ErrZeroWithdrawalAmount               = errors.New("ZeroWithdrawalAmount()")

	protocolErrorsDictionary = map[string]error{
		"0x3c355a89": ErrAlreadyClaimed,
		"0x93105c68": ErrPayerReportAlreadySubmitted,
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

	ErrCodeNotFound = errors.New("error message does not contain a valid error code")
	ErrCodeNotInDic = errors.New("code not found in protocol errors dictionary")
	ErrCompileRegex = errors.New("error compiling regex")
)

type ProtocolError interface {
	error
	Unwrap() error
	IsNoChange() bool
	IsErrInvalidSequenceIDs() bool
	IsErrPayerReportAlreadySubmitted() bool
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
	if e.protocolErr == nil || errors.Is(e.protocolErr, ErrCodeNotFound) ||
		errors.Is(e.protocolErr, ErrCodeNotInDic) ||
		errors.Is(e.protocolErr, ErrCompileRegex) {
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

// IsErrPayerReportAlreadySubmitted returns true if the error is a payer report already submitted error.
func (e BlockchainError) IsErrPayerReportAlreadySubmitted() bool {
	return e.protocolErr != nil && errors.Is(e.protocolErr, ErrPayerReportAlreadySubmitted)
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
