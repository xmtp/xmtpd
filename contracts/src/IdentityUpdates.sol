// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin-contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin-contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title XMTP Identity Updates Contract
contract IdentityUpdates is Initializable, AccessControlUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId);
    event UpgradeAuthorized(address deployer, address newImplementation);

    uint64 private sequenceId;

    /// @notice Initializes the contract with the deployer as admin.
    function initialize() public initializer {
        __UUPSUpgradeable_init();
        __AccessControl_init();
        __Pausable_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }

    /// @notice Pauses the contract.
    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /// @notice Unpauses the contract.
    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /// @notice Adds an identity update to an specific inbox ID.
    /// @param inboxId The inbox ID.
    /// @param update The identity update in bytes.
    function addIdentityUpdate(bytes32 inboxId, bytes calldata update) public whenNotPaused {
        /// @dev Incrementing the sequence ID is safe here due to the extremely large limit of uint64.
        unchecked {
            sequenceId++;
        }

        emit IdentityUpdateCreated(inboxId, update, sequenceId);
    }

    /// @dev Authorizes the upgrade of the contract.
    /// @param newImplementation The address of the new implementation.
    function _authorizeUpgrade(address newImplementation) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        emit UpgradeAuthorized(msg.sender, newImplementation);
    }
}
