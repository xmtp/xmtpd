// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Vm.sol";
import {Test, console} from "forge-std/src/Test.sol";
import {Utils} from "./utils/Utils.sol";
import {XMTP} from "../src/XMTP.sol";

contract XMTPTest is Test, Utils {
    XMTP xmtpImpl;

    address admin = address(this);
    address unauthorized = address(0x1);

    function setUp() public {
        xmtpImpl = new XMTP();
    }

    function testXmtpToken() public {
        
    }
}
