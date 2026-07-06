# CLSVault Keeper

Off-chain keeper for `UniswapV3Strategy`. Polls pool tick and strategy position, quotes swaps via Uniswap V3 Quoter V2 (`eth_call`), decides whether to rebalance / swap / mint, and calls `maintain()` or `collect()` on-chain.

## Prerequisites

- Go 1.22+
- Strategy deployed and `initialize(vault, owner)` called
- Keeper wallet must be the strategy **owner**
- Uniswap V3 Quoter V2 deployed on the target chain

## Setup

```bash
cd keeper

go mod tidy
go run ./cmd/keeper
```

## Environment

| Variable | Description                                        |
|----------|----------------------------------------------------|
| `RPC_URL` | Ethereum JSON-RPC URL                              |
| `PRIVATE_KEY` | Owner private key (only use this on test env)      |
| `STRATEGY_ADDRESS` | UniswapV3Strategy address                          |
| `QUOTER_ADDRESS` | Uniswap V3 Quoter V2 address (**required**)        |
| `POLL_INTERVAL` | Seconds between checks (default 12)                |
| `CHAIN_ID` | Chain ID (auto-detect if blank)                    |
| `DRY_RUN` | If `1` or `true`, log decisions without sending tx |
| `MIN_IDLE_AMOUNT0/1` | Optional idle mint threshold override              |

### Quoter V2 addresses

| Chain | Address |
|-------|---------|
| Ethereum Mainnet | `0x61fFE014bA17989E743c5F6cB21bF9697530B21e` |
| Sepolia | `0x82825d0554fA07f7FC52Ab63c961F330fdEFa8E8` |

## Decision logic

1. **`maintain()`** when:
   - `needRebalance()` is true (tick out of range or no position)
   - Binary search finds a beneficial swap to rebalance idle tokens for LP efficiency
   - Idle tokens above threshold can be minted into position

2. **`collect()`** when:
   - Uncollected fees exist and no maintain action needed

### Swap quoting

Every swap estimate goes through the Quoter V2 `quoteExactInputSingle` via `eth_call`. No local swap math — the Quoter simulates the actual pool including tick-crossing and price impact.

The binary search in `strategy.FindBestSwap`:

1. Determines which token is in excess (higher per-unit liquidity contribution)
2. Quotes the full excess amount — if price crosses `sqrtLower`/`sqrtUpper`, clamps the search upper bound to the maximum amountIn that stays within range
3. Binary-searches on `[0, hi]` using the Quoter's `sqrtPriceX96After` to compute post-swap liquidity
4. Returns the `SwapSuggestion` only if `amountIn >= minSwap`

### Price-bound clamping

When the full excess swap would push the price outside `[sqrtLower, sqrtUpper]`:
- `zeroForOne`: swap is clamped so `sqrtPriceAfter >= sqrtLower`
- `oneForZero`: swap is clamped so `sqrtPriceAfter <= sqrtUpper`

This avoids wasting iterations on infeasible amounts and ensures the LP range is respected.

## Contract entrypoints used

| Contract | Methods |
|----------|---------|
| UniswapV3Strategy | `needRebalance()`, `getPosition()`, `getPool()`, `getMinSwapAmount0/1()`, `getHalfRangeTicks()`, `getFee()` → view |
| | `maintain(zeroForOne, amountIn, tickLower, tickUpper)`, `collect()` → tx |
| UniswapV3Pool | `slot0()`, `tickSpacing()` → view |
| ERC20 (token0/1) | `balanceOf(strategyAddr)` → view |
| IQuoterV2 | `quoteExactInputSingle(params)` → `eth_call` |

## Package structure

```
keeper/
├── cmd/keeper/         # Entrypoint
├── internal/
│   ├── config/         # Env-based config loading
│   ├── contracts/      # abigen bindings (strategy, pool, erc20, quoter)
│   ├── keeper/         # Keeper: chain reads, Evaluate, tx sending
│   ├── math/           # Tick math, liquidity math (no swap math)
│   └── strategy/       # FindBestSwap: binary search with QuoteFunc
└── .env.example
```

## Testing

```bash
# Unit tests (no RPC required — uses spot-price QuoteFunc mock)
go test ./internal/...

# Quoter V2 integration test (requires RPC)
export RPC_URL="https://..."
export QUOTER_ADDRESS="0x61fFE014bA17989E743c5F6cB21bF9697530B21e"
export TOKEN0="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
export TOKEN1="0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
# FEE defaults to 3000; optional AMOUNT_IN for custom amount
go test ./internal/keeper/ -run TestQuoter -v
```
