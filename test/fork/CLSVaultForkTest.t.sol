// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {console2} from "forge-std/console2.sol";
import {CLSVault} from "../../src/vault/CLSVault.sol";
import {UniswapV3Strategy} from "../../src/strategy/UniswapV3Strategy.sol";

contract CLSVaultForkTest is Test {
    address private constant POOL = 0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640;
    address private constant NPM = 0xC36442b4a4522E871399CD717aBDD847Ab11FE88;
    address private constant WETH = 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2;
    address private constant USDC = 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48;

    CLSVault public vault;
    UniswapV3Strategy public strategy;

    function setUp() public {
        // vm.createSelectFork(vm.envString("MAINNET_RPC"), 19_000_000);
        // strategy = new UniswapV3Strategy(POOL, NPM);
        // vault = new CLSVault(address(strategy), strategy.getToken0(), strategy.getToken1());
    }
}
