// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title XMTP Group Messages Contract
contract GroupMessages is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    /// @notice Emitted when a message is sent.
    /// @param groupId The group ID.
    /// @param message The message in bytes. Contains the full mls group message payload.
    /// @param sequenceId The unique sequence ID of the message.
    event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId);

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
    /// @param _admin The address of the admin.
    function initialize(address _admin) public initializer {
        require(_admin != address(0), ZeroAdminAddress());

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        minPayloadSize = 78;
        maxPayloadSize = 4_194_304;

        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
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

    // Messaging functionality
    /// @notice Adds a message to the group.
    /// @param groupId The group ID.
    /// @param message The message in bytes.
    /// @dev Ensures the message length is within the allowed range and increments the sequence ID.
    function addMessage(bytes32 groupId, bytes calldata message) public whenNotPaused {
        require(
            message.length >= minPayloadSize && message.length <= maxPayloadSize,
            InvalidPayloadSize(message.length, minPayloadSize, maxPayloadSize)
        );

        // Increment sequence ID safely using unchecked to save gas.
        unchecked {
            sequenceId++;
        }

        emit MessageSent(groupId, message, sequenceId);
    }

    /// @notice Sets the minimum payload size
    /// @param _minPayloadSize The new minimum payload size
    /// @dev Ensures the new minimum is less than the maximum
    function setMinPayloadSize(uint256 _minPayloadSize) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_minPayloadSize < maxPayloadSize, InvalidMinPayloadSize());
        require(_minPayloadSize > 0, InvalidMinPayloadSize());
        uint256 oldSize = minPayloadSize;
        minPayloadSize = _minPayloadSize;
        emit MinPayloadSizeUpdated(oldSize, _minPayloadSize);
    }

    /// @notice Sets the maximum payload size
    /// @param _maxPayloadSize The new maximum payload size
    /// @dev Ensures the new maximum is greater than the minimum
    function setMaxPayloadSize(uint256 _maxPayloadSize) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_maxPayloadSize > minPayloadSize, InvalidMaxPayloadSize());
        require(_maxPayloadSize <= 4_194_304, InvalidMaxPayloadSize());
        uint256 oldSize = maxPayloadSize;
        maxPayloadSize = _maxPayloadSize;
        emit MaxPayloadSizeUpdated(oldSize, _maxPayloadSize);
    }

    // Upgradeability
    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
