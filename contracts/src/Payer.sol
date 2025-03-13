// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { IERC20Metadata } from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { AccessControlUpgradeable } from "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import { IPayer } from "./interfaces/IPayer.sol";
import { IPayerReport } from "./interfaces/IPayerReport.sol";
import { INodes } from "./interfaces/INodes.sol";

/**
 * @title  Payer
 * @notice Implementation for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process.
 */
contract Payer is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable, IPayer{
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ============ Constants ============ */

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    uint256 public constant DEFAULT_MINIMUM_REGISTRATION_AMOUNT_MICRO_DOLLARS = 10_000_000;
    uint256 public constant DEFAULT_MINIMUM_DEPOSIT_AMOUNT_MICRO_DOLLARS = 10_000_000;
    uint256 public constant DEFAULT_WITHDRAWAL_LOCK_PERIOD = 3 days;
    uint256 public constant ABSOLUTE_MINIMUM_WITHDRAWAL_LOCK_PERIOD = 1 days;
    uint256 public constant DEFAULT_MAX_BACKDATED_TIME = 1 days;
    uint256 public constant ABSOLUTE_MINIMUM_MAX_BACKDATED_TIME = 6 hours;
    string internal constant USDC_SYMBOL = "USDC";

    /* ============ UUPS Storage ============ */

    /// @custom:storage-location erc7201:xmtp.storage.Payer
    struct PayerStorage {
        /// @dev Contracts to interact with.
        IERC20 usdcToken;
        address distributionContract;
        address nodesContract;
        address payerReportContract;

        /// @dev Configuration parameters
        uint256 minimumRegistrationAmountMicroDollars;
        uint256 minimumDepositAmountMicroDollars;
        uint256 withdrawalLockPeriod;
        uint256 maxBackdatedTime;

        /// @dev State variables
        uint256 lastFeeTransferTimestamp;
        uint256 pendingFees;
        uint256 totalAmountDeposited;
        uint256 totalDebtAmount;

        /// @dev Mappings
        mapping(address => Payer) payers;
        mapping(address => Withdrawal) withdrawals;
        EnumerableSet.AddressSet totalPayers;
        EnumerableSet.AddressSet activePayers;
    }

    // keccak256(abi.encode(uint256(keccak256("xmtp.storage.Payer")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 internal constant PayerStorageLocation = 0xd0335f337c570f3417b0f0d20340c88da711d60e810b5e9b3ecabe9ccfcdce5a;

    function _getPayerStorage() internal pure returns (PayerStorage storage $) {
        assembly {
            $.slot := PayerStorageLocation
        }
    }

    /* ============ Modifiers ============ */

    /**
     * @dev Modifier to check if caller is an active node operator
     */
    modifier onlyNodeOperator() {
        if (!_getIsActiveNodeOperator(msg.sender)) {
            revert UnauthorizedNodeOperator();
        }
        _;
    }

    /**
     * @dev Modifier to check if caller is the payer report contract
     */
    modifier onlyPayerReport() {
        if (msg.sender != _getPayerStorage().payerReportContract) {
            revert Unauthorized();
        }
        _;
    }

    /**
     * @dev Modifier to check if address is an active payer
     */
    modifier onlyPayer(address payer) {
        require(_payerExists(payer), PayerDoesNotExist());
        require(msg.sender == payer, Unauthorized());
        _;
    }

    /* ============ Initialization ============ */

    /**
     * @notice Initializes the contract with the deployer as admin.
     * @param  _initialAdmin The address of the admin.
     * @dev    There's a chicken-egg problem here with PayerReport and Distribution contracts.
     *         We need to deploy these contracts first, then set their addresses
     *         in the Payer contract.
     */
    function initialize(
        address _initialAdmin,
        address _usdcToken,
        address _nodesContract
    ) public initializer {
        if (_initialAdmin == address(0) || _usdcToken == address(0) || _nodesContract == address(0)) {
            revert InvalidAddress();
        }

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        PayerStorage storage $ = _getPayerStorage();

        $.minimumRegistrationAmountMicroDollars = DEFAULT_MINIMUM_REGISTRATION_AMOUNT_MICRO_DOLLARS;
        $.minimumDepositAmountMicroDollars = DEFAULT_MINIMUM_DEPOSIT_AMOUNT_MICRO_DOLLARS;
        $.withdrawalLockPeriod = DEFAULT_WITHDRAWAL_LOCK_PERIOD;
        $.maxBackdatedTime = DEFAULT_MAX_BACKDATED_TIME;

        _setUsdcTokenContract(_usdcToken);
        _setNodesContract(_nodesContract);

        require(_grantRole(DEFAULT_ADMIN_ROLE, _initialAdmin), FailedToGrantRole(DEFAULT_ADMIN_ROLE, _initialAdmin));
        require(_grantRole(ADMIN_ROLE, _initialAdmin), FailedToGrantRole(ADMIN_ROLE, _initialAdmin));
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
    function deposit(uint256 amount) external whenNotPaused onlyPayer(msg.sender) {
        PayerStorage storage $ = _getPayerStorage();

        require(amount >= $.minimumDepositAmountMicroDollars, InsufficientAmount());

        if ($.withdrawals[msg.sender].requestTimestamp != 0) revert PayerInWithdrawal();

        _deposit(msg.sender, amount);

        _updatePayerBalance(msg.sender, amount);

        $.payers[msg.sender].latestDepositTimestamp = block.timestamp;
    }

    /**
     * @inheritdoc IPayer
     */
    function donate(address payer, uint256 amount) external whenNotPaused {
        require(amount > 0, InsufficientAmount());
        _revertIfPayerDoesNotExist(payer);

        PayerStorage storage $ = _getPayerStorage();

        if ($.withdrawals[payer].requestTimestamp != 0) revert PayerInWithdrawal();

        _deposit(msg.sender, amount);

        _updatePayerBalance(payer, amount);

        $.payers[payer].latestDonationTimestamp = block.timestamp;

        emit Donation(msg.sender, payer, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deactivatePayer(address payer) external whenNotPaused onlyNodeOperator {
        _revertIfPayerDoesNotExist(payer);

        PayerStorage storage $ = _getPayerStorage();

        $.payers[payer].isActive = false;

        // Deactivating a payer only removes them from the active payers set
        require($.activePayers.remove(payer), FailedToDeactivatePayer());

        emit PayerDeactivated(payer);
    }

    /**
     * @inheritdoc IPayer
     */
    function deletePayer(address payer) external whenNotPaused onlyRole(ADMIN_ROLE) {
        _revertIfPayerDoesNotExist(payer);

        PayerStorage storage $ = _getPayerStorage();

        if ($.payers[payer].balance > 0 || $.payers[payer].debtAmount > 0) {
            revert PayerHasBalanceOrDebt();
        }

        if ($.withdrawals[payer].requestTimestamp != 0) revert PayerInWithdrawal();

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
    function requestWithdrawal(uint256 amount) external whenNotPaused() onlyPayer(msg.sender) {
        if (_withdrawalExists(msg.sender)) revert WithdrawalAlreadyRequested();

        PayerStorage storage $ = _getPayerStorage();

        require($.payers[msg.sender].debtAmount == 0, PayerHasDebt());
        require($.payers[msg.sender].balance >= amount, InsufficientBalance());

        // Balance to be withdrawn is deducted from the payer's balance,
        // it can't be used to settle payments.
        _decreasePayerBalance(msg.sender, amount);
        _decreaseTotalAmountDeposited(amount);

        uint256 withdrawableTimestamp = block.timestamp + $.withdrawalLockPeriod;

        $.withdrawals[msg.sender] = Withdrawal({
            requestTimestamp: block.timestamp,
            withdrawableTimestamp: withdrawableTimestamp,
            amount: amount
        });

        emit WithdrawalRequested(
            msg.sender, 
            block.timestamp, 
            withdrawableTimestamp, 
            amount
        );
    }

    /**
     * @inheritdoc IPayer
     */
    function cancelWithdrawal() external whenNotPaused() onlyPayer(msg.sender) {
        _revertIfWithdrawalNotExists(msg.sender);

        PayerStorage storage $ = _getPayerStorage();

        Withdrawal memory withdrawal = $.withdrawals[msg.sender];

        delete $.withdrawals[msg.sender];

        _updatePayerBalance(msg.sender, withdrawal.amount);

        emit WithdrawalCancelled(msg.sender, withdrawal.requestTimestamp);
    }

    /**
     * @inheritdoc IPayer
     */
    function finalizeWithdrawal() external whenNotPaused() onlyPayer(msg.sender) {
        _revertIfWithdrawalNotExists(msg.sender);

        PayerStorage storage $ = _getPayerStorage();

        Withdrawal memory withdrawal = $.withdrawals[msg.sender];

        delete $.withdrawals[msg.sender];

        $.usdcToken.safeTransfer(msg.sender, withdrawal.amount);

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
    function settleUsage(
        address originatorNode,
        uint256 reportIndex,
        address[] calldata payerList,
        uint256[] calldata amounts
    ) external whenNotPaused onlyPayerReport
    {
        // TODO: Implement usage settlement logic
    }

    /**
     * @inheritdoc IPayer
     */
    function transferFeesToDistribution() external whenNotPaused onlyRole(ADMIN_ROLE) {
        // TODO: Implement fee transfer logic
        //       Update lastFeeTransferTimestamp
        //       Update pendingFees
        //       Transfer fees to distribution contract
        //       Emit FeesTransferred event
    }

    /* ========== Administrative Functions ========== */

    /**
     * @inheritdoc IPayer
     */
    function setDistributionContract(address _newDistributionContract) external onlyRole(ADMIN_ROLE) {
        _setDistributionContract(_newDistributionContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setPayerReportContract(address _newPayerReportContract) external onlyRole(ADMIN_ROLE) {
        _setPayerReportContract(_newPayerReportContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setNodesContract(address _newNodesContract) external onlyRole(ADMIN_ROLE) {
        _setNodesContract(_newNodesContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setUsdcToken(address _newUsdcToken) external onlyRole(ADMIN_ROLE) {
        _setUsdcTokenContract(_newUsdcToken);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMinimumDeposit(uint256 _newMinimumDeposit) external onlyRole(ADMIN_ROLE) {
        require(_newMinimumDeposit > DEFAULT_MINIMUM_DEPOSIT_AMOUNT_MICRO_DOLLARS, InvalidMinimumDeposit());

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldMinimumDeposit = $.minimumDepositAmountMicroDollars;
        $.minimumDepositAmountMicroDollars = _newMinimumDeposit;

        emit MinimumDepositSet(oldMinimumDeposit, _newMinimumDeposit);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMinimumRegistrationAmount(uint256 _newMinimumRegistrationAmount) external onlyRole(ADMIN_ROLE) {
        require(_newMinimumRegistrationAmount > DEFAULT_MINIMUM_REGISTRATION_AMOUNT_MICRO_DOLLARS, InvalidMinimumRegistrationAmount());

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldMinimumRegistrationAmount = $.minimumRegistrationAmountMicroDollars;
        $.minimumRegistrationAmountMicroDollars = _newMinimumRegistrationAmount;

        emit MinimumRegistrationAmountSet(oldMinimumRegistrationAmount, _newMinimumRegistrationAmount);
    }

    /**
     * @inheritdoc IPayer
     */
    function setWithdrawalLockPeriod(uint256 _newWithdrawalLockPeriod) external onlyRole(ADMIN_ROLE) {
        require(_newWithdrawalLockPeriod >= ABSOLUTE_MINIMUM_WITHDRAWAL_LOCK_PERIOD, InvalidWithdrawalLockPeriod());

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldWithdrawalLockPeriod = $.withdrawalLockPeriod;
        $.withdrawalLockPeriod = _newWithdrawalLockPeriod;

        emit WithdrawalLockPeriodSet(oldWithdrawalLockPeriod, _newWithdrawalLockPeriod);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMaxBackdatedTime(uint256 _newMaxBackdatedTime) external onlyRole(ADMIN_ROLE) {
        require(_newMaxBackdatedTime >= ABSOLUTE_MINIMUM_MAX_BACKDATED_TIME, InvalidMaxBackdatedTime());

        PayerStorage storage $ = _getPayerStorage();

        uint256 oldMaxBackdatedTime = $.maxBackdatedTime;
        $.maxBackdatedTime = _newMaxBackdatedTime;

        emit MaxBackdatedTimeSet(oldMaxBackdatedTime, _newMaxBackdatedTime);
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
    function getActivePayers(uint256 offset, uint256 limit) external view returns (Payer[] memory payers, bool hasMore) {
        PayerStorage storage $ = _getPayerStorage();
        
        uint256 totalCount = $.activePayers.length();

        if (offset >= totalCount) revert OutOfBounds();

        uint256 count = totalCount - offset;
        if (count > limit) {
            count = limit;
            hasMore = true;
        } else {
            hasMore = false;
        }

        payers = new Payer[](count);

        for (uint256 i = 0; i < count; i++) {
            address payerAddress = $.activePayers.at(offset + i);
            payers[i] = $.payers[payerAddress];
        }
        
        return (payers, hasMore);
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
    function getPayersInDebt(uint256 offset, uint256 limit) external view returns (
        address[] memory debtors,
        uint256[] memory debtAmounts,
        uint256 totalCount
    ) {
        // TODO: Implement payers in debt retrieval logic
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

    /**
     * @inheritdoc IPayer
     */
    function getMaxBackdatedTime() external view returns (uint256 maxTime) {
        return _getPayerStorage().maxBackdatedTime;
    }

    /* ============ Internal ============ */

    function _deposit(address payer, uint256 amount) internal {
        PayerStorage storage $ = _getPayerStorage();

        $.usdcToken.safeTransferFrom(payer, address(this), amount);
    }

    /**
     * @notice Updates a payer's balance by the specified amount.
     *         If the payer has debt, the debt is subtracted from the amount.
     *         If the payer has no debt, the amount is added to the balance.
     * @param  payerAddress The address of the payer.
     * @param  amount The amount to update by.
     */
    function _updatePayerBalance(address payerAddress, uint256 amount) internal {
        Payer storage payer = _getPayerStorage().payers[payerAddress];

        if (payer.debtAmount == 0) {
            payer.balance += amount;
            _increaseTotalAmountDeposited(amount);
        } else {
            if (payer.debtAmount < amount) {
                uint256 debtToRemove = payer.debtAmount;
                amount -= debtToRemove;
                payer.debtAmount = 0;
                payer.balance += amount;
                _increaseTotalAmountDeposited(amount);
                _decreaseTotalDebtAmount(debtToRemove);
            } else {
                payer.debtAmount -= amount;
                _decreaseTotalDebtAmount(amount);
            }
        }

        emit PayerBalanceUpdated(payerAddress, payer.balance, payer.debtAmount);
    }

    /**
    * @notice Decreases a payer's balance by the specified amount.
    * @param  payerAddress The address of the payer.
    * @param  amount The amount to decrease by.
    */
    function _decreasePayerBalance(address payerAddress, uint256 amount) internal {
        Payer storage payer = _getPayerStorage().payers[payerAddress];

        payer.balance -= amount;

        emit PayerBalanceUpdated(payerAddress, payer.balance, payer.debtAmount);
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

    function _setDistributionContract(address _newDistributionContract) internal {
        PayerStorage storage $ = _getPayerStorage();

        // TODO: Add check to ensure the new distribution contract is valid
        //       Wait until Distribution contract is implemented
        // IDistribution distribution = IDistribution(_newDistributionContract);
        // require(distribution.supportsInterface(type(IDistribution).interfaceId), InvalidDistributionContract());

        require (_newDistributionContract != address(0), InvalidAddress());

        $.distributionContract = _newDistributionContract;

        emit DistributionContractSet(_newDistributionContract);
    }

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

    function _increaseTotalAmountDeposited(uint256 amount) internal {
        _getPayerStorage().totalAmountDeposited += amount;
    }

    // TODO: Check for underflow
    function _decreaseTotalAmountDeposited(uint256 amount) internal {
        _getPayerStorage().totalAmountDeposited -= amount;
    }

    function _increaseTotalDebtAmount(uint256 amount) internal {
        _getPayerStorage().totalDebtAmount += amount;
    }

    // TODO: Check for underflow
    function _decreaseTotalDebtAmount(uint256 amount) internal {
        _getPayerStorage().totalDebtAmount -= amount;
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