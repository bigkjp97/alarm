package utils

import (
	"flag"
	"fmt"
	"os"
)

var version = "v2.0"

type manual struct {
	printVersion string
	configFile   string
}

func (m *manual) InitFlags() {
	flag.Usage = func() {
		ShowVersion()
		fmt.Fprint(os.Stderr, Usage("Alarm"))
	}

	flag.StringVar(&m.printVersion, "v", version, "Print this builds version information")
	flag.StringVar(&m.configFile, "c", "", "yaml file to load")

	flag.Parse()
}

func RegisterFlags() {
	var m manual
	m.InitFlags()
}

func ShowVersion() {
	fmt.Printf("Version: %s\n", version)
}

func Usage(n string) string {
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
