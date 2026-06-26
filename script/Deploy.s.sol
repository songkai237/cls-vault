// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {UniswapV3Strategy} from "../src/strategy/UniswapV3Strategy.sol";
import {CLSVault} from "../src/vault/CLSVault.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/interfaces/IUniswapV3Pool.sol";

/// @title DeployCLSVault
/// @notice Deploy UniswapV3Strategy + CLSVault and wire initialize(vault, owner).
///
/// Required env:
///   PRIVATE_KEY          — deployer (pays gas; becomes strategy owner if STRATEGY_OWNER unset)
///
/// Network defaults (by chain id, or override with DEPLOY_NETWORK=mainnet|sepolia):
///   Mainnet (1)  — WETH/USDC 0.05% pool
///   Sepolia (11155111) — USDC/WETH 0.3% pool
///
/// Optional env overrides:
///   POOL, NPM, SWAP_ROUTER, HALF_RANGE_TICKS, MIN_SWAP_AMOUNT0/1, STRATEGY_OWNER
///
/// Usage:
///   make deploy-dry          # mainnet simulate
///   make deploy              # mainnet broadcast
///   make deploy-sepolia-dry  # sepolia simulate
///   make deploy-sepolia      # sepolia broadcast
contract DeployCLSVault is Script {
    uint256 internal constant MAINNET_CHAIN_ID = 1;
    uint256 internal constant SEPOLIA_CHAIN_ID = 11155111;

    // Mainnet WETH/USDC 0.05%
    address internal constant MAINNET_POOL = 0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640;
    address internal constant MAINNET_NPM = 0xC36442b4a4522E871399CD717aBDD847Ab11FE88;
    address internal constant MAINNET_SWAP_ROUTER = 0xE592427A0AEce92De3Edee1F18E0157C05861564;

    // Sepolia USDC/WETH 0.3% (Circle USDC + canonical WETH)
    address internal constant SEPOLIA_POOL = 0x6Ce0896eAE6D4BD668fDe41BB784548fb8F59b50;
    address internal constant SEPOLIA_NPM = 0x1238536071E1c677A632429e3655c799b22cDA52;
    address internal constant SEPOLIA_SWAP_ROUTER = 0x3bFA4769FB09eefC5a80d6E87c3B9C650f7Ae48E;

    int24 internal constant DEFAULT_HALF_RANGE = 600;
    uint256 internal constant DEFAULT_MIN_SWAP0 = 100 * 1e6;
    uint256 internal constant DEFAULT_MIN_SWAP1 = 0.01 ether;

    address internal owner = 0x8fD7cA167763e78289B418f3f675B2a87cc769E7;

    struct NetworkDefaults {
        string name;
        address pool;
        address npm;
        address swapRouter;
    }

    function run() external {
        NetworkDefaults memory defaults = _resolveNetworkDefaults();
        address pool = _envAddress("POOL", defaults.pool);
        address npm = _envAddress("NPM", defaults.npm);
        address swapRouter = _envAddress("SWAP_ROUTER", defaults.swapRouter);
        int24 halfRange = int24(_envInt("HALF_RANGE_TICKS", DEFAULT_HALF_RANGE));
        uint256 minSwap0 = _envUint("MIN_SWAP_AMOUNT0", DEFAULT_MIN_SWAP0);
        uint256 minSwap1 = _envUint("MIN_SWAP_AMOUNT1", DEFAULT_MIN_SWAP1);

        vm.startBroadcast(owner);

        address token0 = IUniswapV3Pool(pool).token0();
        address token1 = IUniswapV3Pool(pool).token1();

        console2.log("=== CLSVault Deploy ===");
        console2.log("Network", defaults.name);
        console2.log("Chain ID", block.chainid);
        console2.log("Strategy owner", owner);
        console2.log("Pool", pool);
        console2.log("NPM", npm);
        console2.log("SwapRouter", swapRouter);
        console2.log("Half range ticks", uint256(int256(halfRange)));
        console2.log("Min swap amount0", minSwap0);
        console2.log("Min swap amount1", minSwap1);
        console2.log("Token0", token0);
        console2.log("Token1", token1);

        UniswapV3Strategy strategy = new UniswapV3Strategy(
            pool, npm, swapRouter, halfRange, minSwap0, minSwap1
        );
        CLSVault vault = new CLSVault(address(strategy), token0, token1);
        strategy.initialize(address(vault), owner);

        vm.stopBroadcast();

        console2.log("=== Deployed ===");
        console2.log("Strategy", address(strategy));
        console2.log("Vault", address(vault));
        console2.log("");
        console2.log("Keeper .env:");
        console2.log("  STRATEGY_ADDRESS=", address(strategy));
        console2.log("  PRIVATE_KEY=<owner key if owner != deployer>");
    }

    function _resolveNetworkDefaults() internal view returns (NetworkDefaults memory) {
        uint256 chainId = _resolveChainId();

        if (chainId == MAINNET_CHAIN_ID) {
            return NetworkDefaults({
                name: "mainnet",
                pool: MAINNET_POOL,
                npm: MAINNET_NPM,
                swapRouter: MAINNET_SWAP_ROUTER
            });
        }
        if (chainId == SEPOLIA_CHAIN_ID) {
            return NetworkDefaults({
                name: "sepolia",
                pool: SEPOLIA_POOL,
                npm: SEPOLIA_NPM,
                swapRouter: SEPOLIA_SWAP_ROUTER
            });
        }

        revert("Unsupported chain; set DEPLOY_NETWORK=mainnet|sepolia or POOL/NPM/SWAP_ROUTER");
    }

    function _resolveChainId() internal view returns (uint256) {
        try vm.envString("DEPLOY_NETWORK") returns (string memory network) {
            bytes32 h = keccak256(bytes(network));
            if (h == keccak256("mainnet")) return MAINNET_CHAIN_ID;
            if (h == keccak256("sepolia")) return SEPOLIA_CHAIN_ID;
            revert("DEPLOY_NETWORK must be mainnet or sepolia");
        } catch {}

        return block.chainid;
    }

    function _envAddress(string memory key, address defaultAddr) internal view returns (address) {
        try vm.envAddress(key) returns (address v) {
            return v;
        } catch {
            return defaultAddr;
        }
    }

    function _envInt(string memory key, int256 defaultVal) internal view returns (int256) {
        try vm.envInt(key) returns (int256 v) {
            return v;
        } catch {
            return defaultVal;
        }
    }

    function _envUint(string memory key, uint256 defaultVal) internal view returns (uint256) {
        try vm.envUint(key) returns (uint256 v) {
            return v;
        } catch {
            return defaultVal;
        }
    }
}
