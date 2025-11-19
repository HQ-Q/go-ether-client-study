// contract_event.go 展示了如何监听和查询以太坊智能合约中的事件（特别是 ItemSet 事件）
// 主要功能包括：
// 1. 实时监听新产生的 ItemSet 事件（推荐使用 WebSocket 连接）
// 2. 查询指定区块范围内的历史 ItemSet 事件

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// StoreABI 是一个 Solidity 合约的 ABI 描述 JSON 字符串
// 包含了合约的构造函数、状态变量 items 和 version，以及 setItem 函数和 ItemSet 事件定义
// 合约包含以下内容：
// - 构造函数：接受一个 string 类型的 _version 参数
// - 状态变量 items：映射类型 mapping(bytes32 => bytes32)
// - 状态变量 version：string 类型
// - 函数 setItem：设置键值对
// - 事件 ItemSet：当调用 setItem 时触发
var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

// ItemSetEvent 定义了与 Solidity 中 ItemSet 事件相对应的 Go 结构体
// 该结构体用于解析从区块链上获取的事件日志数据
type ItemSetEvent struct {
	// Key 对应合约事件中的 key 参数（bytes32 类型）
	// 使用 abi:"key" 标签确保正确映射到合约事件字段
	Key [32]byte `abi:"key"`

	// Value 对应合约事件中的 value 参数（bytes32 类型）
	// 使用 abi:"value" 标签确保正确映射到合约事件字段
	Value [32]byte `abi:"value"`
}

// main 函数是程序入口点
// 建立与以太坊节点的连接，设置合约地址和 ABI，然后启动事件监听
func main() {
	// 使用 WebSocket 连接到 Ethereum 节点（推荐用于实时监听）
	// wss 协议相比 http 协议更适合实时通信场景
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}

	// 指定目标合约地址（在 Sepolia 测试网络上部署的 Store 合约）
	contractAddress := common.HexToAddress("0x183AdfEe585d04Db1Ab151840D6399009beC2bC4")

	// 将合约 ABI 字符串解析为 Go 可操作的对象
	// abi.ABI 对象提供了对合约函数和事件的操作方法
	contractABI, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatalf("解析ABI失败：%v", err)
	}

	// ------------------------------
	// 场景1：实时监听新产生的 ItemSet 事件（推荐 WS 协议）
	// ------------------------------
	fmt.Println("开始监听ItemSet事件...")
	listenNewEvents(context.Background(), client, contractABI, contractAddress)

	//// ------------------------------
	//// 场景2：查询历史 ItemSet 事件（指定区块范围）
	//// ------------------------------
	//// 示例：从第1个区块到最新区块之间查找所有 ItemSet 事件
	//fromBlock := uint64(1)
	//toBlock := uint64(0) // 0 表示最新区块
	//fmt.Printf("查询区块%d到最新区块的历史事件...\n", fromBlock)
	//historyEvents, err := queryHistoryEvents(context.Background(), client, contractABI, contractAddress, fromBlock, toBlock)
	//if err != nil {
	//	log.Fatalf("查询历史事件失败：%v", err)
	//}
	//for i, event := range historyEvents {
	//	fmt.Printf("历史事件%d：Key=%x, Value=%x\n", i+1, event.Key, event.Value)
	//}
	//
	//// 阻塞主协程，以便持续监听事件（在生产环境中可以通过信号处理实现优雅退出）
	//select {}
}

// listenNewEvents 监听来自指定合约的新 ItemSet 事件
// 使用 WebSocket 连接实时接收事件通知
// 参数:
// ctx - 上下文，用于控制函数执行生命周期
// client - 与以太坊节点建立的连接客户端
// contractABI - 合约的 ABI 对象，用于解析事件
// contractAddr - 合约地址，指定要监听的合约
func listenNewEvents(ctx context.Context, client *ethclient.Client, contractABI abi.ABI, contractAddr common.Address) {
	// 构造 FilterQuery 来筛选特定合约及事件的数据
	query := ethereum.FilterQuery{
		// Addresses 指定要监听的合约地址列表
		// 只关注这个合约地址的日志，过滤掉其他合约的事件
		Addresses: []common.Address{contractAddr},

		// Topics 用于指定要监听的事件签名
		// ItemSet 事件的 ID 作为第一个 Topic，确保只接收该类型的事件
		// contractABI.Events["ItemSet"].ID 获取事件的 Keccak256 哈希值
		Topics: [][]common.Hash{{contractABI.Events["ItemSet"].ID}},
	}

	// 创建一个通道用来接收日志数据
	// types.Log 是 go-ethereum 中定义的日志结构体
	logs := make(chan types.Log)

	// 订阅符合条件的日志更新
	// SubscribeFilterLogs 方法会持续监听新区块中的相关事件
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatalf("订阅事件失败：%v", err)
	}

	// 在函数返回前取消订阅，释放资源
	defer sub.Unsubscribe()

	// 不断循环等待新的日志到来
	for {
		select {
		// 处理订阅过程中可能出现的错误
		case err := <-sub.Err():
			log.Printf("订阅异常：%v，重新订阅...", err)
			// 如果发生错误，则尝试重新订阅（提高健壮性）
			// 忽略重新订阅可能产生的错误，继续尝试监听
			sub, _ = client.SubscribeFilterLogs(ctx, query, logs)
			continue

		// 处理接收到的新日志事件
		case log := <-logs:
			// 解析日志数据并填充到 ItemSetEvent 结构体中
			// UnpackIntoInterface 方法根据 ABI 信息将日志数据解析到指定结构体
			var event ItemSetEvent
			if err := contractABI.UnpackIntoInterface(&event, "ItemSet", log.Data); err != nil {
				fmt.Printf("解析事件失败：%v", err)
				continue
			}

			// 打印事件相关信息
			// log.BlockNumber 提供事件所在区块号
			// log.TxHash.Hex() 提供事件所在的交易哈希
			// event.Key 和 event.Value 是解析出的事件参数
			fmt.Printf(
				"收到新事件：区块号=%d, 交易哈希=%s, Key=%x, Value=%x\n",
				log.BlockNumber, log.TxHash.Hex(), event.Key, event.Value,
			)
		}
	}
}

// queryHistoryEvents 查询指定区块范围内发生的 ItemSet 事件
// 适用于需要检索过去事件的场景，如数据同步、审计等
// 参数:
// ctx - 上下文
// client - 与以太坊节点建立的连接客户端
// contractABI - 合约的 ABI 对象，用于解析事件
// contractAddr - 合约地址
// fromBlock - 起始区块号（包含）
// toBlock - 结束区块号（包含），0 表示最新区块
// 返回值:
// []ItemSetEvent - 解析后的历史事件列表
// error - 查询过程中的错误信息
func queryHistoryEvents(
	ctx context.Context,
	client *ethclient.Client,
	contractABI abi.ABI,
	contractAddr common.Address,
	fromBlock, toBlock uint64,
) ([]ItemSetEvent, error) {
	// 构造带有区块范围的查询参数
	query := ethereum.FilterQuery{
		FromBlock: bigInt(fromBlock),                                   // 起始区块
		ToBlock:   bigInt(toBlock),                                     // 截止区块（0代表最新）
		Addresses: []common.Address{contractAddr},                      // 指定合约地址
		Topics:    [][]common.Hash{{contractABI.Events["ItemSet"].ID}}, // 指定事件签名
	}

	// 发起请求获取日志记录
	// FilterLogs 方法一次性获取指定区块范围内的所有匹配日志
	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs: %v", err)
	}

	// 解析每一条日志，并将其转换成 ItemSetEvent 对象
	var events []ItemSetEvent
	for _, log := range logs {
		var event ItemSetEvent
		// 使用 ABI 信息解析日志数据
		if err := contractABI.UnpackIntoInterface(&event, "ItemSet", log.Data); err != nil {
			return nil, fmt.Errorf("unpack log: %v", err)
		}
		// 将解析成功的事件添加到结果列表中
		events = append(events, event)
	}

	return events, nil
}

// bigInt 工具函数：将 uint64 数值转化为 *big.Int 类型
// 适配 FilterQuery 中 FromBlock / ToBlock 的要求
// Ethereum 区块号使用 *big.Int 类型表示，因此需要转换
func bigInt(n uint64) *big.Int {
	return new(big.Int).SetUint64(n)
}
