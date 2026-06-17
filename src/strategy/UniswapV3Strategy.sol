// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IStragtegy} from "../interfaces/IStragtegy.sol";

import {IUniswapV3Pool} from "@uniswap/v3-core/interfaces/IUniswapV3Pool.sol";
import {INonfungiblePositionManager} from "@uniswap/v3-periphery/interfaces/INonfungiblePositionManager.sol";
import {LiquidityAmounts} from "../libraries/LiquidityAmounts.sol";
import {TickMath} from "../libraries/TickMath.sol";


contract UniswapV3Strategy is IStragtegy {
    address private immutable pool;
    address private immutable token0;
    address private immutable token1;
    uint24 private immutable fee;
    INonfungiblePositionManager private immutable npm;

    uint256 private tokenID; // this version only support one position

    constructor(address _pool, address _npm) {
        pool = _pool;
        npm = INonfungiblePositionManager(_npm);
        IUniswapV3Pool pool_ = IUniswapV3Pool(_pool);
        token0 = pool_.token0();
        token1 = pool_.token1();
        fee = pool_.fee();
    }

    function depositLiquidity(
        int24 tickLower,
        int24 tickUpper,
         uint256 amount0Desired,
        uint256 amount1Desired,
        uint256 amount0Min,
        uint256 amount1Min,
        uint256 deadline
    ) external {
        if (tokenID == 0) {
            INonfungiblePositionManager.MintParams memory params = INonfungiblePositionManager.MintParams({
                token0: token0,
                token1: token1,
                fee: fee,
                tickLower: tickLower,
                tickUpper: tickUpper,
                amount0Desired: amount0Desired,
                amount1Desired: amount1Desired,
                amount0Min: amount0Min,
                amount1Min: amount1Min,
                recipient: address(this),
                deadline: deadline
            });
            (uint256 _tokenID, uint128 _liquidity, uint256 _amount0, uint256 _amount1) = npm.mint(params);
            tokenID = _tokenID;
        } else {
            INonfungiblePositionManager.IncreaseLiquidityParams memory params = INonfungiblePositionManager.IncreaseLiquidityParams({
                tokenId: tokenID,
                amount0Desired: amount0Desired,
                amount1Desired: amount1Desired,
                amount0Min: amount0Min,
                amount1Min: amount1Min,
                deadline: deadline
            });
            (uint128 _liquidity, uint256 _amount0, uint256 _amount1) = npm.increaseLiquidity(params);
        }
    }

    function withdrawLiquidity(uint128 liquidity) external override {
        npm.decreaseLiquidity(
            INonfungiblePositionManager.DecreaseLiquidityParams({
                tokenId: tokenID,
                liquidity: liquidity,
                amount0Min: 0,
                amount1Min: 0,
                deadline: block.timestamp
            })
        );
    }

    function collectFees() external override {
        npm.collect(
            INonfungiblePositionManager.CollectParams({
                tokenId: tokenID,
                recipient: address(this),
                amount0Max: type(uint128).max,
                amount1Max: type(uint128).max
            })
        );
    }
    
    function rebalance(uint256, uint256) external pure override {
        revert("UniswapV3Strategy: rebalance not implemented");
    }

    /// view functions ///
    function getPosition()
        external
        view
        returns (
            uint256 amount0,
            uint256 amount1,
            uint128 liquidity,
            int24 tickLower,
            int24 tickUpper,
            uint128 tokensOwed0,
            uint128 tokensOwed1
        )
    {
        if (tokenID == 0) {
            return (0, 0, 0, 0, 0, 0, 0);
        }

        (, , , , , tickLower, tickUpper, liquidity, , , tokensOwed0, tokensOwed1) = npm.positions(tokenID);

        (uint160 sqrtPriceX96,,,,,,) = IUniswapV3Pool(pool).slot0();

        (amount0, amount1) = LiquidityAmounts.getAmountsForLiquidity(
            sqrtPriceX96,
            TickMath.getSqrtRatioAtTick(tickLower),
            TickMath.getSqrtRatioAtTick(tickUpper),
            liquidity
        );
    }

    function getPositionValue() external view override returns (uint256 amount0, uint256 amount1) {
        (
            uint256 principal0,
            uint256 principal1,
            ,
            ,
            ,
            uint128 tokensOwed0,
            uint128 tokensOwed1
        ) = this.getPosition();
        return (principal0 + tokensOwed0, principal1 + tokensOwed1);
    }

    function getPool() external view returns (address) {
        return pool;
    }

    function getToken0() external view returns (address) {
        return token0;
    }

    function getToken1() external view returns (address) {
        return token1;
    }

    function getFee() external view returns (uint24) {
        return fee;
    }

}