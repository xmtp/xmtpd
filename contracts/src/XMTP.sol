// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

// Import OpenZeppelin contracts.
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract XMTP is ERC20, ERC20Burnable, Ownable {
    // Only this contract may call mintRebates.
    address public minterContract;
    // The staking contract is now exempt from locked-token checks (both sending and receiving).
    address public stakingContract;
    
    // --- Aggregated Lock Data Structures ---
    // For each account and unlock timestamp (the start of a weekly bucket), track the total tokens locked.
    mapping(address => mapping(uint256 => uint256)) private _locks;
    // For each account, store the list of unlock timestamps (weekly buckets) in ascending order.
    mapping(address => uint256[]) private _lockTimestamps;
    // To avoid iterating over expired buckets repeatedly, keep a per-account cursor.
    mapping(address => uint256) private _lockCursor;

    constructor() ERC20("XMTP", "XMTP") Ownable(msg.sender) {
        // Mint initial supply: 1 billion tokens with 18 decimals.
        _mint(msg.sender, 1_000_000_000 * 10 ** decimals());
    }

    /**
     * @notice Sets the address allowed to mint via mintRebates.
     */
    function setMinterContract(address _minterContract) external onlyOwner {
        require(_minterContract != address(0), "Invalid address");
        minterContract = _minterContract;
    }

    /**
     * @notice Sets the staking contract that is exempt from locked-token restrictions.
     */
    function setStakingContract(address _stakingContract) external onlyOwner {
        require(_stakingContract != address(0), "Invalid address");
        stakingContract = _stakingContract;
    }

    /**
     * @notice Mint new tokens for fee rebates.
     *
     * Tokens minted via this function are locked until the start of the weekly bucket (7-day period)
     * that is 26 weeks in the future. All tokens minted within the same bucket are aggregated into one
     * lock entry.
     */
    function mintRebates(address to, uint256 amount) external {
        require(msg.sender == minterContract, "Caller is not authorized to mint");
        _mint(to, amount);

        // Compute the unlock timestamp by rounding the current time down to a 7-day bucket
        // and then adding 26 weeks.
        uint256 unlockTimestamp = _computeUnlockTimestamp(block.timestamp);

        // If this bucket is new for the recipient, record its timestamp.
        if (_locks[to][unlockTimestamp] == 0) {
            _lockTimestamps[to].push(unlockTimestamp);
        }
        // Aggregate the new tokens into the appropriate weekly bucket.
        _locks[to][unlockTimestamp] += amount;
    }

    /**
     * @dev Override _beforeTokenTransfer to enforce that, except for minting, burning,
     *      or any transfer involving the staking contract (either as sender or receiver),
     *      the sender has enough unlocked tokens.
     */
    function _update(
        address from,
        address to,
        uint256 amount
    ) internal virtual override {
        
        // Skip locked-token checks for minting, burning, or if either party is the staking contract.
        if (from == address(0) || to == address(0) || from == stakingContract || to == stakingContract) {
            super._update(from, to, amount);
            return;
        }
        // Compute the total tokens still locked for the sender.
        uint256 locked = _lockedAmount(from);
        require(balanceOf(from) - locked >= amount, "Transfer amount exceeds unlocked balance");

        // Always call the underlying ERC20 transfer function.
        super._update(from, to, amount);
    }

    /**
     * @dev Computes the total locked amount for an account by iterating over its weekly lock buckets.
     *      Buckets that have expired (block.timestamp >= unlock time) are skipped and cleared.
     */
    function _lockedAmount(address account) internal returns (uint256 totalLocked) {
        uint256[] storage timestamps = _lockTimestamps[account];
        uint256 cursor = _lockCursor[account];
        uint256 len = timestamps.length;
        for (uint256 i = cursor; i < len; i++) {
            uint256 unlockTime = timestamps[i];
            if (block.timestamp >= unlockTime) {
                // This bucket has expired; update the cursor and clear storage.
                _lockCursor[account] = i + 1;
                delete _locks[account][unlockTime];
            } else {
                totalLocked += _locks[account][unlockTime];
            }
        }
        // If all buckets have expired, clear the array and reset the cursor.
        if (_lockCursor[account] >= len && len > 0) {
            delete _lockTimestamps[account];
            _lockCursor[account] = 0;
        }
        return totalLocked;
    }

    function lockedAmount(address account) external returns (uint256) {
        return _lockedAmount(account);
    }

    /**
     * @dev Computes the unlock timestamp for a deposit made at `timestamp` by rounding down to
     *      the start of the current 7-day period and then adding 26 weeks.
     *
     * This is a simplified approximation that treats one week as exactly 7 days.
     */
    function _computeUnlockTimestamp(uint256 timestamp) internal pure returns (uint256) {
        // Round down to the start of the current 7-day bucket.
        uint256 currentBucket = timestamp / 7 days;
        // Add 26 buckets (i.e., 26 weeks).
        uint256 unlockBucket = currentBucket + 26;
        // Return the start time of the unlock bucket.
        return unlockBucket * 7 days;
    }
}