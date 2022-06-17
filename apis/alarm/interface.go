package alarm

import (
	"alarm/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
)

type Msg struct {
	AlarmLog
	Wiki AlarmWiki `json:"wiki"`
}

// 告警状态
type Status struct {
	Code            string    // 告警编号
	Last_check_time time.Time `json:"last_check_time"` // 最后查询时间
	Last_alarm_time time.Time `json:"last_alarm_time"` // 最后告警时间
	Alarm_count     int64     `json:"alarm_count"`     // 告警计数
	Send_count      int64     `json:"send_count"`      // 发送计数
}

// 更新状态缓存
func setStatus(cache *cache.Cache, s Status, sch chan Status, exp time.Duration) {
	// 全局的管道
	cache.Set(s.Code, s, exp)
	sch <- s
}

// 获取状态缓存
func getStatus(c *utils.Config, r *redis.Pool, cache *cache.Cache, k string) (Status, error) {
	fmt.Println(k)
	if val, found := cache.Get(k); found {
		fmt.Println(found)
		fmt.Println(val.(Status))
		return val.(Status), nil
	}

	return statusFromRedis(c, r, k)
}

// 状态缓存保存到redis
func statusToRedis(c *utils.Config, r *redis.Pool, s_ch chan Status) {
	// var wg sync.WaitGroup
	// defer wg.Done()
	rc := r.Get()
	defer rc.Close()
	for status := range s_ch {
		jstr, _ := json.Marshal(status)
		if _, err := rc.Do("HSET", c.Redis.Stats_name, status.Code, jstr); err != nil {
			// todo: 错误重试
			// log.Error("func: srv.cache2Redis(); redis hset:%s error:%v", status.Code, err)
		}
	}
}

// 从redis读取状态缓存
func statusFromRedis(c *utils.Config, r *redis.Pool, k string) (Status, error) {
	// defer wg.Done()
	var status Status

	rc := r.Get()
	defer rc.Close()

	val, err := redis.Bytes(rc.Do("HGET", c.Redis.Stats_name, k))
	if err != nil {
		// todo: 错误重试
		// log.Error("func: srv.cacheFromRedis(); redis hget:%s error:%v", k, err)
		return status, err
	}

	err = json.Unmarshal(val, &status)
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