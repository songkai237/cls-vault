// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {console2} from "forge-std/console2.sol";
import {UniswapV3Strategy} from "../../src/strategy/UniswapV3Strategy.sol";
import {UniswapV3StrategyHarness} from "../helpers/UniswapV3StrategyHarness.sol";
import {ISwapRouter} from "../../src/interfaces/ISwapRouter.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/interfaces/IUniswapV3Pool.sol";

contract UniswapV3StrategyForkTest is Test {
    address private constant POOL = 0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640;
    address private constant NPM = 0xC36442b4a4522E871399CD717aBDD847Ab11FE88;
    address private constant SWAP_ROUTER = 0xE592427A0AEce92De3Edee1F18E0157C05861564;
    address private constant WETH = 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2;
    address private constant USDC = 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48;

    uint256 private constant MIN_SWAP_USDC = 100 * 1e6;
    uint256 private constant MIN_SWAP_WETH = 0.01 ether;

    int24 private constant HALF_RANGE = 600;

    uint256 private constant WETH_DEPOSIT = 10 ether;
    uint256 private constant USDC_DEPOSIT = 30_000 * 1e6;

    UniswapV3Strategy internal strategy;
    UniswapV3StrategyHarness internal harness;

    address internal vault = makeAddr("vault");
    address internal owner = makeAddr("owner");

    function setUp() public {
        vm.createSelectFork(vm.envString("MAINNET_RPC_URL"), 19_000_000);

        strategy = new UniswapV3Strategy(POOL, NPM, SWAP_ROUTER, HALF_RANGE, MIN_SWAP_USDC, MIN_SWAP_WETH);
        strategy.initialize(vault, owner);

        harness = new UniswapV3StrategyHarness(POOL, NPM, SWAP_ROUTER, HALF_RANGE, MIN_SWAP_USDC, MIN_SWAP_WETH);
    }

    function _fundStrategy() internal {
        deal(WETH, address(strategy), WETH_DEPOSIT);
        deal(USDC, address(strategy), USDC_DEPOSIT);
    }

    function _currentTick() internal view returns (int24 tick) {
        (, tick,,,,,) = IUniswapV3Pool(POOL).slot0();
    }

    function _expectedTickRange(int24 currentTick) internal view returns (int24 tickLower, int24 tickUpper) {
        (tickLower, tickUpper) = harness.exposedComputeTickRange(currentTick);
    }

    /// @dev Push pool tick above `upperBound` by swapping WETH into USDC.
    function _pushTickAbove(int24 upperBound) internal {
        deal(WETH, address(this), 50_000 ether);
        IERC20(WETH).approve(SWAP_ROUTER, type(uint256).max);

        uint256 amountIn = 500 ether;
        while (_currentTick() <= upperBound && amountIn <= 10_000 ether) {
            ISwapRouter(SWAP_ROUTER).exactInputSingle(
                ISwapRouter.ExactInputSingleParams({
                    tokenIn: WETH,
                    tokenOut: USDC,
                    fee: 500,
                    recipient: address(this),
                    deadline: block.timestamp,
                    amountIn: amountIn,
                    amountOutMinimum: 0,
                    sqrtPriceLimitX96: 0
                })
            );
            amountIn += 500 ether;
        }
        console2.log("pushed tick above", upperBound);
        console2.log("current tick", _currentTick());
    }

    function testFirstDepositMintsPositionCenteredOnCurrentTick() public {
        _fundStrategy();

        int24 tickBefore = _currentTick();
        (int24 expectedLower, int24 expectedUpper) = _expectedTickRange(tickBefore);

        vm.prank(vault);
        strategy.deposit();

        (,, uint128 liquidity, int24 tickLower, int24 tickUpper,,) = strategy.getPosition();
        assertGt(liquidity, 0);
        assertEq(tickLower, expectedLower);
        assertEq(tickUpper, expectedUpper);
        assertGe(tickBefore, tickLower);
        assertLe(tickBefore, tickUpper);
    }

    function testSecondDepositInRangeIncreasesLiquidityWithoutChangingRange() public {
        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        (,, uint128 liquidityBefore, int24 lowerBefore, int24 upperBefore,,) = strategy.getPosition();

        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        (,, uint128 liquidityAfter, int24 lowerAfter, int24 upperAfter,,) = strategy.getPosition();
        assertGt(liquidityAfter, liquidityBefore);
        assertEq(lowerAfter, lowerBefore);
        assertEq(upperAfter, upperBefore);
    }

    function testDepositRebalancesWhenOutOfRange() public {
        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        (,,, int24 oldLower, int24 oldUpper,,) = strategy.getPosition();
        _pushTickAbove(oldUpper);
        int24 tickAfterSwap = _currentTick();
        assertGt(tickAfterSwap, oldUpper);

        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        (,, uint128 liquidity, int24 newLower, int24 newUpper,,) = strategy.getPosition();
        assertGt(liquidity, 0);
        assertFalse(newLower == oldLower && newUpper == oldUpper);

        (int24 expectedLower, int24 expectedUpper) = _expectedTickRange(tickAfterSwap);
        assertEq(newLower, expectedLower);
        assertEq(newUpper, expectedUpper);
        assertGe(tickAfterSwap, newLower);
        assertLe(tickAfterSwap, newUpper);
    }

    function testOwnerRebalanceReturnsNewRangeWithoutMutatingPosition() public {
        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        (,, uint128 liquidityBefore, int24 posLowerBefore, int24 posUpperBefore,,) = strategy.getPosition();
        int24 tick = _currentTick();

        vm.prank(owner);
        (int24 tickLower, int24 tickUpper) = strategy.rebalance();

        (int24 expectedLower, int24 expectedUpper) = _expectedTickRange(tick);
        assertEq(tickLower, expectedLower);
        assertEq(tickUpper, expectedUpper);

        (,, uint128 liquidityAfter, int24 posLowerAfter, int24 posUpperAfter,,) = strategy.getPosition();
        assertEq(liquidityAfter, liquidityBefore);
        assertEq(posLowerAfter, posLowerBefore);
        assertEq(posUpperAfter, posUpperBefore);
    }

    function testPartialWithdrawReturnsProRataTokens() public {
        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        (uint256 total0, uint256 total1) = strategy.getTotalAssets();
        uint256 totalShares = 1000;
        uint256 withdrawShares = 250;

        uint256 expected0 = total0 * withdrawShares / totalShares;
        uint256 expected1 = total1 * withdrawShares / totalShares;

        address recipient = makeAddr("recipient");
        vm.prank(vault);
        (uint256 amount0, uint256 amount1) = strategy.withdraw(withdrawShares, totalShares, recipient);

        assertApproxEqAbs(amount0, expected0, 1e4);
        assertApproxEqAbs(amount1, expected1, 1e15);
        assertEq(IERC20(USDC).balanceOf(recipient), amount0);
        assertEq(IERC20(WETH).balanceOf(recipient), amount1);
    }

    function testFullWithdrawClosesPosition() public {
        _fundStrategy();
        vm.prank(vault);
        strategy.deposit();

        address recipient = makeAddr("recipient");
        vm.prank(vault);
        (uint256 amount0, uint256 amount1) = strategy.withdraw(100, 100, recipient);

        assertGt(amount0 + amount1, 0);
        (,, uint128 liquidity,,,,) = strategy.getPosition();
        assertEq(liquidity, 0);
        assertEq(IERC20(USDC).balanceOf(address(strategy)), 0);
        assertEq(IERC20(WETH).balanceOf(address(strategy)), 0);
    }
}
