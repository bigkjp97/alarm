package main

/**
author: KJP
mark: 方法使用驼峰，变量使用下划线
*/

import (
	"alarm/pkg/conn"
	"alarm/pkg/utils"

	// embed time zone data
	_ "time/tzdata"
)

func init() {

}

func main() {
	cfg := utils.RegisterFlags()

	// 解析配置文件
	c := utils.ParseConfig(&utils.Config{}, cfg)

	// redis连接
	conn.NewRedisConn(c)
}
