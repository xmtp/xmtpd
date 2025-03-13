// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { IERC165 } from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title  IPayer
 * @notice Interface for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process.
 */
interface IPayer is IERC165 {
    /* ============ Structs ============ */

    /**
     * @dev   Struct to store payer information.
     * @param balance                The current USDC balance of the payer.
     * @param isActive               Indicates whether the payer is active.
     * @param creationTimestamp      The timestamp when the payer was first registered.
     * @param latestDepositTimestamp The timestamp of the most recent deposit.
     * @param debtAmount             The amount of fees owed but not yet settled.
     */
    struct Payer {
        uint256 balance;
        uint256 debtAmount;
        uint256 creationTimestamp;
        uint256 latestDepositTimestamp;
        uint256 latestDonationTimestamp;
        bool isActive;
    }

    /**
     * @dev   Struct to store withdrawal request information.
     * @param requestTimestamp      The timestamp when the withdrawal was requested.
     * @param withdrawableTimestamp The timestamp when the withdrawal can be finalized.
     * @param amount                The amount requested for withdrawal.
     */
    struct Withdrawal {
        uint256 requestTimestamp;
        uint256 withdrawableTimestamp;
        uint256 amount;
    }

    /* ============ Events ============ */

    /// @dev Emitted when a new payer is registered.
    event PayerRegistered(address indexed payer, uint256 amount);

    /// @dev Emitted when a payer is deactivated by an owner.
    event PayerDeactivated(address indexed payer);

    /// @dev Emitted when a payer is permanently deleted from the system.
    event PayerDeleted(address indexed payer, uint256 timestamp);

    /// @dev Emitted when a deposit is made to a payer's account.
    event Deposit(address indexed payer, uint256 amount);

    /// @dev Emitted when a user donates to a payer's account.
    event Donation(address indexed donor, address indexed payer, uint256 amount);

    /// @dev Emitted when a payer balance is updated.
    event PayerBalanceUpdated(address indexed payer, uint256 newBalance, uint256 newDebtAmount);

    /// @dev Emitted when a payer initiates a withdrawal request.
    event WithdrawalRequested(
        address indexed payer, uint256 indexed requestTimestamp, uint256 withdrawableTimestamp, uint256 amount
    );

    /// @dev Emitted when a payer cancels a withdrawal request.
    event WithdrawalCancelled(address indexed payer, uint256 indexed requestTimestamp);

    /// @dev Emitted when a payer's withdrawal is finalized.
    event WithdrawalFinalized(address indexed payer, uint256 indexed requestTimestamp);

    /// @dev Emitted when usage is settled and fees are calculated.
    event UsageSettled(uint256 fees, address indexed payer, uint256 indexed originatorNode, uint256 timestamp);

    /// @dev Emitted when fees are transferred to the distribution contract.
    event FeesTransferred(uint256 amount);

    /// @dev Emitted when the distribution contract address is updated.
    event DistributionContractSet(address indexed newDistributionContract);

    /// @dev Emitted when the nodes contract address is updated.
    event NodesContractSet(address indexed newNodesContract);

    /// @dev Emitted when the payer report contract address is updated.
    event PayerReportContractSet(address indexed newPayerReportContract);

    /// @dev Emitted when the USDC token address is updated.
    event UsdcTokenSet(address indexed newUsdcToken);

    /// @dev Emitted when the minimum deposit amount is updated.
    event MinimumDepositSet(uint256 oldMinimumDeposit, uint256 newMinimumDeposit);

    /// @dev Emitted when the minimum registration amount is updated.
    event MinimumRegistrationAmountSet(uint256 oldMinimumRegistrationAmount, uint256 newMinimumRegistrationAmount);

    /// @dev Emitted when the upgrade is authorized.
    event UpgradeAuthorized(address indexed upgrader, address indexed newImplementation);

    /// @dev Emitted when the withdrawal lock period is updated.
    event WithdrawalLockPeriodSet(uint256 oldWithdrawalLockPeriod, uint256 newWithdrawalLockPeriod);

    /// @dev Emitted when the maximum backdated time is updated.
    event MaxBackdatedTimeSet(uint256 oldMaxBackdatedTime, uint256 newMaxBackdatedTime);

    /* ============ Custom Errors ============ */

    /// @dev Error thrown when a call is unauthorized.
    error Unauthorized();

    /// @dev Error thrown when caller is not an authorized node operator.
    error UnauthorizedNodeOperator();

    /// @dev Error thrown when contract is not the distribution contract.
    error InvalidDistributionContract();

    /// @dev Error thrown when contract is not the payer report contract.
    error InvalidPayerReportContract();

    /// @dev Error thrown when contract is not the nodes contract.
    error InvalidNodesContract();

    /// @dev Error thrown when contract is not the USDC token contract.
    error InvalidUsdcTokenContract();

    /// @dev Error thrown when an address is invalid (usually zero address).
    error InvalidAddress();

    /// @dev Error thrown when the amount is insufficient.
    error InsufficientAmount();

    /// @dev Error thrown when balance is insufficient.
    error InsufficientBalance();

    /// @dev Error thrown when the minimum deposit is invalid.
    error InvalidMinimumDeposit();

    /// @dev Error thrown when the minimum registration amount is invalid.
    error InvalidMinimumRegistrationAmount();

    /// @dev Error thrown when the withdrawal lock period is invalid.
    error InvalidWithdrawalLockPeriod();

    /// @dev Error thrown when the maximum backdated time is invalid.
    error InvalidMaxBackdatedTime();

    /// @dev Error thrown when a withdrawal is not in the requested state.
    error WithdrawalNotRequested();

    /// @dev Error thrown when a withdrawal is already in progress.
    error WithdrawalAlreadyRequested();

    /// @dev Error thrown when a withdrawal is not in the requested state.
    error WithdrawalNotExists();

    /// @dev Error thrown when a lock period has not yet elapsed.
    error LockPeriodNotElapsed();

    /// @dev Error thrown when arrays have mismatched lengths.
    error ArrayLengthMismatch();

    /// @dev Error thrown when trying to backdate settlement too far.
    error InvalidSettlementTime();

    /// @dev Error thrown when payer already exists.
    error PayerAlreadyRegistered();

    /// @dev Error thrown when payer does not exist.
    error PayerDoesNotExist();

    /// @dev Error thrown when trying to delete a payer with balance or debt.
    error PayerHasBalanceOrDebt();

    /// @dev Error thrown when payer has debt.
    error PayerHasDebt();

    /// @dev Error thrown when trying to delete a payer in withdrawal state.
    error PayerInWithdrawal();

    /// @dev Error thrown when a payer is not active.
    error PayerIsNotActive();

    /// @notice Error thrown when granting a role has failed.
    error FailedToGrantRole(bytes32 role, address account);

    /// @notice Error thrown when registering a payer has failed.
    error FailedToRegisterPayer();

    /// @notice Error thrown when deactivating a payer has failed.
    error FailedToDeactivatePayer();

    /// @notice Error thrown when deleting a payer has failed.
    error FailedToDeletePayer();

    /* ============ Payer Management ============ */

    /**
     * @notice Registers the caller as a new payer upon depositing the minimum required USDC.
     *         The caller must approve this contract to spend USDC beforehand.
     * @param  amount The amount of USDC to deposit (must be at least the minimum required).
     *
     * Emits `PayerRegistered`.
     */
    function register(uint256 amount) external;

    /**
     * @notice Allows the caller to deposit USDC into their own payer account.
     *         The caller must approve this contract to spend USDC beforehand.
     * @param  amount The amount of USDC to deposit.
     *
     * Emits `Deposit`.
     */
    function deposit(uint256 amount) external;

    /**
     * @notice Allows anyone to donate USDC to an existing payer's account.
     *         The sender must approve this contract to spend USDC beforehand.
     * @param  payer  The address of the payer receiving the donation.
     * @param  amount The amount of USDC to donate.
     *
     * Emits `Donation`.
     */
    function donate(address payer, uint256 amount) external;

    /**
     * @notice Deactivates a payer, signaling XMTP nodes they should not accept messages from them.
     *         Only callable by authorized node operators.
     * @param  payer The address of the payer to deactivate.
     *
     * Emits `PayerDeactivated`.
     */
    function deactivatePayer(address payer) external;

    /**
     * @notice Permanently deletes a payer from the system.
     * @dev    Can only delete payers with zero balance and zero debt who are not in withdrawal.
     *         Only callable by authorized node operators.
     * @param  payer The address of the payer to delete.
     *
     * Emits `PayerDeleted`.
     */
    function deletePayer(address payer) external;

    /* ============ Payer Balance Management ============ */

    /**
     * @notice Initiates a withdrawal request for the caller.
     *         - Sets the payer into withdrawal mode (no further usage allowed).
     *         - Records a timestamp for the withdrawal lock period.
     * @param  amount The amount to withdraw (can be less than or equal to current balance).
     *
     * Emits `WithdrawalRequest`.
     */
    function requestWithdrawal(uint256 amount) external;

    /**
     * @notice Cancels a previously requested withdrawal, removing withdrawal mode.
     * @dev    Only callable by the payer who initiated the withdrawal.
     *
     * Emits `WithdrawalCancelled`.
     */
    function cancelWithdrawal() external;

    /**
     * @notice Finalizes a payer's withdrawal after the lock period has elapsed.
     *         - Accounts for any pending usage during the lock.
     *         - Returns the unspent balance to the payer.
     *
     * Emits `WithdrawalFinalized`.
     */
    function finalizeWithdrawal() external;

    /**
     * @notice Checks if a payer is currently in withdrawal mode and the timestamp
     *         when they initiated the withdrawal.
     * @param  payer                 The address to check.
     * @return withdrawal            The withdrawal status.
     */
    function getWithdrawalStatus(address payer)
        external
        view
        returns (Withdrawal memory withdrawal);

    /* ============ Usage Settlement ============ */

    /**
     * @notice Settles usage for a contiguous batch of (payer, amount) entries.
     * Assumes that the PayerReport contract has already verified the aggregated Merkle proof.
     *
     * @param  originatorNode The node that submitted the report.
     * @param  reportIndex    The index of the report.
     * @param  payers         A contiguous array of payer addresses.
     * @param  amounts        A contiguous array of usage amounts corresponding to each payer.
     */
    function settleUsage(
        address originatorNode,
        uint256 reportIndex,
        address[] calldata payers,
        uint256[] calldata amounts
    ) external; /* onlyPayerReport */

    /**
     * @notice Transfers all pending fees to the designated distribution contract.
     * @dev    Uses a single storage write for updating accumulated fees.
     *
     * Emits `FeesTransferred`.
     */
    function transferFeesToDistribution() external;

    /* ============ Administrative Functions ============ */

    /**
     * @notice Sets the address of the distribution contract.
     * @param  distributionContract The address of the new distribution contract.
     *
     * Emits `DistributionContractUpdated`.
     */
    function setDistributionContract(address distributionContract) external;

    /**
     * @notice Sets the address of the payer report contract.
     * @param  payerReportContract The address of the new payer report contract.
     *
     * Emits `PayerReportContractUpdated`.
     */
    function setPayerReportContract(address payerReportContract) external;

    /**
     * @notice Sets the address of the nodes contract for operator verification.
     * @param  nodesContract The address of the new nodes contract.
     *
     * Emits `NodesContractUpdated`.
     */
    function setNodesContract(address nodesContract) external;

    /**
     * @notice Sets the address of the USDC token contract.
     * @param  usdcToken The address of the new USDC token contract.
     *
     * Emits `UsdcTokenUpdated`.
     */
    function setUsdcToken(address usdcToken) external;

    /**
     * @notice Sets the minimum deposit amount required for registration.
     * @param  newMinimumDeposit The new minimum deposit amount.
     *
     * Emits `MinimumDepositUpdated`.
     */
    function setMinimumDeposit(uint256 newMinimumDeposit) external;

    /**
     * @notice Sets the minimum deposit amount required for registration.
     * @param  newMinimumRegistrationAmount The new minimum deposit amount.
     *
     * Emits `MinimumRegistrationAmountUpdated`.
     */
    function setMinimumRegistrationAmount(uint256 newMinimumRegistrationAmount) external;

    /**
     * @notice Sets the withdrawal lock period.
     * @param  newWithdrawalLockPeriod The new withdrawal lock period.
     *
     * Emits `WithdrawalLockPeriodUpdated`.
     */
    function setWithdrawalLockPeriod(uint256 newWithdrawalLockPeriod) external;

    /**
     * @notice Sets the maximum backdated time for settlements.
     * @param  newMaxBackdatedTime The new maximum backdated time.
     *
     * Emits `MaxBackdatedTimeUpdated`.
     */
    function setMaxBackdatedTime(uint256 newMaxBackdatedTime) external;

    /**
     * @notice Pauses the contract functions in case of emergency.
     *
     * Emits `Paused()`.
     */
    function pause() external;

    /**
     * @notice Unpauses the contract.
     *
     * Emits `Unpaused()`.
     */
    function unpause() external;

    /* ============ Getters ============ */

    /**
     * @notice Returns the payer information.
     * @param  payer The address of the payer.
     * @return payerInfo The payer information.
     */
    function getPayer(address payer) external view returns (Payer memory payerInfo);

    /**
     * @notice Checks if a given address is an active payer.
     * @param  payer    The address to check.
     * @return isActive True if the address is an active payer, false otherwise.
     */
    function getIsActivePayer(address payer) external view returns (bool isActive);

    /**
     * @notice Returns a paginated list of payers with outstanding debt.
     * @param  offset      Number of payers to skip before starting to return results.
     * @param  limit       Maximum number of payers to return.
     * @return debtors     Array of payer addresses with debt.
     * @return debtAmounts Corresponding debt amounts for each payer.
     * @return totalCount  Total number of payers with debt (regardless of pagination).
     */
    function getPayersInDebt(uint256 offset, uint256 limit)
        external
        view
        returns (address[] memory debtors, uint256[] memory debtAmounts, uint256 totalCount);

    /**
     * @notice Returns the total number of registered payers.
     * @return count The total number of registered payers.
     */
    function getTotalPayerCount() external view returns (uint256 count);

    /**
     * @notice Returns the number of active payers.
     * @return count The number of active payers.
     */
    function getActivePayerCount() external view returns (uint256 count);

    /**
     * @notice Returns the timestamp of the last fee transfer to the rewards contract.
     * @return timestamp The last fee transfer timestamp.
     */
    function getLastFeeTransferTimestamp() external view returns (uint256 timestamp);

    /**
     * @notice Returns the total value locked in the contract (all payer balances).
     * @return tvl The total value locked in USDC.
     */
    function getTotalValueLocked() external view returns (uint256 tvl);

    /**
     * @notice Returns the total outstanding debt amount across all payers.
     * @return totalDebt The total debt amount in USDC.
     */
    function getTotalDebtAmount() external view returns (uint256 totalDebt);

    /**
     * @notice Returns the actual USDC balance held by the contract.
     * @dev    This can be used to verify the contract's accounting is accurate.
     * @return balance The USDC token balance of the contract.
     */
    function getContractBalance() external view returns (uint256 balance);

    /**
     * @notice Retrieves the address of the current distribution contract.
     * @return distributionContract The address of the distribution contract.
     */
    function getDistributionContract() external view returns (address distributionContract);

    /**
     * @notice Retrieves the address of the current nodes contract.
     * @return nodesContract The address of the nodes contract.
     */
    function getNodesContract() external view returns (address nodesContract);

    /**
     * @notice Retrieves the address of the current payer report contract.
     * @return payerReportContract The address of the payer report contract.
     */
    function getPayerReportContract() external view returns (address payerReportContract);

    /**
     * @notice Retrieves the minimum deposit amount required to register as a payer.
     * @return minimumDeposit The minimum deposit amount in USDC.
     */
    function getMinimumDeposit() external view returns (uint256 minimumDeposit);

    /**
     * @notice Retrieves the minimum deposit amount required to register as a payer.
     * @return minimumRegistrationAmount The minimum deposit amount in USDC.
     */
    function getMinimumRegistrationAmount() external view returns (uint256 minimumRegistrationAmount);

    /**
     * @notice Retrieves the current total balance of a given payer.
     * @param  payer   The address of the payer.
     * @return balance The current balance of the payer.
     */
    function getPayerBalance(address payer) external view returns (uint256 balance);

    /**
     * @notice Returns the duration of the lock period required before a withdrawal
     *         can be finalized.
     * @return lockPeriod The lock period in seconds.
     */
    function getWithdrawalLockPeriod() external view returns (uint256 lockPeriod);

    /**
     * @notice Retrieves the total pending fees that have not yet been transferred
     *         to the distribution contract.
     * @return fees The total pending fees in USDC.
     */
    function getPendingFees() external view returns (uint256 fees);

    /**
     * @notice Returns the maximum allowed time difference for backdated settlements.
     * @return maxTime The maximum allowed time difference in seconds.
     */
    function getMaxBackdatedTime() external view returns (uint256 maxTime);
}