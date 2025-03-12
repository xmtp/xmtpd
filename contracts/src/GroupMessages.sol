// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { AccessControlUpgradeable } from "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

// TODO: IGroupMessages
// TODO: Abstract PayloadBroadcaster.

/// @title XMTP Group Messages Contract
contract GroupMessages is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    /* ============ Events ============ */

    /**
     * @notice Emitted when a message is sent.
     * @param  groupId    The group ID.
     * @param  message    The message in bytes. Contains the full mls group message payload.
     * @param  sequenceId The unique sequence ID of the message.
     */
    event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId); // TODO: indexed groupId and sequenceId.

    /**
     * @notice Emitted when an upgrade is authorized.
     * @param  upgrader          The EOA authorizing the upgrade.
     * @param  newImplementation The address of the new implementation.
     */
    event UpgradeAuthorized(address upgrader, address newImplementation); // TODO: both indexed.

    /**
     * @notice Emitted when the minimum payload size is updated.
     * @param  oldSize The old minimum payload size.
     * @param  newSize The new minimum payload size.
     */
    event MinPayloadSizeUpdated(uint256 oldSize, uint256 newSize);

    /**
     * @notice Emitted when the maximum payload size is updated.
     * @param  oldSize The old maximum payload size.
     * @param  newSize The new maximum payload size.
     */
    event MaxPayloadSizeUpdated(uint256 oldSize, uint256 newSize);

    /* ============ Custom Errors ============ */

    error ZeroAdminAddress();
    error InvalidPayloadSize(uint256 actualSize, uint256 minSize, uint256 maxSize);
    error InvalidMaxPayloadSize();
    error InvalidMinPayloadSize();
    error ZeroImplementationAddress();

    /* ============ Constants ============ */

    uint256 public constant ABSOLUTE_MIN_PAYLOAD_SIZE = 78;
    uint256 public constant ABSOLUTE_MAX_PAYLOAD_SIZE = 4_194_304;

    /* ============ UUPS Storage ============ */

    /// @custom:storage-location erc7201:xmtp.storage.GroupMessages
    struct GroupMessagesStorage {
        uint256 minPayloadSize;
        uint256 maxPayloadSize;
        uint64 sequenceId;
    }

    // keccak256(abi.encode(uint256(keccak256("xmtp.storage.GroupMessages")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 internal constant GROUP_MESSAGE_STORAGE_LOCATION =
        0x5d34bcd3bd75a3e15b8380222f0e4a5877bc3f258e24e1caa87a1298d2a61000;

    function _getGroupMessagesStorage() internal pure returns (GroupMessagesStorage storage $) {
        // slither-disable-next-line assembly
        assembly {
            $.slot := GROUP_MESSAGE_STORAGE_LOCATION
        }
    }

    /* ============ Initialization ============ */

    /**
     * @notice Initializes the contract with the deployer as admin.
     * @param  admin The address of the admin.
     */
    function initialize(address admin) external initializer {
        require(admin != address(0), ZeroAdminAddress());

        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        GroupMessagesStorage storage $ = _getGroupMessagesStorage();

        $.minPayloadSize = ABSOLUTE_MIN_PAYLOAD_SIZE;
        $.maxPayloadSize = ABSOLUTE_MAX_PAYLOAD_SIZE;

        // slither-disable-next-line unused-return
        _grantRole(DEFAULT_ADMIN_ROLE, admin); // Will return false if the role is already granted.
    }

    /* ============ Pausable functionality ============ */

    /**
     * @notice Pauses the contract, restricting certain actions.
     * @dev    Callable only by accounts with the DEFAULT_ADMIN_ROLE.
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpauses the contract, allowing normal operations.
     * @dev    Callable only by accounts with the DEFAULT_ADMIN_ROLE.
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /* ============ Messaging functionality ============ */

    /**
     * @notice Adds a message to the group.
     * @param  groupId The group ID.
     * @param  message The message in bytes.
     * @dev    Ensures the message length is within the allowed range and increments the sequence ID.
     */
    function addMessage(bytes32 groupId, bytes calldata message) external whenNotPaused {
        GroupMessagesStorage storage $ = _getGroupMessagesStorage();

        if (message.length < $.minPayloadSize || message.length > $.maxPayloadSize) {
            revert InvalidPayloadSize(message.length, $.minPayloadSize, $.maxPayloadSize);
        }

        // Increment sequence ID safely using unchecked to save gas.
        unchecked {
            emit MessageSent(groupId, message, ++$.sequenceId);
        }
    }

    /* ============ Payload Size Setters ============ */

    /**
     * @notice Sets the minimum payload size.
     * @param  minPayloadSizeRequest The new minimum payload size.
     * @dev    Ensures the new minimum is less than the maximum.
     */
    function setMinPayloadSize(uint256 minPayloadSizeRequest) external onlyRole(DEFAULT_ADMIN_ROLE) {
        GroupMessagesStorage storage $ = _getGroupMessagesStorage();

        if (minPayloadSizeRequest > $.maxPayloadSize || minPayloadSizeRequest < ABSOLUTE_MIN_PAYLOAD_SIZE) {
            revert InvalidMinPayloadSize();
        }

        uint256 oldSize = $.minPayloadSize;

        emit MinPayloadSizeUpdated(oldSize, $.minPayloadSize = minPayloadSizeRequest);
    }

    /**
     * @notice Sets the maximum payload size.
     * @param  maxPayloadSizeRequest The new maximum payload size.
     * @dev    Ensures the new maximum is greater than the minimum.
     */
    function setMaxPayloadSize(uint256 maxPayloadSizeRequest) external onlyRole(DEFAULT_ADMIN_ROLE) {
        GroupMessagesStorage storage $ = _getGroupMessagesStorage();

        if (maxPayloadSizeRequest < $.minPayloadSize || maxPayloadSizeRequest > ABSOLUTE_MAX_PAYLOAD_SIZE) {
            revert InvalidMaxPayloadSize();
        }

        uint256 oldSize = $.maxPayloadSize;

        emit MaxPayloadSizeUpdated(oldSize, $.maxPayloadSize = maxPayloadSizeRequest);
    }

    /* ============ Getters ============ */

    /// @notice Minimum valid payload size (in bytes).
    function minPayloadSize() external view returns (uint256 size) {
        return _getGroupMessagesStorage().minPayloadSize;
    }

    /// @notice Maximum valid payload size (in bytes).
    function maxPayloadSize() external view returns (uint256 size) {
        return _getGroupMessagesStorage().maxPayloadSize;
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
