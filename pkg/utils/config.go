package utils

import (
	"github.com/jinzhu/configor"
	"log"
	"time"
)

type config struct {
	app struct {
		service_name string        `default:"app"`
		log_file     string        `default:"app"`
		log_level    string        `default:"INFO"`
		interval     time.Duration `default: 60`
		channel_size uint16        `default:10`
	}

	database struct {
		db_type  string `default:"mysql"`
		db_host  string `default:"localhost"`
		db_port  string `default:"3306"`
		db_name  string `default:"dbname"`
		username string `default:"user"`
		password string `default:"password"`
		charset  string `default:"utf8mb4"`
	}

	redis struct {
		network      string        `default:"tcp"`
		address      string        `default:"127.0.0.1:6379"`
		password     string        `default: ""`
		max_idle     int           `default:5`   // 最大的空闲连接数
		idle_timeout time.Duration `default:240` // 最大的空闲连接等待时间，单位 秒
		db           int           `default:0`
		queue_name   string        `default:"queue"`
		touch_name   string        `default:"touch"`
		stats_name   string        `default:"stats"` // Redis 查询状态保存
	}
}

func parseConfig(c *config, cfg string) {
	err := configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadInterval: 10 * time.Second}).Load(config, cfg)
	if err != nil {
		log.Fatalf("load config file %s error: %v", cfg, err)
	}
}
