// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {UniswapV3Strategy} from "../../src/strategy/UniswapV3Strategy.sol";

/// @dev Exposes internal helpers for unit testing.
contract UniswapV3StrategyHarness is UniswapV3Strategy {
    constructor(address pool, address npm, address swapRouter, int24 halfRangeTicks)
        UniswapV3Strategy(pool, npm, swapRouter, halfRangeTicks)
    {}

    function exposedComputeTickRange(int24 currentTick) external view returns (int24 tickLower, int24 tickUpper) {
        return _computeTickRange(currentTick);
    }

    function exposedFloor(int24 tick, int24 spacing) external pure returns (int24) {
        return _floor(tick, spacing);
    }

    function exposedNeedRebalance(int24 currentTick, int24 currentTickLower, int24 currentTickUpper)
        external
        view
        returns (bool)
    {
        return _needRebalance(currentTick, currentTickLower, currentTickUpper);
    }

    function setTokenID(uint256 id) external {
        tokenID = id;
    }
}
