// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IStrategy} from "../interfaces/IStrategy.sol";

import {IUniswapV3Pool} from "@uniswap/v3-core/interfaces/IUniswapV3Pool.sol";
import {INonfungiblePositionManager} from "@uniswap/v3-periphery/interfaces/INonfungiblePositionManager.sol";
import {ISwapRouter} from "../interfaces/ISwapRouter.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import {LiquidityAmounts} from "../libraries/LiquidityAmounts.sol";
import {TickMath} from "../libraries/TickMath.sol";
import {FullMath} from "../libraries/FullMath.sol";

contract UniswapV3Strategy is IStrategy {
    using SafeERC20 for IERC20;

    error UniswapV3Strategy__AlreadyInitialized();
    error UniswapV3Strategy__OnlyVault();
    error UniswapV3Strategy__InvalidTickRange();
    error UniswapV3Strategy__OnlyOwner();

    address private immutable pool;
    address private immutable token0;
    address private immutable token1;
    uint24 private immutable fee;
    address private immutable swapRouter;
    int24 private immutable halfRangeTicks;
    uint256 private immutable minSwapAmount0;
    uint256 private immutable minSwapAmount1;

    address private vault;
    address private owner;
    bool private initialized;
    INonfungiblePositionManager private immutable npm;

    uint256 internal tokenID;

    constructor(
        address _pool,
        address _npm,
        address _swapRouter,
        int24 _halfRangeTicks,
        uint256 _minSwapAmount0,
        uint256 _minSwapAmount1
    ) {
        if (_halfRangeTicks <= 0) {
            revert UniswapV3Strategy__InvalidTickRange();
        }
        pool = _pool;
        npm = INonfungiblePositionManager(_npm);
        swapRouter = _swapRouter;
        halfRangeTicks = _halfRangeTicks;
        minSwapAmount0 = _minSwapAmount0;
        minSwapAmount1 = _minSwapAmount1;

        IUniswapV3Pool pool_ = IUniswapV3Pool(_pool);
        token0 = pool_.token0();
        token1 = pool_.token1();
        fee = pool_.fee();

        IERC20(token0).forceApprove(_npm, type(uint256).max);
        IERC20(token1).forceApprove(_npm, type(uint256).max);
        IERC20(token0).forceApprove(_swapRouter, type(uint256).max);
        IERC20(token1).forceApprove(_swapRouter, type(uint256).max);
    }

    modifier onlyVault() {
        if (msg.sender != vault) {
            revert UniswapV3Strategy__OnlyVault();
        }
        _;
    }

    modifier onlyOwner() {
        if (msg.sender != owner) {
            revert UniswapV3Strategy__OnlyOwner();
        }
        _;
    }

    function initialize(address _vault, address _owner) external {
        if (initialized) {
            revert UniswapV3Strategy__AlreadyInitialized();
        }
        vault = _vault;
        owner = _owner;
        initialized = true;
    }

    function deposit() external onlyVault {
        int24 tickLower;
        int24 tickUpper;
        (int24 currentTick, int24 currentTickLower, int24 currentTickUpper, uint128 liquidity) =
            _getRebalanceParams();

        if (_needRebalance(currentTick, currentTickLower, currentTickUpper)) {
            if (tokenID != 0 && liquidity > 0) {
                _withdrawLiquidity(liquidity);
            }
            if (tokenID != 0) {
                _collect();
                npm.burn(tokenID);
                tokenID = 0;
            }
            (tickLower, tickUpper) = _rebalance(currentTick);
        } else {
            tickLower = currentTickLower;
            tickUpper = currentTickUpper;
        }

        _basicSwap(tickLower, tickUpper);

        uint256 amount0 = IERC20(token0).balanceOf(address(this));
        uint256 amount1 = IERC20(token1).balanceOf(address(this));

        if (amount0 > 0 || amount1 > 0) {
            _depositLiquidity(tickLower, tickUpper, amount0, amount1, 0, 0, block.timestamp + 600);
        }
    }

    function withdraw(uint256 shares, uint256 totalShares, address recipient)
        external
        onlyVault
        returns (uint256 amount0, uint256 amount1)
    {
        if (shares == 0 || totalShares == 0) {
            return (0, 0);
        }

        if (shares == totalShares) {
            return _withdrawAll(recipient);
        }

        (uint256 total0, uint256 total1) = _getTotalAssets();
        amount0 = FullMath.mulDiv(total0, shares, totalShares);
        amount1 = FullMath.mulDiv(total1, shares, totalShares);

        uint256 idle0 = IERC20(token0).balanceOf(address(this));
        uint256 idle1 = IERC20(token1).balanceOf(address(this));

        uint128 liquidity;
        uint128 tokensOwed0;
        uint128 tokensOwed1;
        if (tokenID != 0) {
            (,, liquidity,,, tokensOwed0, tokensOwed1) = _getPosition();
        }

        uint256 available0 = idle0 + tokensOwed0;
        uint256 available1 = idle1 + tokensOwed1;
        bool needsLp = amount0 > available0 || amount1 > available1;

        if (needsLp && tokenID != 0 && liquidity > 0) {
            uint128 liquidityToWithdraw = uint128(FullMath.mulDiv(uint256(liquidity), shares, totalShares));
            if (liquidityToWithdraw > 0) {
                _withdrawLiquidity(liquidityToWithdraw);
            }
            _collect();
        } else if (tokenID != 0 && (tokensOwed0 > 0 || tokensOwed1 > 0)) {
            _collect();
        }

        uint256 balance0 = IERC20(token0).balanceOf(address(this));
        uint256 balance1 = IERC20(token1).balanceOf(address(this));
        if (amount0 > balance0) amount0 = balance0;
        if (amount1 > balance1) amount1 = balance1;

        IERC20(token0).safeTransfer(recipient, amount0);
        IERC20(token1).safeTransfer(recipient, amount1);
    }

    // function withdrawLiquidity(uint128 liquidity) external override onlyVault {
    //     _withdrawLiquidity(liquidity);
    // }

    function collect() external override onlyOwner {
        _collect();
    }

    function rebalance() external override onlyOwner returns (int24 tickLower, int24 tickUpper) {
        (tickLower, tickUpper) = _rebalance();
    }

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
        return _getPosition();
    }

    function getPositionValue() external view override returns (uint256 amount0, uint256 amount1) {
        return _getPositionValue();
    }

    function getTotalValue() external view returns (uint256 value) {
        (uint256 amountInPool0, uint256 amountInPool1) = _getPositionValue();
        uint256 amount0 = amountInPool0 + IERC20(token0).balanceOf(address(this));
        uint256 amount1 = amountInPool1 + IERC20(token1).balanceOf(address(this));
        value = _valueInToken1(amount0) + amount1;
    }

    function getTotalValue(uint256 amount0, uint256 amount1) external view returns (uint256 value) {
        value = _valueInToken1(amount0) + amount1;
    }

    function getTotalAssets() external view returns (uint256 amount0, uint256 amount1) {
        return _getTotalAssets();
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

    function _getPositionValue() internal view returns (uint256 amount0, uint256 amount1) {
        (uint256 principal0, uint256 principal1,,,, uint128 tokensOwed0, uint128 tokensOwed1) = _getPosition();
        return (principal0 + tokensOwed0, principal1 + tokensOwed1);
    }

    function _getTotalAssets() internal view returns (uint256 amount0, uint256 amount1) {
        (uint256 pool0, uint256 pool1) = _getPositionValue();
        return (pool0 + IERC20(token0).balanceOf(address(this)), pool1 + IERC20(token1).balanceOf(address(this)));
    }

    function _withdrawAll(address recipient) internal returns (uint256 amount0, uint256 amount1) {
        if (tokenID != 0) {
            (,, uint128 liquidity,,,,) = _getPosition();
            if (liquidity > 0) {
                _withdrawLiquidity(liquidity);
            }
            _collect();
            npm.burn(tokenID);
            tokenID = 0;
        }

        amount0 = IERC20(token0).balanceOf(address(this));
        amount1 = IERC20(token1).balanceOf(address(this));

        if (amount0 > 0) {
            IERC20(token0).safeTransfer(recipient, amount0);
        }
        if (amount1 > 0) {
            IERC20(token1).safeTransfer(recipient, amount1);
        }
    }

    function _needRebalance(int24 currentTick, int24 currentTickLower, int24 currentTickUpper)
        internal
        view
        returns (bool)
    {
        if (tokenID == 0) {
            return true;
        }
        if (currentTick < currentTickLower || currentTick > currentTickUpper) {
            return true;
        }
        return false;
    }

    function _valueInToken1(uint256 amount0) internal view returns (uint256) {
        (uint160 sqrtPriceX96,,,,,,) = IUniswapV3Pool(pool).slot0();
        uint256 v0 = FullMath.mulDiv(amount0, uint256(sqrtPriceX96), 1 << 96);
        v0 = FullMath.mulDiv(v0, uint256(sqrtPriceX96), 1 << 96);
        return v0;
    }

    function _depositLiquidity(
        int24 tickLower,
        int24 tickUpper,
        uint256 amount0Desired,
        uint256 amount1Desired,
        uint256 amount0Min,
        uint256 amount1Min,
        uint256 deadline
    ) internal {
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
            (uint256 _tokenID,,,) = npm.mint(params);
            tokenID = _tokenID;
        } else {
            INonfungiblePositionManager.IncreaseLiquidityParams memory params = INonfungiblePositionManager
                .IncreaseLiquidityParams({
                tokenId: tokenID,
                amount0Desired: amount0Desired,
                amount1Desired: amount1Desired,
                amount0Min: amount0Min,
                amount1Min: amount1Min,
                deadline: deadline
            });
            npm.increaseLiquidity(params);
        }
    }

    function _getPosition()
        internal
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

        (,,,,, tickLower, tickUpper, liquidity,,, tokensOwed0, tokensOwed1) = npm.positions(tokenID);

        (uint160 sqrtPriceX96,,,,,,) = IUniswapV3Pool(pool).slot0();

        (amount0, amount1) = LiquidityAmounts.getAmountsForLiquidity(
            sqrtPriceX96, TickMath.getSqrtRatioAtTick(tickLower), TickMath.getSqrtRatioAtTick(tickUpper), liquidity
        );
    }

    function _withdrawLiquidity(uint128 liquidity) internal {
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

    function _collect() internal {
        npm.collect(
            INonfungiblePositionManager.CollectParams({
                tokenId: tokenID,
                recipient: address(this),
                amount0Max: type(uint128).max,
                amount1Max: type(uint128).max
            })
        );
    }

    function _getRebalanceParams()
        internal
        view
        returns (int24 currentTick, int24 currentTickLower, int24 currentTickUpper, uint128 liquidity)
    {
        (,, liquidity, currentTickLower, currentTickUpper,,) = _getPosition();
        (, currentTick,,,,,) = IUniswapV3Pool(pool).slot0();
    }

    function _rebalance() internal returns (int24 tickLower, int24 tickUpper) {
        (int24 currentTick,,,) = _getRebalanceParams();
        return _rebalance(currentTick);
    }

    function _rebalance(int24 currentTick)
        internal
        returns (int24 tickLower, int24 tickUpper)
    {
        (tickLower, tickUpper) = _computeTickRange(currentTick);
    }

    function _computeTickRange(int24 currentTick) internal view returns (int24 tickLower, int24 tickUpper) {
        int24 spacing = IUniswapV3Pool(pool).tickSpacing();
        int24 center = _floor(currentTick, spacing);
        tickLower = _floor(center - halfRangeTicks, spacing);
        tickUpper = _floor(center + halfRangeTicks, spacing);

        if (tickLower >= tickUpper) {
            tickUpper = tickLower + spacing;
        }
    }

    function _floor(int24 tick, int24 spacing) internal pure returns (int24) {
        int24 compressed = tick / spacing;
        if (tick < 0 && tick % spacing != 0) {
            compressed--;
        }
        return compressed * spacing;
    }

    /// @dev Swap excess token inventory so mint uses both sides efficiently (closed-form at spot price).
    function _basicSwap(int24 tickLower, int24 tickUpper) internal {
        uint256 amount0 = IERC20(token0).balanceOf(address(this));
        uint256 amount1 = IERC20(token1).balanceOf(address(this));
        if (amount0 == 0 && amount1 == 0) {
            return;
        }

        (uint160 sqrtPriceX96,,,,,,) = IUniswapV3Pool(pool).slot0();
        uint160 sqrtLower = TickMath.getSqrtRatioAtTick(tickLower);
        uint160 sqrtUpper = TickMath.getSqrtRatioAtTick(tickUpper);

        uint128 liquidity =
            LiquidityAmounts.getLiquidityForAmounts(sqrtPriceX96, sqrtLower, sqrtUpper, amount0, amount1);
        if (liquidity == 0) {
            return;
        }

        (uint256 need0, uint256 need1) =
            LiquidityAmounts.getAmountsForLiquidity(sqrtPriceX96, sqrtLower, sqrtUpper, liquidity);

        if (amount0 > need0) {
            uint256 excess0 = amount0 - need0;
            if (excess0 >= minSwapAmount0) {
                _swapExactInput(token0, token1, excess0);
            }
        } else if (amount1 > need1) {
            uint256 excess1 = amount1 - need1;
            if (excess1 >= minSwapAmount1) {
                _swapExactInput(token1, token0, excess1);
            }
        }
    }

    function _swapExactInput(address tokenIn, address tokenOut, uint256 amountIn) internal {
        if (amountIn == 0) {
            return;
        }

        ISwapRouter(swapRouter).exactInputSingle(
            ISwapRouter.ExactInputSingleParams({
                tokenIn: tokenIn,
                tokenOut: tokenOut,
                fee: fee,
                recipient: address(this),
                deadline: block.timestamp,
                amountIn: amountIn,
                amountOutMinimum: 0,
                sqrtPriceLimitX96: 0
            })
        );
    }
}
