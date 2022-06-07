package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// 配置文件结构
type config struct {
	app struct {
		service_name string
		log_file     string
		log_level    string
		interval     time.Duration
		channel_size uint16
	}

	database struct {
		db_type  string
		db_host  string
		db_port  string
		db_name  string
		username string
		password string
		charset  string
	}

	redis struct {
		network      string
		address      string
		auth         string
		max_idle     int
		idle_timeout time.Duration
		db_select    int
		queue_name   string
		touch_name   string
		stats_name   string
	}
}

// 解析配置文件
func parseConfig(c *config, cfg string) {
	f, err := os.Open(cfg)
	if err != nil {
		log.Fatalf("load config file %s error: %v", cfg, err)
	}
	yaml.NewDecoder(f).Decode(c)
	fmt.Println(c)
}

func ParseConfig(cfg string) {
	var c config
	parseConfig(&c, cfg)
}
