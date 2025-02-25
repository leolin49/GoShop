package main

import (
	"flag"
	"goshop/configs"
	"goshop/pkg/mq"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
)

type GatewayServer struct {
	service.Service
	MQClient *mq.RabbitMQ
}

var (
	server *GatewayServer
	once   sync.Once
)

func GatewayServerGetInstance() *GatewayServer {
	once.Do(func() {
		server = &GatewayServer{}
		server.Derived = server
	})
	return server
}

func (s *GatewayServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[GatewayServer] parse config error.")
		return false
	}
	// rpc clients.
	rpcClientsStart()

	// rabbit client.
	var err error
	cfg := configs.GetConf()
	s.MQClient, err = mq.NewRabbitMQWorkClient("checkout-queue", cfg.GetRabbitMQUrl())
	if err != nil {
		glog.Errorln("[GatewayServer] rabbit mq client start error.")
		return false
	}

	if !httpServerStart() {
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
	flag.Set("v", "2")
	flag.Parse()
	GatewayServerGetInstance().Main()
	glog.Infoln("[GatewayServer] server closed.")
}
