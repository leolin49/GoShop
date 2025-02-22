package main

import (
	authpb "goshop/api/protobuf/auth"
	"goshop/pkg/service"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	auth_client authpb.AuthServiceClient 
	auth_conn   *grpc.ClientConn
)

func AuthClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := service.ServiceRecover("auth-service")
	if err != nil || addr == "" {
		glog.Errorln("[Authserver] consul service recover failed.")
		return false
	}
	auth_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Authserver] new login rpc client error: ", err.Error())
		return false
	}
	auth_client = authpb.NewAuthServiceClient(auth_conn)
	glog.Infoln("[Authserver] connect [auth-service] server successful on: ", addr)
	return true
}

func AuthClient() authpb.AuthServiceClient { return auth_client }

func AuthClientClose() { auth_conn.Close() }
