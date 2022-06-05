package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	// embed time zone data
	_ "time/tzdata"
)

type Config struct {
	config.Config   `yaml:",inline"`
	printVersion    bool
	printConfig     bool
	logConfig       bool
	dryRun          bool
	configFile      string
	configExpandEnv bool
	inspect         bool
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&c.printVersion, "version", false, "Print this builds version information")
	f.BoolVar(&c.printConfig, "print-config-stderr", false, "Dump the entire Alarm config object to stderr")
	f.BoolVar(&c.logConfig, "log-config-reverse-order", false, "Dump the entire Alarm config object at Info log "+
		"level with the order reversed, reversing the order makes viewing the entries easier in Grafana.")
	f.BoolVar(&c.inspect, "inspect", false, "Allows for detailed inspection of pipeline stages")
	f.StringVar(&c.configFile, "config.file", "", "yaml file to load")
	f.BoolVar(&c.configExpandEnv, "config.expand-env", false, "Expands ${var} in config according to the values of the environment variables.")
	c.Config.RegisterFlags(f)
}

func main() {
	// Load config, merging config file and CLI flags
	var config Config
	if err := cfg.DefaultUnmarshal(&config, os.Args[1:], flag.CommandLine); err != nil {
		fmt.Println("Unable to parse config:", err)
		os.Exit(1)
	}

	// Handle -version CLI flag
	if config.printVersion {
		fmt.Println(version.Print("alarm"))
		os.Exit(0)
	}

	// Init the logger which will honor the log level set in cfg.Server
	if reflect.DeepEqual(&config.Config.ServerConfig.Config.LogLevel, &logging.Level{}) {
		fmt.Println("Invalid log level")
		os.Exit(1)
	}
	util_log.InitLogger(&config.Config.ServerConfig.Config, prometheus.DefaultRegisterer)

}
