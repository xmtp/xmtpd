pragma solidity 0.8.28;

contract NodeSlots {
    struct Bid {
        uint256 nodeID;
        // Mapping from address to the amount bid
        mapping(address => uint256) delegations;
        // Total amount of delegations
        uint256 totalDelegations;
        // The fee that the node operator receives for the slot
        // which is paid out 100% before the delegators receive their funds.
        uint256 nodeOperatorFee;
    }

    struct ActiveNode {
        uint256 expiresAt;
        Bid bid;
    }

    // Mapping from slot ID to the winning bid
    mapping(uint16 => ActiveNode) private currentNodes;

    mapping(uint16 => Bid) public pendingBids;
}
