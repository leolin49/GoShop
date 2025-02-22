package main

import (
	"context"
	authpb "goshop/api/protobuf/auth"
	errorcode "goshop/pkg/error"
	"goshop/pkg/util"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type AuthRpcService struct {
	authpb.UnimplementedAuthServiceServer
}

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49300")
	if err != nil {
		glog.Fatalf("[AuthServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(rpcServer, new(AuthRpcService))
	glog.Infoln("[AuthServer] Starting rpc server on :49300")
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[AuthServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *AuthRpcService) DeliverToken(ctx context.Context, req *authpb.ReqDeliverToken) (*authpb.RspDeliverToken, error) {
	tokenString, err := util.JwtGenerateToken(req.UserId)
	if err != nil {
		return nil, err
	}
	return &authpb.RspDeliverToken{
		ErrorCode: errorcode.Ok,
		Token: tokenString,
	}, nil
}

func (s *AuthRpcService) VerifyToken(ctx context.Context, req *authpb.ReqVerifyToken) (*authpb.RspVerifyToken, error) {
	userId, err := util.JwtExtractTokenUserId(req.Token)
	if err != nil {
		return nil, err
	}
	return &authpb.RspVerifyToken{
		ErrorCode: errorcode.Ok,
		UserId: userId,
	}, nil
}

