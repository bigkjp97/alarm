package server

type DBServer struct {
	Type     string `yaml:"db_type"`
	Host     string `yaml:"db_host"`
	Port     int    `yaml:"db_port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	DBname   string `yaml:"db_name"`
	Charset  string `yaml:"charset"`
}
