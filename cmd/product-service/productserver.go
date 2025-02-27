package main

import (
	"flag"
	"goshop/configs"
	"goshop/models"
	"goshop/pkg/mysql"
	"goshop/pkg/redis"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type ProductServer struct {
	service.Service
}

var (
	serverId string
	consul   *service.ConsulClient
	db       *gorm.DB
	rdb      redis.IRdb
	server   *ProductServer
	once     sync.Once
)

func ProductServerGetInstance() *ProductServer {
	once.Do(func() {
		server = &ProductServer{}
		server.Derived = server
	})
	return server
}

func (s *ProductServer) Init() bool {
	var err error

	if !configs.ParseConfig() {
		glog.Errorln("[ProductServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[ProductServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery("product-service/" + serverId)
	if err != nil {
		glog.Errorln("[ProductServer] recover config from consul error: ", err.Error())
		return false
	}

	// Rpc server
	if !rpcServerStart() {
		glog.Errorln("[ProductServer] rpc server start error.")
		return false
	}

	// MySQL cluster connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[ProductServer] mysql cluster init failed: ", err.Error())
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.Product{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&cfg.RedisCfg); err != nil {
		glog.Errorln("[ProductServer] redis database init failed: ", err.Error())
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		"product-service",
		cfg.ProductCfg.Host,
		cfg.ProductCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[ProductServer] consul register error.")
		return false
	}

	return true
}

func (s *ProductServer) Reload() {
}

func (s *ProductServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *ProductServer) Final() bool {
	return true
}

func main() {
	defer func() {
		_ = consul.ServiceDeregister(serverId)
		glog.Flush()
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	{
		flag.Set("v", "2")
		flag.StringVar(&serverId, "node", "node1", "the name of the service instance")
		flag.Parse()
	}

	ProductServerGetInstance().Main()
	glog.Infoln("[ProductServer] server closed.")
}
