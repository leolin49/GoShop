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

type CheckoutServer struct {
	service.Service
}

var (
	serviceName string = "checkout-service"
	serverId    string
	serverName  string
	consul      *service.ConsulClient
	db          *gorm.DB
	rdb         redis.IRdb
	server      *CheckoutServer
	once        sync.Once
)

func CheckoutServerGetInstance() *CheckoutServer {
	once.Do(func() {
		server = &CheckoutServer{}
		server.Derived = server
	})
	return server
}

func (s *CheckoutServer) Init() bool {
	var err error
	serverName = serviceName + "/" + serverId

	if !configs.ParseConfig() {
		glog.Errorln("[CheckoutServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[CheckoutServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery(serverName)
	if err != nil {
		glog.Errorln("[CheckoutServer] recover config from consul error: ", err.Error())
		return false
	}

	if !rpcServerStart(cfg) {
		glog.Errorln("[CheckoutServer] rpc server start error.")
		return false
	}

	// rpc clients
	rpcClientsStart()

	// MySQL connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[CartServer] mysql database init error.")
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.PaymentLog{})

	// redis
	if rdb, err = redis.NewRedisClient(&configs.GetConf().RedisCfg); err != nil {
		glog.Errorln("[CheckoutServer] redis database init error: ", err.Error())
		return false
	}

	// RabbitMQ consume
	go rabbitConsumer(cfg.CheckoutCfg.MqName, cfg.GetRabbitMQUrl())

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		serviceName,
		cfg.CheckoutCfg.Host,
		cfg.CheckoutCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[CheckoutServer] consul register error.")
		return false
	}
	return true
}

func (s *CheckoutServer) Reload() {
}

func (s *CheckoutServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *CheckoutServer) Final() bool {
	return true
}

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
	CheckoutServerGetInstance().Main()
	glog.Infoln("[CheckoutServer] server closed.")
}
