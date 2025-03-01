package main

import (
	"context"
	"fmt"
	authpb "goshop/api/protobuf/auth"
	loginpb "goshop/api/protobuf/login"
	"goshop/configs"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"goshop/pkg/util"
	"net"
	"sync"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type LoginRpcService struct {
	loginpb.UnimplementedLoginServiceServer
	// TODO grpc retry
	mu 			sync.Mutex
	reqCounter 	uint
	reqModulo  	uint
}

func (s *LoginRpcService) RegisterUser(ctx context.Context, req *loginpb.ReqRegisterUser) (*loginpb.RspRegisterUser, error) {
	var err error
	// Create the user record.
	Md5Pwd := util.MD5WithSaltFun(req.Password, nil)
	user := &models.User{
		Email:    req.Email,
		Name:     req.Username,
		Password: Md5Pwd,
	}
	err = models.NewUserQuery(db).CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &loginpb.RspRegisterUser{
		ErrorCode: errorcode.Ok,
		UserId:    uint32(user.ID),
	}, nil
}

func (s *LoginRpcService) LoginUser(ctx context.Context, req *loginpb.ReqLoginUser) (*loginpb.RspLoginUser, error) {
	user_id, pwd, err := models.NewUserQueryRead(db).GetIdAndPwdByEmail(req.Email)
	if err != nil {
		return &loginpb.RspLoginUser{
			ErrorCode: errorcode.UnknowError,
		}, err
	}
	// Check the password.
	if req.Password != pwd {
		return &loginpb.RspLoginUser{
			ErrorCode: errorcode.LoginPasswordError,
		}, nil
	}
	// Delivery Jwt Token.
	ret, err := AuthClient().DeliverDoubleToken(ctx, &authpb.ReqDeliverDoubleToken{
		UserId: user_id,
	})
	if err != nil {
		return nil, err
	}

	return &loginpb.RspLoginUser{
		ErrorCode:    errorcode.Ok,
		AccessToken:  ret.AccessToken,
		RefreshToken: ret.RefreshToken,
	}, nil
}

func rpcServerStart(cfg *configs.Config) bool {
	addr := fmt.Sprintf("%s:%s", cfg.LoginCfg.Host, cfg.LoginCfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("[LoginServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	loginpb.RegisterLoginServiceServer(rpcServer, new(LoginRpcService))
	glog.Infof("[LoginServer] Starting rpc server on [%s]\n", addr)
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[LoginServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}
