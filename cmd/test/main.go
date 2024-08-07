package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 定义一个 defer 方法
	defer func() {
		fmt.Println("Deferred function executed")
	}()

	// 创建一个通道来接收信号通知
	sigs := make(chan os.Signal, 1)

	// 通知接收特定的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 使用 goroutine 来处理信号
	go func() {
		sig := <-sigs
		fmt.Println("Received signal:", sig)
		os.Exit(0)
	}()

	fmt.Println("Running... Press Ctrl+C to exit")
	for {
		time.Sleep(1 * time.Second)
	}
}
