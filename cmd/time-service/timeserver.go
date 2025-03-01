package main

import (
	"flag"
	"goshop/configs"
	"goshop/pkg/redis"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

type TimeServer struct {
	service.Service
}

var (
	serviceName string = "time-service"
	serverName  string
	serverId    string
	consul      *service.ConsulClient
	rdb         redis.IRdb
	server      *TimeServer
	once        sync.Once
)

func TimeServerGetInstance() *TimeServer {
	once.Do(func() {
		server = &TimeServer{}
		server.Derived = server
	})
	return server
}

func (s *TimeServer) Init() bool {
	var err error
	serverName = serviceName + "/" + serverId

	// Parse config
	if !configs.ParseConfig() {
		glog.Errorln("[TimeServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[TimeServer] new consul client failed: ", err.Error())
		return false
	}

	cfg, err := consul.ConfigQuery(serverName)
	if err != nil {
		glog.Errorln("[TimeServer] recover config from consul error: ", err.Error())
		return false
	}

	// Rpc client
	rpcClientsStart()

	// Redis connect
	if rdb, err = redis.NewRedisClient(&cfg.RedisCfg); err != nil {
		glog.Errorln("[TimeServer] redis database init failed: ", err.Error())
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		serviceName,
		cfg.TimeCfg.Host,
		cfg.TimeCfg.Port,
		"5s",
		"5s",
	) {
		glog.Errorln("[TimeServer] consul register error.")
		return false
	}

	go clockTrigger()

	startAllTicker()

	if !httpServerStart(cfg) {
		glog.Errorln("[GatewayServer] http server start error.")
		return false
	}

	return true
}

func (s *TimeServer) Reload() {
}

func (s *TimeServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *TimeServer) Final() bool {
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

	TimeServerGetInstance().Main()
	glog.Infoln("[TimeServer] server closed.")
}
