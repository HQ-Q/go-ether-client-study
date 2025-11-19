package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//queryBlockInfo()
	transferEth()
}

// 给指定用户转账eth
func transferEth() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxx")
	if err != nil {
		fmt.Println("连接失败", err)
		return
	}
	defer client.Close()

	//私钥
	privateKey, err := crypto.HexToECDSA("xxxxx")
	if err != nil {
		fmt.Println("私钥转换失败", err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("转换公钥失败")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//获取最新nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("获取nonce失败", err)
		return
	}
	fmt.Println("最新nonce:", nonce)

	//gasLimit
	gasLimit := uint64(21000)
	//计算gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("获取gasPrice失败", err)
		return
	}
	fmt.Println("gasPrice:", gasPrice.String())
	//收款地址
	toAddress := common.HexToAddress("0xac787ff5df204282fc4a9216e2c5e5fc3d703574")
	amount := big.NewInt(664000000000000000) //1 eth
	//构建交易
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("获取chainID失败", err)
		return
	}
	//签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("交易签名失败", err)
		return
	}
	//发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("发送交易失败", err)
		return
	}
	fmt.Printf("交易发送成功！TxHash：%s\n", signedTx.Hash().Hex())
}

// 查询区块信息
func queryBlockInfo() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/qLVks-ioeg_c5r6FDc8VT")
	if err != nil {
		fmt.Println("连接失败", err)
		return
	}
	defer client.Close()

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		fmt.Println("获取区块失败", err)
		return
	}
	fmt.Println("区块号：", block.Number().Uint64())
	fmt.Println("区块哈希：", block.Hash().Hex())
	fmt.Println("区块大小：", block.Size())
	fmt.Println("区块时间：", block.Time())
	fmt.Println("区块交易数：", block.Transactions().Len())
}
