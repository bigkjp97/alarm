package main

import (
	"alarm/pkg/utils"

	// embed time zone data
	_ "time/tzdata"
)

func init() {

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
	utils.RegisterFlags()

}
