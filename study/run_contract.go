package main

import (
	"context"
	"crypto/ecdsa"
	"eth-client-study/study/store"
	"eth-client-study/utils"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	runSetItemsByAbi()
}

// 通过abi调用
func runSetItemsByAbi() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	privateKey, err := crypto.HexToECDSA(utils.GetEnv("PRIVATE_KEY1"))
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//获取最新的nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	//估算gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//准备交易 calldata
	abiStr := "[\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_version\",\n        \"type\": \"string\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"constructor\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes32\",\n        \"name\": \"key\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes32\",\n        \"name\": \"value\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"ItemSet\",\n    \"type\": \"event\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"items\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"key\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"value\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"setItem\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"version\",\n    \"outputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"\",\n        \"type\": \"string\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  }\n]"
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatal(err)
	}
	//函数名称
	methodName := "setItem"
	var key [32]byte
	var value [32]byte

	copy(key[:], "demo_save_key_use_abi")
	copy(value[:], "demo_save_value_use_abi_11111")
	input, err := contractABI.Pack(methodName, key, value)

	tx := types.NewTransaction(nonce, common.HexToAddress("0x183AdfEe585d04Db1Ab151840D6399009beC2bC4"), big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(11155111)), privateKey)
	if err != nil {
		log.Fatal("交易签名失败：", err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败：", err)
	}
	fmt.Printf("交易发送成功！TxHash：%s\n", signedTx.Hash().Hex())

	_, err = waitForReceipt2(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// 查询刚刚设置的值
	callInput, err := contractABI.Pack("items", key)
	if err != nil {
		log.Fatal(err)
	}

	to := common.HexToAddress("0x183AdfEe585d04Db1Ab151840D6399009beC2bC4")
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	// 解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result:", string(result))

	var unpacked [32]byte
	contractABI.UnpackIntoInterface(&unpacked, "items", result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("is value saving in contract equals to origin value:", unpacked == value)

}

func runSetItemsNoUseAbi() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	privateKey, err := crypto.HexToECDSA(utils.GetEnv("PRIVATE_KEY1"))
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//获取最新的nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	//估算gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	methodSignature := []byte("setItem(bytes32,bytes32)")
	methodSelector := crypto.Keccak256(methodSignature)[:4]

	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("demo_save_key_no_use_abi"))
	copy(value[:], []byte("demo_save_value_no_use_abi_11111"))

	// 组合调用数据
	var input []byte
	input = append(input, methodSelector...)
	input = append(input, key[:]...)
	input = append(input, value[:]...)

	tx := types.NewTransaction(nonce, common.HexToAddress("0x183AdfEe585d04Db1Ab151840D6399009beC2bC4"), big.NewInt(0), 300000, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(11155111)), privateKey)
	if err != nil {
		log.Fatal("交易签名失败：", err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败：", err)
	}
	fmt.Printf("交易发送成功！TxHash：%s\n", signedTx.Hash().Hex())

	_, err = waitForReceipt2(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	itemsSignature := []byte("items(bytes32)")
	itemsSelector := crypto.Keccak256(itemsSignature)[:4]

	var callInput []byte
	callInput = append(callInput, itemsSelector...)
	callInput = append(callInput, key[:]...)

	to := common.HexToAddress("0x183AdfEe585d04Db1Ab151840D6399009beC2bC4")
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	var unpacked [32]byte
	copy(unpacked[:], result)
	fmt.Println("is value saving in contract equals to origin value:", unpacked == value)

}

func waitForReceipt2(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}

func runSetItems() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	contractAddr := "0x183AdfEe585d04Db1Ab151840D6399009beC2bC4"
	storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}
	version, err := storeContract.Version(&bind.CallOpts{})
	if err != nil {

		log.Fatal(err)
	}
	log.Println("version:", version)

	//私钥
	privateKey, err := crypto.HexToECDSA(utils.GetEnv("PRIVATE_KEY1"))
	if err != nil {
		log.Fatal(err)
	}
	var key [32]byte
	var value [32]byte
	copy(key[:], "demo_save_key")
	copy(value[:], "demo_save_value11111")
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tx:", tx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("receipt:", receipt.Status)
	log.Println("等待交易确认...")
	for receipt.Status != 1 {
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("交易已确认")
	log.Println("等待交易确认结束")

	//读取items
	items, err := storeContract.Items(&bind.CallOpts{}, key)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("items:", items)
	log.Println("交易结束")
	//字节数组转字符串
	log.Println("items:", string(items[:]))
}
