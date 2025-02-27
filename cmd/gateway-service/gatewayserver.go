package main

import (
	"flag"
	"goshop/configs"
	"goshop/pkg/mq"
	service "goshop/pkg/service"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

type GatewayServer struct {
	service.Service
	MQClient *mq.RabbitMQ
}

var (
	serverId string
	localIp  string
	consul   *service.ConsulClient
	server   *GatewayServer
	once     sync.Once
)

func GatewayServerGetInstance() *GatewayServer {
	once.Do(func() {
		server = &GatewayServer{}
		server.Derived = server
	})
	return server
}

func (s *GatewayServer) Init() bool {
	var err error

	// FIXME: Here, only get consul config.
	if !configs.ParseConfig() {
		glog.Errorln("[GatewayServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[GatewayServer] new consul client failed: ", err.Error())
		return false
	}

	// Get config from consul
	cfg, err := consul.ConfigQuery("gateway-service/" + serverId)
	if err != nil {
		glog.Errorln("[GatewayServer] recover config from consul error: ", err.Error())
		return false
	}

	// Service register
	if !consul.ServiceRegister(
		serverId,
		"gateway-service",
		localIp,
		cfg.GatewayCfg.Port,
		"5s",
		"5s",
	) {
		glog.Errorln("[GatewayServer] consul register error.")
		return false
	}

	// rpc clients.
	rpcClientsStart()

	// rabbit client.
	s.MQClient, err = mq.NewRabbitMQWorkClient(cfg.GatewayCfg.MqName, cfg.GetRabbitMQUrl())
	if err != nil {
		glog.Errorln("[GatewayServer] rabbit mq client start error.")
		return false
	}

	if !httpServerStart(cfg) {
		glog.Errorln("[GatewayServer] http server start error.")
		return false
	}
	return true
}

func (s *GatewayServer) Reload() {
}

func (s *GatewayServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *GatewayServer) Final() bool {
	return true
}

func main() {
	defer func() {
		rpcClientsClose()
		glog.Flush()
	}()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	localIp = os.Getenv("LOCALIP")

	{
		flag.Set("v", "2")
		flag.StringVar(&serverId, "node", "node1", "the name of the service instance")
		flag.Parse()
	}

	GatewayServerGetInstance().Main()

	glog.Infoln("[GatewayServer] server closed.")
}
