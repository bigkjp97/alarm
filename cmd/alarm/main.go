package main

/**
author: KJP
mark: 方法使用驼峰，变量使用下划线
*/

import (
	"time"
	// "alarm/apis/alarm"
	"alarm/apis/alarm"
	"alarm/pkg/conn"
	"alarm/pkg/utils"

	// embed time zone data
	_ "time/tzdata"

	"github.com/patrickmn/go-cache"
)

func init() {

}

// 启动告警程序
func StartAlarm() {
	// 先判断是否是节假日

	// 获取查询接口

	// 获取监控项

	// 判断告警
}

func main() {
	cfg := utils.RegisterFlags()

	// 解析配置文件
	c := utils.ParseConfig(&utils.Config{}, cfg)

	// redis连接
	p, err := conn.NewRedisConn(c)
	if err != nil {
		return
	}
	conn.TestFunc(p)

	dbconn, err := conn.NewDBConn(c)
	if err != nil {
		return
	}
	cache := cache.New(5*time.Minute, 10*time.Minute)

	var sch chan alarm.Status
	server := alarm.Server{Cfg: c, Pool: p, Cache: cache, Dbconn: dbconn}
	server.RunCheck(sch)

	// conn.TestTables(c)

	// alarm.GetHoliday("20220615", dbconn)
}
