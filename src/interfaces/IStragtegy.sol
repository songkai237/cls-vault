// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

interface IStragtegy {
    function depositLiquidity(
        int24 tickLower,
        int24 tickUpper,
        uint256 amount0Desired,
        uint256 amount1Desired,
        uint256 amount0Min,
        uint256 amount1Min,
        uint256 deadline
    ) external;
    function withdrawLiquidity(uint128 liquidity) external;
    function collectFees() external;
    function getPositionValue() external view returns (uint256 amount0, uint256 amount1);
    function rebalance(uint256 amount0, uint256 amount1) external;
}