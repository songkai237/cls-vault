package contracts

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const erc20ABIJSON = `[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

type ERC20 struct {
	contract *bind.BoundContract
}

func NewERC20(token common.Address, client *ethclient.Client) (*ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		return nil, fmt.Errorf("parse erc20 abi: %w", err)
	}
	return &ERC20{
		contract: bind.NewBoundContract(token, parsed, client, client, client),
	}, nil
}

func (e *ERC20) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	if err := e.contract.Call(opts, &out, "balanceOf", account); err != nil {
		return nil, err
	}
	bal, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected balance type %T", out[0])
	}
	return bal, nil
}
