package main

import (
	"flag"
	"goshop/configs"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
)

type AuthServer struct {
	service.Service
}

var (
	serverId string
	server   *AuthServer
	once     sync.Once
	consul   *service.ConsulClient
)

func AuthServerGetInstance() *AuthServer {
	once.Do(func() {
		server = &AuthServer{}
		server.Derived = server
	})
	return server
}

func (s *AuthServer) Init() bool {
	var err error

	if !configs.ParseConfig() {
		glog.Errorln("[AuthServer] parse config error.")
		return false
	}

	// Consul client
	consul, err = service.NewConsulClient(&configs.GetConf().ConsulCfg)
	if err != nil {
		glog.Errorln("[AuthServer] new consul client failed: ", err.Error())
		return false
	}
	cfg, err := consul.ConfigQuery("auth-service/" + serverId)
	if err != nil {
		glog.Errorln("[AuthServer] recover config from consul error: ", err.Error())
		return false
	}

	if !rpcServerStart(cfg) {
		glog.Errorln("[AuthServer] rpc server start error.")
		return false
	}

	// Consul register
	if !consul.ServiceRegister(
		serverId,
		"auth-service",
		cfg.AuthCfg.Host,
		cfg.AuthCfg.Port,
		"1s",
		"5s",
	) {
		glog.Errorln("[AuthServer] consul register error.")
		return false
	}
	return true
}

func (s *AuthServer) Reload() {
}

func (s *AuthServer) MainLoop() {
	time.Sleep(time.Second)
}

func (s *AuthServer) Final() bool {
	return true
}

func main() {
	defer func() {
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

	AuthServerGetInstance().Main()
	glog.Infoln("[AuthServer] server closed.")
}
