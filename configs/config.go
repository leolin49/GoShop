package configs

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
)

var cfg *Config

type Config struct {
	MysqlCfg        MySQLConfig        `yaml:"mysql"`
	ConsulCfg       ConsulConfig       `yaml:"consul"`
	GatewayCfg      GatewayConfig      `yaml:"gateway-service"`
	LoginCfg        ServiceConfig      `yaml:"login-service"`
	AuthCfg         ServiceConfig      `yaml:"auth-service"`
	ProductCfg      ServiceConfig      `yaml:"product-service"`
	CartCfg         ServiceConfig      `yaml:"cart-service"`
	PayCfg          ServiceConfig      `yaml:"pay-service"`
	CheckoutCfg     ServiceConfig      `yaml:"checkout-service"`
	OrderCfg        ServiceConfig      `yaml:"order-service"`
	RabbitMqCfg     RabbitMQConfig     `yaml:"rabbitmq"`
	RedisCfg        RedisConfig        `yaml:"redis"`
	MysqlClusterCfg MySQLClusterConfig `yaml:"mysql-cluster"`
}

type MySQLConfig struct {
	DSNFormat string `yaml:"dsnformat"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DataBase  string `yaml:"database"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Charset   string `yaml:"charset"`
}

func (cfg *MySQLConfig) GetDSN() string {
	return fmt.Sprintf(cfg.DSNFormat,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DataBase,
		cfg.Charset,
	)
}

type MySQLClusterConfig struct {
	Master   MySQLConfig   `yaml:"master"`
	Replicas []MySQLConfig `yaml:"replicas"`
}

type ConsulConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Scheme string `yaml:"scheme"`
}

type RabbitMQConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	UrlFormat string `yaml:"urlformat"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Network  string `yaml:"network"`
	Db       int    `yaml:"db"`
	Protocol int    `yaml:"protocol"`
}

type GatewayConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	MqName string `yaml:"mqname"`
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

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisCfg.Host, c.RedisCfg.Port)
}

func (c *Config) GetRabbitMQUrl() string {
	return fmt.Sprintf(c.RabbitMqCfg.UrlFormat,
		c.RabbitMqCfg.Username,
		c.RabbitMqCfg.Password,
		c.RabbitMqCfg.Host,
		c.RabbitMqCfg.Port,
	)
}

func (c *Config) GetConsulAddr() string {
	return fmt.Sprintf("%s:%s", c.ConsulCfg.Host, c.ConsulCfg.Port)
}
