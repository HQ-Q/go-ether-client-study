package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		fmt.Println("链接失败", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("xxxxxxxxxx")
	if err != nil {
		fmt.Println("私钥转换失败", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("公钥转换失败")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress:", fromAddress)
	nonceAt, err := client.PendingNonceAt(context.Background(), fromAddress)

	if err != nil {
		fmt.Println("获取nonce失败", err)
	}
	fmt.Println("交易nonce:", nonceAt)

	value := big.NewInt(1000000000000000000) //1 eth
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//收款账户地址 0xfa73ee972cb6a7af855846635ad65427a7009d4e
	toAddress := common.HexToAddress("0xfa73ee972cb6a7af855846635ad65427a7009d4e")
	var data []byte
	tx := types.NewTransaction(nonceAt, toAddress, value, gasLimit, gasPrice, data)
	fmt.Println("tx:", tx.Hash().Hex())
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
