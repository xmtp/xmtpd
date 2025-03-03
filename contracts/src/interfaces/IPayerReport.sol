// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title IPayerReport
 * @notice Updated interface for a Rewards (PayerReport) contract that:
 *  1) Accepts fee transfers from the Payer contract.
 *  2) Stores 12-hour usage reports submitted by originator nodes. Each report
 *     now includes a commitment (Merkle root) to the detailed (payer, amount)
 *     data that can be verified later.
 *  3) Allows nodes to attest to the report’s correctness.
 *  4) Finalizes reports upon majority attestation, triggering settlement of usage.
 *  5) Enables batch settlement of usage via aggregated Merkle proofs.
 *  6) Periodically distributes accumulated rewards to active node operators.
 *
 * "Active slots" data is retrieved from the Auction contract (not shown here).
 */
interface IPayerReport {

    //==============================================================
    //                          EVENTS
    //==============================================================

    /**
     * @dev Emitted when an originator node submits a usage report.
     * The report includes the Merkle root of the detailed (payer, amount) data.
     */
    event PayerReportSubmitted(
        address indexed originatorNode,
        uint256 indexed reportIndex,
        uint256 startingSequenceID,
        uint256 endingSequenceID,
        uint256 lastMessageTimestamp,
        uint256 reportTimestamp,
        bytes32 reportMerkleRoot,
        address[] payers,
        uint256[] amountsSpent
    );

    /**
     * @dev Emitted when a node attests to the correctness of a report.
     */
    event PayerReportAttested(
        address indexed originatorNode,
        bytes32 reportMerkleRoot
    );

    /**
     * @dev Emitted when a usage report is confirmed.
     */
    event PayerReportConfirmed(
        address indexed originatorNode,
        bytes32 reportMerkleRoot
    );

    //==============================================================
    //                     USAGE REPORT LOGIC
    //==============================================================

    /**
     * @notice Submits a usage report for a node covering messages from
     *         startingSequenceID to endingSequenceID.
     * @param originatorNode The node’s address/ID.
     * @param startingSequenceID The first message included in the report.
     * @param endingSequenceID The last message included in the report.
     * @param lastMessageTimestamp The timestamp of the last message in the report.
     * @param reportTimestamp The time the report was generated.
     * @param reportMerkleRoot The Merkle root of the detailed (payer, amount) data.
     * @param payers An array of payer addresses included in this usage window.
     * @param amountsSpent The usage cost each payer owes.
     *
     * Emits a PayerReportSubmitted event.
     */
    function submitPayerReport(
        address originatorNode,
        uint256 startingSequenceID,
        uint256 endingSequenceID,
        uint256 lastMessageTimestamp,
        uint256 reportTimestamp,
        bytes32 reportMerkleRoot,
        address[] calldata payers,
        uint256[] calldata amountsSpent
    ) external;

    /**
     * @notice Allows nodes to attest to the correctness of a submitted usage report.
     * @param originatorNode The node that submitted the report.
     * @param reportIndex The index of the report.
     *
     * Emits a PayerReportAttested event.
     */
    function attestPayerReport(address originatorNode, uint256 reportIndex) external;

    /**
     * @notice Finalizes a usage report once majority attestation is reached.
     * Settlement happens in batches, by calling the settleUsage function in Payer contract.
     * Emits a PayerReportConfirmed event.
     * @param originatorNode The node that submitted the report.
     * @param reportIndex The index of the report to confirm.
     */
    function confirmPayerReport(
        address originatorNode,
        uint256 reportIndex
    ) external;

    /**
     * @notice Returns a list of all payer reports for a given originator node.
     * @param originatorNode The address of the originator node.
     * @return startingSequenceID The first sequence ID in the report.
     * @return reportsMerkleRoot The Merkle root for the detailed usage data.
     */
    function listPayerReports(address originatorNode) external view returns (uint256[] memory startingSequenceID, bytes32[] memory reportsMerkleRoot);

    /**
     * @notice Returns summary info about a specific usage report.
     * @param originatorNode The node that submitted the report.
     * @param reportIndex The report's index.
     * @return startingSequenceID The first sequence ID in the report.
     * @return endingSequenceID The last sequence ID in the report.
     * @return lastMessageTimestamp The timestamp of the last message in the report.
     * @return reportTimestamp The time the report was generated.
     * @return attestationCount The number of attestations received.
     * @return isConfirmed Whether the report is finalized.
     * @return reportMerkleRoot The Merkle root for the detailed usage data.
     */
    function getPayerReport(
        address originatorNode,
        uint256 reportIndex
    )
        external
        view
        returns (
            uint256 startingSequenceID,
            uint256 endingSequenceID,
            uint256 lastMessageTimestamp,
            uint256 reportTimestamp,
            uint256 attestationCount,
            bool isConfirmed,
            bytes32 reportMerkleRoot
        );
}
