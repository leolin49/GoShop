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

type OrderServer struct {
	service.Service
}

var (
	serverId    string
	serviceName string = "order-service"
	serverName  string
	server      *OrderServer
	consul      *service.ConsulClient
	db          *gorm.DB
	rdb         redis.IRdb
	once        sync.Once
)

func OrderServerGetInstance() *OrderServer {
	once.Do(func() {
		server = &OrderServer{}
		server.Derived = server
	})
	return server
}

func (s *OrderServer) Init() bool {
	var err error
	serverName = serviceName + "/" + serverId

	if !configs.ParseConfig() {
		glog.Errorln("[OrderServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[OrderServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery(serverName)
	if err != nil {
		glog.Errorln("[OrderServer] recover config from consul error: ", err.Error())
		return false
	}

	if !rpcServerStart(cfg) {
		glog.Errorln("[OrderServer] rpc server start error.")
		return false
	}

	// MySQL connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[OrderServer] mysql database init failed: ", err.Error())
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.Order{}, &models.OrderItem{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&configs.GetConf().RedisCfg); err != nil {
		glog.Errorln("[OrderServer] redis database init failed: ", err.Error())
		return false
	}
	// Consul register
	if !consul.ServiceRegister(
		serverId,
		serviceName,
		cfg.OrderCfg.Host,
		cfg.OrderCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[OrderServer] consul register error.")
		return false
	}
	return true
}

func (s *OrderServer) Reload() {
}

func (s *OrderServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *OrderServer) Final() bool {
	return true
}

func main() {
	defer func() {
		_ = consul.ServiceDeregister(serverName)
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
	OrderServerGetInstance().Main()
	glog.Infoln("[OrderServer] server closed.")
}
