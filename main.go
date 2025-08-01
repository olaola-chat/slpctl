package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// FunctionMap 存储所有可用功能
var FunctionMap = make(map[string]Function)

type Function interface {
	// 初始化功能的参数
	InitArgs(flagset *flag.FlagSet)
	// 执行功能
	Execute() error
	// 显示功能帮助信息
	Help()
}

// 初始化所有功能
func initFunctions() {
	FunctionMap["state"] = &FunctionState{}
	FunctionMap["codec"] = &FunctionCodec{}
	FunctionMap["helm"] = &FunctionHelm{}
	// 可以在这里添加更多功能
}

func main() {
	// 初始化所有功能
	initFunctions()

	// 检查参数数量
	if len(os.Args) < 2 {
		showMainHelp()
		os.Exit(1)
	}

	// 获取功能类型
	functionType := os.Args[1]
	function, exists := FunctionMap[functionType]
	if !exists {
		fmt.Printf("未知功能类型: %s\n\n", functionType)
		showMainHelp()
		os.Exit(1)
	}

	// 处理功能特定的参数
	fmt.Printf("执行功能: %s\n", functionType)
	flagset := flag.NewFlagSet(functionType, flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Printf("Usage of %s:\n", functionType)
		flagset.PrintDefaults()
		function.Help()
	}

	function.InitArgs(flagset)

	// 解析命令行参数（跳过第一个参数: 功能类型）
	err := flagset.Parse(os.Args[2:])
	if err != nil {
		// 处理帮助请求
		if err == flag.ErrHelp {
			function.Help()
			os.Exit(0)
		}
		fmt.Printf("参数错误: %v\n", err)
		flagset.Usage()
		os.Exit(1)
	}

	// 执行功能
	if err = function.Execute(); err != nil {
		log.Fatalf("执行失败: %v", err)
		os.Exit(1)
	}
	return
}

// 显示主帮助信息
func showMainHelp() {
	fmt.Println("slpctl 工具 usage:")
	fmt.Println("  slpctl <功能> [参数...]")
	fmt.Println("")
	fmt.Println("可用功能:")
	for name := range FunctionMap {
		fmt.Printf("  %s - \n", name)
	}
	fmt.Println("")
	fmt.Println("使用 'slpctl <功能> -h' 查看具体功能的帮助信息")
}
