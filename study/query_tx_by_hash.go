package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 查询交易
func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	// 获取当前网络的链ID
	chainID, err := client.ChainID(context.Background())
	hashStr := "0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5"
	hash := common.HexToHash(hashStr)
	tx, isPending, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tx:", tx.Hash().Hex(), "isPending:", isPending)
	// 打印交易哈希值
	fmt.Println("交易哈希值：", tx.Hash().Hex())

	// 打印交易金额（wei单位）
	fmt.Println("交易金额：", tx.Value().String())

	// 打印交易的Gas限制
	fmt.Println("Gas限制：", tx.Gas())

	// 打印交易的Gas价格
	fmt.Println("Gas价格：", tx.GasPrice().Uint64())

	// 打印交易的nonce值
	fmt.Println("Nonce值：", tx.Nonce())

	// 打印交易的数据字段（通常是合约调用数据）
	fmt.Println("交易数据：", tx.Data())

	// 打印交易接收方地址
	fmt.Println("接收方地址：", tx.To().Hex())

	// 打印交易所在的链ID
	fmt.Println("链ID：", tx.ChainId().Uint64())

	// 使用EIP155签名者从交易中恢复发送方地址
	if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
		// 打印发送方地址
		fmt.Println("发送方地址：", sender.Hex())
	} else {
		// 如果恢复发送方地址失败，则打印错误信息
		fmt.Println("发送方地址获取失败：", err)
		log.Fatal(err)
	}

	// 根据交易哈希获取交易回执
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())

	// 检查获取交易回执是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 打印交易状态（1表示成功，0表示失败）
	fmt.Println("交易状态：", receipt.Status)

	// 打印交易产生的事件日志
	fmt.Println("事件日志：", receipt.Logs)

}
