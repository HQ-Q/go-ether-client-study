package main

import (
	"context"
	"eth-client-study/study/store"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	contractAddr := "0x183AdfEe585d04Db1Ab151840D6399009beC2bC4"
	storeContract, err := store.NewStoreFilterer(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个通道来接收事件
	itemSetCh := make(chan *store.StoreItemSet)

	// 监听ItemSet事件
	sub, err := storeContract.WatchItemSet(&bind.WatchOpts{Context: context.Background()}, itemSetCh)
	if err != nil {
		log.Fatal("监听ItemSet事件失败:", err)
	}
	defer sub.Unsubscribe()

	log.Println("开始监听ItemSet事件...")

	// 循环监听事件
	for {
		select {
		case itemSet := <-itemSetCh:
			log.Printf("收到ItemSet事件:")
			log.Printf("  交易哈希: %s", itemSet.Raw.TxHash.Hex())
			log.Printf("  区块号: %d", itemSet.Raw.BlockNumber)
			log.Printf("  Key: %s", string(itemSet.Key[:]))
			log.Printf("  Value: %s", string(itemSet.Value[:]))
		case err := <-sub.Err():
			log.Fatal("订阅错误:", err)
		}
	}
}
