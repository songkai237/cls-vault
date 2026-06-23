// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {UniswapV3Strategy} from "../../src/strategy/UniswapV3Strategy.sol";
import {UniswapV3StrategyHarness} from "../helpers/UniswapV3StrategyHarness.sol";
import {MockPool} from "../mocks/MockPool.sol";
import {MockERC20} from "../mocks/MockERC20.sol";

contract UniswapV3StrategyTest is Test {
    int24 internal constant HALF_RANGE = 600;
    int24 internal constant TICK_SPACING = 10;

    MockERC20 internal token0;
    MockERC20 internal token1;
    MockPool internal pool;
    UniswapV3StrategyHarness internal strategy;

    address internal vault = makeAddr("vault");
    address internal owner = makeAddr("owner");
    address internal npm = makeAddr("npm");
    address internal swapRouter = makeAddr("swapRouter");

    function setUp() public {
        token0 = new MockERC20("Token0", "T0");
        token1 = new MockERC20("Token1", "T1");
        pool = new MockPool(address(token0), address(token1), 500, TICK_SPACING, 0);
        strategy = new UniswapV3StrategyHarness(address(pool), npm, swapRouter, HALF_RANGE);
        strategy.initialize(vault, owner);
    }

    function testConstructorRevertsWhenHalfRangeZero() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__InvalidTickRange.selector);
        new UniswapV3StrategyHarness(address(pool), npm, swapRouter, 0);
    }

    function testConstructorRevertsWhenHalfRangeNegative() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__InvalidTickRange.selector);
        new UniswapV3StrategyHarness(address(pool), npm, swapRouter, -100);
    }

    function testInitializeSetsVaultAndOwner() public {
        UniswapV3StrategyHarness fresh =
            new UniswapV3StrategyHarness(address(pool), npm, swapRouter, HALF_RANGE);
        fresh.initialize(vault, owner);
    }

    function testInitializeRevertsWhenAlreadyInitialized() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__AlreadyInitialized.selector);
        strategy.initialize(vault, owner);
    }

    function testDepositRevertsWhenNotVault() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__OnlyVault.selector);
        strategy.deposit();
    }

    function testWithdrawRevertsWhenNotVault() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__OnlyVault.selector);
        strategy.withdraw(1, 10, address(this));
    }

    function testRebalanceRevertsWhenNotOwner() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__OnlyOwner.selector);
        strategy.rebalance();
    }

    function testCollectRevertsWhenNotOwner() public {
        vm.expectRevert(UniswapV3Strategy.UniswapV3Strategy__OnlyOwner.selector);
        strategy.collect();
    }

    function testFloorPositiveTick() public view {
        assertEq(strategy.exposedFloor(201537, TICK_SPACING), 201530);
        assertEq(strategy.exposedFloor(201530, TICK_SPACING), 201530);
    }

    function testFloorNegativeTick() public view {
        assertEq(strategy.exposedFloor(-5, TICK_SPACING), -10);
        assertEq(strategy.exposedFloor(-10, TICK_SPACING), -10);
    }

    function testComputeTickRangeAlignsToSpacing() public view {
        int24 currentTick = 201537;
        (int24 tickLower, int24 tickUpper) = strategy.exposedComputeTickRange(currentTick);

        assertEq(tickLower % TICK_SPACING, 0);
        assertEq(tickUpper % TICK_SPACING, 0);
        assertGt(tickUpper, tickLower);

        int24 center = strategy.exposedFloor(currentTick, TICK_SPACING);
        assertEq(tickLower, strategy.exposedFloor(center - HALF_RANGE, TICK_SPACING));
        assertEq(tickUpper, strategy.exposedFloor(center + HALF_RANGE, TICK_SPACING));
    }

    function testComputeTickRangeEnsuresUpperAboveLower() public {
        (int24 tickLower, int24 tickUpper) = strategy.exposedComputeTickRange(0);
        assertGt(tickUpper, tickLower);
    }

    function testNeedRebalanceTrueWhenNoPosition() public view {
        assertTrue(strategy.exposedNeedRebalance(100, -1000, 1000));
    }

    function testNeedRebalanceTrueWhenTickBelowRange() public {
        strategy.setTokenID(1);
        assertTrue(strategy.exposedNeedRebalance(-100, 0, 1000));
    }

    function testNeedRebalanceTrueWhenTickAboveRange() public {
        strategy.setTokenID(1);
        assertTrue(strategy.exposedNeedRebalance(2000, 0, 1000));
    }

    function testNeedRebalanceFalseWhenTickInRange() public {
        strategy.setTokenID(1);
        assertFalse(strategy.exposedNeedRebalance(500, 0, 1000));
        assertFalse(strategy.exposedNeedRebalance(0, 0, 1000));
        assertFalse(strategy.exposedNeedRebalance(1000, 0, 1000));
    }

    function testOwnerRebalanceReturnsComputedRangeWithoutPositionChange() public {
        pool.setCurrentTick(201_537);

        vm.prank(owner);
        (int24 tickLower, int24 tickUpper) = strategy.rebalance();

        (int24 expectedLower, int24 expectedUpper) = strategy.exposedComputeTickRange(201_537);
        assertEq(tickLower, expectedLower);
        assertEq(tickUpper, expectedUpper);

        (,, uint128 liquidity,,,,) = strategy.getPosition();
        assertEq(liquidity, 0);
    }
}
