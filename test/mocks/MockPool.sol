// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

/// @dev Minimal Uniswap V3 pool stub for unit tests.
contract MockPool {
    address public immutable token0;
    address public immutable token1;
    uint24 public immutable fee;
    int24 public immutable tickSpacing;
    int24 public currentTick;

    constructor(address _token0, address _token1, uint24 _fee, int24 _tickSpacing, int24 _currentTick) {
        token0 = _token0;
        token1 = _token1;
        fee = _fee;
        tickSpacing = _tickSpacing;
        currentTick = _currentTick;
    }

    function setCurrentTick(int24 tick) external {
        currentTick = tick;
    }

    function slot0()
        external
        view
        returns (
            uint160 sqrtPriceX96,
            int24 tick,
            uint16 observationIndex,
            uint16 observationCardinality,
            uint16 observationCardinalityNext,
            uint8 feeProtocol,
            bool unlocked
        )
    {
        sqrtPriceX96 = 0;
        tick = currentTick;
        observationIndex = 0;
        observationCardinality = 0;
        observationCardinalityNext = 0;
        feeProtocol = 0;
        unlocked = true;
    }
}
