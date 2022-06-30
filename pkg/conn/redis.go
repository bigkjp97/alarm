package conn

import (
	"alarm/pkg/utils"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 测试空闲连接方法
func testOnBorrowFunc(conn redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := conn.Do("PING")
	return err
}

// 建立连接方法
func dialFunc(c *utils.Config) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", c.Redis.Address)
	if err != nil {
		return nil, err
	}

	// Redis验证
	if c.Redis.Auth != "" {
		if _, err := conn.Do("AUTH", c.Redis.Auth); err != nil {
			conn.Close()
			return nil, err
		}
	}

	// 选择数据库0
	if _, err := conn.Do("SELECT", "0"); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// 建立连接池
func NewRedisConn(c *utils.Config) (*redis.Pool, error) {
	// go引用结构体指针
	pool := &redis.Pool{
		// 最大空闲连接数
		MaxIdle: c.Redis.Max_idle,

		// 最大空闲连接等待时间
		IdleTimeout: c.Redis.Idle_timeout * time.Second,

		// 测试空闲连接的健康（空闲超过1分钟）
		TestOnBorrow: testOnBorrowFunc,

		// 建立连接
		Dial: func() (redis.Conn, error) { return dialFunc(c) },
	}

	return pool, nil
}

func TestFunc(p *redis.Pool) {
	con := p.Get()
	defer con.Close()

	_, err := con.Do("SET", "TESTALARM", 10000)
	if err != nil {
		fmt.Println(err)
	}

	res, err := redis.Int(con.Do("GET", "TESTALARM"))
	if err != nil {
		fmt.Println("GET TESTALARM FAILED: ", err)
	}

	fmt.Println(res)
	p.Close()
}
