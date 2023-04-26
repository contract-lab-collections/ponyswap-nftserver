package blockserv

import (
	"fmt"
	"math/big"
	"ponytaapi/pkg/blockchain/ethrpc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

/*
interface ERC20 {
    event Approval(address indexed owner, address indexed spender, uint value);
    event Transfer(address indexed from, address indexed to, uint value);

    function name() external pure returns (string memory);
    function symbol() external pure returns (string memory);
    function decimals() external pure returns (uint8);
    function totalSupply() external view returns (uint);
    function balanceOf(address owner) external view returns (uint);
    function allowance(address owner, address spender) external view returns (uint);

    function approve(address spender, uint value) external returns (bool);
    function transfer(address to, uint value) external returns (bool);
    function transferFrom(address from, address to, uint value) external returns (bool);
}
*/

type Erc20TokenCC struct {
	W3Rpc *ethrpc.Web3RPC
}

func NewErc20TokenCC(erc20address, privateKey string) (*Erc20TokenCC, error) {
	w3, err := NewW3RPC()
	if err != nil {
		return nil, err
	}

	res := &Erc20TokenCC{W3Rpc: w3}

	if err := res.W3Rpc.LoadContract(erc20address, ERC20_ABI, privateKey); err != nil {
		return nil, fmt.Errorf("load exchange contract object err: %s", err.Error())
	}

	return res, nil
}

func (t *Erc20TokenCC) BalanceOf(address string) (*big.Int, error) {
	// bytes32
	input := []interface{}{
		common.HexToAddress(address),
	}

	balance := new(big.Int)
	output := []interface{}{&balance}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "balanceOf", input...)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (t *Erc20TokenCC) Allowance(owner, spender string) (*big.Int, error) {
	// bytes32
	input := []interface{}{
		common.HexToAddress(owner),
		common.HexToAddress(spender),
	}

	balance := new(big.Int)
	output := []interface{}{&balance}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "allowance", input...)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// function decimals() external pure returns (uint8);
func (t *Erc20TokenCC) Decimals() (uint8, error) {
	input := []interface{}{}
	output := []interface{}{}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "decimals", input...)
	if err != nil {
		return 0, err
	}

	rt := output[0].(uint8)
	return rt, nil
}

// function symbol() external pure returns (string memory);
func (t *Erc20TokenCC) Symbol() (string, error) {
	input := []interface{}{}

	output := []interface{}{}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "symbol", input...)
	if err != nil {
		return "", err
	}

	rt := output[0].(string)

	return rt, nil
}

// function name() external pure returns (string memory);
func (t *Erc20TokenCC) Name() (string, error) {
	input := []interface{}{}

	output := []interface{}{}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "name", input...)
	if err != nil {
		return "", err
	}

	rt := output[0].(string)

	return rt, nil
}

func GetErc20TokenInfo(address string) (string, string, uint8, error) {
	cc, err := NewErc20TokenCC(address, "")
	if err != nil {
		return "", "", 0, err
	}

	symbol, err := cc.Symbol()
	if err != nil {
		return "", "", 0, err
	}

	decimals, err := cc.Decimals()
	if err != nil {
		return "", "", 0, err
	}

	name, err := cc.Name()
	if err != nil {
		return "", "", 0, err
	}

	return name, symbol, decimals, nil
}
