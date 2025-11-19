package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// 转移20代币
func main() {

	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	//账户私钥
	privateKey, err := crypto.HexToECDSA("xxxxxxxx")
	if err != nil {
		log.Fatal("privateKey err:", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	//接收地址
	toAddress := common.HexToAddress("0xac787ff5df204282fc4a9216e2c5e5fc3d703574")
	//合约地址
	tokenAddress := common.HexToAddress("0x2f8C29909a2697E4E0449662302aAa1750f2cF98")

	//转账签名
	methodID := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(methodID)
	//转账函数hash
	transferFnSig := hash.Sum(nil)[:4]
	fmt.Println("transferFnSig:", hexutil.Encode(transferFnSig)) // 0xa9059cbb
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println("paddedAddress:", hexutil.Encode(paddedAddress)) //
	amount := new(big.Int)
	amount.SetString("1000000000000000000000000000", 10) // 十亿 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println("paddedAmount:", hexutil.Encode(paddedAmount)) //
	data := append(transferFnSig, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal("Gas估算失败：", err)
	}
	fmt.Println("估算GasLimit：", gasLimit)

	// 7. 获取推荐GasPrice（Legacy交易用固定GasPrice）
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取GasPrice失败：", err)
	}
	fmt.Println("推荐GasPrice：", gasPrice.String(), "wei")
	// 8. 构造交易
	// 注意：To是代币合约地址，Value是0（ERC20转账不转ETH）
	tx := types.NewTransaction(
		nonce,         // 交易序号
		tokenAddress,  // 目标地址（代币合约）
		big.NewInt(0), // 转账ETH金额（ERC20转账设为0）
		gasLimit+1000, // GasLimit（估算值+1000缓冲，避免Gas不足）
		gasPrice,      // 固定GasPrice（Legacy交易）
		data,          // 合约调用数据
	)
	// 9. 签名交易
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		chainID = big.NewInt(11155111) // 降级使用硬编码链ID
		log.Println("获取链ID失败：", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("交易签名失败：", err)
	}

	// 10. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败：", err)
	}

	fmt.Printf("交易发送成功！TxHash：%s\n", signedTx.Hash().Hex())
}
