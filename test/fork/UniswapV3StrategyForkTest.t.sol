// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {console2} from "forge-std/console2.sol";
import {UniswapV3Strategy} from "../../src/strategy/UniswapV3Strategy.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/interfaces/IUniswapV3Pool.sol";

contract UniswapV3StrategyForkTest is Test {
    address private constant POOL = 0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640;
    address private constant NPM = 0xC36442b4a4522E871399CD717aBDD847Ab11FE88;
    address private constant WETH = 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2;
    address private constant USDC = 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48;

    uint256 private constant WETHAmount = 10 * 1e18;
    uint256 private constant USDCAmount = 30000 * 1e6;
    int24 private tick;
    int24 private tickLower;
    int24 private tickUpper;


    UniswapV3Strategy strategy;

    address public user = makeAddr("user");

    function setUp() public {
        vm.createSelectFork(vm.envString("MAINNET_RPC"), 19_000_000); 
        strategy = new UniswapV3Strategy(POOL, NPM);

        (uint160 sqrtPriceX96, int24 _tick, , , , , ) = IUniswapV3Pool(POOL).slot0();
        tick = _tick;

        console2.log("current tick", tick);
    }

    function base() public {
        deal(WETH, address(strategy), WETHAmount);
        deal(USDC, address(strategy), USDCAmount);

        vm.startPrank(address(strategy));
        IERC20(WETH).approve(NPM, WETHAmount);
        IERC20(USDC).approve(NPM, USDCAmount);
        vm.stopPrank();

        tickLower = tick - 18005;
        tickUpper = tick + 19005;    

        strategy.depositLiquidity(
            tickLower, 
            tickUpper, 
            10 * 1e18, 
            30000 * 1e6, 
            0, 
            0, 
            block.timestamp + 1000
        );
    }

    function testMintAndGetPosition() public {
        base();
        (
            uint256 amount0,
            uint256 amount1,
            uint128 liquidity,
            int24 tLower,
            int24 tUpper,
            ,
        ) = strategy.getPosition();
        assertGt(liquidity, 0);
        assertEq(tLower, tickLower);
        assertEq(tUpper, tickUpper);
        assertGt(amount0 + amount1, 0);

        console2.log("amount0", amount0);
        console2.log("amount1", amount1);
        console2.log("liquidity", liquidity);
        console2.log("tLower", tLower);
        console2.log("tUpper", tUpper);
    }

    function testWithdrawAndCollectLiquidity() public {
        base();
        (
            uint256 amount0,
            uint256 amount1,
            uint128 liquidity,
            int24 tLower,
            int24 tUpper,
            ,
        ) = strategy.getPosition();
        strategy.withdrawLiquidity(liquidity);
        (
            amount0,
            amount1,
            liquidity,
            tLower,
            tUpper,
            ,
        ) = strategy.getPosition();
        assertEq(amount0, 0);
        assertEq(amount1, 0);
        assertEq(liquidity, 0);

        strategy.collectFees();
        assertApproxEqAbs(IERC20(WETH).balanceOf(address(strategy)), WETHAmount, 1e15);
        assertApproxEqAbs(IERC20(USDC).balanceOf(address(strategy)), USDCAmount, 1e4);
    }
}