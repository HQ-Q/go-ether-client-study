package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 订阅区块
func main() {
	wsUrl := "wss://eth-sepolia.g.alchemy.com/v2/xxxxxxxx"
	// 2. 创建客户端时添加超时上下文，避免无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(ctx, wsUrl)
	if err != nil {
		fmt.Printf("创建以太坊客户端失败：%v\n", err)
		fmt.Println("请检查：1) WebSocket URL是否完整 2) 网络是否通畅 3) API密钥是否有效")
		return
	}
	defer func() {
		// 确保client非nil时才关闭
		if client != nil {
			client.Close()
		}
	}()

	// 验证客户端是否正常连接
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Printf("验证连接失败：%v\n", err)
		return
	}
	fmt.Printf("成功连接到网络，链ID：%d\n", chainID.Uint64())

	// 3. 创建区块头通道（缓冲区10，避免阻塞）
	headers := make(chan *types.Header, 10)

	// 4. 订阅新区块头
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		fmt.Printf("订阅新区块失败：%v\n", err)
		return
	}
	defer sub.Unsubscribe() // 程序退出时取消订阅
	fmt.Println("开始监听新区块...（按Ctrl+C退出）")

	for {
		select {
		case err := <-sub.Err():
			fmt.Println("Subscribe error:", err)
			return
		case header := <-headers:
			// 成功收到新区块头
			fmt.Printf("\n==================== 新区块 ====================\n")
			fmt.Printf("区块号：%d\n", header.Number.Int64())
			fmt.Printf("区块Hash：%s\n", header.Hash().Hex())

			// 6. 获取完整区块数据（添加超时，避免阻塞）
			blockCtx, blockCancel := context.WithTimeout(context.Background(), 5*time.Second)
			block, err := client.BlockByHash(blockCtx, header.Hash())
			blockCancel() // 及时释放上下文
			if err != nil {
				fmt.Printf("获取完整区块失败（可能包含未支持的交易类型）：%v\n", err)
				fmt.Printf("仅输出区块头信息：区块号=%d，Hash=%s，时间戳=%s\n",
					header.Number.Int64(),
					header.Hash().Hex(),
					time.Unix(int64(header.Time), 0).Format("2006-01-02 15:04:05"))
				continue
			}
			// 输出区块详细信息（避免nil指针）
			fmt.Printf("区块号（完整区块）：%d\n", block.Number().Uint64())
			fmt.Printf("时间戳：%s\n", time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("难度：%d\n", block.Difficulty().Uint64())
			fmt.Printf("Nonce：%d\n", block.Nonce())
			fmt.Printf("区块大小：%d bytes\n", block.Size())
			fmt.Printf("Gas限制：%d\n", block.GasLimit())
			fmt.Printf("Gas使用：%d\n", block.GasUsed())
			fmt.Printf("矿工地址：%s\n", block.Coinbase().Hex())
			fmt.Printf("交易数量：%d\n", len(block.Transactions()))
			fmt.Printf("叔块数量：%d\n", len(block.Uncles()))
			fmt.Println("===============================================")

		}
	}

	fmt.Println("Done")
}
