package utils

import (
	"flag"
	"fmt"
)

var version = "v2.0"

type Manual struct {
	printVersion string
	configFile   string
}

func (c *Manual) RegisterFlags(f *flag.FlagSet, a []string) {
	f.StringVar(&c.printVersion, "v", version, "Print this builds version information")
	f.StringVar(&c.configFile, "c", "", "yaml file to load")
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

func ShowVersion() {
	fmt.Printf("Version: %s\n", version)
}
