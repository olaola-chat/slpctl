package main

import (
	"flag"
	"fmt"
	"github.com/olaola-chat/slpctl/helmgen"
)

// 功能1: 示例功能A - 文件处理
type FunctionHelm struct {
	jsonFolder string
	jsonFile   string
	outputDir  string
}

func (f *FunctionHelm) InitArgs(flagset *flag.FlagSet) {
}

func (f *FunctionHelm) Execute() error {
	helmgen.GenHelmCode()
	return nil
}

func (f *FunctionHelm) Help() {
	fmt.Println("功能: 快速生成项目重启的helm中的tag命令")
}
