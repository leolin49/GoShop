package main

import (
	"flag"
	"goshop/configs"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
)

type GatewayServer struct {
	service.Service
}

var server *GatewayServer
var once sync.Once

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
		glog.Flush()
		LoginClientClose()
	}()
	flag.Set("v", "2")
	flag.Parse()
	GatewayServerGetInstance().Main()
	glog.Infoln("[GatewayServer] server closed.")
}
