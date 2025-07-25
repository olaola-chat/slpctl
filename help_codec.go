package main

import (
	"flag"
	"fmt"
	"github.com/olaola-chat/slpctl/codecgen"
)

// 功能1: 示例功能A - 文件处理
type FunctionCodec struct {
	tablename        string
	s                int64
	h                int64
	d                string
	primaryAlisaName string
	m                string
}

// flag.String("m", "slp", "给个项目的go.mod的包名")
func (f *FunctionCodec) InitArgs(flagset *flag.FlagSet) {
	fmt.Println("FunctionCodec.InitArgs")
	flagset.StringVar(&f.tablename, "t", "", "会根据这个表明生成对应的cache文件")
	flagset.Int64Var(&f.s, "s", 0, "cache 的缓存过期时间，单位s")
	flagset.Int64Var(&f.h, "h", 0, "cache 的缓存过期时间，单位小时,默认3")
	flagset.StringVar(&f.d, "d", "passive", "redis的那个模块的db,按业务区分。目前提供 story,property,block,user...")
	flagset.StringVar(&f.primaryAlisaName, "uq", "id", "默认id，但你的表如果唯一索引锁uid，这里你就可以用uid")
	flagset.StringVar(&f.m, "m", "slp", "给个项目的go.mod的包名")
}

func (f *FunctionCodec) Execute() error {
	if f.tablename == "" {
		return fmt.Errorf("-t 不能为空;会根据这个表明生成对应的cache文件")
	}
	//pTableName string, ps, ph int64, pd string, pName, pmode string
	codecgen.CodecExec(f.tablename, f.s, f.h, f.d, f.primaryAlisaName, f.m)
	return nil
}

func (f *FunctionCodec) Help() {
	fmt.Println("功能: 状态机")
	fmt.Println("  描述: 用户快速生成状态机的基础代码")
	fmt.Println("  参数:")
	fmt.Println("    -j 状态机的配置json文件的目录")
	fmt.Println("    -f 游戏状态机的默认配置文件名称")
	fmt.Println("    -o 状态机代码输出目录")
}
