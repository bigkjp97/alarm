package conn

import (
	"alarm/pkg/utils"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewRedisConn() *redis.Pool {
	var c utils.Config
	return &redis.Pool{
		MaxIdle:     c.Redis.Max_idle,                   // 最大的空闲连接数
		IdleTimeout: c.Redis.Idle_timeout * time.Second, // 最大的空闲连接等待时间

		// 建立连接
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Redis.Network, c.Redis.Address)
			if err != nil {
				log.Fatalf("redis connection error %v", err)
			}
			if c.Redis.Auth != "" {
				if _, err := conn.Do("AUTH", c.Redis.Auth); err != nil {
					conn.Close()
					log.Fatalf("redis connection error %v", err)
					return nil, err
				}
			}
			conn.Do("SELECT", c.Redis.Db_select) // 选择数据库
			return conn, nil
		},

		// 测试连接可用性
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}
