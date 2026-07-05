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

// IQuoterV2QuoteExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterV2QuoteExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	AmountIn          *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterV2QuoteExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterV2QuoteExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Amount            *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterV2MetaData contains all meta data concerning the IQuoterV2 contract.
var IQuoterV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"quoteExactInput\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\",\"internalType\":\"uint160[]\"},{\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteExactInputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIQuoterV2.QuoteExactInputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteExactOutput\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\",\"internalType\":\"uint160[]\"},{\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteExactOutputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIQuoterV2.QuoteExactOutputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"}]",
}

// IQuoterV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use IQuoterV2MetaData.ABI instead.
var IQuoterV2ABI = IQuoterV2MetaData.ABI

// IQuoterV2 is an auto generated Go binding around an Ethereum contract.
type IQuoterV2 struct {
	IQuoterV2Caller     // Read-only binding to the contract
	IQuoterV2Transactor // Write-only binding to the contract
	IQuoterV2Filterer   // Log filterer for contract events
}

// IQuoterV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type IQuoterV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IQuoterV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IQuoterV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IQuoterV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IQuoterV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IQuoterV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IQuoterV2Session struct {
	Contract     *IQuoterV2        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IQuoterV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IQuoterV2CallerSession struct {
	Contract *IQuoterV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IQuoterV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IQuoterV2TransactorSession struct {
	Contract     *IQuoterV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IQuoterV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type IQuoterV2Raw struct {
	Contract *IQuoterV2 // Generic contract binding to access the raw methods on
}

// IQuoterV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IQuoterV2CallerRaw struct {
	Contract *IQuoterV2Caller // Generic read-only contract binding to access the raw methods on
}

// IQuoterV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IQuoterV2TransactorRaw struct {
	Contract *IQuoterV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIQuoterV2 creates a new instance of IQuoterV2, bound to a specific deployed contract.
func NewIQuoterV2(address common.Address, backend bind.ContractBackend) (*IQuoterV2, error) {
	contract, err := bindIQuoterV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IQuoterV2{IQuoterV2Caller: IQuoterV2Caller{contract: contract}, IQuoterV2Transactor: IQuoterV2Transactor{contract: contract}, IQuoterV2Filterer: IQuoterV2Filterer{contract: contract}}, nil
}

// NewIQuoterV2Caller creates a new read-only instance of IQuoterV2, bound to a specific deployed contract.
func NewIQuoterV2Caller(address common.Address, caller bind.ContractCaller) (*IQuoterV2Caller, error) {
	contract, err := bindIQuoterV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IQuoterV2Caller{contract: contract}, nil
}

// NewIQuoterV2Transactor creates a new write-only instance of IQuoterV2, bound to a specific deployed contract.
func NewIQuoterV2Transactor(address common.Address, transactor bind.ContractTransactor) (*IQuoterV2Transactor, error) {
	contract, err := bindIQuoterV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IQuoterV2Transactor{contract: contract}, nil
}

// NewIQuoterV2Filterer creates a new log filterer instance of IQuoterV2, bound to a specific deployed contract.
func NewIQuoterV2Filterer(address common.Address, filterer bind.ContractFilterer) (*IQuoterV2Filterer, error) {
	contract, err := bindIQuoterV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IQuoterV2Filterer{contract: contract}, nil
}

// bindIQuoterV2 binds a generic wrapper to an already deployed contract.
func bindIQuoterV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IQuoterV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IQuoterV2 *IQuoterV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IQuoterV2.Contract.IQuoterV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IQuoterV2 *IQuoterV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IQuoterV2.Contract.IQuoterV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IQuoterV2 *IQuoterV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IQuoterV2.Contract.IQuoterV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IQuoterV2 *IQuoterV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IQuoterV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IQuoterV2 *IQuoterV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IQuoterV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IQuoterV2 *IQuoterV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IQuoterV2.Contract.contract.Transact(opts, method, params...)
}

// QuoteExactInput is a paid mutator transaction binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Transactor) QuoteExactInput(opts *bind.TransactOpts, path []byte, amountIn *big.Int) (*types.Transaction, error) {
	return _IQuoterV2.contract.Transact(opts, "quoteExactInput", path, amountIn)
}

// QuoteExactInput is a paid mutator transaction binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Session) QuoteExactInput(path []byte, amountIn *big.Int) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactInput(&_IQuoterV2.TransactOpts, path, amountIn)
}

// QuoteExactInput is a paid mutator transaction binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2TransactorSession) QuoteExactInput(path []byte, amountIn *big.Int) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactInput(&_IQuoterV2.TransactOpts, path, amountIn)
}

// QuoteExactInputSingle is a paid mutator transaction binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Transactor) QuoteExactInputSingle(opts *bind.TransactOpts, params IQuoterV2QuoteExactInputSingleParams) (*types.Transaction, error) {
	return _IQuoterV2.contract.Transact(opts, "quoteExactInputSingle", params)
}

// QuoteExactInputSingle is a paid mutator transaction binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Session) QuoteExactInputSingle(params IQuoterV2QuoteExactInputSingleParams) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactInputSingle(&_IQuoterV2.TransactOpts, params)
}

// QuoteExactInputSingle is a paid mutator transaction binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2TransactorSession) QuoteExactInputSingle(params IQuoterV2QuoteExactInputSingleParams) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactInputSingle(&_IQuoterV2.TransactOpts, params)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Transactor) QuoteExactOutput(opts *bind.TransactOpts, path []byte, amountOut *big.Int) (*types.Transaction, error) {
	return _IQuoterV2.contract.Transact(opts, "quoteExactOutput", path, amountOut)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Session) QuoteExactOutput(path []byte, amountOut *big.Int) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactOutput(&_IQuoterV2.TransactOpts, path, amountOut)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2TransactorSession) QuoteExactOutput(path []byte, amountOut *big.Int) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactOutput(&_IQuoterV2.TransactOpts, path, amountOut)
}

// QuoteExactOutputSingle is a paid mutator transaction binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Transactor) QuoteExactOutputSingle(opts *bind.TransactOpts, params IQuoterV2QuoteExactOutputSingleParams) (*types.Transaction, error) {
	return _IQuoterV2.contract.Transact(opts, "quoteExactOutputSingle", params)
}

// QuoteExactOutputSingle is a paid mutator transaction binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2Session) QuoteExactOutputSingle(params IQuoterV2QuoteExactOutputSingleParams) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactOutputSingle(&_IQuoterV2.TransactOpts, params)
}

// QuoteExactOutputSingle is a paid mutator transaction binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_IQuoterV2 *IQuoterV2TransactorSession) QuoteExactOutputSingle(params IQuoterV2QuoteExactOutputSingleParams) (*types.Transaction, error) {
	return _IQuoterV2.Contract.QuoteExactOutputSingle(&_IQuoterV2.TransactOpts, params)
}
