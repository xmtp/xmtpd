// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {AccessControlUpgradeable} from "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {Initializable} from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import {PausableUpgradeable} from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin-contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {UUPSUpgradeable} from "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {INodes} from "./interfaces/INodes.sol";
import {IPayer} from "./interfaces/IPayer.sol";
import {IPayerReport} from "./interfaces/IPayerReport.sol";

/**
 * @title  Payer
 * @notice Implementation for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process.
 */
contract Payer is
    Initializable,
    AccessControlUpgradeable,
    UUPSUpgradeable,
    PausableUpgradeable,
    ReentrancyGuardUpgradeable,
    IPayer
{
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ============ Constants ============ */

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    string internal constant USDC_SYMBOL = "USDC";
    uint64 private constant DEFAULT_MINIMUM_REGISTRATION_AMOUNT_MICRO_DOLLARS = 10_000_000; // 10 USD
    uint64 private constant DEFAULT_MINIMUM_DEPOSIT_AMOUNT_MICRO_DOLLARS = 10_000_000;      // 10 USD
    uint64 private constant DEFAULT_MAX_TOLERABLE_DEBT_AMOUNT_MICRO_DOLLARS = 50_000_000;   // 50 USD
    uint32 private constant DEFAULT_MINIMUM_TRANSFER_FEES_PERIOD = 6 hours;
    uint32 private constant ABSOLUTE_MINIMUM_TRANSFER_FEES_PERIOD = 1 hours;
    uint32 private constant DEFAULT_WITHDRAWAL_LOCK_PERIOD = 3 days;
    uint32 private constant ABSOLUTE_MINIMUM_WITHDRAWAL_LOCK_PERIOD = 1 days;

    /* ============ UUPS Storage ============ */

    /// @custom:storage-location erc7201:xmtp.storage.Payer
    struct PayerStorage {
        /// @dev Contracts to interact with.
        IERC20 usdcToken;
        address distributionContract;
        address nodesContract;
        address payerReportContract;
        /// @dev Configuration parameters
        uint64 minimumRegistrationAmountMicroDollars;
        uint64 minimumDepositAmountMicroDollars;
        uint64 maxTolerableDebtAmountMicroDollars;
        uint32 withdrawalLockPeriod;
        uint32 transferFeesPeriod;
        /// @dev State variables
        uint256 lastFeeTransferTimestamp;
        uint256 totalAmountDeposited;
        uint256 totalDebtAmount;
        uint256 pendingFees;
        uint256 collectedFees;
        /// @dev Mappings
        mapping(address => Payer) payers;
        mapping(address => Withdrawal) withdrawals;
        EnumerableSet.AddressSet totalPayers;
        EnumerableSet.AddressSet activePayers;
        EnumerableSet.AddressSet debtPayers;
    }

    // keccak256(abi.encode(uint256(keccak256("xmtp.storage.Payer")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 internal constant PAYER_STORAGE_LOCATION = 0xd0335f337c570f3417b0f0d20340c88da711d60e810b5e9b3ecabe9ccfcdce5a;

    function _getPayerStorage() internal pure returns (PayerStorage storage $) {
        // slither-disable-next-line assembly
        assembly {
            $.slot := PAYER_STORAGE_LOCATION
        }
    }

    /* ============ Modifiers ============ */

    /**
     * @dev Modifier to check if caller is an active node operator.
     */
    modifier onlyNodeOperator() {
        if (!_getIsActiveNodeOperator(msg.sender)) {
            revert UnauthorizedNodeOperator();
        }
        _;
    }

    /**
     * @dev Modifier to check if caller is the payer report contract.
     */
    modifier onlyPayerReport() {
        if (msg.sender != _getPayerStorage().payerReportContract) {
            revert Unauthorized();
        }
        _;
    }

    /**
     * @dev Modifier to check if address is an active payer.
     */
    modifier onlyPayer(address payer) {
        require(_payerExists(payer), PayerDoesNotExist());
        require(msg.sender == payer, Unauthorized());
        _;
    }

    /* ============ Initialization ============ */

    /**
     * @notice Initializes the contract with the deployer as admin.
     * @param  initialAdmin The address of the admin.
     * @dev    There's a chicken-egg problem here with PayerReport and Distribution contracts.
     *         We need to deploy these contracts first, then set their addresses
     *         in the Payer contract.
     */
    function initialize(address initialAdmin, address usdcToken, address nodesContract) public initializer {
        if (initialAdmin == address(0) || usdcToken == address(0) || nodesContract == address(0)) {
            revert InvalidAddress();
        }

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        PayerStorage storage $ = _getPayerStorage();

        $.minimumRegistrationAmountMicroDollars = DEFAULT_MINIMUM_REGISTRATION_AMOUNT_MICRO_DOLLARS;
        $.minimumDepositAmountMicroDollars = DEFAULT_MINIMUM_DEPOSIT_AMOUNT_MICRO_DOLLARS;
        $.withdrawalLockPeriod = DEFAULT_WITHDRAWAL_LOCK_PERIOD;
        $.maxTolerableDebtAmountMicroDollars = DEFAULT_MAX_TOLERABLE_DEBT_AMOUNT_MICRO_DOLLARS;
        $.transferFeesPeriod = DEFAULT_MINIMUM_TRANSFER_FEES_PERIOD;

        _setUsdcTokenContract(usdcToken);
        _setNodesContract(nodesContract);

        require(_grantRole(DEFAULT_ADMIN_ROLE, initialAdmin), FailedToGrantRole(DEFAULT_ADMIN_ROLE, initialAdmin));
        require(_grantRole(ADMIN_ROLE, initialAdmin), FailedToGrantRole(ADMIN_ROLE, initialAdmin));
    }

    /* ============ Payers Management ============ */

    /**
     * @inheritdoc IPayer
     */
    function register(uint256 amount) external whenNotPaused {
        PayerStorage storage $ = _getPayerStorage();

        require(amount >= $.minimumRegistrationAmountMicroDollars, InsufficientAmount());
        require(!_payerExists(msg.sender), PayerAlreadyRegistered());

        _deposit(msg.sender, amount);

        // New payer registration
        $.payers[msg.sender] = Payer({
            balance: amount,
            debtAmount: 0,
            creationTimestamp: block.timestamp,
            latestDepositTimestamp: block.timestamp,
            latestDonationTimestamp: 0,
            isActive: true
        });

        // Add new payer to active and total payers sets
        require($.activePayers.add(msg.sender), FailedToRegisterPayer());
        require($.totalPayers.add(msg.sender), FailedToRegisterPayer());

        _increaseTotalAmountDeposited(amount);

        emit PayerRegistered(msg.sender, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deposit(uint256 amount) external whenNotPaused nonReentrant onlyPayer(msg.sender) {
        PayerStorage storage $ = _getPayerStorage();

        require(amount >= $.minimumDepositAmountMicroDollars, InsufficientAmount());
        require($.withdrawals[msg.sender].requestTimestamp == 0, PayerInWithdrawal());

        _deposit(msg.sender, amount);

        _updatePayerBalance(msg.sender, amount);

        emit PayerBalanceUpdated(msg.sender, $.payers[msg.sender].balance, $.payers[msg.sender].debtAmount);
    }

    /**
     * @inheritdoc IPayer
     */
    function donate(address payer, uint256 amount) external whenNotPaused {
        require(amount > 0, InsufficientAmount());
        _revertIfPayerDoesNotExist(payer);

        PayerStorage storage $ = _getPayerStorage();

        require($.withdrawals[payer].requestTimestamp == 0, PayerInWithdrawal());

        _deposit(msg.sender, amount);

        _updatePayerBalance(payer, amount);

        $.payers[payer].latestDonationTimestamp = block.timestamp;

        emit PayerBalanceUpdated(payer, $.payers[payer].balance, $.payers[payer].debtAmount);
        emit Donation(msg.sender, payer, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deactivatePayer(address payer) external whenNotPaused onlyNodeOperator {
        _revertIfPayerDoesNotExist(payer);

        _deactivatePayer(payer);
    }

    /**
     * @inheritdoc IPayer
     */
    function deletePayer(address payer) external whenNotPaused onlyRole(ADMIN_ROLE) {
        _revertIfPayerDoesNotExist(payer);

        PayerStorage storage $ = _getPayerStorage();

        require($.withdrawals[payer].requestTimestamp == 0, PayerInWithdrawal());

        if ($.payers[payer].balance > 0 || $.payers[payer].debtAmount > 0) {
            revert PayerHasBalanceOrDebt();
        }

        // Delete all payer data
        delete $.payers[payer];
        require($.totalPayers.remove(payer), FailedToDeletePayer());
        require($.activePayers.remove(payer), FailedToDeletePayer());

        emit PayerDeleted(payer, block.timestamp);
    }

    /* ========== Payers Balance Management ========= */

    /**
     * @inheritdoc IPayer
     */
    function requestWithdrawal(uint256 amount) external whenNotPaused onlyPayer(msg.sender) {
        if (_withdrawalExists(msg.sender)) revert WithdrawalAlreadyRequested();

        PayerStorage storage $ = _getPayerStorage();

        require($.payers[msg.sender].debtAmount == 0, PayerHasDebt());
        require($.payers[msg.sender].balance >= amount, InsufficientBalance());

        // Balance to be withdrawn is deducted from the payer's balance,
        // it can't be used to settle payments.
        $.payers[msg.sender].balance -= amount;
        _decreaseTotalAmountDeposited(amount);

        uint256 withdrawableTimestamp = block.timestamp + $.withdrawalLockPeriod;

        $.withdrawals[msg.sender] = Withdrawal({
            requestTimestamp: block.timestamp,
            withdrawableTimestamp: withdrawableTimestamp,
            amount: amount
        });

        emit WithdrawalRequested(msg.sender, block.timestamp, withdrawableTimestamp, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function cancelWithdrawal() external whenNotPaused onlyPayer(msg.sender) {
        _revertIfWithdrawalNotExists(msg.sender);

        PayerStorage storage $ = _getPayerStorage();

        Withdrawal memory withdrawal = $.withdrawals[msg.sender];

        delete $.withdrawals[msg.sender];

        $.payers[msg.sender].balance += withdrawal.amount;
        _increaseTotalAmountDeposited(withdrawal.amount);

        emit WithdrawalCancelled(msg.sender, withdrawal.requestTimestamp);
    }

    /**
     * @inheritdoc IPayer
     */
    function finalizeWithdrawal() external whenNotPaused nonReentrant onlyPayer(msg.sender) {
        _revertIfWithdrawalNotExists(msg.sender);

        PayerStorage storage $ = _getPayerStorage();

        Withdrawal memory withdrawal = $.withdrawals[msg.sender];

        delete $.withdrawals[msg.sender];

        uint256 finalWithdrawalAmount = withdrawal.amount;

        if ($.payers[msg.sender].debtAmount > 0) {
            finalWithdrawalAmount = _settleDebts(msg.sender, withdrawal.amount);
        }

        if (finalWithdrawalAmount > 0) {
            $.usdcToken.safeTransfer(msg.sender, finalWithdrawalAmount);
        }

        emit WithdrawalFinalized(msg.sender, withdrawal.requestTimestamp);
    }

    /**
     * @inheritdoc IPayer
     */
    function getWithdrawalStatus(address payer) external view returns (Withdrawal memory withdrawal) {
        _revertIfPayerDoesNotExist(payer);

        return _getPayerStorage().withdrawals[payer];
    }

    /* ============ Usage Settlement ============ */

    /**
     * @inheritdoc IPayer
     */
    function settleUsage(uint256 originatorNode, address[] calldata payerList, uint256[] calldata usageAmountsList)
        external
        whenNotPaused
        nonReentrant
        onlyPayerReport
    {
        require(payerList.length == usageAmountsList.length, InvalidPayerListLength());

        PayerStorage storage $ = _getPayerStorage();

        uint256 fees = 0;

        for (uint256 i = 0; i < payerList.length; i++) {
            address payer = payerList[i];
            uint256 usage = usageAmountsList[i];

            // This should never happen, as PayerReport has already verified the payers and amounts.
            // Payers in payerList should always exist and be active.
            if (!_payerExists(payer) || !_payerIsActive(payer)) {
                continue;
            }

            Payer memory storedPayer = $.payers[payer];

            if (storedPayer.balance < usage) {
                uint256 debt = usage - storedPayer.balance;

                $.pendingFees += storedPayer.balance;
                fees += storedPayer.balance;

                storedPayer.balance = 0;
                storedPayer.debtAmount = debt;
                $.payers[payer] = storedPayer;

                _addDebtor(payer);
                _increaseTotalDebtAmount(debt);

                if (debt > $.maxTolerableDebtAmountMicroDollars) _deactivatePayer(payer);

                emit PayerBalanceUpdated(payer, storedPayer.balance, storedPayer.debtAmount);

                continue;
            }

            $.pendingFees += usage;
            fees += usage;

            storedPayer.balance -= usage;

            $.payers[payer] = storedPayer;

            emit PayerBalanceUpdated(payer, storedPayer.balance, storedPayer.debtAmount);
        }

        emit UsageSettled(originatorNode, block.timestamp, fees);
    }

    /**
     * @inheritdoc IPayer
     */
    function transferFeesToDistribution() external whenNotPaused nonReentrant {
        PayerStorage storage $ = _getPayerStorage();

        /// @dev slither marks this as a security issue because validators can modify block.timestamp.
        ///      However, in this scenario it's fine, as we'd just send fees a earlier than expected.
        ///      It would be a bigger issue if we'd rely on timestamp for randomness or calculations.
        // slither-disable-next-line timestamp
        require(block.timestamp - $.lastFeeTransferTimestamp >= $.transferFeesPeriod, InsufficientTimePassed());

        require($.pendingFees > 0, InsufficientAmount());

        uint256 feesToTransfer = $.pendingFees;

        $.usdcToken.safeTransfer($.distributionContract, feesToTransfer);

        $.lastFeeTransferTimestamp = block.timestamp;
        $.collectedFees += $.pendingFees;
        $.pendingFees = 0;

        emit FeesTransferred(block.timestamp, feesToTransfer);
    }

    /* ========== Administrative Functions ========== */

    /**
     * @inheritdoc IPayer
     */
    function setDistributionContract(address newDistributionContract) external onlyRole(ADMIN_ROLE) {
        _setDistributionContract(newDistributionContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setPayerReportContract(address newPayerReportContract) external onlyRole(ADMIN_ROLE) {
        _setPayerReportContract(newPayerReportContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setNodesContract(address newNodesContract) external onlyRole(ADMIN_ROLE) {
        _setNodesContract(newNodesContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setUsdcToken(address newUsdcToken) external onlyRole(ADMIN_ROLE) {
        _setUsdcTokenContract(newUsdcToken);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMinimumDeposit(uint64 newMinimumDeposit) external onlyRole(ADMIN_ROLE) {
        require(newMinimumDeposit > DEFAULT_MINIMUM_DEPOSIT_AMOUNT_MICRO_DOLLARS, InvalidMinimumDeposit());

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldMinimumDeposit = $.minimumDepositAmountMicroDollars;
        $.minimumDepositAmountMicroDollars = newMinimumDeposit;

        emit MinimumDepositSet(oldMinimumDeposit, newMinimumDeposit);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMinimumRegistrationAmount(uint64 newMinimumRegistrationAmount) external onlyRole(ADMIN_ROLE) {
        require(
            newMinimumRegistrationAmount > DEFAULT_MINIMUM_REGISTRATION_AMOUNT_MICRO_DOLLARS,
            InvalidMinimumRegistrationAmount()
        );

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldMinimumRegistrationAmount = $.minimumRegistrationAmountMicroDollars;
        $.minimumRegistrationAmountMicroDollars = newMinimumRegistrationAmount;

        emit MinimumRegistrationAmountSet(oldMinimumRegistrationAmount, newMinimumRegistrationAmount);
    }

    /**
     * @inheritdoc IPayer
     */
    function setWithdrawalLockPeriod(uint32 newWithdrawalLockPeriod) external onlyRole(ADMIN_ROLE) {
        require(newWithdrawalLockPeriod >= ABSOLUTE_MINIMUM_WITHDRAWAL_LOCK_PERIOD, InvalidWithdrawalLockPeriod());

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldWithdrawalLockPeriod = $.withdrawalLockPeriod;
        $.withdrawalLockPeriod = newWithdrawalLockPeriod;

        emit WithdrawalLockPeriodSet(oldWithdrawalLockPeriod, newWithdrawalLockPeriod);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMaxTolerableDebtAmount(uint64 newMaxTolerableDebtAmountMicroDollars) external onlyRole(ADMIN_ROLE) {
        require(newMaxTolerableDebtAmountMicroDollars > 0, InvalidMaxTolerableDebtAmount());

        PayerStorage storage $ = _getPayerStorage();

        uint64 oldMaxTolerableDebtAmount = $.maxTolerableDebtAmountMicroDollars;
        $.maxTolerableDebtAmountMicroDollars = newMaxTolerableDebtAmountMicroDollars;

        emit MaxTolerableDebtAmountSet(oldMaxTolerableDebtAmount, newMaxTolerableDebtAmountMicroDollars);
    }

    /**
     * @inheritdoc IPayer
     */
    function setTransferFeesPeriod(uint32 newTransferFeesPeriod) external onlyRole(ADMIN_ROLE) {
        require(newTransferFeesPeriod >= ABSOLUTE_MINIMUM_TRANSFER_FEES_PERIOD, InvalidTransferFeesPeriod());

        PayerStorage storage $ = _getPayerStorage();

        uint32 oldTransferFeesPeriod = $.transferFeesPeriod;
        $.transferFeesPeriod = newTransferFeesPeriod;

        emit TransferFeesPeriodSet(oldTransferFeesPeriod, newTransferFeesPeriod);
    }

    /**
     * @inheritdoc IPayer
     */
    function pause() external onlyRole(ADMIN_ROLE) {
        _pause();
    }

    /**
     * @inheritdoc IPayer
     */
    function unpause() external onlyRole(ADMIN_ROLE) {
        _unpause();
    }

    /* ============ Getters ============ */

    /**
     * @inheritdoc IPayer
     */
    function getPayer(address payer) external view returns (Payer memory payerInfo) {
        _revertIfPayerDoesNotExist(payer);

        return _getPayerStorage().payers[payer];
    }

    /**
     * @inheritdoc IPayer
     */
    function getActivePayers(uint256 offset, uint256 limit)
        external
        view
        returns (Payer[] memory payers, bool hasMore)
    {
        PayerStorage storage $ = _getPayerStorage();

        (address[] memory payerAddresses, bool _hasMore) = _getPaginatedAddresses($.activePayers, offset, limit);

        payers = new Payer[](payerAddresses.length);
        for (uint256 i = 0; i < payerAddresses.length; i++) {
            payers[i] = $.payers[payerAddresses[i]];
        }

        return (payers, _hasMore);
    }

    /**
     * @inheritdoc IPayer
     */
    function getIsActivePayer(address payer) public view returns (bool isActive) {
        _revertIfPayerDoesNotExist(payer);

        return _getPayerStorage().payers[payer].isActive;
    }

    /**
     * @inheritdoc IPayer
     */
    function getPayerBalance(address payer) external view returns (uint256 balance) {
        _revertIfPayerDoesNotExist(payer);

        return _getPayerStorage().payers[payer].balance;
    }

    /**
     * @inheritdoc IPayer
     */
    function getPayersInDebt(uint256 offset, uint256 limit)
        external
        view
        returns (Payer[] memory payers, bool hasMore)
    {
        PayerStorage storage $ = _getPayerStorage();

        (address[] memory payerAddresses, bool _hasMore) = _getPaginatedAddresses($.debtPayers, offset, limit);

        payers = new Payer[](payerAddresses.length);
        for (uint256 i = 0; i < payerAddresses.length; i++) {
            payers[i] = $.payers[payerAddresses[i]];
        }

        return (payers, _hasMore);
    }

    /**
     * @inheritdoc IPayer
     */
    function getTotalPayerCount() external view returns (uint256 count) {
        return _getPayerStorage().totalPayers.length();
    }

    /**
     * @inheritdoc IPayer
     */
    function getActivePayerCount() external view returns (uint256 count) {
        return _getPayerStorage().activePayers.length();
    }

    /**
     * @inheritdoc IPayer
     */
    function getLastFeeTransferTimestamp() external view returns (uint256 timestamp) {
        return _getPayerStorage().lastFeeTransferTimestamp;
    }

    /**
     * @inheritdoc IPayer
     */
    function getTotalValueLocked() external view returns (uint256 totalValueLocked) {
        PayerStorage storage $ = _getPayerStorage();

        if ($.totalDebtAmount > $.totalAmountDeposited) return 0;

        return $.totalAmountDeposited - $.totalDebtAmount;
    }

    /**
     * @inheritdoc IPayer
     */
    function getTotalDebtAmount() external view returns (uint256 totalDebt) {
        return _getPayerStorage().totalDebtAmount;
    }

    /**
     * @inheritdoc IPayer
     */
    function getContractBalance() external view returns (uint256 balance) {
        return _getPayerStorage().usdcToken.balanceOf(address(this));
    }

    /**
     * @inheritdoc IPayer
     */
    function getDistributionContract() external view returns (address distributionContractAddress) {
        return _getPayerStorage().distributionContract;
    }

    /**
     * @inheritdoc IPayer
     */
    function getNodesContract() external view returns (address nodesContractAddress) {
        return _getPayerStorage().nodesContract;
    }

    /**
     * @inheritdoc IPayer
     */
    function getPayerReportContract() external view returns (address payerReportContractAddress) {
        return _getPayerStorage().payerReportContract;
    }

    /**
     * @inheritdoc IPayer
     */
    function getMinimumDeposit() external view returns (uint256 minimumDeposit) {
        return _getPayerStorage().minimumDepositAmountMicroDollars;
    }

    /**
     * @inheritdoc IPayer
     */
    function getMinimumRegistrationAmount() external view returns (uint256 minimumRegistrationAmount) {
        return _getPayerStorage().minimumRegistrationAmountMicroDollars;
    }

    /**
     * @inheritdoc IPayer
     */
    function getWithdrawalLockPeriod() external view returns (uint256 lockPeriod) {
        return _getPayerStorage().withdrawalLockPeriod;
    }

    /**
     * @inheritdoc IPayer
     */
    function getPendingFees() external view returns (uint256 fees) {
        return _getPayerStorage().pendingFees;
    }

    /* ============ Internal ============ */

    /**
     * @notice Deposits USDC from a payer to the contract.
     * @param  payer The address of the payer.
     * @param  amount The amount to deposit.
     */
    function _deposit(address payer, uint256 amount) internal {
        PayerStorage storage $ = _getPayerStorage();

        $.usdcToken.safeTransferFrom(payer, address(this), amount);
    }

    /**
     * @notice Updates a payer's balance, handling debt settlement if applicable.
     * @param  payerAddress The address of the payer.
     * @param  amount The amount to add to the payer's balance.
     * @return leftoverAmount Amount remaining after debt settlement (if any).
     */
    function _updatePayerBalance(address payerAddress, uint256 amount) internal returns (uint256 leftoverAmount) {
        PayerStorage storage $ = _getPayerStorage();

        if ($.payers[payerAddress].debtAmount > 0) {
            return _settleDebts(payerAddress, amount);
        } else {
            $.payers[payerAddress].balance += amount;
            _increaseTotalAmountDeposited(amount);
            return amount;
        }
    }

    /**
     * @notice Settles debts for a payer, updating their balance and total amounts.
     * @param  payer The address of the payer.
     * @param  amount The amount to settle debts for.
     * @return amountAfterSettlement The amount remaining after debt settlement.
     */
    function _settleDebts(address payer, uint256 amount) internal returns (uint256 amountAfterSettlement) {
        PayerStorage storage $ = _getPayerStorage();

        Payer memory storedPayer = $.payers[payer];

        if (storedPayer.debtAmount < amount) {
            uint256 debtToRemove = storedPayer.debtAmount;
            amount -= debtToRemove;

            storedPayer.debtAmount = 0;
            storedPayer.balance += amount;

            _removeDebtor(payer);
            _increaseTotalAmountDeposited(amount);
            _decreaseTotalDebtAmount(debtToRemove);

            amountAfterSettlement = amount;
        } else {
            storedPayer.debtAmount -= amount;

            _decreaseTotalDebtAmount(amount);

            amountAfterSettlement = 0;
        }

        $.payers[payer] = storedPayer;

        return amountAfterSettlement;
    }

    /**
     * @notice Checks if a payer exists.
     * @param  payer The address of the payer to check.
     * @return exists True if the payer exists, false otherwise.
     */
    function _payerExists(address payer) internal view returns (bool exists) {
        return _getPayerStorage().payers[payer].creationTimestamp != 0;
    }

    /**
     * @notice Checks if a payer is active.
     * @param  payer The address of the payer to check.
     * @return isActive True if the payer is active, false otherwise.
     */
    function _payerIsActive(address payer) internal view returns (bool isActive) {
        return _getPayerStorage().payers[payer].isActive;
    }

    /**
     * @notice Deactivates a payer.
     * @param  payer The address of the payer to deactivate.
     */
    function _deactivatePayer(address payer) internal {
        PayerStorage storage $ = _getPayerStorage();

        $.payers[payer].isActive = false;

        // Deactivating a payer only removes them from the active payers set
        require($.activePayers.remove(payer), FailedToDeactivatePayer());

        emit PayerDeactivated(payer);
    }

    /**
     * @notice Reverts if a payer does not exist.
     * @param  payer The address of the payer to check.
     */
    function _revertIfPayerDoesNotExist(address payer) internal view {
        require(_payerExists(payer), PayerDoesNotExist());
    }

    /**
     * @notice Checks if a withdrawal exists.
     * @param  payer The address of the payer to check.
     * @return exists True if the withdrawal exists, false otherwise.
     */
    function _withdrawalExists(address payer) internal view returns (bool exists) {
        return _getPayerStorage().withdrawals[payer].requestTimestamp != 0;
    }

    /**
     * @notice Reverts if a withdrawal does not exist.
     * @param  payer The address of the payer to check.
     */
    function _revertIfWithdrawalNotExists(address payer) internal view {
        require(_withdrawalExists(payer), WithdrawalNotExists());
    }

    /**
     * @notice Removes a payer from the debt payers set.
     * @param  payer The address of the payer to remove.
     */
    function _removeDebtor(address payer) internal {
        PayerStorage storage $ = _getPayerStorage();

        if ($.debtPayers.contains(payer)) {
            require($.debtPayers.remove(payer), FailedToRemoveDebtor());
        }
    }

    /**
     * @notice Adds a payer to the debt payers set.
     * @param  payer The address of the payer to add.
     */
    function _addDebtor(address payer) internal {
        PayerStorage storage $ = _getPayerStorage();

        if (!$.debtPayers.contains(payer)) {
            require($.debtPayers.add(payer), FailedToAddDebtor());
        }
    }

    /**
     * @notice Checks if a given address is an active node operator.
     * @param  operator The address to check.
     * @return isActiveNodeOperator True if the address is an active node operator, false otherwise.
     */
    function _getIsActiveNodeOperator(address operator) internal view returns (bool isActiveNodeOperator) {
        INodes nodes = INodes(_getPayerStorage().nodesContract);

        require(address(nodes) != address(0), Unauthorized());

        // TODO: Implement this in Nodes contract
        // return nodes.isActiveNodeOperator(operator);
        return true;
    }

    /**
     * @notice Sets the Distribution contract.
     * @param  _newDistributionContract The address of the new Distribution contract.
     */
    function _setDistributionContract(address _newDistributionContract) internal {
        PayerStorage storage $ = _getPayerStorage();

        // TODO: Add check to ensure the new distribution contract is valid
        //       Wait until Distribution contract is implemented
        // IDistribution distribution = IDistribution(_newDistributionContract);
        // require(distribution.supportsInterface(type(IDistribution).interfaceId), InvalidDistributionContract());

        require(_newDistributionContract != address(0), InvalidAddress());

        $.distributionContract = _newDistributionContract;

        emit DistributionContractSet(_newDistributionContract);
    }

    /**
     * @notice Sets the PayerReport contract.
     * @param  _newPayerReportContract The address of the new PayerReport contract.
     */
    function _setPayerReportContract(address _newPayerReportContract) internal {
        PayerStorage storage $ = _getPayerStorage();

        IPayerReport payerReport = IPayerReport(_newPayerReportContract);

        try payerReport.supportsInterface(type(IPayerReport).interfaceId) returns (bool supported) {
            require(supported, InvalidPayerReportContract());
        } catch {
            revert InvalidPayerReportContract();
        }

        $.payerReportContract = _newPayerReportContract;

        emit PayerReportContractSet(_newPayerReportContract);
    }

    /**
     * @notice Sets the Nodes contract.
     * @param  _newNodesContract The address of the new Nodes contract.
     */
    function _setNodesContract(address _newNodesContract) internal {
        PayerStorage storage $ = _getPayerStorage();

        try INodes(_newNodesContract).supportsInterface(type(INodes).interfaceId) returns (bool supported) {
            require(supported, InvalidNodesContract());
        } catch {
            revert InvalidNodesContract();
        }

        $.nodesContract = _newNodesContract;

        emit NodesContractSet(_newNodesContract);
    }

    /**
     * @notice Sets the USDC token contract.
     * @param  _newUsdcToken The address of the new USDC token contract.
     */
    function _setUsdcTokenContract(address _newUsdcToken) internal {
        PayerStorage storage $ = _getPayerStorage();

        try IERC20Metadata(_newUsdcToken).symbol() returns (string memory symbol) {
            require(keccak256(bytes(symbol)) == keccak256(bytes(USDC_SYMBOL)), InvalidUsdcTokenContract());
        } catch {
            revert InvalidUsdcTokenContract();
        }

        $.usdcToken = IERC20(_newUsdcToken);

        emit UsdcTokenSet(_newUsdcToken);
    }

    /**
     * @notice Increases the total amount deposited by a given amount.
     * @param  amount The amount to increase the total amount deposited by.
     */
    function _increaseTotalAmountDeposited(uint256 amount) internal {
        _getPayerStorage().totalAmountDeposited += amount;
    }

    /**
     * @notice Decreases the total amount deposited by a given amount.
     * @param  amount The amount to decrease the total amount deposited by.
     */
    function _decreaseTotalAmountDeposited(uint256 amount) internal {
        PayerStorage storage $ = _getPayerStorage();

        $.totalAmountDeposited = amount > $.totalAmountDeposited ? 0 : $.totalAmountDeposited - amount;
    }

    /**
     * @notice Increases the total debt amount by a given amount.
     * @param  amount The amount to increase the total debt amount by.
     */
    function _increaseTotalDebtAmount(uint256 amount) internal {
        _getPayerStorage().totalDebtAmount += amount;
    }

    /**
     * @notice Decreases the total debt amount by a given amount.
     * @param  amount The amount to decrease the total debt amount by.
     */
    function _decreaseTotalDebtAmount(uint256 amount) internal {
        PayerStorage storage $ = _getPayerStorage();

        $.totalDebtAmount = amount > $.totalDebtAmount ? 0 : $.totalDebtAmount - amount;
    }

    /**
     * @notice Internal helper for paginated access to EnumerableSet.AddressSet.
     * @param  addressSet The EnumerableSet to paginate.
     * @param  offset The starting index.
     * @param  limit Maximum number of items to return.
     * @return addresses Array of addresses from the set.
     * @return hasMore Whether there are more items after this page.
     */
    function _getPaginatedAddresses(EnumerableSet.AddressSet storage addressSet, uint256 offset, uint256 limit)
        internal
        view
        returns (address[] memory addresses, bool hasMore)
    {
        uint256 totalCount = addressSet.length();

        if (offset >= totalCount) revert OutOfBounds();

        uint256 count = totalCount - offset;
        if (count > limit) {
            count = limit;
            hasMore = true;
        } else {
            hasMore = false;
        }

        addresses = new address[](count);

        for (uint256 i = 0; i < count; i++) {
            addresses[i] = addressSet.at(offset + i);
        }

        return (addresses, hasMore);
    }

    /* ============ Upgradeability ============ */

    /**
     * @dev   Authorizes the upgrade of the contract.
     * @param newImplementation The address of the new implementation.
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), InvalidAddress());
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
