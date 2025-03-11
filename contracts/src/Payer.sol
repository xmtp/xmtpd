// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./interfaces/IPayer.sol";
import "./interfaces/INodes.sol";

interface INodesCaller {
    // TODO: Node has to implement ERC721Enumerable or an ad-hoc function (senderIsNodeOperator)
    function senderIsActiveNodeOperator(address sender) external view returns (bool isNodeOperator);
    function tokenOfOwnerByIndex(address owner, uint256 index) external view returns (uint256);
}

/**
 * @title Payer
 * @notice Implementation for managing payer USDC deposits, usage settlements,
 *         and a secure withdrawal process with optimized storage using
 *         Merkle trees and EIP-1283 gas optimizations.
 */
contract Payer is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable, IPayer{
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @dev Roles
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    /// @dev USDC token contract
    IERC20 public usdcToken;

    /// @dev Distribution contract
    address public distributionContract;

    /// @dev Nodes contract
    address public nodesContract;

    /// @dev Payer report contract
    address public payerReportContract;

    /// @dev Minimum deposit amount in micro-USDC
    uint256 public minimumDepositMicroDollars = 10_000_000;

    /// @dev Pending fees
    uint256 public pendingFees;

    /// @dev Last fee transfer timestamp
    uint256 public lastFeeTransferTimestamp;

    /// @dev Withdrawal lock period
    uint256 public withdrawalLockPeriod = 3 days;

    /// @dev Maximum backdated time
    uint256 public maxBackdatedTime = 1 days;

    /// @dev Mapping of payer address to their information
    mapping(address => Payer) private payers;

    /// @dev Mapping of payer address to withdrawal information
    mapping(address => Withdrawal) private withdrawals;

    /// @dev Set of all payer addresses
    EnumerableSet.AddressSet private totalPayers;

    /// @dev Set of active payer addresses
    EnumerableSet.AddressSet private activePayers;

    /// @dev Total value locked
    uint256 public totalValueLocked;

    /// @dev Total debt amount
    uint256 public totalDebtAmount;

    //==============================================================
    //                          Modifiers
    //==============================================================

    /**
     * @dev Modifier to check if caller is an active node operator
     */
    modifier onlyNodeOperator() {
        if (!getIsActiveNodeOperator(msg.sender)) {
            revert UnauthorizedNodeOperator();
        }
        _;
    }

    /**
     * @dev Modifier to check if caller is the payer report contract
     */
    modifier onlyPayerReport() {
        if (msg.sender != payerReportContract) {
            revert NotPayerReportContract();
        }
        _;
    }

    /**
     * @dev Modifier to check if address is an active payer
     */
    modifier onlyPayer(address payer) {
        require(_payerExists(payer), PayerDoesNotExist(payer));
        _;
    }

    //==============================================================
    //                          Initialization
    //==============================================================

    /// @notice Initializes the contract with the deployer as admin.
    /// @param _initialAdmin The address of the admin.
    function initialize(
        address _initialAdmin,
        address _usdcToken,
        address _distributionContract,
        address _nodesContract,
        uint256 _withdrawalLockPeriod,
        uint256 _maxBackdatedTime
    ) public initializer {
        if (_initialAdmin == address(0) || 
            _usdcToken == address(0) || 
            _nodesContract == address(0) ||
            _withdrawalLockPeriod == 0 ||
            _maxBackdatedTime == 0) {
            revert InvalidAddress();
        }

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _initialAdmin);
        _grantRole(ADMIN_ROLE, _initialAdmin);

        usdcToken = IERC20(_usdcToken);
        distributionContract = _distributionContract;
        nodesContract = _nodesContract;
        withdrawalLockPeriod = _withdrawalLockPeriod;
        maxBackdatedTime = _maxBackdatedTime;
    }

    //==============================================================
    //                   Payers Management
    //==============================================================

    /**
     * @inheritdoc IPayer
     */
    function register(uint256 amount) external whenNotPaused {
        require(amount >= minimumDepositMicroDollars, InsufficientAmount());
        require(!_payerExists(msg.sender), PayerAlreadyRegistered(msg.sender));

        // Transfer USDC from the sender to this contract
        usdcToken.safeTransferFrom(msg.sender, address(this), amount);

        // New payer registration
        payers[msg.sender] = Payer({
            balance: amount,
            isActive: true,
            creationTimestamp: block.timestamp,
            latestDepositTimestamp: block.timestamp,
            debtAmount: 0
        });

        // Add new payer to active and total payers sets
        activePayers.add(msg.sender);
        totalPayers.add(msg.sender);

        // Update counters
        totalValueLocked += amount;
        
        emit PayerRegistered(msg.sender, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deposit(uint256 amount) external whenNotPaused onlyPayer(msg.sender) {
        require(amount > 0, InsufficientAmount());

        if (withdrawals[msg.sender].requestTimestamp != 0) {
            revert PayerInWithdrawal();
        }

        // Transfer USDC from sender to this contract
        usdcToken.safeTransferFrom(msg.sender, address(this), amount);

        // Update payer record
        payers[msg.sender].balance += amount;
        payers[msg.sender].latestDepositTimestamp = block.timestamp;
        totalValueLocked += amount;

        emit Deposit(msg.sender, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function donate(address payer, uint256 amount) external whenNotPaused {
        require(amount > 0, InsufficientAmount());
        require(_payerExists(payer), PayerDoesNotExist(payer));

        if (withdrawals[payer].requestTimestamp != 0) {
            revert PayerInWithdrawal();
        }

        // Transfer USDC from sender to this contract
        usdcToken.safeTransferFrom(msg.sender, address(this), amount);

        // Update payer record
        payers[payer].balance += amount;
        
        // Update TVL
        totalValueLocked += amount;

        emit Donation(msg.sender, payer, amount);
    }

    /**
     * @inheritdoc IPayer
     */
    function deactivatePayer(address payer) external whenNotPaused onlyRole(ADMIN_ROLE) {
        require(_payerExists(payer), PayerDoesNotExist(payer));
        payers[payer].isActive = false;
        activePayers.remove(payer);
        emit PayerDeactivated(payer);
    }

    /**
     * @inheritdoc IPayer
     */
    function deletePayer(address payer) external whenNotPaused onlyRole(ADMIN_ROLE) {
        require(_payerExists(payer), PayerDoesNotExist(payer));

        if (payers[payer].balance > 0 || payers[payer].debtAmount > 0) {
            revert PayerHasBalanceOrDebt();
        }

        if (withdrawals[payer].requestTimestamp != 0) {
            revert PayerInWithdrawal();
        }

        // Delete payer data
        delete payers[payer];

        // Remove from totalPayers set
        totalPayers.remove(payer);
        activePayers.remove(payer);
        
        emit PayerDeleted(payer, block.timestamp);
    }

    //==============================================================
    //                  Payers Balance Management
    //==============================================================

    /**
     * @inheritdoc IPayer
     */
    function getPayerBalance(address payer) external view returns (uint256 balance) {
        require(_payerExists(payer), PayerDoesNotExist(payer));
        return payers[payer].balance;
    }

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

    //==============================================================
    //                    Usage Settlement
    //==============================================================

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

    function calculateFees(uint256 amount) external view returns (uint256 fees) {
        // TODO: Implement fee calculation logic
    }

    /**
     * @inheritdoc IPayer
     */
    function transferFeesToDistribution() external whenNotPaused onlyRole(ADMIN_ROLE) {
        // TODO: Implement fee transfer logic
    }

    //==============================================================
    //                 Administrative Functions
    //==============================================================

    /**
     * @inheritdoc IPayer
     */
    function setDistributionContract(address _newDistributionContract) external onlyRole(ADMIN_ROLE) {
        require (_newDistributionContract != address(0), InvalidAddress());
        // TODO: Add check to ensure the new distribution contract is valid
        distributionContract = _newDistributionContract;
        emit DistributionContractUpdated(_newDistributionContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setNodesContract(address _newNodesContract) external onlyRole(ADMIN_ROLE) {
        require (_newNodesContract != address(0), InvalidAddress());
        // TODO: Add check to ensure the new nodes contract is valid
        nodesContract = _newNodesContract;
        emit NodesContractUpdated(_newNodesContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setPayerReportContract(address _newPayerReportContract) external onlyRole(ADMIN_ROLE) {
        require (_newPayerReportContract != address(0), InvalidAddress());
        // TODO: Add check to ensure the new payer report contract is valid
        payerReportContract = _newPayerReportContract;
        emit PayerReportContractUpdated(_newPayerReportContract);
    }

    /**
     * @inheritdoc IPayer
     */
    function setMinimumDeposit(uint256 _newMinimumDeposit) external onlyRole(ADMIN_ROLE) {
        uint256 oldMinimumDeposit = minimumDepositMicroDollars;
        minimumDepositMicroDollars = _newMinimumDeposit;
        emit MinimumDepositUpdated(oldMinimumDeposit, _newMinimumDeposit);
    }

    /**
     * @inheritdoc IPayer
     */
    function getContractBalance() external view returns (uint256 balance) {
        return usdcToken.balanceOf(address(this));
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

    // Upgradeability
    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }

    /**
     * @inheritdoc IPayer
     */
    function getIsActiveNodeOperator(address operator) public view returns (bool) {
        INodesCaller nodes = INodesCaller(nodesContract);
        require(address(nodes) != address(0), Unauthorized());

        return nodes.senderIsActiveNodeOperator(operator);
    }

    //==============================================================
    //                      Getters
    //==============================================================

    /**
     * @dev Returns the payer information.
     * @param payer The address of the payer.
     * @return payerInfo The payer information.
     */
    function getPayer(address payer) external view returns (Payer memory payerInfo) {
        require(_payerExists(payer), PayerDoesNotExist(payer));
        return payers[payer];
    }

    /**
     * @inheritdoc IPayer
     */
    function getIsActivePayer(address payer) public view returns (bool isActive) {
        return activePayers.contains(payer);
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
        return totalPayers.length();
    }

    /**
     * @inheritdoc IPayer
     */
    function getActivePayerCount() external view returns (uint256 count) {
        return activePayers.length();
    }

    /**
     * @inheritdoc IPayer
     */
    function getLastFeeTransferTimestamp() external view returns (uint256 timestamp) {
        return lastFeeTransferTimestamp;
    }

    function getTotalValueLocked() external view returns (uint256 tvl) {
        // TODO: TVL should subtract the total debt amount
        return totalValueLocked;
    }

    function getTotalDebtAmount() external view returns (uint256 totalDebt) {
        return totalDebtAmount;
    }

    //==============================================================
    //                    Internal Functions
    //==============================================================

    /**
     * @dev Checks if a payer exists.
     * @param payer The address of the payer to check.
     * @return exists True if the payer exists, false otherwise.
     */
    function _payerExists(address payer) internal view returns (bool exists) {
        return payers[payer].creationTimestamp != 0;
    }
}