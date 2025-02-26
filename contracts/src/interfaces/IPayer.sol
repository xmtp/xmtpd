// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

/**
 * @title IPayer
 * @notice Interface for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process with optimized storage using
 *         batch hash commitments and EIP-1283 gas optimizations.
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
    struct WithdrawalRequest {
        uint256 requestTimestamp;
        uint256 withdrawableTimestamp;
        uint256 amount;
    }

    /**
     * @dev Struct to store usage reporting data.
     * @notice This struct is used primarily for event emission and off-chain tracking.
     * @dev For storage efficiency, complete reports are NOT stored on-chain; only their
     *      cryptographic hash commitments are stored.
     * @param payer The address of the payer being charged.
     * @param nodeId The ID of the node that provided service.
     * @param fees The amount charged for this usage.
     * @param timestamp When the usage occurred.
     */
    struct UsageReport {
        address payer;
        uint256 nodeId;
        uint256 fees;
        uint256 timestamp;
    }

    /**
     * @dev Struct to store batch commitment information.
     * @param batchHash Hash of the complete batch data (payers, fees, timestamp, nodeId).
     * @param totalFees The total fees collected in this batch.
     * @param nodeId The node operator who submitted this batch.
     * @param timestamp When this batch was processed.
     * @param blockNumber The block in which this commitment was stored.
     */
    struct BatchCommitment {
        bytes32 batchHash;
        uint256 totalFees;
        uint256 nodeId;
        uint256 timestamp;
        uint256 blockNumber;
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
    event DepositMade(address indexed from, address indexed payer, uint256 amount);

    /// @dev Emitted when a user donates to a payer's account.
    event DonationMade(address indexed donor, address indexed payer, uint256 amount);

    /// @dev Emitted when a payer initiates a withdrawal request.
    event WithdrawRequested(address indexed payer, uint256 requestTimestamp, uint256 withdrawableTimestamp, uint256 amount);

    /// @dev Emitted when a payer cancels a withdrawal request.
    event WithdrawCancelled(address indexed payer);

    /// @dev Emitted when a payer's withdrawal is finalized.
    event WithdrawFinalized(address indexed payer, uint256 amountReturned);

    /// @dev Emitted when usage is settled and fees are calculated.
    event UsageSettled(uint256 fees, address indexed payer, uint256 indexed nodeId, uint256 timestamp);

    /// @dev Emitted when batch usage is settled.
    event BatchUsageSettled(bytes32 batchHash, uint256 totalFees, uint256 indexed nodeId, uint256 timestamp);

    /// @dev Emitted when fees are transferred to the rewards contract.
    event FeesTransferred(uint256 amount);

    /// @dev Emitted when old batch commitments are cleared from storage.
    event CommitmentsCleared(uint256 count, uint256 oldestTimestamp);

    /// @dev Emitted when the rewards contract address is updated.
    event RewardsContractUpdated(address indexed newRewardsContract);

    /// @dev Emitted when the nodes contract address is updated.
    event NodesContractUpdated(address indexed newNodesContract);

    /// @dev Emitted when the minimum deposit amount is updated.
    event MinimumDepositUpdated(uint256 newMinimumDeposit);

    /// @dev Emitted when the retention period for commitments is updated.
    event CommitmentRetentionPeriodUpdated(uint256 newPeriod);

    /// @dev Emitted when the contract is paused.
    event Paused();

    /// @dev Emitted when the contract is unpaused.
    event Unpaused();

    //==============================================================
    //                             ERRORS
    //==============================================================

    /// @dev Error thrown when an operation is attempted while the contract is paused.
    error ContractPaused();

    /// @dev Error thrown when caller is not an authorized node operator.
    error NotAuthorizedNodeOperator();

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

    /// @dev Error thrown when a batch hash is invalid.
    error InvalidBatchHash();

    /// @dev Error thrown when a batch commitment already exists.
    error CommitmentAlreadyExists();

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
     * Emits `DepositMade`.
     */
    function deposit(uint256 amount) external;

    /**
     * @notice Allows anyone to donate USDC to an existing payer's account.
     *         The sender must approve this contract to spend USDC beforehand.
     * @param payer The address of the payer receiving the donation.
     * @param amount The amount of USDC to donate.
     *
     * Emits `DonationMade`.
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
     * @notice Allows a payer to delete their own account from the system.
     * @dev Can only be called if the payer has zero balance and zero debt and is not in withdrawal.
     *
     * Emits `PayerDeleted`.
     */
    function deleteMyAccount() external;

    /**
     * @notice Checks if a given address is an active payer.
     * @param payer The address to check.
     * @return True if the address is an active payer, false otherwise.
     */
    function isActivePayer(address payer) external view returns (bool);

    /**
     * @notice Retrieves the minimum deposit amount required to register as a payer.
     * @return The minimum deposit amount in USDC.
     */
    function getMinimumDeposit() external view returns (uint256);

    /**
     * @notice Updates the minimum deposit amount required for registration.
     * @param newMinimumDeposit The new minimum deposit amount.
     * 
     * Emits `MinimumDepositUpdated`.
     */
    function updateMinimumDeposit(uint256 newMinimumDeposit) external;

    //==============================================================
    //                      PAYER BALANCE MANAGEMENT
    //==============================================================

    /**
     * @notice Retrieves the current total balance of a given payer.
     * @param payer The address of the payer.
     * @return The current balance of the payer.
     */
    function getPayerBalance(address payer) external view returns (uint256);

    /**
     * @notice Initiates a withdrawal request for the caller.
     *         - Sets the payer into withdrawal mode (no further usage allowed).
     *         - Records a timestamp for the withdrawal lock period.
     * @param amount The amount to withdraw (can be less than or equal to current balance).
     *
     * Emits `WithdrawRequested`.
     */
    function requestWithdraw(uint256 amount) external;

    /**
     * @notice Cancels a previously requested withdrawal, removing withdrawal mode.
     * @dev Only callable by the payer who initiated the withdrawal.
     *
     * Emits `WithdrawCancelled`.
     */
    function cancelWithdraw() external;

    /**
     * @notice Finalizes a payer's withdrawal after the lock period has elapsed.
     *         - Accounts for any pending usage during the lock.
     *         - Returns the unspent balance to the payer.
     *
     * Emits `WithdrawFinalized`.
     */
    function finalizeWithdraw() external;

    /**
     * @notice Checks if a payer is currently in withdrawal mode and the timestamp
     *         when they initiated the withdrawal.
     * @param payer The address to check.
     * @return inWithdrawal True if in withdrawal mode, false otherwise.
     * @return requestTimestamp The timestamp when `requestWithdraw()` was called.
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
    function getLockPeriod() external view returns (uint256);

    //==============================================================
    //                       USAGE SETTLEMENT
    //==============================================================

    /**
     * @notice Called by node operators to settle usage and calculate fees owed.
     * @dev This function is EIP-1283 optimized by using accumulators for multiple state updates.
     * @param fees The total USDC fees computed from this usage period.
     * @param payer The address of the payer being charged.
     * @param nodeId The ID of the node operator submitting the usage.
     * @param timestamp The timestamp when the usage occurred (can be backdated).
     * @param commitmentHash Optional hash for additional off-chain data verification.
     *
     * Emits `UsageSettled`.
     */
    function settleUsage(
        uint256 fees,
        address payer,
        uint256 nodeId,
        uint256 timestamp,
        bytes32 commitmentHash
    ) external;

    /**
     * @notice Called by node operators to settle usage for multiple payers in a batch.
     * @dev Uses EIP-1283 optimizations for storage efficiency and a simple hash commitment
     *      for batch verification.
     * @param payers Array of payer addresses being charged.
     * @param fees Array of USDC fees corresponding to each payer.
     * @param timestamp When this batch of usage occurred (can be backdated).
     * @param nodeId The ID of the node operator submitting the usage.
     *
     * Emits `BatchUsageSettled` and multiple `UsageSettled` events.
     */
    function settleUsageBatch(
        address[] calldata payers,
        uint256[] calldata fees,
        uint256 timestamp,
        uint256 nodeId
    ) external;

    /**
     * @notice Verifies if a specific batch matches the stored commitment hash.
     * @param batchHash The commitment hash stored on-chain.
     * @param payers Array of payer addresses to verify.
     * @param fees Array of fee amounts to verify.
     * @param timestamp The timestamp to verify.
     * @param nodeId The node ID to verify.
     * @return True if the batch data matches the stored hash, false otherwise.
     */
    function verifyBatch(
        bytes32 batchHash,
        address[] calldata payers,
        uint256[] calldata fees,
        uint256 timestamp,
        uint256 nodeId
    ) external pure returns (bool);

    /**
     * @notice Computes the hash for a batch of usage data.
     * @param payers Array of payer addresses.
     * @param fees Array of fee amounts.
     * @param timestamp The timestamp of the batch.
     * @param nodeId The node ID that generated the batch.
     * @return The computed batch hash.
     */
    function computeBatchHash(
        address[] calldata payers,
        uint256[] calldata fees,
        uint256 timestamp,
        uint256 nodeId
    ) external pure returns (bytes32);

    /**
     * @notice Retrieves information about a batch commitment.
     * @param batchHash The hash of the batch to query.
     * @return exists Whether this commitment exists.
     * @return totalFees Total fees in this batch.
     * @return nodeId The node that created this batch.
     * @return timestamp When the batch was processed.
     * @return blockNumber The block where this commitment was stored.
     */
    function getBatchCommitmentInfo(bytes32 batchHash) external view returns (
        bool exists,
        uint256 totalFees,
        uint256 nodeId,
        uint256 timestamp,
        uint256 blockNumber
    );

    /**
     * @notice Clears old batch commitments to optimize storage (uses EIP-1283 refunds).
     * @dev Only clears commitments older than the retention period.
     * @param hashes Array of batch hashes to check and potentially clear.
     * @param maxToClear Maximum number of commitments to clear in this call.
     * @return cleared The number of commitments that were cleared.
     *
     * Emits `CommitmentsCleared`.
     */
    function clearOldCommitments(bytes32[] calldata hashes, uint256 maxToClear) external returns (uint256 cleared);

    /**
     * @notice Sets the retention period for batch commitments.
     * @param newPeriod The new retention period in seconds.
     *
     * Emits `CommitmentRetentionPeriodUpdated`.
     */
    function setCommitmentRetentionPeriod(uint256 newPeriod) external;

    /**
     * @notice Retrieves the total pending fees that have not yet been transferred
     *         to the rewards contract.
     * @return pending The total pending fees in USDC.
     */
    function pendingFees() external view returns (uint256 pending);

    /**
     * @notice Transfers all pending fees to the designated rewards contract for
     *         distribution using EIP-1283 optimizations.
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

    /**
     * @notice Returns the current commitment retention period.
     * @return The retention period in seconds.
     */
    function getCommitmentRetentionPeriod() external view returns (uint256);

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
     * @notice Returns the total amount of fees collected since contract deployment.
     * @return amount The total amount of fees ever collected.
     */
    function getLifetimeFeesCollected() external view returns (uint256 amount);

    /**
     * @notice Returns the current amount locked in withdrawal requests.
     * @return amount The total amount in pending withdrawal requests.
     */
    function getTotalWithdrawalRequests() external view returns (uint256 amount);

    /**
     * @notice Returns the timestamp of the last fee transfer to the rewards contract.
     * @return timestamp The last fee transfer timestamp.
     */
    function getLastFeeTransferTimestamp() external view returns (uint256 timestamp);

    /**
     * @notice Returns historical usage statistics for a specific node.
     * @param nodeId The ID of the node to query.
     * @return totalFees Total fees generated by this node.
     * @return lastSettlementTime The last time this node settled usage.
     */
    function getNodeStatistics(uint256 nodeId) external view returns (
        uint256 totalFees,
        uint256 lastSettlementTime
    );

    /**
     * @notice Returns a list of payers with outstanding debt.
     * @param limit Maximum number of payers to return.
     * @return debtors Array of payer addresses with debt.
     * @return debtAmounts Corresponding debt amounts for each payer.
     */
    function getPayersInDebt(uint256 limit) external view returns (
        address[] memory debtors,
        uint256[] memory debtAmounts
    );

    /**
     * @notice Returns the actual USDC balance held by the contract.
     * @dev This can be used to verify the contract's accounting is accurate.
     * @return balance The USDC token balance of the contract.
     */
    function getContractBalance() external view returns (uint256 balance);

    /**
     * @notice Returns the total amount of fees transferred to rewards.
     * @return amount The total amount of fees sent to rewards contract.
     */
    function getTotalFeesTransferred() external view returns (uint256 amount);

    /**
     * @notice Returns the number of active batch commitments stored.
     * @return count The number of active batch commitments.
     */
    function getActiveBatchCommitmentCount() external view returns (uint256 count);

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
     * @return True if the address is an active node operator, false otherwise.
     */
    function isActiveNodeOperator(address operator) external view returns (bool);

    /**
     * @notice Checks if the contract is currently paused.
     * @return True if the contract is paused, false otherwise.
     */
    function isPaused() external view returns (bool);
}
