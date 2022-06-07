package main

import (
	"alarm/pkg/utils"

	// embed time zone data
	_ "time/tzdata"
)

func init() {

}

func main() {
	cfg := utils.RegisterFlags()
	utils.ParseConfig(cfg)
}
