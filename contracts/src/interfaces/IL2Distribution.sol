// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { IERC165 } from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title  IDistribution
 * @notice Interface for distributing rewards.
 */
interface IL2Distribution is IERC165 {
    /* ============ Distribution ============ */

    /**
     * @notice Distributes the rewards.
     * @dev    Only callable by the Payer contract. When called:
     *         - The protocol fee is deducted and kept in the treasury.
     *         - The Node Operators are paid for their operational costs.
     *         - Any excess is split between:
     *           - The possibility of minting XMTP rebates to incentivize Payers.
     *           - Used to buy back and burn XMTP tokens.
     */
    function distributeRewards() external; /* onlyPayerContract */

    /* ============ Administrative Functions ============ */

    /**
     * @notice Refreshes the node operator list.
     * @dev    Only callable by the admin or active nodes.
     *         This function is used to refresh the node operator list.
     *         It is called when a new node is added, disabled or removed.
     */
    function refreshNodeOperatorList() external;

    /**
     * @notice Sets the address of the nodes contract for operator verification.
     * @param  nodesContract The address of the new nodes contract.
     *
     * Emits `NodesContractUpdated`.
     */
    function setNodesContract(address nodesContract) external;

    /**
     * @notice Sets the address of the payer contract.
     * @param  payerContract The address of the new payer contract.
     *
     * Emits `PayerContractUpdated`.
     */
    function setPayerContract(address payerContract) external;

    /**
     * @notice Sets the protocol fee.
     * @param  newProtocolFee New protocol fee (in basis points, e.g., 100 = 1%).
     */
    function setProtocolFee(uint256 newProtocolFee) external;

    /**
     * @notice Sets the rebates percentage which will be minted as XMTP rebates.
     * @param  newRebatesPercentage New rebates percentage.
     *
     * Emits `RebatesPercentageUpdated`.
     */
    function setRebatesPercentage(uint256 newRebatesPercentage) external;

    /**
     * @notice Sets the address of the USDC token contract.
     * @param  usdcToken The address of the new USDC token contract.
     *
     * Emits `UsdcTokenUpdated`.
     */
    function setUsdcToken(address usdcToken) external;

    /* ============ Getters ============ */

    /**
     * @notice Returns the current address of the nodes contract.
     * @return nodesContract The current address of the nodes contract.
     */
    function getNodesContract() external view returns (address nodesContract);

    /**
     * @notice Returns the current address of the payer contract.
     * @return payerContract The current address of the payer contract.
     */
    function getPayerContract() external view returns (address payerContract);

    /**
     * @notice Returns the current protocol fee (in basis points).
     * @return protocolFee The current protocol fee.
     */
    function getProtocolFee() external view returns (uint256 protocolFee);

    /**
     * @notice Returns the current rebates percentage (in basis points).
     * @return rebatesPercentage The current rebates percentage.
     */
    function getRebatesPercentage() external view returns (uint256 rebatesPercentage);

    /**
     * @notice Returns the current address of the USDC token contract.
     * @return usdcToken The current address of the USDC token contract.
     */
    function getUsdcToken() external view returns (address usdcToken);
}
