// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Vm.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Test, console} from "forge-std/src/Test.sol";
import {Utils} from "./utils/Utils.sol";
import {XMTP} from "../src/XMTP.sol";
import {XMTPStaking} from "../src/Staking.sol";

contract USDC is ERC20 {
    constructor() ERC20("USDC", "USDC") {
        _mint(msg.sender, 1000000000000000000);
    }
}

contract XMTPTest is Test, Utils {
    XMTP xmtpImpl;
    USDC usdcImpl;
    XMTPStaking stakingImpl;

    address admin = address(this);
    address unauthorized = address(0x1);

    function setUp() public {
        xmtpImpl = new XMTP();
        usdcImpl = new USDC();
        stakingImpl = new XMTPStaking(xmtpImpl, usdcImpl, "XMTP Staking", "XMTP-STAKING");
    }

    function testXmtpToken() public {
        address staker1 = vm.randomAddress();
        address staker2 = vm.randomAddress();
        address staker3 = vm.randomAddress();
        xmtpImpl.transfer(staker1, 1000);
        xmtpImpl.transfer(staker2, 2000);
        xmtpImpl.transfer(staker3, 3000);

        vm.startPrank(staker1);
        require(xmtpImpl.balanceOf(staker1) == 1000, "Staker 1 should have 1000 XMTP");
        require(xmtpImpl.lockedAmount(staker1) == 0, "Staker 1 should have no locked XMTP");
        xmtpImpl.approve(address(stakingImpl), 1000);
        stakingImpl.deposit(1000, staker1);

        vm.startPrank(staker2);
        xmtpImpl.approve(address(stakingImpl), 2000);
        stakingImpl.deposit(2000, staker2);

        vm.assertEq(stakingImpl.withdrawableRewardOf(staker1), 0);
        vm.assertEq(stakingImpl.withdrawableRewardOf(staker2), 0);

        vm.startPrank(admin);
        usdcImpl.approve(address(stakingImpl), 1000000000000000000);
        stakingImpl.depositReward(100);

        vm.assertEq(stakingImpl.withdrawableRewardOf(staker1), 33);
        vm.assertEq(stakingImpl.withdrawableRewardOf(staker2), 66);
        
        vm.startPrank(staker3);
        xmtpImpl.approve(address(stakingImpl), 3000);
        stakingImpl.deposit(100, staker3);

        vm.assertEq(stakingImpl.withdrawableRewardOf(staker3), 0);
    }
}
