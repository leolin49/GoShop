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
	server *AuthServer
	once   sync.Once
)

func AuthServerGetInstance() *AuthServer {
	once.Do(func() {
		server = &AuthServer{}
		server.Derived = server
	})
	return server
}

func (s *AuthServer) Init() bool {
	if !configs.ParseConfig() {
		glog.Errorln("[AuthServer] parse config error.")
		return false
	}
	if !rpcServerStart() {
		glog.Errorln("[AuthServer] rpc server start error.")
		return false
	}
	// Consul register
	if !service.ServiceRegister(
		"1",
		"auth-service",
		configs.GetConf().AuthCfg.Host,
		configs.GetConf().AuthCfg.Port,
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
		glog.Flush()
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	flag.Set("v", "2")
	flag.Parse()

	AuthServerGetInstance().Main()
	glog.Infoln("[AuthServer] server closed.")
}
