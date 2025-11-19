package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 查询账户余额
func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		fmt.Println("连接失败", err)
	}
	defer client.Close()

	account := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")
	// 获取账户余额
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		fmt.Println("查询余额失败", err)
	}

	fmt.Println("余额:", balance)
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	price := new(big.Float).Quo(fbalance, big.NewFloat(1e18))
	fmt.Println("余额:", price, "ETH")
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println("待处理余额:", pendingBalance)

	//查询指定区块余额
	balance, err = client.BalanceAt(context.Background(), account, big.NewInt(5532993))

	fmt.Println("区块余额:", balance)
}
