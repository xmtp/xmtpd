// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

/**
 * @title  IPayerReport
 * @notice Interface for the PayerReport contract handling usage reports and batch settlements.
 */
interface IPayerReport {
    /* ============ Structs ============ */

    /**
     * @notice A struct containing the usage report details.
     * @param  originatorNode       The address of the originator node.
     * @param  startingSequenceID   The starting sequence ID of the report.
     * @param  endingSequenceID     The ending sequence ID of the report.
     * @param  lastMessageTimestamp The timestamp of the last message in the report.
     * @param  reportTimestamp      The timestamp of the report.
     * @param  reportMerkleRoot     The Merkle root of the report.
     * @param  leafCount            The number of leaves in the report. A leaf is a single (payer, amount) pair.
     */
    struct PayerReport {
        address originatorNode;
        uint256 startingSequenceID;
        uint256 endingSequenceID;
        uint256 lastMessageTimestamp;
        uint256 reportTimestamp;
        bytes32 reportMerkleRoot;
        uint16 leafCount;
    }

    /* ============ Events ============ */

    /**
     * @dev Emitted when an originator node submits a usage report.
     * The report includes the Merkle root of the detailed (payer, amount) data.
     * Note: Payers and amounts are not stored on-chain, only emitted in the event.
     * Nodes listen to this event to get full details of the report.
     */
    event PayerReportSubmitted(
        address indexed originatorNode,
        uint256 indexed reportIndex,
        bytes32 indexed reportMerkleRoot,
        uint256 startingSequenceID,
        uint256 endingSequenceID,
        uint256 lastMessageTimestamp,
        uint256 reportTimestamp,
        uint16 leafCount
    );

    /**
     * @dev Emitted when a node attests to the correctness of a report.
     */
    event PayerReportAttested(
        address indexed originatorNode, uint256 indexed reportIndex, bytes32 indexed reportMerkleRoot
    );

    /**
     * @dev Emitted when a usage report is confirmed.
     */
    event PayerReportConfirmed(
        address indexed originatorNode, uint256 indexed reportIndex, bytes32 indexed reportMerkleRoot
    );

    /**
     * @dev Emitted when a batch of usage is settled.
     */
    event PayerReportPartiallySettled(
        address indexed originatorNode,
        uint256 indexed reportIndex,
        bytes32 indexed reportMerkleRoot,
        address[] payers,
        uint256[] amounts,
        uint16 offset
    );

    /**
     * @dev Emitted when a usage report is fully settled.
     */
    event PayerReportFullySettled(
        address indexed originatorNode, uint256 indexed reportIndex, bytes32 indexed reportMerkleRoot
    );

    /* ============ Usage Report Logic ============ */

    /**
     * @notice Submits a usage report for a node covering messages from
     *         startingSequenceID to endingSequenceID.
     * @param  payerReport A struct containing the usage report details.
     *
     * Emits a PayerReportSubmitted event.
     */
    function submitPayerReport(PayerReport calldata payerReport) external;

    /**
     * @notice Allows nodes to attest to the correctness of a submitted usage report.
     * @param originatorNode The node that submitted the report.
     * @param reportIndex The index of the report.
     *
     * Emits a PayerReportAttested event.
     */
    function attestPayerReport(address originatorNode, uint256 reportIndex) external;

    /**
     * @notice Returns a list of all payer reports for a given originator node.
     * @param  originatorNode      The address of the originator node.
     * @return startingSequenceIDs The array of starting sequence IDs for each report.
     * @return reportsMerkleRoots  The array of Merkle roots for each report.
     */
    function listPayerReports(address originatorNode)
        external
        view
        returns (uint256[] memory startingSequenceIDs, bytes32[] memory reportsMerkleRoots);

    /**
     * @notice Returns summary info about a specific usage report.
     * @param  originatorNode The node that submitted the report.
     * @param  reportIndex    The index of the report.
     * @return payerReport    A PayerReport struct with the report details.
     */
    function getPayerReport(address originatorNode, uint256 reportIndex)
        external
        view
        returns (PayerReport memory payerReport);

    /**
     * @notice Settles a contiguous batch of usage data from a confirmed report.
     * Verifies an aggregated Merkle proof that the provided (payer, amount)
     * batch is included in the report's committed Merkle root, then calls the
     * settleUsage function in the Payer contract.
     *
     * @param originatorNode The node that submitted the report.
     * @param reportIndex    The index of the report.
     * @param offset         The index of the batch in the report's data (managed off-chain).
     * @param payers         A contiguous array of payer addresses.
     * @param amounts        A contiguous array of usage amounts corresponding to each payer.
     * @param proof          An aggregated Merkle proof containing branch hashes.
     *
     * Emits a UsageSettled event.
     */
    function settleUsageBatch(
        address originatorNode,
        uint256 reportIndex,
        uint16 offset,
        address[] calldata payers,
        uint256[] calldata amounts,
        bytes32[] calldata proof
    ) external;

    /**
     * @notice Sets the maximum batch size for usage settlements.
     * @param  maxBatchSize The new maximum batch size.
     */
    function setMaxBatchSize(uint256 maxBatchSize) external;

    /**
     * @notice Returns the current maximum batch size.
     * @return batchSize The current maximum batch size.
     */
    function getMaxBatchSize() external view returns (uint256 batchSize);
}
