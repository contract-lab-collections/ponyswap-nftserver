package ethrpc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Web3RPC struct {
	Client    *ethclient.Client
	Contract  *bind.BoundContract
	Account   *common.Address
	PrivKey   *ecdsa.PrivateKey
	Abi       *abi.ABI
	CcAddress common.Address
}

const (
	TX_FAULT   = 0
	TX_SUCCESS = 1
)

func NewWeb3RPC(web3url string) (*Web3RPC, error) {
	// Initialize an RPC connection
	web3Client, err := ethclient.Dial(web3url)
	if err != nil {
		return nil, err
	}

	w3 := &Web3RPC{
		Client: web3Client,
	}

	return w3, nil
}

// Query the latest account balance, error return 0
func (t *Web3RPC) BalanceAt(address string) (*big.Int, error) {
	balance, err := t.Client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("query account balance err:%s", err.Error())
	}

	return balance, nil
}

func (t *Web3RPC) GetBlockTimeByNumber(height int64) (int64, error) {
	blockNumber := big.NewInt(height)
	block, err := t.Client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return 0, err
	}

	return int64(block.Time()), nil
}

// Query pending or completed transactions
func (t *Web3RPC) GetPendingTx(txHash string) (*types.Receipt, error) {
	// 1 * 10 second
	for i := 0; i < 1; i++ {
		tx, isPending, err := t.Client.TransactionByHash(context.Background(), common.HexToHash(txHash))
		if err != nil {
			time.Sleep(3 * time.Second)
			tx, isPending, err = t.Client.TransactionByHash(context.Background(), common.HexToHash(txHash))
			if err != nil {
				time.Sleep(3 * time.Second)
				tx, isPending, err = t.Client.TransactionByHash(context.Background(), common.HexToHash(txHash))
				if err != nil {
					return nil, err
				}
			}
		}

		if isPending {
			time.Sleep(time.Second * 3)
			continue
		}

		receipt, err := t.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return nil, err
		}

		if receipt.Status == TX_FAULT {
			return nil, errors.New("tx failed")
		}

		return receipt, nil

	}

	return nil, errors.New("tx not found")
}

func (t *Web3RPC) GetTxReceipt(txHash string) (*types.Receipt, error) {

	receipt, err := t.Client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, errors.Wrap(err, "TransactionReceipt")
	}

	if receipt.Status == TX_FAULT {
		return nil, errors.New("tx failed")
	}

	return receipt, nil

	// return nil, errors.New("tx not found")
}
