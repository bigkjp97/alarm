package main

import (
	"alarm/pkg/utils"
	"flag"
	"fmt"
	"os"

	// embed time zone data
	_ "time/tzdata"
)

func init() {
	flag.Usage = func() {
		utils.ShowVersion()
		fmt.Fprint(os.Stderr, utils.Usage("Alarm"))
	}

	flag.Parse()
}

func main() {
	// Load config, merging config file and CLI flags
	// var config Config
	// if err := utils.DefaultUnmarshal(&config, os.Args[1:], flag.CommandLine); err != nil {
	// 	fmt.Println("Unable to parse config:", err)
	// 	os.Exit(1)
	// }

	// // Handle -version CLI flag
	// if config.printVersion {
	// 	fmt.Println("v2.0")
	// 	os.Exit(0)
	// }
	// f := flag.NewFlagSet("Alarm", flag.ExitOnError)
	// var m Manual

	// m.RegisterFlags(f, os.Args[2:])
}
