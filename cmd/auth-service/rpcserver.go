package main

import (
	"context"
	"fmt"
	authpb "goshop/api/protobuf/auth"
	"goshop/configs"
	errorcode "goshop/pkg/error"
	"goshop/pkg/util"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type AuthRpcService struct {
	authpb.UnimplementedAuthServiceServer
}

func rpcServerStart(cfg *configs.Config) bool {
	addr := fmt.Sprintf("%s:%s", cfg.AuthCfg.Host, cfg.AuthCfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("[AuthServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(rpcServer, new(AuthRpcService))
	glog.Infof("[AuthServer] Starting rpc server on [%s]\n", addr)
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[AuthServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *AuthRpcService) DeliverDoubleToken(ctx context.Context, req *authpb.ReqDeliverDoubleToken) (*authpb.RspDeliverDoubleToken, error) {
	accessToken, refreshToken, err := util.JwtDoubleToken(req.UserId, 30*60, 1*24*60*60)
	if err != nil {
		return nil, err
	}
	return &authpb.RspDeliverDoubleToken{
		ErrorCode:    errorcode.Ok,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthRpcService) VerifyToken(ctx context.Context, req *authpb.ReqVerifyToken) (*authpb.RspVerifyToken, error) {
	var (
		userId uint32
		err    error
	)
	if req.IsAccess {
		userId, err = util.JwtExtractAccessTokenUserId(req.Token)
	} else {
		userId, err = util.JwtExtractRefreshTokenUserId(req.Token)
	}
	if err != nil {
		return nil, err
	}
	return &authpb.RspVerifyToken{
		ErrorCode: errorcode.Ok,
		UserId:    userId,
	}, nil
}
