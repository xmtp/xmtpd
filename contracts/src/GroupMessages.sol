// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title XMTP Group Messages Contract
contract GroupMessages is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId);
    event UpgradeAuthorized(address deployer, address newImplementation);

    error InvalidIdentityUpdateSize(uint256 actualSize, uint256 minSize, uint256 maxSize);

    uint256 private constant MIN_PAYLOAD_SIZE = 78;
    uint256 private constant MAX_PAYLOAD_SIZE = 4_194_304;

    error InvalidMessage();

    uint64 private sequenceId;

    /// @dev Reserved storage gap for future upgrades
    uint256[50] private __gap;

    /// @notice Initializes the contract with the deployer as admin.
    function initialize(address _admin) public initializer {
        require(_admin != address(0), "Admin address cannot be zero");
        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
    }

    /// @notice Pauses the contract.
    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /// @notice Unpauses the contract.
    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /// @notice Adds a message to the group.
    /// @param groupId The group ID.
    /// @param message The message in bytes.
    function addMessage(bytes32 groupId, bytes calldata message) public whenNotPaused {
        require(
            message.length >= MIN_PAYLOAD_SIZE && message.length <= MAX_PAYLOAD_SIZE,
            InvalidMessage()
        );

        /// @dev Incrementing the sequence ID is safe here due to the extremely large limit of uint64.
        unchecked {
            sequenceId++;
        }

        emit MessageSent(groupId, message, sequenceId);
    }

    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "New implementation cannot be zero address");
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
