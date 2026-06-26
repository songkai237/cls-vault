# CLSVault Keeper

Off-chain keeper for `UniswapV3Strategy`. Polls pool tick and strategy position, decides whether to rebalance / swap / mint, and calls `maintain()` or `collect()` on-chain.

## Prerequisites

- Go 1.22+
- Strategy deployed and `initialize(vault, owner)` called
- Keeper wallet must be the strategy **owner**

## Setup

```bash
# Generate Go bindings after contract changes
make keeper-bindings

cd keeper
cp .env.example .env
# edit .env

go mod tidy
go run ./cmd/keeper
```

## Environment

| Variable | Description |
|----------|-------------|
| `RPC_URL` | Ethereum JSON-RPC URL |
| `PRIVATE_KEY` | Owner private key |
| `STRATEGY_ADDRESS` | UniswapV3Strategy address |
| `POLL_INTERVAL` | Seconds between checks (default 12) |
| `DRY_RUN` | If true, log decisions only |
| `MIN_IDLE_AMOUNT0/1` | Optional idle mint threshold |

## Decision logic

1. **`maintain()`** when:
   - `needRebalance()` is true (tick out of range or no position)
   - Idle excess exceeds min swap (closed-form, mirrors `_basicSwap`)
   - Idle tokens above threshold can be minted into position

2. **`collect()`** when:
   - Uncollected fees exist and no maintain action needed

`maintain()` runs the full on-chain flow: rebalance range → closed-form swap → mint/increase liquidity.

## Contract entrypoints used

- View: `needRebalance()`, `getPosition()`, `getPool()`, `getMinSwapAmount0/1()`, `getHalfRangeTicks()`
- Pool view: `slot0()`, `tickSpacing()`
- Tx: `maintain()`, `collect()`
