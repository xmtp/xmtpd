// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title XMTP Identity Updates Contract
contract IdentityUpdates is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    /// @notice Emitted when an identity update is sent.
    /// @param inboxId The inbox ID.
    /// @param update The identity update in bytes. Contains the full mls identity update payload.
    /// @param sequenceId The unique sequence ID of the identity update.
    event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId);

    /// @notice Emitted when an upgrade is authorized.
    /// @param upgrader The EOA authorizing the upgrade.
    /// @param newImplementation The address of the new implementation.
    event UpgradeAuthorized(address upgrader, address newImplementation);

    /// @notice Emitted when the minimum payload size is updated.
    /// @param oldSize The old minimum payload size.
    /// @param newSize The new minimum payload size.
    event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize);

    /// @notice Emitted when the maximum payload size is updated.
    /// @param oldSize The old maximum payload size.
    /// @param newSize The new maximum payload size.
    event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize);

    // Custom errors
    error ZeroAdminAddress();
    error InvalidPayloadSize(uint256 actualSize, uint256 minSize, uint256 maxSize);
    error InvalidMaxPayloadSize();
    error InvalidMinPayloadSize();

    /// @dev Minimum valid payload size (in bytes).
    // slither-disable-next-line constable-states
    uint256 public minPayloadSize;

    /// @dev Maximum valid payload size (in bytes).
    // slither-disable-next-line constable-states
    uint256 public maxPayloadSize;

    // State variables
    // slither-disable-next-line unused-state,constable-states
    uint64 private sequenceId;

    /// @dev Reserved storage gap for future upgrades
    // slither-disable-next-line unused-state,naming-convention
    uint256[50] private __gap;

    // Initialization
    /// @notice Initializes the contract with the deployer as admin.
    /// @param admin The address of the admin.
    function initialize(address admin) public initializer {
        require(admin != address(0), ZeroAdminAddress());

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        minPayloadSize = 104;
        maxPayloadSize = 4_194_304;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
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

    // IdentityUpdate functionality
    /// @notice Adds an identity update to an specific inbox ID.
    /// @param inboxId The inbox ID.
    /// @param update The identity update in bytes.
    function addIdentityUpdate(bytes32 inboxId, bytes calldata update) public whenNotPaused {
        require(
            update.length >= minPayloadSize && update.length <= maxPayloadSize,
            InvalidPayloadSize(update.length, minPayloadSize, maxPayloadSize)
        );

        // Increment sequence ID safely using unchecked to save gas.
        unchecked {
            sequenceId++;
        }

        emit IdentityUpdateCreated(inboxId, update, sequenceId);
    }

    /// @notice Sets the minimum payload size
    /// @param minPayloadSizeRequest The new minimum payload size
    /// @dev Ensures the new minimum is less than the maximum
    function setMinPayloadSize(uint256 minPayloadSizeRequest) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(minPayloadSizeRequest < maxPayloadSize, InvalidMinPayloadSize());
        require(minPayloadSizeRequest > 0, InvalidMinPayloadSize());
        uint256 oldSize = minPayloadSize;
        minPayloadSize = minPayloadSizeRequest;
        emit MinPayloadSizeUpdated(oldSize, minPayloadSizeRequest);
    }

    /// @notice Sets the maximum payload size
    /// @param maxPayloadSizeRequest The new maximum payload size
    /// @dev Ensures the new maximum is greater than the minimum
    function setMaxPayloadSize(uint256 maxPayloadSizeRequest) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(maxPayloadSizeRequest > minPayloadSize, InvalidMaxPayloadSize());
        require(maxPayloadSizeRequest <= 4_194_304, InvalidMaxPayloadSize());
        uint256 oldSize = maxPayloadSize;
        maxPayloadSize = maxPayloadSizeRequest;
        emit MaxPayloadSizeUpdated(oldSize, maxPayloadSizeRequest);
    }

    // Upgradeability
    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
