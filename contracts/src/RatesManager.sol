// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/access/extensions/AccessControlDefaultAdminRules.sol";

contract RatesManager is AccessControlDefaultAdminRules {
    // Create a dedicated role for managing rates
    bytes32 public constant RATES_ADMIN_ROLE = keccak256("RATES_ADMIN_ROLE");

    // Rates struct holds the fees and the start time of the rates
    struct Rates {
        uint64 messageFee;
        uint64 storageFee;
        uint64 congestionFee;
        uint64 startTime;
    }

    // All Rates appended here
    Rates[] private allRates;

    // Event emitted when new Rates are added
    event RatesAdded(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime);

    constructor()
        AccessControlDefaultAdminRules(
            1 days, // adminTransferDelay
            msg.sender // initialAdmin
        )
    {
        // Setup RATES_ADMIN_ROLE so that the default admin can grant it
        _setRoleAdmin(RATES_ADMIN_ROLE, DEFAULT_ADMIN_ROLE);

        // Grant the deployer the RATES_ADMIN_ROLE
        _grantRole(RATES_ADMIN_ROLE, msg.sender);
    }

    /**
     * @dev Add new Rates. Can only be called by addresses with RATES_ADMIN_ROLE.
     *      The array only grows; we do not allow removal or updating.
     */
    function addRates(uint64 messageFee, uint64 storageFee, uint64 congestionFee, uint64 startTime)
        external
        onlyRole(RATES_ADMIN_ROLE)
    {
        // Enforce chronological order
        require(
            allRates.length == 0 || startTime > allRates[allRates.length - 1].startTime,
            "startTime must be greater than the last startTime"
        );

        allRates.push(
            Rates({messageFee: messageFee, storageFee: storageFee, congestionFee: congestionFee, startTime: startTime})
        );

        emit RatesAdded(messageFee, storageFee, congestionFee, startTime);
    }

    /**
     * @dev Returns a slice of the Rates list for pagination.
     * @param fromIndex Index from which to start (must be < allRates.length).
     * @return rates - The subset of Rates.
     * @return hasMore - True if there are more items beyond this slice.
     */
    function getRates(uint256 fromIndex) external view returns (Rates[] memory rates, bool hasMore) {
        if (allRates.length == 0 && fromIndex == 0) {
            return (new Rates[](0), false);
        }

        require(fromIndex < allRates.length, "fromIndex out of range");

        uint256 pageSize = 50; // Fixed page size
        uint256 toIndex = fromIndex + pageSize;
        if (toIndex > allRates.length) {
            toIndex = allRates.length;
        }

        uint256 resultSize = toIndex - fromIndex;
        Rates[] memory tempRates = new Rates[](resultSize);
        for (uint256 i = 0; i < resultSize; i++) {
            tempRates[i] = allRates[fromIndex + i];
        }

        bool moreData = (toIndex < allRates.length);
        return (tempRates, moreData);
    }

    /**
     * @dev Returns the total number of Rates stored.
     */
    function getRatesCount() external view returns (uint256) {
        return allRates.length;
    }
}
