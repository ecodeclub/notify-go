package tool

import (
	"math/rand"
	"time"
)

// RandIntN 生成指定范围内的随机整数
func RandIntN(min, max int) int {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(max-min+1) + min

	return randomNumber
}
