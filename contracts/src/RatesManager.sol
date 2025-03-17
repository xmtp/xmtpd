// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { AccessControlUpgradeable } from "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

// TODO: PAGE_SIZE should be a default, but overridden by the caller.
// TODO: Nodes should filter recent events to build rates array, without requiring contract to maintain it.

contract RatesManager is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    /* ============ Events ============ */

    // Event emitted when new Rates are added
    event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime);

    /**
     * @notice Emitted when an upgrade is authorized.
     * @param  upgrader          The EOA authorizing the upgrade.
     * @param  newImplementation The address of the new implementation.
     */
    event UpgradeAuthorized(address upgrader, address newImplementation); // TODO: index both.

    /* ============ Custom Errors ============ */

    error ZeroAdminAddress();
    error InvalidStartTime();
    error FromIndexOutOfRange();
    error ZeroImplementationAddress();

    /* ============ Structs ============ */

    // Rates struct holds the fees and the start time of the rates.
    struct Rates {
        uint64 messageFee;
        uint64 storageFee;
        uint64 congestionFee;
        uint64 targetRatePerMinute;
        uint64 startTime;
    }

    /* ============ Constants ============ */

    bytes32 public constant RATES_MANAGER_ROLE = keccak256("RATES_MANAGER_ROLE");
    uint256 public constant PAGE_SIZE = 50; // Fixed page size for reading rates

    /* ============ UUPS Storage ============ */

    /// @custom:storage-location erc7201:xmtp.storage.RatesManager
    struct RatesManagerStorage {
        Rates[] allRates; // All Rates appended here.
    }

    // keccak256(abi.encode(uint256(keccak256("xmtp.storage.RatesManager")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 internal constant RATES_MANAGER_STORAGE_LOCATION =
        0x6ad1a01bf62225c91223b2956030efc848b0def7d19ed478ca6dd31490e2d000;

    function _getRatesManagerStorage() internal pure returns (RatesManagerStorage storage $) {
        // slither-disable-next-line assembly
        assembly {
            $.slot := RATES_MANAGER_STORAGE_LOCATION
        }
    }

    /* ============ Initialization ============ */

    /**
     * @notice Initializes the contract with the deployer as admin.
     * @param  admin The address of the admin.
     */
    function initialize(address admin) public initializer {
        require(admin != address(0), ZeroAdminAddress());

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        _setRoleAdmin(RATES_MANAGER_ROLE, DEFAULT_ADMIN_ROLE);

        // slither-disable-next-line unused-return
        _grantRole(DEFAULT_ADMIN_ROLE, admin); // Will return false if the role is already granted.

        // slither-disable-next-line unused-return
        _grantRole(RATES_MANAGER_ROLE, admin); // Will return false if the role is already granted.
    }

    /* ============ Pausable functionality ============ */

    /**
     * @notice Pauses the contract, restricting certain actions.
     * @dev    Callable only by accounts with the DEFAULT_ADMIN_ROLE.
     */
    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpauses the contract, allowing normal operations.
     * @dev    Callable only by accounts with the DEFAULT_ADMIN_ROLE.
     */
    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /* ============ RatesManager functionality ============ */

    /**
     * @dev Add new Rates. Can only be called by addresses with RATES_ADMIN_ROLE.
     *      The array only grows; we do not allow removal or updating.
     */
    function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 targetRatePerMinute, uint64 startTime)
        external
        onlyRole(RATES_MANAGER_ROLE)
    {
        RatesManagerStorage storage $ = _getRatesManagerStorage();

        // Enforce chronological order.
        if ($.allRates.length > 0 && startTime <= $.allRates[$.allRates.length - 1].startTime) {
            revert InvalidStartTime();
        }

        $.allRates.push(
            Rates({
                messageFee: messageFee,
                storageFee: storageFee,
                congestionFee: congestionFee,
                startTime: startTime,
                targetRatePerMinute: targetRatePerMinute
            })
        );

        emit RatesAdded(messageFee, storageFee, congestionFee, targetRatePerMinute, startTime);
    }

    /**
     * @dev    Returns a slice of the Rates list for pagination.
     * @param  fromIndex Index from which to start (must be < allRates.length).
     * @return rates     The subset of Rates.
     * @return hasMore   True if there are more items beyond this slice.
     */
    function getRates(uint256 fromIndex) external view returns (Rates[] memory rates, bool hasMore) {
        RatesManagerStorage storage $ = _getRatesManagerStorage();

        // TODO: Fix unexpected behavior that an out of bounds query is not an error when the list is empty.
        if ($.allRates.length == 0 && fromIndex == 0) return (new Rates[](0), false);

        require(fromIndex < $.allRates.length, FromIndexOutOfRange());

        uint256 toIndex = _min(fromIndex + PAGE_SIZE, $.allRates.length);

        rates = new Rates[](toIndex - fromIndex);

        for (uint256 i; i < rates.length; ++i) {
            rates[i] = $.allRates[fromIndex + i];
        }

        hasMore = toIndex < $.allRates.length;
    }

    /**
     * @dev Returns the total number of Rates stored.
     */
    function getRatesCount() external view returns (uint256 count) {
        return _getRatesManagerStorage().allRates.length;
    }

    /* ============ Internal Functions ============ */

    /// @dev Returns the minimum of two numbers.
    function _min(uint256 a, uint256 b) internal pure returns (uint256 min) {
        return a < b ? a : b;
    }

    /* ============ Upgradeability ============ */

    /**
     * @dev   Authorizes the upgrade of the contract.
     * @param newImplementation The address of the new implementation.
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        // TODO: Consider reverting if there is no code at the new implementation address.
        require(newImplementation != address(0), ZeroImplementationAddress());
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
