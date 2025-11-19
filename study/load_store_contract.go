package main

import (
	"context"
	"eth-client-study/study/store"
	"eth-client-study/utils"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractAddress = "0x183AdfEe585d04Db1Ab151840D6399009beC2bC4"

// 加载store合约
func main() {

	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	newStore, err := store.NewStore(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	version, err := newStore.Version(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("合约版本:", version)
	//私钥
	privateKey, err := crypto.HexToECDSA(utils.GetEnv("PRIVATE_KEY1"))
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}
	transactOpts.GasLimit = 3000000
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	transactOpts.GasPrice = gasPrice
	//调用setItem函数
	item, err := newStore.SetItem(transactOpts, [32]byte{0x01}, [32]byte{0x02})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("交易hash:", item.Hash().Hex())
	items, err := newStore.Items(&bind.CallOpts{}, [32]byte{0x01})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("获取的key:", items)

}
