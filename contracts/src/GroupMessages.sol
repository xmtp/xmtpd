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

    // Custom errors
    error ZeroAdminAddress();
    error InvalidPayloadSize(uint256 actualSize, uint256 minSize, uint256 maxSize);

    /// @dev Minimum valid payload size (in bytes).
    uint256 public constant MIN_PAYLOAD_SIZE = 78;

    /// @dev Maximum valid payload size (4 MB).
    uint256 public constant MAX_PAYLOAD_SIZE = 4_194_304;

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
            message.length >= MIN_PAYLOAD_SIZE && message.length <= MAX_PAYLOAD_SIZE,
            InvalidPayloadSize(message.length, MIN_PAYLOAD_SIZE, MAX_PAYLOAD_SIZE)
        );

        // Increment sequence ID safely using unchecked to save gas.
        unchecked {
            sequenceId++;
        }

        emit MessageSent(groupId, message, sequenceId);
    }

    // Upgradeability
    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
