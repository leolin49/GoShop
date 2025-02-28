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

type LoginServer struct {
	service.Service
}

var (
	serviceName string = "login-service"
	serverName  string
	serverId    string
	consul      *service.ConsulClient
	db          *gorm.DB
	rdb         redis.IRdb
	server      *LoginServer
	once        sync.Once
)

func LoginServerGetInstance() *LoginServer {
	once.Do(func() {
		server = &LoginServer{}
		server.Derived = server
	})
	return server
}

func (s *LoginServer) Init() bool {
	var err error

	// Parse config
	if !configs.ParseConfig() {
		glog.Errorln("[LoginServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[LoginServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery(serverName)
	if err != nil {
		glog.Errorln("[LoginServer] recover config from consul error: ", err.Error())
		return false
	}

	// Rpc server
	if !rpcServerStart(cfg) {
		glog.Errorln("[LoginServer] rpc server start error.")
		return false
	}

	// Rpc client
	rpcClientsStart()

	// MySQL cluster connect
	if db, err = mysql.DBClusterInit(&cfg.MysqlClusterCfg); err != nil {
		glog.Errorln("[LoginServer] mysql cluster init failed: ", err.Error())
		return false
	}
	// MySQL table migrate
	db.AutoMigrate(&models.User{})

	// Redis connect
	if rdb, err = redis.NewRedisClient(&cfg.RedisCfg); err != nil {
		glog.Errorln("[LoginServer] redis database init failed: ", err.Error())
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		serviceName,
		cfg.LoginCfg.Host,
		cfg.LoginCfg.Port,
		"5s",
		"5s",
	) {
		glog.Errorln("[LoginServer] consul register error.")
		return false
	}

	return true
}

func (s *LoginServer) Reload() {
}

func (s *LoginServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *LoginServer) Final() bool {
	return true
}

func main() {
	defer func() {
		rpcClientsClose()
		_ = consul.ServiceDeregister(serviceName)
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

	LoginServerGetInstance().Main()
	glog.Infoln("[LoginServer] server closed.")
}
