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
	flag.StringVar(&f.jsonFolder, "j", "./rpc/server/internal/room_game/state/json", "状态机的配置json文件的目录")
	flag.StringVar(&f.jsonFile, "f", "", "游戏状态机的配置文件名称")
	flag.StringVar(&f.outputDir, "o", "./rpc/server/internal/room_game", "状态机代码输出目录")
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
	fmt.Println("功能: 状态机代码生成")
	fmt.Println("  描述: 基于JSON配置生成游戏状态机基础代码")
	fmt.Println("  参数:")
	fmt.Println("    -j <目录>  状态机配置JSON文件的目录 (默认: ./rpc/server/internal/room_game/state/json)")
	fmt.Println("    -f <文件>  状态机配置文件名 (必须指定)")
	fmt.Println("    -o <目录>  生成代码的输出目录 (默认: ./rpc/server/internal/room_game)")
	fmt.Println("  示例:")
	fmt.Println("    slpctl state -f game_state.json -o ./output")
}
