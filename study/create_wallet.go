// create_wallet.go - 以太坊钱包密钥对生成示例程序
package main

import (
	"crypto/ecdsa" // 提供椭圆曲线数字签名算法支持
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil" // 提供十六进制编码和解码功能
	"github.com/ethereum/go-ethereum/crypto"         // 提供以太坊加密相关功能
	"golang.org/x/crypto/sha3"
)

// main 函数是程序入口点，演示如何生成以太坊钱包的私钥和公钥
func main() {
	// 使用 crypto.GenerateKey() 生成一个新的椭圆曲线私钥
	// 该私钥基于 secp256k1 曲线，是以太坊默认的加密算法
	privateKey, err := crypto.GenerateKey() // 生成私钥
	if err != nil {
		panic(err) // 如果生成失败则抛出异常
	}

	// 将私钥转换为字节切片格式，便于后续编码处理
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// 使用 hexutil.Encode 将字节切片编码为十六进制字符串
	// [2:] 用于去除 "0x" 前缀，只显示纯十六进制内容
	fmt.Println("私钥：", hexutil.Encode(privateKeyBytes)[2:])

	// 通过 privateKey.Public() 获取与私钥对应的公钥
	publicKey := privateKey.Public() // 获取公钥

	// 类型断言：将接口类型的 publicKey 转换为具体的 *ecdsa.PublicKey 类型
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		// 如果类型断言失败，则输出错误信息并终止程序
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 将公钥转换为字节切片格式
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// 输出公钥的十六进制表示，[4:] 去除前缀（通常是 0x04）
	// 0x04 是未压缩公钥的标识符
	fmt.Println("公钥：", hexutil.Encode(publicKeyBytes)[4:])

	// 使用 `crypto.PubkeyToAddress` 将公钥转换为以太坊地址
	// 这是官方推荐的方法，会自动处理地址生成的所有步骤
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address:", address)

	// 手动计算地址：使用 Keccak256 哈希算法计算公钥哈希
	// 创建一个新的 Keccak256 哈希实例
	hash := sha3.NewLegacyKeccak256()
	// 写入公钥数据（跳过第一个字节，即未压缩公钥标识符0x04）
	hash.Write(publicKeyBytes[1:])
	// 计算哈希值并取最后20字节作为地址（以太坊地址长度为20字节）
	// [12:] 表示跳过前12字节，保留后20字节（32-12=20）
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}
