package main

import (
	token "eth-client-study/study/erc20"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 查询代币余额
func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		fmt.Println("连接失败", err)
	}
	defer client.Close()

	tokenAddress := common.HexToAddress("0x2f8C29909a2697E4E0449662302aAa1750f2cF98")
	instance, err := token.NewErc20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	accountAddress := common.HexToAddress("0xFA73Ee972cB6A7af855846635Ad65427a7009d4e")
	balance, err := instance.BalanceOf(&bind.CallOpts{}, accountAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("账户余额:", balance.String())
	//账户余额除以1e18
	floatBalance := new(big.Float)
	floatBalance.SetString(balance.String())
	price := new(big.Float).Quo(floatBalance, big.NewFloat(1e18))
	fmt.Println("账户余额:", price)

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("代币名称:", name)

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("代币符号:", symbol)

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("代币精度:", decimals)

}
