package utils

import (
	"flag"
	"fmt"
	"os"
)

var version = "v2.0"

// 命令行结构
type manual struct {
	printVersion bool
	configFile   string
}

func (m *manual) initFlags() {
	// 重写Usage
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usagePrint("Alarm"))
	}

	// 存放命令行参数变量
	flag.BoolVar(&m.printVersion, "v", false, "Print this builds version information")
	flag.StringVar(&m.configFile, "c", "test", "yaml file to load")

	// 解析命令行
	flag.Parse()

	// 如果输入-v，则变量为true，判断后打印版本
	if m.printVersion {
		showVersion()
	}
}

// 显示版本
func showVersion() {
	fmt.Printf("Version: %s\n", version)
}

// 显示使用说明
func usagePrint(n string) string {
	return fmt.Sprintf(`
	 ------------------------------
	 Usage: %s [options...]

	 Options:
	 -c    Config file. (default: "conf/config.yaml")
	 -v    Show version and exit.

	 Example:

	   %s -c conf/config.yaml
	`, n, n)
}

// 注册命令行
func RegisterFlags() string {
	var m manual
	m.initFlags()
	return m.configFile
}
