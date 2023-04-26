package chainlink

import (
	"fmt"
	"math/big"
	"ponytaapi/internal/service/blockserv"
	"ponytaapi/pkg/blockchain/ethrpc"
	"ponytaapi/pkg/utils"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

type AggregatorCC struct {
	W3Rpc *ethrpc.Web3RPC
}

// const PancakeswapFactoryAddress = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
// const PancakeswapPairAbi = `
// [{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"},
// {"constant":true,"inputs":[],"name":"token0","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},
// {"constant":true,"inputs":[],"name":"token1","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]
// `

func NewAggregatorCC(address string) (*AggregatorCC, error) {
	// https://eth-mainnet.g.alchemy.com/v2/d4QY1BuMJPkipt5d0zQg0xabWWgfG94X
	// https://mainnet.infura.io/v3/3316f4fb520340e1bd369cadf0d4a8b6
	w3, err := blockserv.NewSpecifiedRPC("https://eth-mainnet.g.alchemy.com/v2/d4QY1BuMJPkipt5d0zQg0xabWWgfG94X")
	if err != nil {
		return nil, err
	}

	res := &AggregatorCC{W3Rpc: w3}

	if err := res.W3Rpc.LoadContract(address, Aggregator_ABI_V3, ""); err != nil {
		return nil, fmt.Errorf("load contract object err: %s", err.Error())
	}

	return res, nil
}

func GetErc20PriceByAddress(address string) (float64, error) {
	cc, err := NewAggregatorCC(address)
	if err != nil {
		return 0, err
	}

	decimals, err := cc.CcDecimals()
	if err != nil {
		return 0, errors.Wrap(err, "get decimals")
	}

	basePirce, err := cc.CcLatestAnswer()
	if err != nil {
		return 0, errors.Wrap(err, "get LatestAnswer")
	}

	price, err := utils.TokenValueTo(basePirce.String(), decimals)
	if err != nil {
		return 0, err
	}

	price, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", price), 64)
	return price, nil
}

// function decimals() external view returns (uint8);
func (t *AggregatorCC) CcDecimals() (uint8, error) {
	output := []interface{}{}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "decimals")
	if err != nil {
		return 0, err
	}

	decimals := output[0].(uint8)

	return decimals, nil
}

// function latestAnswer() external view returns (int256);
func (t *AggregatorCC) CcLatestAnswer() (*big.Int, error) {
	output := []interface{}{}

	err := t.W3Rpc.Contract.Call(&bind.CallOpts{}, &output, "latestAnswer")
	if err != nil {
		return nil, err
	}

	price := output[0].(*big.Int)

	return price, nil
}
