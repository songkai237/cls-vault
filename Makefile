build:
	forge build

deploy-dry:
	forge script script/Deploy.s.sol:DeployCLSVault --rpc-url mainnet -vvvv

deploy:
	forge script script/Deploy.s.sol:DeployCLSVault --rpc-url mainnet --broadcast -vvvv

deploy-sepolia-dry:
	forge script script/Deploy.s.sol:DeployCLSVault --rpc-url sepolia --account test -vvvv

deploy-sepolia:
	forge script script/Deploy.s.sol:DeployCLSVault --rpc-url sepolia --account test --broadcast -vvvv

keeper-bindings:
	forge build
	mkdir -p keeper/internal/contracts
	jq '.abi' out/UniswapV3Strategy.sol/UniswapV3Strategy.json > keeper/internal/contracts/strategy.abi.json
	abigen --abi keeper/internal/contracts/strategy.abi.json --pkg contracts --type UniswapV3Strategy --out keeper/internal/contracts/strategy.go
	jq '.abi' lib/v3-core/out/UniswapV3Pool.sol/UniswapV3Pool.json > keeper/internal/contracts/pool.abi.json
	abigen --abi keeper/internal/contracts/pool.abi.json --pkg contracts --type UniswapV3Pool --out keeper/internal/contracts/pool.go

keeper-run:
	cd keeper && go run ./cmd/keeper

test-unit:
	forge test --match-path "test/unit/*"

test-fork:
	forge test --match-path "test/fork/*"