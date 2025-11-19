// query_receipt.go - 查询以太坊区块交易收据示例程序
package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/ethclient" // 提供与以太坊节点交互的客户端
	"github.com/ethereum/go-ethereum/rpc"       // 支持通过RPC调用传递块号或块哈希
)

// 查询区块交易收据
func main() {
	// 使用 Alchemy 提供的 Sepolia 测试网络端点创建一个以太坊客户端连接
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err) // 如果建立连接失败，则记录错误并退出程序
	}
	defer client.Close() // 程序结束前关闭客户端连接
	//blockNumber := big.NewInt(5671744)
	//block, err := client.BlockByNumber(context.Background(), blockNumber)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 定义要查询的区块哈希（Sepolia测试网上的某个区块）
	//blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	// 根据区块哈希获取该区块内所有交易的收据信息
	//receiptByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(block.Hash(), false))
	// 根据区块编号获取该区块内所有交易的收据信息
	receiptByNumber, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(5671744))

	//fmt.Println("receiptByNumber:", receiptByNumber) // 打印基于区块编号获取的收据列表
	//fmt.Println("receiptByHash:", receiptByHash)     // 打印基于区块哈希获取的收据列表

	if err != nil {
		log.Fatal(err) // 若获取过程中出现错误则终止程序
	}

	//jsonReceiptByHash, _ := json.Marshal(receiptByHash)
	//println("receiptByNumber:", string(jsonReceiptByHash))
	jsonReceiptByNumber, _ := json.Marshal(receiptByNumber)
	println("receiptByNumber:", string(jsonReceiptByNumber))

	// 遍历打印每个交易收据中的关键字段
	//for i, receipt := range receiptByHash {
	//fmt.Println("PostState:", receipt.PostState)                 // EIP-658后不再使用，通常为空
	//fmt.Println("Status:", receipt.Status)                       // 交易状态：1表示成功，0表示失败
	//fmt.Println("CumulativeGasUsed:", receipt.CumulativeGasUsed) // 区块中累计消耗的gas总量
	//fmt.Println("Bloom:", receipt.Bloom)                         // 日志布隆过滤器，用于快速检索事件日志
	//fmt.Println("Logs:", receipt.Logs)                           // 交易产生的事件日志数组
	//fmt.Println("TxHash:", receipt.TxHash)                       // 交易哈希值
	//fmt.Println("ContractAddress:", receipt.ContractAddress)     // 合约部署地址（如果是合约创建交易）
	//fmt.Println("GasUsed:", receipt.GasUsed)                     // 此次交易实际使用的gas数量
	//fmt.Println("EffectiveGasPrice:", receipt.EffectiveGasPrice) // 实际支付的gas价格（包括基础费用和小费）
	//fmt.Println("BlobGasUsed:", receipt.BlobGasUsed)             // blob交易使用的gas量（适用于EIP-4844）
	//fmt.Println("BlobGasPrice:", receipt.BlobGasPrice)           // blob数据的gas价格（适用于EIP-4844）
	//fmt.Println("BlockHash:", receipt.BlockHash)                 // 包含此交易的区块哈希
	//fmt.Println("BlockNumber:", receipt.BlockNumber)             // 包含此交易的区块编号
	//fmt.Println("TransactionIndex:", receipt.TransactionIndex)   // 交易在区块中的索引位置
	//fmt.Println("------------------")                            // 分隔线以便阅读输出结果
	//bytes, _ := json.Marshal(receipt)
	//fmt.Println("RECEIPT:", i, "------", string(bytes))
	//}

}
