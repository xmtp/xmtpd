// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/introspection/ERC165Checker.sol";
import "./interfaces/INodes.sol";

/// @title XMTP Governance
/// @notice This contract provides governance functions for the XMTP network.
contract XMTPGovernance is Ownable {
    /// @notice The XMTP Node Registry contract.
    INodes public xmtpNodeRegistry;

    /// @notice The interface ID for the INodes contract.
    bytes4 private constant _INTERFACE_ID_NODES = type(INodes).interfaceId;

    /// @notice The error thrown when an unauthorized call is made.
    error Unauthorized();

    /// @notice Initializes the XMTP Governance contract.
    constructor(address _owner) Ownable(_owner) {}
    
    /// @notice Sets the XMTP Node Registry contract.
    /// @param _xmtpNodeRegistry The address of the XMTP Node Registry contract.
    function setXMTPNodeRegistry(address _xmtpNodeRegistry) external onlyOwner {
        require(_xmtpNodeRegistry != address(0), "Invalid address");

        require(
            ERC165Checker.supportsInterface(_xmtpNodeRegistry, _INTERFACE_ID_NODES),
            "Contract does not support INodes interface"
        );

        xmtpNodeRegistry = INodes(_xmtpNodeRegistry);
    }

    /// @notice Modifier that ensures only owners of an "active" Node NFT can call the function.
    /// @param tokenId The token ID of the Node NFT.
    modifier onlyActiveNode(uint256 tokenId) {
        require(xmtpNodeRegistry.ownerOf(tokenId) == msg.sender, Unauthorized());
        INodes.Node memory nodeData = xmtpNodeRegistry.getNode(tokenId);
        require(nodeData.isActive, Unauthorized());
        _;
    }
}
