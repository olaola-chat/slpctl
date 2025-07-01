package main

import (
	"flag"
	"fmt"
	"github.com/olaola-chat/slpctl/stategen"
	"log"
)

// 功能1: 示例功能A - 文件处理
type FunctionState struct {
	jsonFolder string
	jsonFile   string
	outputDir  string
}

func (f *FunctionState) InitArgs() {
	flag.StringVar(&f.jsonFolder, "j", "./rpc/server/internal/room_game/state/json", "游戏状态机的默认目录")
	flag.StringVar(&f.jsonFile, "f", "", "游戏状态机的默认配置文件名称")
	flag.StringVar(&f.outputDir, "o", "./rpc/server/internal/room_game", "输出目录")
}

func (f *FunctionState) Execute() error {
	if f.jsonFile == "" {
		return fmt.Errorf("源文件和目标文件路径必须指定")
	}

	jPath := fmt.Sprintf("%s/%s", f.jsonFolder, f.jsonFile)
	generator, nErr := stategen.NewGameGenerator(jPath, f.outputDir)
	if nErr != nil {
		log.Fatalf("生成失败了: %v", nErr)
		return nErr
	}
	if err := generator.Generate(); err != nil {
		log.Fatalf("生成失败: %v", err)
		return err
	}

	fmt.Printf("游戏代码已成功生成到目录: %s\n", f.outputDir)
	return nil
}

func (f *FunctionState) Help() {
	fmt.Println("功能: 状态机")
	fmt.Println("  描述: 用户快速生成状态机的基础代码")
	fmt.Println("  参数:")
	fmt.Println("    -j 状态机的配置json文件的目录")
	fmt.Println("    -f 游戏状态机的默认配置文件名称")
	fmt.Println("    -o 状态机代码输出目录")
}
