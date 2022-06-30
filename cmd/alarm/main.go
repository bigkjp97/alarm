package main

/**
author: KJP
mark: 方法使用驼峰，变量使用下划线
*/

import (
	"alarm/apis/alarm"
	"time"
	// "alarm/apis/query"
	"alarm/pkg/conn"
	"alarm/pkg/utils"

	// embed time zone data
	"github.com/patrickmn/go-cache"
	_ "time/tzdata"
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
	// conn.TestFunc(p)

	dbconn, err := conn.NewDBConn(c)
	if err != nil {
		return
	}
	cache := cache.New(5*time.Minute, 10*time.Minute)
	sch := make(chan alarm.Status, c.App.Channel_size)

	server := alarm.Server{Cfg: c, Pool: p, Cache: cache, Dbconn: dbconn, Status_ch: sch}

	for {
		server.RunCheck()
		time.Sleep(10 * time.Second)
	}

	// conn.TestTables(c)

	// alarm.GetHoliday("20220615", dbconn)

	// var l query.Loki
	// var r query.Result
	// l.Query(&r,"http://10.89.5.130:3100/",`{job="varlogs"}`)
}
