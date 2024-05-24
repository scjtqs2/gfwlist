// Package 入口
package main

import (
	"bytes"
	"flag"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	rules "github.com/Dreamacro/clash/rule"
	conf2 "github.com/scjtqs2/gfwlist/conf"
	"github.com/scjtqs2/gfwlist/gfw"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	h bool
	// debug   bool
	d       bool
	Version string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&d, "d", true, "是否每次运行都重新下载gfwlist.txt")
	flag.StringVar(&gfw.GfwlistUrl, "gfw", "https://pagure.io/gfwlist/raw/master/f/gfwlist.txt", "gfwlist download url for https://github.com/gfwlist/gfwlist")
	flag.Parse()
}

// Help cli命令行-h的帮助提示
func Help() {
	fmt.Printf(`clash 规则添加gfwlist命令工具
version: %s
Usage:
server [OPTIONS]
Options:
`, Version)
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	if h {
		Help()
	}
	if d {
		_ = os.RemoveAll(gfw.Gfwlist)
	}
	defer os.RemoveAll(gfw.Gfwlist)
	domainList, err := gfw.LoadGfwList()
	if err != nil {
		panic(err)
	}
	// 生成clash规则
	makeCashRule(domainList)
	// 生成qx规则
	makeQuantumultXRule(domainList)
}

// 写入 clash的规则
func makeCashRule(domainList []string) {

	tmpRules := make([]C.Rule, 0)
	for _, s := range domainList {
		tmpRules = append(tmpRules, rules.NewDomainSuffix(s, "GFW"))
	}
	rRule := conf2.TransRule(tmpRules)
	buf, _ := yaml.Marshal(rRule)
	conf := bytes.NewBufferString("payload:\n")
	conf.Write(buf)
	err := writefile("Rules/Clash/gfwlist.yml", conf.Bytes())
	if err != nil {
		panic(err)
	}
}

func makeQuantumultXRule(domainList []string) {
	conf := bytes.NewBuffer(nil)
	for _, s := range domainList {
		conf.WriteString(fmt.Sprintf("HOST-SUFFIX,%s,GFWLIST\n", s))
	}
	err := writefile("Rules/QuantumultX/gfwlist.conf", conf.Bytes())
	if err != nil {
		panic(err)
	}
}

// writefile 写入文件
func writefile(file string, buf []byte) error {
	return os.WriteFile(file, buf, 0o644)
}
