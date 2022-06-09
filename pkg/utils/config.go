package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// 配置文件结构
type Config struct {
	App struct {
		Service_name string
		Log_file     string
		Log_level    string
		Interval     time.Duration
		Channel_size int
	}

	Database struct {
		Db_type  string
		Db_host  string
		Db_port  string
		Db_name  string
		Username string
		Password string
		Charset  string
	}

	Redis struct {
		Network      string
		Address      string
		Auth         string
		Max_idle     int
		Idle_timeout time.Duration
		Db_select    int
		Queue_name   string
		Touch_name   string
		Stats_name   string
	}
}

// 解析配置文件
func parseConfig(c *Config, cfg string) {
	f, err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Fatalf("load config file %s error: %v", cfg, err)
	}

	err = yaml.Unmarshal([]byte(f), c)
	if err != nil {
		log.Fatalf("unmarshal config file %s error: %v", cfg, err)
	}

	y, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("marshal config file %s error: %v", cfg, err)
	}

	fmt.Println(string(y))
}

func ParseConfig(cfg string) {
	var c Config
	parseConfig(&c, cfg)
}
