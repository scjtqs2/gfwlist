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
	// 生成surge规则
	makeSurgeRule(domainList)
	makeSurfboardRule(domainList)
	makeSsRule(domainList)
	makeQuantumultRule(domainList)
}

// 写入 clash的规则
func makeCashRule(domainList []string) {

	tmpRules := make([]C.Rule, 0)
	for _, s := range domainList {
		tmpRules = append(tmpRules, rules.NewDomainSuffix(s, "🈲 GFW"))
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

// makeQuantumultXRule 生成qx用的gfw规则
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

// makeSurgeRule 生成surge的gfw规则
func makeSurgeRule(domainList []string) {
	conf := bytes.NewBuffer(nil)
	for _, s := range domainList {
		conf.WriteString(fmt.Sprintf("DOMAIN-SUFFIX,%s\n", s))
	}
	err := writefile("Rules/Surge/gfwlist.list", conf.Bytes())
	if err != nil {
		panic(err)
	}
}

// makeSurfboardRule 生成 surfboard的规则
func makeSurfboardRule(domainList []string) {
	conf := bytes.NewBuffer(nil)
	for _, s := range domainList {
		conf.WriteString(fmt.Sprintf("DOMAIN-SUFFIX,%s\n", s))
	}
	err := writefile("Rules/Surfboard/gfwlist.conf", conf.Bytes())
	if err != nil {
		panic(err)
	}
}

// makeSsRule 生成ss的订阅规则
func makeSsRule(domainList []string) {
	preText := `[General]
bypass-system = true
skip-proxy = 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, localhost, *.local, e.crashlytics.com, captive.apple.com
bypass-tun = 10.0.0.0/8,100.64.0.0/10,127.0.0.0/8,169.254.0.0/16,172.16.0.0/12,192.0.0.0/24,192.0.2.0/24,192.88.99.0/24,192.168.0.0/16,198.18.0.0/15,198.51.100.0/24,203.0.113.0/24,224.0.0.0/4,255.255.255.255/32
dns-server = system, 114.114.114.114, 112.124.47.27, 8.8.8.8, 8.8.4.4

[Rule]
`
	conf := bytes.NewBufferString(preText)
	for _, s := range domainList {
		conf.WriteString(fmt.Sprintf("DOMAIN-SUFFIX,%s,PROXY\n", s))
	}
	err := writefile("Rules/Shadowrocket/gfwlist.conf", conf.Bytes())
	if err != nil {
		panic(err)
	}
}

// makeQuantumultRule 生成 quantumult 的规则
func makeQuantumultRule(domainList []string) {
	conf := bytes.NewBuffer(nil)
	for _, s := range domainList {
		conf.WriteString(fmt.Sprintf("DOMAIN-SUFFIX,%s,选择GFWLIST的策略\n", s))
	}
	err := writefile("Rules/Quantumult/gfwlist.conf", conf.Bytes())
	if err != nil {
		panic(err)
	}
}

// writefile 写入文件
func writefile(file string, buf []byte) error {
	return os.WriteFile(file, buf, 0o644)
}
