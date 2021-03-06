package alarm

import (
	// "alarm/pkg/utils"
	"encoding/json"
	"fmt"

	// "fmt"

	"github.com/gomodule/redigo/redis"
	// "github.com/patrickmn/go-cache"
)

type Msg struct {
	AlarmLog
	Wiki AlarmWiki `json:"wiki"`
}

// 状态缓存保存到redis
func (server *Server) statusToRedis() {
	// var wg sync.WaitGroup
	// defer wg.Done()
	fmt.Println("正在缓存当前状态...")
	rc := server.Pool.Get()
	defer rc.Close()
	for status := range server.Status_ch {
		jstr, _ := json.Marshal(status)
		if _, err := rc.Do("HSET", server.Cfg.Redis.Stats_name, status.Code, jstr); err != nil {
			// todo: 错误重试
			// log.Error("func: srv.cache2Redis(); redis hset:%s error:%v", status.Code, err)
			fmt.Println("缓存状态失败")
		}
	}
}

// 从redis读取状态缓存
func (server *Server) statusFromRedis(k string) (Status, error) {
	// defer wg.Done()
	var status Status

	rc := server.Pool.Get()
	defer rc.Close()

	val, err := redis.Bytes(rc.Do("HGET", server.Cfg.Redis.Stats_name, k))
	if err != nil {
		// todo: 错误重试
		// log.Error("func: srv.cacheFromRedis(); redis hget:%s error:%v", k, err)
		fmt.Println("空哈希")
		return status, err
	}

	err = json.Unmarshal(val, &status)
	fmt.Println("有状态")
	return status, err
}

// // 告警通知记录日志
// func alarmToRedis() {
// 	defer wg.Done()
// 	rc := srv.RedisCli.Get()
// 	defer rc.Close()

// 	for jstr := range srv.alarmCh {
// 		log.Info("func: srv.alarm2Redis(); send notice to redis: %v", jstr)
// 		if _, err := rc.Do("LPUSH", conf.Redis.QueueName, jstr); err != nil {
// 			// todo: 错误重试
// 			log.Error("func: srv.alarm2Redis(); redis lpush error:%v", err)
// 		}
// 		rc.Do("INCR", conf.Redis.TouchName)
// 	}
// }

// // 告警通知记录日志
// func alarmToDB(l AlarmLog) error {
// 	if result := srv.DbConn.Table(l.TableName()).Create(&l); result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }
