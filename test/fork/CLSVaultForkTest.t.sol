// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {console2} from "forge-std/console2.sol";
import {CLSVault} from "../../src/vault/CLSVault.sol";
import {UniswapV3Strategy} from "../../src/strategy/UniswapV3Strategy.sol";
import {UniswapV3StrategyHarness} from "../helpers/UniswapV3StrategyHarness.sol";
import {ISwapRouter} from "../../src/interfaces/ISwapRouter.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/interfaces/IUniswapV3Pool.sol";

contract CLSVaultForkTest is Test {
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

    CLSVault internal vault;
    UniswapV3Strategy internal strategy;
    UniswapV3StrategyHarness internal harness;

    address internal owner = makeAddr("owner");
    address internal user = makeAddr("user");

    function setUp() public {
        vm.createSelectFork(vm.envString("MAINNET_RPC"), 19_000_000);

        strategy = new UniswapV3Strategy(POOL, NPM, SWAP_ROUTER, HALF_RANGE, MIN_SWAP_USDC, MIN_SWAP_WETH);
        vault = new CLSVault(address(strategy), USDC, WETH);
        strategy.initialize(address(vault), owner);

        harness = new UniswapV3StrategyHarness(POOL, NPM, SWAP_ROUTER, HALF_RANGE, MIN_SWAP_USDC, MIN_SWAP_WETH);
    }

    function _depositAsUser() internal returns (uint256 shares) {
        deal(WETH, user, WETH_DEPOSIT);
        deal(USDC, user, USDC_DEPOSIT);

        vm.startPrank(user);
        IERC20(WETH).approve(address(vault), type(uint256).max);
        IERC20(USDC).approve(address(vault), type(uint256).max);
        vault.deposit(USDC_DEPOSIT, WETH_DEPOSIT);
        shares = vault.balanceOf(user);
        vm.stopPrank();
    }

    function _currentTick() internal view returns (int24 tick) {
        (, tick,,,,,) = IUniswapV3Pool(POOL).slot0();
    }

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

    function testDepositMintsSharesAndCreatesPosition() public {
        uint256 shares = _depositAsUser();

        assertGt(shares, 0);
        assertEq(vault.balanceOf(user), shares);

        (,, uint128 liquidity, int24 tickLower, int24 tickUpper,,) = strategy.getPosition();
        assertGt(liquidity, 0);

        int24 tick = _currentTick();
        (int24 expectedLower, int24 expectedUpper) = harness.exposedComputeTickRange(tick);
        assertEq(tickLower, expectedLower);
        assertEq(tickUpper, expectedUpper);
    }

    function testSecondDepositIncreasesSharesAndLiquidity() public {
        _depositAsUser();
        (,, uint128 liquidityBefore,,,,) = strategy.getPosition();
        uint256 sharesBefore = vault.balanceOf(user);

        _depositAsUser();

        assertGt(vault.balanceOf(user), sharesBefore);
        (,, uint128 liquidityAfter,,,,) = strategy.getPosition();
        assertGt(liquidityAfter, liquidityBefore);
    }

    function testPartialWithdrawReturnsTokens() public {
        uint256 shares = _depositAsUser();

        uint256 wethBefore = IERC20(WETH).balanceOf(user);
        uint256 usdcBefore = IERC20(USDC).balanceOf(user);

        vm.prank(user);
        (uint256 amount0, uint256 amount1) = vault.withdraw(shares / 2);

        assertGt(amount0 + amount1, 0);
        assertEq(IERC20(USDC).balanceOf(user), usdcBefore + amount0);
        assertEq(IERC20(WETH).balanceOf(user), wethBefore + amount1);
        assertEq(vault.balanceOf(user), shares / 2);
    }

    function testFullWithdrawBurnsAllShares() public {
        uint256 shares = _depositAsUser();

        vm.prank(user);
        vault.withdraw(shares);

        assertEq(vault.balanceOf(user), 0);
        assertEq(vault.totalSupply(), 0);
        (,, uint128 liquidity,,,,) = strategy.getPosition();
        assertEq(liquidity, 0);
    }

    function testDepositRebalancesWhenOutOfRange() public {
        _depositAsUser();

        (,,, int24 oldLower, int24 oldUpper,,) = strategy.getPosition();
        _pushTickAbove(oldUpper);
        int24 tickAfterSwap = _currentTick();
        assertGt(tickAfterSwap, oldUpper);

        deal(WETH, user, WETH_DEPOSIT);
        deal(USDC, user, USDC_DEPOSIT);
        vm.startPrank(user);
        IERC20(WETH).approve(address(vault), type(uint256).max);
        IERC20(USDC).approve(address(vault), type(uint256).max);
        vault.deposit(USDC_DEPOSIT, WETH_DEPOSIT);
        vm.stopPrank();

        (,, uint128 liquidity, int24 newLower, int24 newUpper,,) = strategy.getPosition();
        assertGt(liquidity, 0);
        assertFalse(newLower == oldLower && newUpper == oldUpper);

        (int24 expectedLower, int24 expectedUpper) = harness.exposedComputeTickRange(tickAfterSwap);
        assertEq(newLower, expectedLower);
        assertEq(newUpper, expectedUpper);
    }
}
