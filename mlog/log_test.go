package mlog

import (
	"fmt"
	"sync"
	"testing"
)

var Log *Logger

// go test -v -run TestNewLog
func TestNewLog(t *testing.T) {
	Log = New(Config{
		Path: "./mo7-logs",
		Name: "m",
	})
	go Log.Info("this is info") // 并发写入测试，其实写日志建立 goroutine 开销可能更大

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		Log.Info("this is info")
	}()
	go func() {
		defer wg.Done()
		Log.Warn("this is warn")
	}()
	go func() {
		defer wg.Done()
		Log.Error("this is error")
	}()
	go func() {
		defer wg.Done()
		Log.Debug("this is debug")
	}()

	fmt.Println("等待所有任务完成...")
	wg.Wait()
	fmt.Println("所有任务完成。")
}

// go test -v -run TestClearLog
func TestClearLog(t *testing.T) {
	Log = New(Config{
		Path: "./mo7-logs",
		Name: "m-logtest",
	})

	Log.Clear(ClearOpt{
		Type:   []string{"debug", "warn"},
		Before: 7,
	})
}
