package main

import (
	"flag"
	"goshop/configs"
	"goshop/models"
	"goshop/pkg/mysql"
	service "goshop/pkg/service"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type LoginServer struct {
	service.Service
	db *gorm.DB
}

var (
	server *LoginServer
	once   sync.Once
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
	// Rpc server
	if !rpcServerStart() {
		glog.Errorln("[LoginServer] rpc server start error.")
		return false
	}

	// MySQL connect
	if s.db, err = mysql.DatabaseInit(&configs.GetConf().MysqlCfg); err != nil {
		glog.Errorln("[LoginServer] mysql database init error.")
		return false
	}
	// MySQL table migrate
	s.db.AutoMigrate(&models.User{})

	// Rpc client
	rpcClientsStart()

	// Consul register
	if !service.ServiceRegister(
		"1",
		"login-service",
		configs.GetConf().LoginCfg.Host,
		configs.GetConf().LoginCfg.Port,
		"1s",
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
		AuthClientClose()
		glog.Flush()
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	flag.Set("v", "2")
	flag.Parse()

	LoginServerGetInstance().Main()
	glog.Infoln("[LoginServer] server closed.")
}
