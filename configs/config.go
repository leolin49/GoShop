package configs

import (
	"os"

	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
)

var cfg *Config

type Config struct {
	MysqlCfg   MySQLConfig   `yaml:"mysql"`
	ConsulCfg  ConsulConfig  `yaml:"consul"`
	GatewayCfg ServiceConfig `yaml:"gateway-service"`
	LoginCfg   ServiceConfig `yaml:"login-service"`
	AuthCfg    ServiceConfig `yaml:"auth-service"`
	ProductCfg ServiceConfig `yaml:"product-service"`
	CartCfg    ServiceConfig `yaml:"cart-service"`
	PayCfg		ServiceConfig `yaml:"pay-service"`
	CheckoutCfg	ServiceConfig `yaml:"checkout-service"`
	OrderCfg	ServiceConfig `yaml:"order-service"`
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
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Scheme string `yaml:"scheme"`
}

type ServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func ParseConfig() bool {
	data, err := os.ReadFile(os.Getenv("GOSHOP") + "configs/dev/config.yaml")
	if err != nil {
		glog.Errorln("[Config] Failed to read config file: ", err.Error())
		return false
	}

	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		glog.Errorln("[Config] Failed to parse config file: ", err.Error())
		return false
	}
	return true
}

func GetConf() *Config { return cfg }
