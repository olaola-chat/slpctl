package main

import (
	"flag"
	"fmt"
	"os"
)

// FunctionMap 存储所有可用功能
var FunctionMap = make(map[string]Function)

type Function interface {
	// 初始化功能的参数
	InitArgs()
	// 执行功能
	Execute() error
	// 显示功能帮助信息
	Help()
}

// 初始化所有功能
func initFunctions() {
	FunctionMap["state"] = &FunctionState{}
	FunctionMap["codec"] = &FunctionCodec{}
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
	function.InitArgs()

	// 重新解析命令行参数，跳过第一个参数(功能类型)
	err := flagset.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("参数解析错误: %v functionType=%v \n", err, functionType)
		os.Exit(1)
	}

	// 显示帮助信息
	if flagset.NArg() > 0 && (flagset.Arg(0) == "-h" || flagset.Arg(0) == "--help") {
		function.Help()
		os.Exit(0)
	}

	// 执行功能
	err = function.Execute()
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("功能执行完成")
}

// 显示主帮助信息
func showMainHelp() {
	fmt.Println("多功能工具 usage:")
	fmt.Println("  slpctl state [function-args]")
	fmt.Println("")
	fmt.Println("可用功能:")
	for name := range FunctionMap {
		fmt.Printf("  %s\n", name)
	}
	fmt.Println("")
	fmt.Println("使用 'slpctl tp -h' 查看具体功能的帮助信息")
}
