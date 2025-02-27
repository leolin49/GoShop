package main

import (
	authpb "goshop/api/protobuf/auth"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	authClient authpb.AuthServiceClient
	authConn   *grpc.ClientConn
)

func AuthClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("auth-service")
	if err != nil || addr == "" {
		glog.Errorln("[Authserver] consul service recover failed.")
		return err
	}
	authConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Authserver] new login rpc client error: ", err.Error())
		return err
	}
	authClient = authpb.NewAuthServiceClient(authConn)
	glog.Infoln("[Authserver] connect [auth-service] server successful on: ", addr)
	return nil
}

func AuthClient() authpb.AuthServiceClient {
	return authClient
}

func AuthClientClose() error {
	return authConn.Close()
}
