// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import {IStrategy} from "../interfaces/IStrategy.sol";
import {FullMath} from "../libraries/FullMath.sol";

contract CLSVault is ERC20 {
    using SafeERC20 for IERC20;

    error CLSVault__Amount0OrAmount1MustBeGreaterThanZero();
    error CLSVault__FirstSharesMustBeGreaterThanMinShares();
    error CLSVault__InvalidShares();

    uint256 private constant MIN_SHARES = 1e18; // first deposit should be at least 1e18

    address private immutable token0;
    address private immutable token1;
    address private immutable strategy;

    constructor(address _strategy, address _token0, address _token1) ERC20("CLSVault", "CLS") {
        strategy = _strategy;
        token0 = _token0;
        token1 = _token1;
    }

    function deposit(uint256 amount0, uint256 amount1) external {
        if (amount0 == 0 && amount1 == 0) {
            revert CLSVault__Amount0OrAmount1MustBeGreaterThanZero();
        }

        // deposit price
        uint256 depositValue = IStrategy(strategy).getTotalValue(amount0, amount1);

        // total value from strategy and pool
        uint256 totalValue = IStrategy(strategy).getTotalValue();

        uint256 shares;
        if (totalSupply() == 0) {
            shares = depositValue;
            if (shares < MIN_SHARES) {
                revert CLSVault__FirstSharesMustBeGreaterThanMinShares();
            }
        } else {
            shares = FullMath.mulDiv(depositValue, totalSupply(), totalValue);
        }

        IERC20(token0).safeTransferFrom(msg.sender, strategy, amount0);
        IERC20(token1).safeTransferFrom(msg.sender, strategy, amount1);

        _mint(msg.sender, shares);
        IStrategy(strategy).deposit();
    }

    /// @notice Redeem vault shares for a pro-rata share of token0 and token1.
    function withdraw(uint256 shares)
        external
        returns (uint256 amount0, uint256 amount1)
    {
        if (shares == 0 || shares > balanceOf(msg.sender)) {
            revert CLSVault__InvalidShares();
        }

        uint256 totalShares = totalSupply();
        _burn(msg.sender, shares);

        (amount0, amount1) = IStrategy(strategy).withdraw(shares, totalShares, msg.sender);
    }
}
