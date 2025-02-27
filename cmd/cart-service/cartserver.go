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

type CartServer struct {
	service.Service
}

var (
	serverId string
	consul   *service.ConsulClient
	server   *CartServer
	db       *gorm.DB
	rdb      redis.IRdb
	once     sync.Once
)

func CartServerGetInstance() *CartServer {
	once.Do(func() {
		server = &CartServer{}
		server.Derived = server
	})
	return server
}

func (s *CartServer) Init() bool {
	var err error
	if !configs.ParseConfig() {
		glog.Errorln("[CartServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[CartServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery("cart-service/" + serverId)
	if err != nil {
		glog.Errorln("[CartServer] recover config from consul error: ", err.Error())
		return false
	}

	// RPC server
	if !rpcServerStart(cfg) {
		glog.Errorln("[CartServer] rpc server start error.")
		return false
	}
	// RPC client
	rpcClientsStart()

	// MySQL connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[CartServer] mysql database init error.")
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.Cart{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&configs.GetConf().RedisCfg); err != nil {
		glog.Errorln("[CartServer] redis database init error.")
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		"cart-service",
		cfg.CartCfg.Host,
		cfg.CartCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[CartServer] consul register error.")
		return false
	}
	return true
}

func (s *CartServer) Reload() {}

func (s *CartServer) MainLoop() { time.Sleep(time.Second) }

func (s *CartServer) Final() bool { return true }

func main() {
	defer func() {
		rpcClientsClose()
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
	CartServerGetInstance().Main()
	glog.Infoln("[CartServer] server closed.")
}
