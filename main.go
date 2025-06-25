package main

import (
	"flag"
	"fmt"
	"github.com/olaola-chat/fsm_ctl/codegen"
	"log"
	"os"
)

func main() {
	exec()
}

func exec() {
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
