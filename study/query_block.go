// 定义包名为main，表示这是一个可执行程序
package main

// 导入所需的包
import (
	"context" // 用于控制请求的上下文
	"fmt"     // 用于格式化输入输出
	"log"     // 用于记录日志

	"github.com/ethereum/go-ethereum/ethclient" // 以太坊客户端库
)

// 查询区块信息
func main() {
	// 连接到以太坊Sepolia测试网络的Alchemy节点
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	//client, err := ethclient.Dial("http://localhost:8545") // 备用的本地节点连接（已注释）

	// 检查连接是否出错，如果出错则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 延迟关闭客户端连接
	defer client.Close()

	//blockNumber := big.NewInt(5671744) // 创建一个指定区块号的big.Int（已注释）

	// 获取最新的区块头信息（nil表示最新区块）
	header, err := client.HeaderByNumber(context.Background(), nil)

	// 打印区块号
	fmt.Println(header.Number.Uint64())

	// 打印区块时间戳
	fmt.Println(header.Time)

	// 打印区块难度值
	fmt.Println(header.Difficulty.Uint64())

	// 打印区块哈希值
	fmt.Println(header.Hash().Hex())

	// 检查获取区块头是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 获取完整的区块信息（包括交易数据）
	block, err := client.BlockByNumber(context.Background(), nil)

	// 检查获取区块是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 打印区块号
	fmt.Println(block.Number().Uint64())

	// 打印区块时间戳
	fmt.Println(block.Time())

	// 打印区块难度值
	fmt.Println(block.Difficulty().Uint64())

	// 打印区块哈希值
	fmt.Println(block.Hash().Hex())

	// 打印区块中包含的交易数量
	fmt.Println(len(block.Transactions()))

	// 获取指定区块哈希的交易数量
	count, err := client.TransactionCount(context.Background(), block.Hash())

	// 检查获取交易数量是否有错误，如果有则记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}

	// 打印交易数量
	fmt.Println(count)
}
