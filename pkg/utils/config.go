package utils

import (
	"alarm/pkg/utils/server"
	"fmt"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	DBConfig server.DBServer `yaml:"db_config,omitempty`
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}

	fmt.Printf(string(b))
	return string(b)
}
