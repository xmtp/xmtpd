// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { AccessControlUpgradeable } from "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

// TODO: IIdentityUpdates
// TODO: Abstract PayloadBroadcaster.

/// @title XMTP Identity Updates Contract
contract IdentityUpdates is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    /* ============ Events ============ */

    /**
     * @notice Emitted when an identity update is sent.
     * @param  inboxId    The inbox ID.
     * @param  update     The identity update in bytes. Contains the full mls identity update payload.
     * @param  sequenceId The unique sequence ID of the identity update.
     */
    event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId); // TODO: indexed inboxId and
        // sequenceId.

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
    error FailedToGrantRole(bytes32 role, address account);

    /* ============ Constants ============ */

    uint256 public constant STARTING_MIN_PAYLOAD_SIZE = 78;
    uint256 public constant ABSOLUTE_MAX_PAYLOAD_SIZE = 4_194_304;
    uint256 public constant ABSOLUTE_MIN_PAYLOAD_SIZE = 1;

    /* ============ UUPS Storage ============ */

    /// @custom:storage-location erc7201:xmtp.storage.IdentityUpdates
    struct IdentityUpdatesStorage {
        uint256 minPayloadSize;
        uint256 maxPayloadSize;
        uint64 sequenceId;
    }

    // keccak256(abi.encode(uint256(keccak256("xmtp.storage.IdentityUpdates")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 internal constant IDENTITY_UPDATES_STORAGE_LOCATION =
        0x92f6d7b379434335724ccaa6ce32661f25de0b6cb746fac5f5edaed4b9685e00;

    function _getIdentityUpdatesStorage() internal pure returns (IdentityUpdatesStorage storage $) {
        // slither-disable-next-line assembly
        assembly {
            $.slot := IDENTITY_UPDATES_STORAGE_LOCATION
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

        IdentityUpdatesStorage storage $ = _getIdentityUpdatesStorage();

        $.minPayloadSize = STARTING_MIN_PAYLOAD_SIZE;
        $.maxPayloadSize = ABSOLUTE_MAX_PAYLOAD_SIZE;

        require(_grantRole(DEFAULT_ADMIN_ROLE, admin), FailedToGrantRole(DEFAULT_ADMIN_ROLE, admin));
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

    /* ============ IdentityUpdate functionality ============ */

    /**
     * @notice Adds an identity update to an specific inbox ID.
     * @param  inboxId The inbox ID.
     * @param  update  The identity update in bytes.
     */
    function addIdentityUpdate(bytes32 inboxId, bytes calldata update) external whenNotPaused {
        IdentityUpdatesStorage storage $ = _getIdentityUpdatesStorage();

        require(
            update.length >= $.minPayloadSize && update.length <= $.maxPayloadSize,
            InvalidPayloadSize(update.length, $.minPayloadSize, $.maxPayloadSize)
        );

        // Increment sequence ID safely using unchecked to save gas.
        unchecked {
            emit IdentityUpdateCreated(inboxId, update, ++$.sequenceId);
        }
    }

    /* ============ Payload Size Setters ============ */

    /**
     * @notice Sets the minimum payload size.
     * @param  minPayloadSizeRequest The new minimum payload size.
     * @dev    Ensures the new minimum is less than the maximum.
     */
    function setMinPayloadSize(uint256 minPayloadSizeRequest) external onlyRole(DEFAULT_ADMIN_ROLE) {
        IdentityUpdatesStorage storage $ = _getIdentityUpdatesStorage();

        require(minPayloadSizeRequest <= $.maxPayloadSize, InvalidMinPayloadSize());
        require(minPayloadSizeRequest >= ABSOLUTE_MIN_PAYLOAD_SIZE, InvalidMinPayloadSize());

        uint256 oldSize = $.minPayloadSize;

        emit MinPayloadSizeUpdated(oldSize, $.minPayloadSize = minPayloadSizeRequest);
    }

    /**
     * @notice Sets the maximum payload size.
     * @param  maxPayloadSizeRequest The new maximum payload size.
     * @dev    Ensures the new maximum is greater than the minimum.
     */
    function setMaxPayloadSize(uint256 maxPayloadSizeRequest) external onlyRole(DEFAULT_ADMIN_ROLE) {
        IdentityUpdatesStorage storage $ = _getIdentityUpdatesStorage();

        require(maxPayloadSizeRequest > $.minPayloadSize, InvalidMaxPayloadSize());
        require(maxPayloadSizeRequest <= ABSOLUTE_MAX_PAYLOAD_SIZE, InvalidMaxPayloadSize());

        uint256 oldSize = $.maxPayloadSize;

        emit MaxPayloadSizeUpdated(oldSize, $.maxPayloadSize = maxPayloadSizeRequest);
    }

    /* ============ Getters ============ */

    /// @notice Minimum valid payload size (in bytes).
    function minPayloadSize() external view returns (uint256) {
        return _getIdentityUpdatesStorage().minPayloadSize;
    }

    /// @notice Maximum valid payload size (in bytes).
    function maxPayloadSize() external view returns (uint256) {
        return _getIdentityUpdatesStorage().maxPayloadSize;
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
