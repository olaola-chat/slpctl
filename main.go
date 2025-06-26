package main

import (
	"flag"
	"fmt"
	"github.com/olaola-chat/slpctl/codegen"
	"log"
	"os"
)

func main() {
	opType := flag.String("t", "state", "游戏状态机模板生成")
	switch *opType {
	case "state":
		stateExec()
	}
}

func stateExec() {
	configPath := flag.String("c", "test.json", "游戏的json配置，请参考test.json")
	outputDir := flag.String("o", "./app", "输出目录")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	generator, _ := codegen.NewGameGenerator(*configPath, *outputDir)
	if err := generator.Generate(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("游戏代码已成功生成到目录: %s\n", *outputDir)
}
