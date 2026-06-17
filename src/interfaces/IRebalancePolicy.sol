// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

interface IRebalancePolicy {
    function rebalance(uint256 amount0, uint256 amount1) external;
}