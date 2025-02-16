// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title XMTPVault
 * @notice An ERC4626 vault that accepts deposits of XMTP (the asset)
 * and periodically receives USDC reward tokens.
 *
 * Depositors' shares earn USDC rewards (only for rewards deposited
 * after they deposit) via a dividend‐paying mechanism.
 *
 * This implementation uses a "magnified rewards per share" pattern similar
 * to that popularized in dividend–paying token implementations.
 *
 */
contract XMTPStaking is ERC4626, Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    /// @notice The reward token (USDC) distributed to vault share holders.
    IERC20 public rewardToken;

    // ================================================================
    // Dividend accounting: "magnified rewards per share" pattern.
    // A very high precision factor is used to allow for fractions.
    // ================================================================
    uint256 public constant MAGNITUDE = 2**128;
    uint256 public magnifiedRewardPerShare;
    mapping(address => int256) public magnifiedRewardCorrections;
    mapping(address => uint256) public withdrawnRewards;

    /**
     * @notice Constructor.
     * @param _asset The XMTP token (asset accepted by the vault)
     * @param _rewardToken The USDC token (reward token)
     * @param name_ The ERC20 name for the vault's share token
     * @param symbol_ The ERC20 symbol for the vault's share token
     */
    constructor(
        IERC20 _asset,
        IERC20 _rewardToken,
        string memory name_,
        string memory symbol_
    )
        ERC20(name_, symbol_)
        ERC4626(_asset)
        Ownable(msg.sender)
    {
        rewardToken = _rewardToken;
    }

    // ================================================================
    // Reward (USDC) functions
    // ================================================================

    function depositReward(uint256 amount) external {
        require(totalSupply() > 0, "No deposits in vault");
        rewardToken.safeTransferFrom(msg.sender, address(this), amount);
        // Increase the accumulated rewards per share.
        magnifiedRewardPerShare += (amount * MAGNITUDE) / totalSupply();
    }

    function claimRewards() external nonReentrant {
        uint256 _withdrawableReward = withdrawableRewardOf(msg.sender);
        require(_withdrawableReward > 0, "No rewards to claim");
        withdrawnRewards[msg.sender] += _withdrawableReward;
        rewardToken.safeTransfer(msg.sender, _withdrawableReward);
    }

    function withdrawableRewardOf(address account) public view returns (uint256) {
        return accumulativeRewardOf(account) - withdrawnRewards[account];
    }

    function accumulativeRewardOf(address account) public view returns (uint256) {
        return uint256(
            int256(balanceOf(account) * magnifiedRewardPerShare) + magnifiedRewardCorrections[account]
        ) / MAGNITUDE;
    }

    /**
     * @dev Overrides _deposit
     *
     * When new shares are minted (i.e. when someone deposits XMTP), we subtract
     * from the depositor's correction so that they do not retroactively receive
     * rewards distributed before their deposit.
     *
     * @param caller The address making the deposit.
     * @param receiver The address receiving the shares.
     * @param assets The number of assets deposited.
     * @param shares The number of vault shares minted.
     */

    function _deposit(address caller, address receiver, uint256 assets, uint256 shares) internal override {
        super._deposit(caller, receiver, assets, shares);
        // Correct the reward correction for the new shares
        magnifiedRewardCorrections[receiver] -= int256(shares * magnifiedRewardPerShare);
    }

    /**
     * @dev Overrides _withdraw
     *
     * When shares are burned (i.e. when someone withdraws XMTP), we adjust
     * the dividend correction so that the account keeps the rewards earned on the
     * shares it is burning.
     *
     * @param caller The address initiating the withdrawal.
     * @param receiver The address receiving the assets.
     * @param owner The address owning the shares.
     * @param assets The number of assets withdrawn.
     * @param shares The number of vault shares that will be burned.
     */
    function _withdraw(
        address caller,
        address receiver,
        address owner,
        uint256 assets,
        uint256 shares
    ) internal override {
    super._withdraw(caller, receiver, owner, assets, shares);
    // Add back to ensure the account keeps the rewards earned on the shares burned.
    magnifiedRewardCorrections[owner] += int256(shares * magnifiedRewardPerShare);
}

    // We still override _update to adjust corrections on share transfers.
    function _update(address from, address to, uint256 amount) internal virtual override {
        if (from != address(0) && to != address(0)) {
            int256 magCorrection = int256(magnifiedRewardPerShare * amount);
            magnifiedRewardCorrections[from] += magCorrection;
            magnifiedRewardCorrections[to] -= magCorrection;
        }
        
        super._update(from, to, amount);
    }
}