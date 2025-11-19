package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	// 加载.env文件
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	value := os.Getenv(key)
	if len(value) == 0 {
		panic("key not exist")
	}
	return value
}
