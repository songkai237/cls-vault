// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

interface IStrategy {
    function deposit() external;
    function withdraw(uint256 shares, uint256 totalShares, address recipient)
        external
        returns (uint256 amount0, uint256 amount1);
    function collect() external;
    function getPositionValue() external view returns (uint256 amount0, uint256 amount1);
    function getTotalAssets() external view returns (uint256 amount0, uint256 amount1);
    function rebalance() external returns (int24 tickLower, int24 tickUpper);
    function getTotalValue() external view returns (uint256 value);
    function getTotalValue(uint256 amount0, uint256 amount1) external view returns (uint256 value);
}
