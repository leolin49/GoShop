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

type PayServer struct {
	service.Service
}

var (
	serverId    string
	serviceName string = "pay-service"
	serverName  string
	server      *PayServer
	consul      *service.ConsulClient
	db          *gorm.DB
	rdb         redis.IRdb
	once        sync.Once
)

func PayServerGetInstance() *PayServer {
	once.Do(func() {
		server = &PayServer{}
		server.Derived = server
	})
	return server
}

func (s *PayServer) Init() bool {
	var err error
	serverName = serviceName + "/" + serverId

	if !configs.ParseConfig() {
		glog.Errorln("[PayServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[PayServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery(serverName)
	if err != nil {
		glog.Errorln("[PayServer] recover config from consul error: ", err.Error())
		return false
	}

	if !rpcServerStart(cfg) {
		glog.Errorln("[PayServer] rpc server start error.")
		return false
	}
	// MySQL connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[PayServer] mysql database init failed: ", err.Error())
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.PaymentLog{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&configs.GetConf().RedisCfg); err != nil {
		glog.Errorln("[PayServer] redis database init failed: ", err.Error())
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		serviceName,
		cfg.PayCfg.Host,
		cfg.PayCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[PayServer] consul register error.")
		return false
	}
	return true
}

func (s *PayServer) Reload() {
}

func (s *PayServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *PayServer) Final() bool {
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
	PayServerGetInstance().Main()
	glog.Infoln("[PayServer] server closed.")
}
