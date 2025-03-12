// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { AccessControlUpgradeable } from "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import { IPayer } from "./interfaces/IPayer.sol";
import { INodes } from "./interfaces/INodes.sol";

/**
 * @title Payer
 * @notice Implementation for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process with optimized storage using
 *         Merkle trees and EIP-1283 gas optimizations.
 */
contract Payer is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable, IPayer{
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ============ Constants ============ */

    /// @dev Roles
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    /* ============ UUPS Storage ============ */

    /// @custom:storage-location erc7201:xmtp.storage.Payer
    struct PayerStorage {
        /// @dev Contracts to interact with.
        IERC20 usdcToken;
        address distributionContract;
        address nodesContract;
        address payerReportContract;

        /// @dev Parameters
        uint256 minimumDepositMicroDollars;
        uint256 pendingFees;
        uint256 lastFeeTransferTimestamp;
        uint256 withdrawalLockPeriod;
        uint256 maxBackdatedTime;
        uint256 totalValueLocked;
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
            revert NotPayerReportContract();
        }
        _;
    }

    /**
     * @dev Modifier to check if address is an active payer
     */
    modifier onlyPayer(address payer) {
        require(_payerExists(payer), PayerDoesNotExist());
        _;
    }

    /* ============ Initialization ============ */

    /**
     * @notice Initializes the contract with the deployer as admin.
     * @param  _initialAdmin The address of the admin.
     */
    function initialize(
        address _initialAdmin,
        address _usdcToken,
        address _distributionContract,
        address _nodesContract
    ) public initializer {
        if (_initialAdmin == address(0) || _usdcToken == address(0) || _nodesContract == address(0)) {
            revert InvalidAddress();
        }

        PayerStorage storage $ = _getPayerStorage();

        $.minimumDepositMicroDollars = 10_000_000;
        $.withdrawalLockPeriod = 3 days;
        $.maxBackdatedTime = 1 days;

        $.usdcToken = IERC20(_usdcToken);
        $.distributionContract = _distributionContract;
        $.nodesContract = _nodesContract;

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        require(_grantRole(DEFAULT_ADMIN_ROLE, _initialAdmin), FailedToGrantRole(DEFAULT_ADMIN_ROLE, _initialAdmin));
        require(_grantRole(ADMIN_ROLE, _initialAdmin), FailedToGrantRole(ADMIN_ROLE, _initialAdmin));
    }

    /* ============ Payers Management ============ */

    /**
     * @inheritdoc IPayer
     */
    function register(uint256 amount) external whenNotPaused {
        PayerStorage storage $ = _getPayerStorage();

        require(amount >= $.minimumDepositMicroDollars, InsufficientAmount());
        require(!_payerExists(msg.sender), PayerAlreadyRegistered());

        // Transfer USDC from the sender to this contract
        $.usdcToken.safeTransferFrom(msg.sender, address(this), amount);

        // New payer registration
        $.payers[msg.sender] = Payer({
            balance: amount,
            isActive: true,
            creationTimestamp: block.timestamp,
            latestDepositTimestamp: block.timestamp,
            debtAmount: 0
        });

        // Add new payer to active and total payers sets
        $.activePayers.add(msg.sender);
        $.totalPayers.add(msg.sender);

        // Update counters
        $.totalValueLocked += amount;
        
        emit PayerRegistered(msg.sender, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deposit(uint256 amount) external whenNotPaused onlyPayer(msg.sender) {
        PayerStorage storage $ = _getPayerStorage();

        require(amount > 0, InsufficientAmount());
        _revertIfPayerDoesNotExist(msg.sender);

        if ($.withdrawals[msg.sender].requestTimestamp != 0) {
            revert PayerInWithdrawal();
        }

        // Transfer USDC from sender to this contract
        $.usdcToken.safeTransferFrom(msg.sender, address(this), amount);

        // Update payer record
        $.payers[msg.sender].balance += amount;
        $.payers[msg.sender].latestDepositTimestamp = block.timestamp;
        $.totalValueLocked += amount;

        emit Deposit(msg.sender, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function donate(address payer, uint256 amount) external whenNotPaused {
        PayerStorage storage $ = _getPayerStorage();

        require(amount > 0, InsufficientAmount());
        _revertIfPayerDoesNotExist(payer);

        if ($.withdrawals[payer].requestTimestamp != 0) {
            revert PayerInWithdrawal();
        }

        // Transfer USDC from sender to this contract
        $.usdcToken.safeTransferFrom(msg.sender, address(this), amount);

        // Update payer record
        $.payers[payer].balance += amount;
        
        // Update TVL
        $.totalValueLocked += amount;

        emit Donation(msg.sender, payer, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deactivatePayer(address payer) external whenNotPaused onlyRole(ADMIN_ROLE) {
        PayerStorage storage $ = _getPayerStorage();

        _revertIfPayerDoesNotExist(payer);
        $.payers[payer].isActive = false;
        $.activePayers.remove(payer);
        emit PayerDeactivated(payer);
    }

    /**
     * @inheritdoc IPayer
     */
    function deletePayer(address payer) external whenNotPaused onlyRole(ADMIN_ROLE) {
        PayerStorage storage $ = _getPayerStorage();

        _revertIfPayerDoesNotExist(payer);

        if ($.payers[payer].balance > 0 || $.payers[payer].debtAmount > 0) {
            revert PayerHasBalanceOrDebt();
        }

        if ($.withdrawals[payer].requestTimestamp != 0) {
            revert PayerInWithdrawal();
        }

        // Delete payer data
        delete $.payers[payer];

        // Remove from totalPayers set
        $.totalPayers.remove(payer);
        $.activePayers.remove(payer);
        
        emit PayerDeleted(payer, block.timestamp);
    }

    /* ========== Payers Balance Management ========= */

    /**
     * @inheritdoc IPayer
     */
    function requestWithdrawal(uint256 amount) external whenNotPaused() onlyPayer(msg.sender) {
        // TODO: Implement withdrawal request logic
    }

    /**
     * @inheritdoc IPayer
     */
    function cancelWithdrawal() external whenNotPaused() onlyPayer(msg.sender) {
        // TODO: Implement withdrawal cancellation logic
    }

    /**
     * @inheritdoc IPayer
     */
    function finalizeWithdrawal() external whenNotPaused() onlyPayer(msg.sender) {
        // TODO: Implement withdrawal finalization logic
    }

    /**
     * @inheritdoc IPayer
     */
    function getWithdrawalStatus(address payer) external view returns (Withdrawal memory withdrawal) {
        // TODO: Implement withdrawal status retrieval logic
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
    }

    /* ========== Administrative Functions ========== */

    /**
     * @inheritdoc IPayer
     */
    function setDistributionContract(address _newDistributionContract) external onlyRole(ADMIN_ROLE) {
        PayerStorage storage $ = _getPayerStorage();

        require (_newDistributionContract != address(0), InvalidAddress());
        // TODO: Add check to ensure the new distribution contract is valid
        $.distributionContract = _newDistributionContract;
        emit DistributionContractUpdated(_newDistributionContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setNodesContract(address _newNodesContract) external onlyRole(ADMIN_ROLE) {
        PayerStorage storage $ = _getPayerStorage();

        require (_newNodesContract != address(0), InvalidAddress());
        // TODO: Add check to ensure the new nodes contract is valid
        $.nodesContract = _newNodesContract;
        emit NodesContractUpdated(_newNodesContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setPayerReportContract(address _newPayerReportContract) external onlyRole(ADMIN_ROLE) {
        PayerStorage storage $ = _getPayerStorage();

        require (_newPayerReportContract != address(0), InvalidAddress());
        // TODO: Add check to ensure the new payer report contract is valid
        $.payerReportContract = _newPayerReportContract;
        emit PayerReportContractUpdated(_newPayerReportContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMinimumDeposit(uint256 _newMinimumDeposit) external onlyRole(ADMIN_ROLE) {
        PayerStorage storage $ = _getPayerStorage();

        uint256 oldMinimumDeposit = $.minimumDepositMicroDollars;
        $.minimumDepositMicroDollars = _newMinimumDeposit;
        emit MinimumDepositUpdated(oldMinimumDeposit, _newMinimumDeposit);
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
        PayerStorage storage $ = _getPayerStorage();

        _revertIfPayerDoesNotExist(payer);
        return $.payers[payer];
    }

    /**
     * @inheritdoc IPayer
     */
    function getIsActivePayer(address payer) public view returns (bool isActive) {
        return _getPayerStorage().activePayers.contains(payer);
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
        PayerStorage storage $ = _getPayerStorage();

        return $.totalPayers.length();
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
    function getTotalValueLocked() external view returns (uint256 tvl) {
        // TODO: TVL should subtract the total debt amount
        return _getPayerStorage().totalValueLocked;
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
     * @notice Retrieves the minimum deposit amount required to register as a payer.
     * @return minimumDeposit The minimum deposit amount in USDC.
     */
    function getMinimumDeposit() external view returns (uint256 minimumDeposit) {
        return _getPayerStorage().minimumDepositMicroDollars;
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

    /**
     * @dev   Reverts if a payer does not exist.
     * @param payer The address of the payer to check.
     */
    function _revertIfPayerDoesNotExist(address payer) internal view {
        require(_payerExists(payer), PayerDoesNotExist());
    }

    /**
     * @dev    Checks if a payer exists.
     * @param  payer The address of the payer to check.
     * @return exists True if the payer exists, false otherwise.
     */
    function _payerExists(address payer) internal view returns (bool exists) {
        return _getPayerStorage().payers[payer].creationTimestamp != 0;
    }

    /**
     * @notice Checks if a given address is an active node operator.
     * @param  operator The address to check.
     * @return isActiveNodeOperator True if the address is an active node operator, false otherwise.
     */
    function _getIsActiveNodeOperator(address operator) internal view returns (bool) {
        INodes nodes = INodes(_getPayerStorage().nodesContract);
        require(address(nodes) != address(0), Unauthorized());

        // TODO: Implement this in Nodes contract
        // return nodes.isActiveNodeOperator(operator);
        return true;
    }

    /* ============ Upgradeability ============ */

    /**
     * @dev   Authorizes the upgrade of the contract.
     * @param newImplementation The address of the new implementation.
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}