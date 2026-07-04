// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// UniswapV3StrategyMetaData contains all meta data concerning the UniswapV3Strategy contract.
var UniswapV3StrategyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_npm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_swapRouter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_halfRangeTicks\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"_minSwapAmount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minSwapAmount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collect\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHalfRangeTicks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinSwapAmount0\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinSwapAmount1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPosition\",\"inputs\":[],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"liquidity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tokensOwed0\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"tokensOwed1\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPositionValue\",\"inputs\":[],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken0\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalValue\",\"inputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalValue\",\"inputs\":[],\"outputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maintain\",\"inputs\":[{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"needRebalance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rebalance\",\"inputs\":[],\"outputs\":[{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"UniswapV3Strategy__AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UniswapV3Strategy__InvalidTickRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UniswapV3Strategy__OnlyOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UniswapV3Strategy__OnlyVault\",\"inputs\":[]}]",
}

// UniswapV3StrategyABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV3StrategyMetaData.ABI instead.
var UniswapV3StrategyABI = UniswapV3StrategyMetaData.ABI

// UniswapV3Strategy is an auto generated Go binding around an Ethereum contract.
type UniswapV3Strategy struct {
	UniswapV3StrategyCaller     // Read-only binding to the contract
	UniswapV3StrategyTransactor // Write-only binding to the contract
	UniswapV3StrategyFilterer   // Log filterer for contract events
}

// UniswapV3StrategyCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV3StrategyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3StrategyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV3StrategyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3StrategyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV3StrategyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3StrategySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV3StrategySession struct {
	Contract     *UniswapV3Strategy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UniswapV3StrategyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV3StrategyCallerSession struct {
	Contract *UniswapV3StrategyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// UniswapV3StrategyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV3StrategyTransactorSession struct {
	Contract     *UniswapV3StrategyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// UniswapV3StrategyRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV3StrategyRaw struct {
	Contract *UniswapV3Strategy // Generic contract binding to access the raw methods on
}

// UniswapV3StrategyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV3StrategyCallerRaw struct {
	Contract *UniswapV3StrategyCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV3StrategyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV3StrategyTransactorRaw struct {
	Contract *UniswapV3StrategyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV3Strategy creates a new instance of UniswapV3Strategy, bound to a specific deployed contract.
func NewUniswapV3Strategy(address common.Address, backend bind.ContractBackend) (*UniswapV3Strategy, error) {
	contract, err := bindUniswapV3Strategy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV3Strategy{UniswapV3StrategyCaller: UniswapV3StrategyCaller{contract: contract}, UniswapV3StrategyTransactor: UniswapV3StrategyTransactor{contract: contract}, UniswapV3StrategyFilterer: UniswapV3StrategyFilterer{contract: contract}}, nil
}

// NewUniswapV3StrategyCaller creates a new read-only instance of UniswapV3Strategy, bound to a specific deployed contract.
func NewUniswapV3StrategyCaller(address common.Address, caller bind.ContractCaller) (*UniswapV3StrategyCaller, error) {
	contract, err := bindUniswapV3Strategy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StrategyCaller{contract: contract}, nil
}

// NewUniswapV3StrategyTransactor creates a new write-only instance of UniswapV3Strategy, bound to a specific deployed contract.
func NewUniswapV3StrategyTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV3StrategyTransactor, error) {
	contract, err := bindUniswapV3Strategy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StrategyTransactor{contract: contract}, nil
}

// NewUniswapV3StrategyFilterer creates a new log filterer instance of UniswapV3Strategy, bound to a specific deployed contract.
func NewUniswapV3StrategyFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV3StrategyFilterer, error) {
	contract, err := bindUniswapV3Strategy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StrategyFilterer{contract: contract}, nil
}

// bindUniswapV3Strategy binds a generic wrapper to an already deployed contract.
func bindUniswapV3Strategy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV3StrategyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3Strategy *UniswapV3StrategyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3Strategy.Contract.UniswapV3StrategyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3Strategy *UniswapV3StrategyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.UniswapV3StrategyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3Strategy *UniswapV3StrategyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.UniswapV3StrategyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3Strategy *UniswapV3StrategyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3Strategy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3Strategy *UniswapV3StrategyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3Strategy *UniswapV3StrategyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.contract.Transact(opts, method, params...)
}

// GetFee is a free data retrieval call binding the contract method 0xced72f87.
//
// Solidity: function getFee() view returns(uint24)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0xced72f87.
//
// Solidity: function getFee() view returns(uint24)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetFee() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetFee(&_UniswapV3Strategy.CallOpts)
}

// GetFee is a free data retrieval call binding the contract method 0xced72f87.
//
// Solidity: function getFee() view returns(uint24)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetFee() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetFee(&_UniswapV3Strategy.CallOpts)
}

// GetHalfRangeTicks is a free data retrieval call binding the contract method 0x18e8b986.
//
// Solidity: function getHalfRangeTicks() view returns(int24)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetHalfRangeTicks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getHalfRangeTicks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetHalfRangeTicks is a free data retrieval call binding the contract method 0x18e8b986.
//
// Solidity: function getHalfRangeTicks() view returns(int24)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetHalfRangeTicks() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetHalfRangeTicks(&_UniswapV3Strategy.CallOpts)
}

// GetHalfRangeTicks is a free data retrieval call binding the contract method 0x18e8b986.
//
// Solidity: function getHalfRangeTicks() view returns(int24)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetHalfRangeTicks() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetHalfRangeTicks(&_UniswapV3Strategy.CallOpts)
}

// GetMinSwapAmount0 is a free data retrieval call binding the contract method 0xdfa68f3e.
//
// Solidity: function getMinSwapAmount0() view returns(uint256)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetMinSwapAmount0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getMinSwapAmount0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinSwapAmount0 is a free data retrieval call binding the contract method 0xdfa68f3e.
//
// Solidity: function getMinSwapAmount0() view returns(uint256)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetMinSwapAmount0() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetMinSwapAmount0(&_UniswapV3Strategy.CallOpts)
}

// GetMinSwapAmount0 is a free data retrieval call binding the contract method 0xdfa68f3e.
//
// Solidity: function getMinSwapAmount0() view returns(uint256)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetMinSwapAmount0() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetMinSwapAmount0(&_UniswapV3Strategy.CallOpts)
}

// GetMinSwapAmount1 is a free data retrieval call binding the contract method 0x700f1acb.
//
// Solidity: function getMinSwapAmount1() view returns(uint256)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetMinSwapAmount1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getMinSwapAmount1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinSwapAmount1 is a free data retrieval call binding the contract method 0x700f1acb.
//
// Solidity: function getMinSwapAmount1() view returns(uint256)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetMinSwapAmount1() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetMinSwapAmount1(&_UniswapV3Strategy.CallOpts)
}

// GetMinSwapAmount1 is a free data retrieval call binding the contract method 0x700f1acb.
//
// Solidity: function getMinSwapAmount1() view returns(uint256)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetMinSwapAmount1() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetMinSwapAmount1(&_UniswapV3Strategy.CallOpts)
}

// GetPool is a free data retrieval call binding the contract method 0x026b1d5f.
//
// Solidity: function getPool() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPool is a free data retrieval call binding the contract method 0x026b1d5f.
//
// Solidity: function getPool() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetPool() (common.Address, error) {
	return _UniswapV3Strategy.Contract.GetPool(&_UniswapV3Strategy.CallOpts)
}

// GetPool is a free data retrieval call binding the contract method 0x026b1d5f.
//
// Solidity: function getPool() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetPool() (common.Address, error) {
	return _UniswapV3Strategy.Contract.GetPool(&_UniswapV3Strategy.CallOpts)
}

// GetPosition is a free data retrieval call binding the contract method 0x7398ab18.
//
// Solidity: function getPosition() view returns(uint256 amount0, uint256 amount1, uint128 liquidity, int24 tickLower, int24 tickUpper, uint128 tokensOwed0, uint128 tokensOwed1)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetPosition(opts *bind.CallOpts) (struct {
	Amount0     *big.Int
	Amount1     *big.Int
	Liquidity   *big.Int
	TickLower   *big.Int
	TickUpper   *big.Int
	TokensOwed0 *big.Int
	TokensOwed1 *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getPosition")

	outstruct := new(struct {
		Amount0     *big.Int
		Amount1     *big.Int
		Liquidity   *big.Int
		TickLower   *big.Int
		TickUpper   *big.Int
		TokensOwed0 *big.Int
		TokensOwed1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Amount1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Liquidity = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TickLower = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TickUpper = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TokensOwed0 = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.TokensOwed1 = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPosition is a free data retrieval call binding the contract method 0x7398ab18.
//
// Solidity: function getPosition() view returns(uint256 amount0, uint256 amount1, uint128 liquidity, int24 tickLower, int24 tickUpper, uint128 tokensOwed0, uint128 tokensOwed1)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetPosition() (struct {
	Amount0     *big.Int
	Amount1     *big.Int
	Liquidity   *big.Int
	TickLower   *big.Int
	TickUpper   *big.Int
	TokensOwed0 *big.Int
	TokensOwed1 *big.Int
}, error) {
	return _UniswapV3Strategy.Contract.GetPosition(&_UniswapV3Strategy.CallOpts)
}

// GetPosition is a free data retrieval call binding the contract method 0x7398ab18.
//
// Solidity: function getPosition() view returns(uint256 amount0, uint256 amount1, uint128 liquidity, int24 tickLower, int24 tickUpper, uint128 tokensOwed0, uint128 tokensOwed1)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetPosition() (struct {
	Amount0     *big.Int
	Amount1     *big.Int
	Liquidity   *big.Int
	TickLower   *big.Int
	TickUpper   *big.Int
	TokensOwed0 *big.Int
	TokensOwed1 *big.Int
}, error) {
	return _UniswapV3Strategy.Contract.GetPosition(&_UniswapV3Strategy.CallOpts)
}

// GetPositionValue is a free data retrieval call binding the contract method 0xb44b3a8f.
//
// Solidity: function getPositionValue() view returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetPositionValue(opts *bind.CallOpts) (struct {
	Amount0 *big.Int
	Amount1 *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getPositionValue")

	outstruct := new(struct {
		Amount0 *big.Int
		Amount1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Amount1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPositionValue is a free data retrieval call binding the contract method 0xb44b3a8f.
//
// Solidity: function getPositionValue() view returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetPositionValue() (struct {
	Amount0 *big.Int
	Amount1 *big.Int
}, error) {
	return _UniswapV3Strategy.Contract.GetPositionValue(&_UniswapV3Strategy.CallOpts)
}

// GetPositionValue is a free data retrieval call binding the contract method 0xb44b3a8f.
//
// Solidity: function getPositionValue() view returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetPositionValue() (struct {
	Amount0 *big.Int
	Amount1 *big.Int
}, error) {
	return _UniswapV3Strategy.Contract.GetPositionValue(&_UniswapV3Strategy.CallOpts)
}

// GetToken0 is a free data retrieval call binding the contract method 0xba94a315.
//
// Solidity: function getToken0() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetToken0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getToken0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetToken0 is a free data retrieval call binding the contract method 0xba94a315.
//
// Solidity: function getToken0() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetToken0() (common.Address, error) {
	return _UniswapV3Strategy.Contract.GetToken0(&_UniswapV3Strategy.CallOpts)
}

// GetToken0 is a free data retrieval call binding the contract method 0xba94a315.
//
// Solidity: function getToken0() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetToken0() (common.Address, error) {
	return _UniswapV3Strategy.Contract.GetToken0(&_UniswapV3Strategy.CallOpts)
}

// GetToken1 is a free data retrieval call binding the contract method 0x6f26a710.
//
// Solidity: function getToken1() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetToken1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getToken1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetToken1 is a free data retrieval call binding the contract method 0x6f26a710.
//
// Solidity: function getToken1() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetToken1() (common.Address, error) {
	return _UniswapV3Strategy.Contract.GetToken1(&_UniswapV3Strategy.CallOpts)
}

// GetToken1 is a free data retrieval call binding the contract method 0x6f26a710.
//
// Solidity: function getToken1() view returns(address)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetToken1() (common.Address, error) {
	return _UniswapV3Strategy.Contract.GetToken1(&_UniswapV3Strategy.CallOpts)
}

// GetTotalAssets is a free data retrieval call binding the contract method 0x6e07302b.
//
// Solidity: function getTotalAssets() view returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetTotalAssets(opts *bind.CallOpts) (struct {
	Amount0 *big.Int
	Amount1 *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getTotalAssets")

	outstruct := new(struct {
		Amount0 *big.Int
		Amount1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Amount1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTotalAssets is a free data retrieval call binding the contract method 0x6e07302b.
//
// Solidity: function getTotalAssets() view returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetTotalAssets() (struct {
	Amount0 *big.Int
	Amount1 *big.Int
}, error) {
	return _UniswapV3Strategy.Contract.GetTotalAssets(&_UniswapV3Strategy.CallOpts)
}

// GetTotalAssets is a free data retrieval call binding the contract method 0x6e07302b.
//
// Solidity: function getTotalAssets() view returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetTotalAssets() (struct {
	Amount0 *big.Int
	Amount1 *big.Int
}, error) {
	return _UniswapV3Strategy.Contract.GetTotalAssets(&_UniswapV3Strategy.CallOpts)
}

// GetTotalValue is a free data retrieval call binding the contract method 0x09a045b1.
//
// Solidity: function getTotalValue(uint256 amount0, uint256 amount1) view returns(uint256 value)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetTotalValue(opts *bind.CallOpts, amount0 *big.Int, amount1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getTotalValue", amount0, amount1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalValue is a free data retrieval call binding the contract method 0x09a045b1.
//
// Solidity: function getTotalValue(uint256 amount0, uint256 amount1) view returns(uint256 value)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetTotalValue(amount0 *big.Int, amount1 *big.Int) (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetTotalValue(&_UniswapV3Strategy.CallOpts, amount0, amount1)
}

// GetTotalValue is a free data retrieval call binding the contract method 0x09a045b1.
//
// Solidity: function getTotalValue(uint256 amount0, uint256 amount1) view returns(uint256 value)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetTotalValue(amount0 *big.Int, amount1 *big.Int) (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetTotalValue(&_UniswapV3Strategy.CallOpts, amount0, amount1)
}

// GetTotalValue0 is a free data retrieval call binding the contract method 0xcaa648b4.
//
// Solidity: function getTotalValue() view returns(uint256 value)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) GetTotalValue0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "getTotalValue0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalValue0 is a free data retrieval call binding the contract method 0xcaa648b4.
//
// Solidity: function getTotalValue() view returns(uint256 value)
func (_UniswapV3Strategy *UniswapV3StrategySession) GetTotalValue0() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetTotalValue0(&_UniswapV3Strategy.CallOpts)
}

// GetTotalValue0 is a free data retrieval call binding the contract method 0xcaa648b4.
//
// Solidity: function getTotalValue() view returns(uint256 value)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) GetTotalValue0() (*big.Int, error) {
	return _UniswapV3Strategy.Contract.GetTotalValue0(&_UniswapV3Strategy.CallOpts)
}

// NeedRebalance is a free data retrieval call binding the contract method 0xffb86c6a.
//
// Solidity: function needRebalance() view returns(bool)
func (_UniswapV3Strategy *UniswapV3StrategyCaller) NeedRebalance(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UniswapV3Strategy.contract.Call(opts, &out, "needRebalance")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// NeedRebalance is a free data retrieval call binding the contract method 0xffb86c6a.
//
// Solidity: function needRebalance() view returns(bool)
func (_UniswapV3Strategy *UniswapV3StrategySession) NeedRebalance() (bool, error) {
	return _UniswapV3Strategy.Contract.NeedRebalance(&_UniswapV3Strategy.CallOpts)
}

// NeedRebalance is a free data retrieval call binding the contract method 0xffb86c6a.
//
// Solidity: function needRebalance() view returns(bool)
func (_UniswapV3Strategy *UniswapV3StrategyCallerSession) NeedRebalance() (bool, error) {
	return _UniswapV3Strategy.Contract.NeedRebalance(&_UniswapV3Strategy.CallOpts)
}

// Collect is a paid mutator transaction binding the contract method 0xe5225381.
//
// Solidity: function collect() returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactor) Collect(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Strategy.contract.Transact(opts, "collect")
}

// Collect is a paid mutator transaction binding the contract method 0xe5225381.
//
// Solidity: function collect() returns()
func (_UniswapV3Strategy *UniswapV3StrategySession) Collect() (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Collect(&_UniswapV3Strategy.TransactOpts)
}

// Collect is a paid mutator transaction binding the contract method 0xe5225381.
//
// Solidity: function collect() returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactorSession) Collect() (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Collect(&_UniswapV3Strategy.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Strategy.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_UniswapV3Strategy *UniswapV3StrategySession) Deposit() (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Deposit(&_UniswapV3Strategy.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactorSession) Deposit() (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Deposit(&_UniswapV3Strategy.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _vault, address _owner) returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactor) Initialize(opts *bind.TransactOpts, _vault common.Address, _owner common.Address) (*types.Transaction, error) {
	return _UniswapV3Strategy.contract.Transact(opts, "initialize", _vault, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _vault, address _owner) returns()
func (_UniswapV3Strategy *UniswapV3StrategySession) Initialize(_vault common.Address, _owner common.Address) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Initialize(&_UniswapV3Strategy.TransactOpts, _vault, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _vault, address _owner) returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactorSession) Initialize(_vault common.Address, _owner common.Address) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Initialize(&_UniswapV3Strategy.TransactOpts, _vault, _owner)
}

// Maintain is a paid mutator transaction binding the contract method 0x15ad819d.
//
// Solidity: function maintain(bool zeroForOne, uint256 amountIn, int24 tickLower, int24 tickUpper) returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactor) Maintain(opts *bind.TransactOpts, zeroForOne bool, amountIn *big.Int, tickLower *big.Int, tickUpper *big.Int) (*types.Transaction, error) {
	return _UniswapV3Strategy.contract.Transact(opts, "maintain", zeroForOne, amountIn, tickLower, tickUpper)
}

// Maintain is a paid mutator transaction binding the contract method 0x15ad819d.
//
// Solidity: function maintain(bool zeroForOne, uint256 amountIn, int24 tickLower, int24 tickUpper) returns()
func (_UniswapV3Strategy *UniswapV3StrategySession) Maintain(zeroForOne bool, amountIn *big.Int, tickLower *big.Int, tickUpper *big.Int) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Maintain(&_UniswapV3Strategy.TransactOpts, zeroForOne, amountIn, tickLower, tickUpper)
}

// Maintain is a paid mutator transaction binding the contract method 0x15ad819d.
//
// Solidity: function maintain(bool zeroForOne, uint256 amountIn, int24 tickLower, int24 tickUpper) returns()
func (_UniswapV3Strategy *UniswapV3StrategyTransactorSession) Maintain(zeroForOne bool, amountIn *big.Int, tickLower *big.Int, tickUpper *big.Int) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Maintain(&_UniswapV3Strategy.TransactOpts, zeroForOne, amountIn, tickLower, tickUpper)
}

// Rebalance is a paid mutator transaction binding the contract method 0x7d7c2a1c.
//
// Solidity: function rebalance() returns(int24 tickLower, int24 tickUpper)
func (_UniswapV3Strategy *UniswapV3StrategyTransactor) Rebalance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Strategy.contract.Transact(opts, "rebalance")
}

// Rebalance is a paid mutator transaction binding the contract method 0x7d7c2a1c.
//
// Solidity: function rebalance() returns(int24 tickLower, int24 tickUpper)
func (_UniswapV3Strategy *UniswapV3StrategySession) Rebalance() (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Rebalance(&_UniswapV3Strategy.TransactOpts)
}

// Rebalance is a paid mutator transaction binding the contract method 0x7d7c2a1c.
//
// Solidity: function rebalance() returns(int24 tickLower, int24 tickUpper)
func (_UniswapV3Strategy *UniswapV3StrategyTransactorSession) Rebalance() (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Rebalance(&_UniswapV3Strategy.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0ad58d2f.
//
// Solidity: function withdraw(uint256 shares, uint256 totalShares, address recipient) returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategyTransactor) Withdraw(opts *bind.TransactOpts, shares *big.Int, totalShares *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _UniswapV3Strategy.contract.Transact(opts, "withdraw", shares, totalShares, recipient)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0ad58d2f.
//
// Solidity: function withdraw(uint256 shares, uint256 totalShares, address recipient) returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategySession) Withdraw(shares *big.Int, totalShares *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Withdraw(&_UniswapV3Strategy.TransactOpts, shares, totalShares, recipient)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0ad58d2f.
//
// Solidity: function withdraw(uint256 shares, uint256 totalShares, address recipient) returns(uint256 amount0, uint256 amount1)
func (_UniswapV3Strategy *UniswapV3StrategyTransactorSession) Withdraw(shares *big.Int, totalShares *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _UniswapV3Strategy.Contract.Withdraw(&_UniswapV3Strategy.TransactOpts, shares, totalShares, recipient)
}
