package pkg

import (
	"os"

	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
)

var cfg Config

type Config struct {
	MysqlCfg  MySQLConfig  `yaml:"mysql"`
	ConsulCfg ConsulConfig `yaml:"consul"`
}

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Charset  string `yaml:"charset"`
}

type ConsulConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

func ParseMysqlConfig() {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		glog.Errorln("[MysqlDB] Failed to read config file: ", err.Error())
		return
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		glog.Errorln("[MysqlDB] Failed to parse config file: ", err.Error())
		return
	}
}
