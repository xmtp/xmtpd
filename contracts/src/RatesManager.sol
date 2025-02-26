// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract RatesManager is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    bytes32 public constant RATES_MANAGER_ROLE = keccak256("RATES_MANAGER_ROLE");
    uint256 constant PAGE_SIZE = 50; // Fixed page size for reading rates

    // All Rates appended here
    Rates[] private allRates;

    // EVENTS
    // Event emitted when new Rates are added
    event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime);

    /// @notice Emitted when an upgrade is authorized.
    /// @param upgrader The EOA authorizing the upgrade.
    /// @param newImplementation The address of the new implementation.
    event UpgradeAuthorized(address upgrader, address newImplementation);

    // Custom errors
    error ZeroAdminAddress();
    error InvalidStartTime();
    error FromIndexOutOfRange();

    // Rates struct holds the fees and the start time of the rates
    struct Rates {
        uint64 messageFee;
        uint64 storageFee;
        uint64 congestionFee;
        uint64 startTime;
    }

    /// @dev Reserved storage gap for future upgrades
    // slither-disable-next-line unused-state,naming-convention
    uint256[50] private __gap;

    // Initialization
    /// @notice Initializes the contract with the deployer as admin.
    /// @param admin The address of the admin.
    function initialize(address admin) public initializer {
        if (admin == address(0)) {
            revert ZeroAdminAddress();
        }

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        _setRoleAdmin(RATES_MANAGER_ROLE, DEFAULT_ADMIN_ROLE);
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(RATES_MANAGER_ROLE, admin);
    }

    // Pausable functionality
    /// @notice Pauses the contract, restricting certain actions.
    /// @dev Callable only by accounts with the DEFAULT_ADMIN_ROLE.
    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /// @notice Unpauses the contract, allowing normal operations.
    /// @dev Callable only by accounts with the DEFAULT_ADMIN_ROLE.
    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @dev Add new Rates. Can only be called by addresses with RATES_ADMIN_ROLE.
     *      The array only grows; we do not allow removal or updating.
     */
    function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime)
        external
        onlyRole(RATES_MANAGER_ROLE)
    {
        // Enforce chronological order
        if (allRates.length > 0 && startTime <= allRates[allRates.length - 1].startTime) {
            revert InvalidStartTime();
        }

        allRates.push(
            Rates({messageFee: messageFee, storageFee: storageFee, congestionFee: congestionFee, startTime: startTime})
        );

        emit RatesAdded(messageFee, storageFee, congestionFee, startTime);
    }

    /**
     * @dev Returns a slice of the Rates list for pagination.
     * @param fromIndex Index from which to start (must be < allRates.length).
     * @return rates - The subset of Rates.
     * @return hasMore - True if there are more items beyond this slice.
     */
    function getRates(uint256 fromIndex) external view returns (Rates[] memory rates, bool hasMore) {
        if (allRates.length == 0 && fromIndex == 0) {
            return (new Rates[](0), false);
        }

        if (fromIndex >= allRates.length) {
            revert FromIndexOutOfRange();
        }

        uint256 toIndex = fromIndex + PAGE_SIZE;
        if (toIndex > allRates.length) {
            toIndex = allRates.length;
        }

        uint256 resultSize = toIndex - fromIndex;
        Rates[] memory tempRates = new Rates[](resultSize);
        for (uint256 i = 0; i < resultSize; i++) {
            tempRates[i] = allRates[fromIndex + i];
        }

        bool moreData = (toIndex < allRates.length);
        return (tempRates, moreData);
    }

    /**
     * @dev Returns the total number of Rates stored.
     */
    function getRatesCount() external view returns (uint256) {
        return allRates.length;
    }

    // Upgradeability
    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
