package main

import (
	"flag"
	"fmt"
	"fsm_ctl/codegen"
	"log"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

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
