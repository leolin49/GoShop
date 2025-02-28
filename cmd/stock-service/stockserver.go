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

type StockServer struct {
	service.Service
}

var (
	serviceName string = "stock-service"
	serverId    string
	serverName  string
	consul      *service.ConsulClient
	server      *StockServer
	db          *gorm.DB
	rdb         redis.IRdb
	once        sync.Once
)

func StockServerGetInstance() *StockServer {
	once.Do(func() {
		server = &StockServer{}
		server.Derived = server
	})
	return server
}

func (s *StockServer) Init() bool {
	var err error
	serverName = serviceName + "/" + serverId

	if !configs.ParseConfig() {
		glog.Errorln("[StockServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[StockServer] new consul client failed: ", err.Error())
		return false
	}

	glog.Errorln(serverName)
	cfg, err := consul.ConfigQuery(serverName)
	glog.Errorln(cfg)
	if err != nil {
		glog.Errorln("[StockServer] recover config from consul error: ", err.Error())
		return false
	}

	// RPC server
	if !rpcServerStart(cfg) {
		glog.Errorln("[StockServer] rpc server start error.")
		return false
	}
	// RPC client
	rpcClientsStart()

	// MySQL connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[StockServer] mysql database init error.")
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.Stock{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&cfg.RedisCfg); err != nil {
		glog.Errorln("[StockServer] redis database init error.")
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		serviceName,
		cfg.StockCfg.Host,
		cfg.StockCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[StockServer] consul register error.")
		return false
	}
	return true
}

func (s *StockServer) Reload() {}

func (s *StockServer) MainLoop() { time.Sleep(time.Second) }

func (s *StockServer) Final() bool { return true }

func main() {
	defer func() {
		rpcClientsClose()
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
	StockServerGetInstance().Main()
	glog.Infoln("[StockServer] server closed.")
}
