// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";

contract Environment is Script {
    string public constant XMTP_GROUP_MESSAGES_OUTPUT_JSON = "GroupMessages";
    string public constant XMTP_IDENTITY_UPDATES_OUTPUT_JSON = "IdentityUpdates";
    string public constant XMTP_NODE_REGISTRY_OUTPUT_JSON = "XMTPNodeRegistry";
}
