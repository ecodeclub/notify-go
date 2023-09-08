package logger

import (
	"sync"
	"testing"
)

// 多个logger，根据日志级别选择打印到的日志文件
func Test_Multi(t *testing.T) {
	Default().Info("hello world", String("h", "1"))
	Default().Warn("hello world", String("j", "2"))
	Default().Error("hello world", String("k", "3"))
}

// 单个logger 并且是标准命令行输出
func Test_Std(t *testing.T) {
	ResetDefault(std)
	Default().Info("hello world", String("h", "231"))
}

func Test_Rotate(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(200000)
	for i := 0; i < 200000; i++ {
		go func() {
			wg := wg
			wg.Done()
			Default().Info("demo:", String("app", "start ok"), Int("major version", 3))
			Default().Error("demo:", String("app", "crash"), Int("reason", -1))
		}()
	}
	wg.Wait()
}
