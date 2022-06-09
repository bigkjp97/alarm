package main

/**
author: KJP
mark: 方法使用驼峰，变量使用下划线
*/

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
