package main

import (
	"fmt"

	whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"
)

// whisper 初步使用
// 什么也不用添加
func main() {
	// 使用默认的配置
	wsp := whisper.New(nil)
	if wsp == nil {
		fmt.Println("Fail to create whisper")
	}

	// 空启动
	wsp.Start(nil)
	defer wsp.Stop()
}
