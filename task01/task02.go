package main

import (
	"context"
	"crypto/ecdsa"
	"eth-client-study/task01/counter"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	deployCounterContract()
}

// 部署合约Counter
func deployCounterContract() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/6FDc8VT")

	if err != nil {
		fmt.Println("连接错误", err)
		return
	}
	defer client.Close()

	//私钥
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		fmt.Println("私钥转换失败", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("公钥转换失败:", ok)
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("获取nonce失败", err)
		return
	}
	fmt.Println("最新nonce:", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("获取gasPrice失败", err)
		return
	}
	fmt.Println("gasPrice:", gasPrice)

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("获取chainId失败", err)
		return
	}
	fmt.Println("chainId:", chainId)

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		fmt.Println("获取transactor失败", err)
		return
	}
	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0)
	opts.GasLimit = uint64(3000000)
	opts.GasPrice = gasPrice
	opts.Context = context.Background()
	//部署合约
	address, transaction, counterContract, err := counter.DeployCounter(opts, client)
	if err != nil {
		fmt.Println("部署合约失败", err)
		return
	}
	fmt.Println("部署成功--合约地址:", address.Hex(), "交易Hash:", transaction.Hash().Hex())

	// 等待部署交易确认
	waitForTransaction(client, transaction.Hash())

	count, _ := counterContract.GetCount(&bind.CallOpts{})
	fmt.Println("初始化count:", count)

	//调用合约Increment方法
	// 创建一个绑定的transactor
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		fmt.Println("创建transactor失败", err)
	}
	transactOpts.GasLimit = 3000000
	gasPrice, _ = client.SuggestGasPrice(context.Background())
	transactOpts.GasPrice = gasPrice
	//调用合约Increment方法
	transaction, err = counterContract.Increment(transactOpts)
	if err != nil {
		fmt.Println("调用合约失败", err)
		return
	} else {
		fmt.Println("调用合约Increment方法成功交易Hash:", transaction.Hash().Hex())

		// 等待交易确认
		waitForTransaction(client, transaction.Hash())

		count, _ = counterContract.GetCount(&bind.CallOpts{})
		fmt.Println("调用合约Increment方法成功count:", count)
	}

	transaction, err = counterContract.Decrement(transactOpts)
	if err != nil {
		fmt.Println("调用合约失败", err)
		return
	} else {
		fmt.Println("调用合约Decrement方法成功交易Hash:", transaction.Hash().Hex())

		// 等待交易确认
		waitForTransaction(client, transaction.Hash())

		count, _ = counterContract.GetCount(&bind.CallOpts{})
		fmt.Println("调用合约Decrement方法成功count:", count)
	}

}

// waitForTransaction 等待交易被确认
func waitForTransaction(client *ethclient.Client, txHash common.Hash) {
	fmt.Printf("等待交易 %s 被确认...\n", txHash.Hex())
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			// 交易还未被确认，继续等待
			time.Sleep(1 * time.Second)
			continue
		}
		if receipt.Status == types.ReceiptStatusSuccessful {
			fmt.Printf("交易 %s 已成功确认\n", txHash.Hex())
			break
		} else {
			fmt.Printf("交易 %s 执行失败\n", txHash.Hex())
			break
		}
	}
}
