// 定义包名为main，表示这是一个可执行程序
package main

// 导入所需的包
import (
	"context"  // 用于控制请求的上下文
	"fmt"      // 用于格式化输入输出
	"log"      // 用于记录日志
	"math/big" // 用于处理大整数

	"github.com/ethereum/go-ethereum/common"     // 以太坊通用工具函数
	"github.com/ethereum/go-ethereum/core/types" // 以太坊核心类型定义
	"github.com/ethereum/go-ethereum/ethclient"  // 以太坊客户端库
)

// 主函数，程序入口点
func main() {
	// 连接到以太坊Sepolia测试网络的Alchemy节点
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/qLVks-ioeg_c5r6FDc8VT")
	//client, err := ethclient.Dial("http://localhost:8545") // 备用的本地节点连接（已注释）

	// 检查连接是否出错，如果出错则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 延迟关闭客户端连接
	defer client.Close()

	// 获取当前网络的链ID
	chainID, err := client.ChainID(context.Background())

	// 检查获取链ID是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个指定区块号5671744的big.Int对象
	blockNumber := big.NewInt(5671744)

	// 根据区块号获取对应区块的信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)

	// 检查获取区块是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 遍历区块中的所有交易
	for _, tx := range block.Transactions() {
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

		// 只处理第一个交易就跳出循环
		//break
	}

	// 将十六进制字符串转换为区块哈希对象
	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")

	// 获取指定区块哈希的交易数量
	count, err := client.TransactionCount(context.Background(), blockHash)

	// 检查获取交易数量是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 遍历区块中的所有交易索引
	for idx := uint(0); idx < count; idx++ {
		// 根据区块哈希和交易索引获取具体的交易
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)

		// 检查获取交易是否有错误，如果有则记录错误并终止程序
		if err != nil {
			log.Fatal(err)
		}

		// 打印交易哈希值
		fmt.Println("区块内交易哈希值：", tx.Hash().Hex())

		// 只处理第一个交易就跳出循环
		break
	}

	// 将十六进制字符串转换为交易哈希对象
	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")

	// 根据交易哈希获取交易详情和是否处于待处理状态
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)

	print(tx.Hash().Hex())

	if err != nil {
		log.Fatal(err)
	}

	// 打印交易是否处于待处理状态
	fmt.Println("是否待处理：", isPending)

	// 打印交易哈希值
	fmt.Println("指定交易哈希值：", tx.Hash().Hex())

	// 注意：这里有一行被注释掉的代码，原本应该是打印isPending状态
	//fmt.Println(isPending)       // false
}
