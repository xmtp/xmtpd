// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

/**
 * @title IPayer
 * @notice Interface for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process.
 */
interface IPayer {
    //==============================================================
    //                             STRUCTS
    //==============================================================

    /**
     * @dev Struct to store payer information.
     * @param balance The current USDC balance of the payer.
     * @param isActive Indicates whether the payer is active.
     * @param creationTimestamp The timestamp when the payer was first registered.
     * @param latestDepositTimestamp The timestamp of the most recent deposit.
     * @param debtAmount The amount of fees owed but not yet settled.
     */
    struct Payer {
        uint256 balance;
        bool isActive;
        uint256 creationTimestamp;
        uint256 latestDepositTimestamp;
        uint256 debtAmount;
    }

    /**
     * @dev Struct to store withdrawal request information.
     * @param requestTimestamp The timestamp when the withdrawal was requested.
     * @param withdrawableTimestamp The timestamp when the withdrawal can be finalized.
     * @param amount The amount requested for withdrawal.
     */
    struct Withdrawal {
        uint256 requestTimestamp;
        uint256 withdrawableTimestamp;
        uint256 amount;
    }

    //==============================================================
    //                             EVENTS
    //==============================================================

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

    /// @dev Emitted when a payer initiates a withdrawal request.
    event WithdrawalRequest(address indexed payer, uint256 requestTimestamp, uint256 withdrawableTimestamp, uint256 amount);

    /// @dev Emitted when a payer cancels a withdrawal request.
    event WithdrawalCancelled(address indexed payer);

    /// @dev Emitted when a payer's withdrawal is finalized.
    event WithdrawalFinalized(address indexed payer, uint256 amountReturned);

    /// @dev Emitted when usage is settled and fees are calculated.
    event UsageSettled(uint256 fees, address indexed payer, uint256 indexed nodeId, uint256 timestamp);

    /// @dev Emitted when fees are transferred to the rewards contract.
    event FeesTransferred(uint256 amount);

    /// @dev Emitted when the rewards contract address is updated.
    event RewardsContractUpdated(address indexed newRewardsContract);

    /// @dev Emitted when the nodes contract address is updated.
    event NodesContractUpdated(address indexed newNodesContract);

    /// @dev Emitted when the minimum deposit amount is updated.
    event MinimumDepositUpdated(uint256 newMinimumDeposit);

    /// @dev Emitted when the pause is triggered by `account`.
    event Paused(address account);

    /// @dev Emitted when the pause is lifted by `account`.
    event Unpaused(address account);

    //==============================================================
    //                             ERRORS
    //==============================================================

    /// @dev Error thrown when caller is not an authorized node operator.
    error UnauthorizedNodeOperator();

    /// @dev Error thrown when caller is not the rewards contract.
    error NotRewardsContract();

    /// @dev Error thrown when an address is invalid (usually zero address).
    error InvalidAddress();

    /// @dev Error thrown when the amount is insufficient.
    error InsufficientAmount();

    /// @dev Error thrown when a withdrawal is not in the requested state.
    error WithdrawalNotRequested();

    /// @dev Error thrown when a withdrawal is already in progress.
    error WithdrawalAlreadyRequested();

    /// @dev Error thrown when a lock period has not yet elapsed.
    error LockPeriodNotElapsed();

    /// @dev Error thrown when arrays have mismatched lengths.
    error ArrayLengthMismatch();

    /// @dev Error thrown when trying to backdate settlement too far.
    error InvalidSettlementTime();

    /// @dev Error thrown when trying to delete a payer with balance or debt.
    error PayerHasBalanceOrDebt();

    /// @dev Error thrown when trying to delete a payer in withdrawal state.
    error PayerInWithdrawal();

    //==============================================================
    //                      PAYER REGISTRATION & MANAGEMENT
    //==============================================================

    /**
     * @notice Registers the caller as a new payer upon depositing the minimum required USDC.
     *         The caller must approve this contract to spend USDC beforehand.
     * @param amount The amount of USDC to deposit (must be at least the minimum required).
     *
     * Emits `PayerRegistered`.
     */
    function register(uint256 amount) external;

    /**
     * @notice Allows the caller to deposit USDC into their own payer account.
     *         The caller must approve this contract to spend USDC beforehand.
     * @param amount The amount of USDC to deposit.
     *
     * Emits `Deposit`.
     */
    function deposit(uint256 amount) external;

    /**
     * @notice Allows anyone to donate USDC to an existing payer's account.
     *         The sender must approve this contract to spend USDC beforehand.
     * @param payer The address of the payer receiving the donation.
     * @param amount The amount of USDC to donate.
     *
     * Emits `Donation`.
     */
    function donate(address payer, uint256 amount) external;

    /**
     * @notice Deactivates a payer, preventing them from initiating new transactions.
     *         Only callable by authorized node operators.
     * @param payer The address of the payer to deactivate.
     *
     * Emits `PayerDeactivated`.
     */
    function deactivatePayer(address payer) external;

    /**
     * @notice Permanently deletes a payer from the system.
     * @dev Can only delete payers with zero balance and zero debt who are not in withdrawal.
     *      Only callable by authorized node operators.
     * @param payer The address of the payer to delete.
     *
     * Emits `PayerDeleted`.
     */
    function deletePayer(address payer) external;

    /**
     * @notice Checks if a given address is an active payer.
     * @param payer The address to check.
     * @return isActive True if the address is an active payer, false otherwise.
     */
    function getIsActivePayer(address payer) external view returns (bool isActive);

    /**
     * @notice Retrieves the minimum deposit amount required to register as a payer.
     * @return minimumDeposit The minimum deposit amount in USDC.
     */
    function getMinimumDeposit() external view returns (uint256 minimumDeposit);

    /**
     * @notice Updates the minimum deposit amount required for registration.
     * @param newMinimumDeposit The new minimum deposit amount.
     * 
     * Emits `MinimumDepositUpdated`.
     */
    function setMinimumDeposit(uint256 newMinimumDeposit) external;

    //==============================================================
    //                      PAYER BALANCE MANAGEMENT
    //==============================================================

    /**
     * @notice Retrieves the current total balance of a given payer.
     * @param payer The address of the payer.
     * @return balance The current balance of the payer.
     */
    function getPayerBalance(address payer) external view returns (uint256 balance);

    /**
     * @notice Initiates a withdrawal request for the caller.
     *         - Sets the payer into withdrawal mode (no further usage allowed).
     *         - Records a timestamp for the withdrawal lock period.
     * @param amount The amount to withdraw (can be less than or equal to current balance).
     *
     * Emits `WithdrawalRequest`.
     */
    function requestWithdrawal(uint256 amount) external;

    /**
     * @notice Cancels a previously requested withdrawal, removing withdrawal mode.
     * @dev Only callable by the payer who initiated the withdrawal.
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
     * @param payer The address to check.
     * @return inWithdrawal True if in withdrawal mode, false otherwise.
     * @return requestTimestamp The timestamp when `requestWithdrawal()` was called.
     * @return withdrawableTimestamp When the withdrawal can be finalized.
     * @return amount The amount requested for withdrawal.
     */
    function getWithdrawalStatus(address payer)
        external
        view
        returns (bool inWithdrawal, uint256 requestTimestamp, uint256 withdrawableTimestamp, uint256 amount);

    /**
     * @notice Returns the duration of the lock period required before a withdrawal
     *         can be finalized.
     * @return The lock period in seconds.
     */
    function getWithdrawalLockPeriod() external view returns (uint256);

    //==============================================================
    //                       USAGE SETTLEMENT
    //==============================================================

    /**
     * @notice Settles a contiguous batch of usage data from a PayerReport.
     * Uses an aggregated Merkle proof to verify that the provided batch of
     * (payer, amount) entries is included in the report’s committed Merkle root.
     *
     * @param originatorNode The node that submitted the report.
     * @param reportIndex The index of the report.
     * @param offset The starting index in the list of (payer, amount) entries.
     * @param payers A contiguous array of payer addresses.
     * @param amounts A contiguous array of usage amounts corresponding to each payer.
     * @param proof An array of branch hashes for the aggregated Merkle proof.
     *
     * The contract computes the aggregated hash for the batch and, along with the
     * provided proof and offset, reconstructs the Merkle path to verify inclusion.
     */
    function settleUsage(
        address originatorNode,
        uint256 reportIndex,
        uint256 offset,
        address[] calldata payers,
        uint256[] calldata amounts,
        bytes32[] calldata proof
    ) external;

    /**
     * @notice Retrieves the total pending fees that have not yet been transferred
     *         to the rewards contract.
     * @return pending The total pending fees in USDC.
     */
    function getPendingFees() external view returns (uint256 pending);

    /**
     * @notice Transfers all pending fees to the designated rewards contract for
     *         distribution using EIP-2200 optimizations.
     * @dev Uses a single storage write for updating accumulated fees.
     *
     * Emits `FeesTransferred`.
     */
    function transferFeesToRewards() external;

    /**
     * @notice Returns the maximum allowed time difference for backdated settlements.
     * @return The maximum allowed time difference in seconds.
     */
    function getMaxBackdatedTime() external view returns (uint256);

    //==============================================================
    //                       OBSERVABILITY FUNCTIONS
    //==============================================================

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
     * @notice Returns a paginated list of payers with outstanding debt.
     * @param offset Number of payers to skip before starting to return results.
     * @param limit Maximum number of payers to return.
     * @return debtors Array of payer addresses with debt.
     * @return debtAmounts Corresponding debt amounts for each payer.
     * @return totalCount Total number of payers with debt (regardless of pagination).
     */
    function getPayersInDebt(uint256 offset, uint256 limit) external view returns (
        address[] memory debtors,
        uint256[] memory debtAmounts,
        uint256 totalCount
    );

    /**
     * @notice Returns the actual USDC balance held by the contract.
     * @dev This can be used to verify the contract's accounting is accurate.
     * @return balance The USDC token balance of the contract.
     */
    function getContractBalance() external view returns (uint256 balance);

    //==============================================================
    //                       ADMINISTRATIVE FUNCTIONS
    //==============================================================

    /**
     * @notice Sets the address of the rewards contract.
     * @param _rewardsContract The address of the new rewards contract.
     *
     * Emits `RewardsContractUpdated`.
     */
    function setRewardsContract(address _rewardsContract) external;

    /**
     * @notice Sets the address of the nodes contract for operator verification.
     * @param _nodesContract The address of the new nodes contract.
     *
     * Emits `NodesContractUpdated`.
     */
    function setNodesContract(address _nodesContract) external;

    /**
     * @notice Retrieves the address of the current rewards contract.
     * @return The address of the rewards contract.
     */
    function getRewardsContract() external view returns (address);

    /**
     * @notice Retrieves the address of the current nodes contract.
     * @return The address of the nodes contract.
     */
    function getNodesContract() external view returns (address);

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

    /**
     * @notice Checks if a given address is an active node operator.
     * @param operator The address to check.
     * @return isActiveNodeOperator True if the address is an active node operator, false otherwise.
     */
    function getIsActiveNodeOperator(address operator) external view returns (bool isActiveNodeOperator);
}