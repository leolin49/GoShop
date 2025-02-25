package main

import (
	"context"
	authpb "goshop/api/protobuf/auth"
	loginpb "goshop/api/protobuf/login"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type LoginRpcService struct {
	loginpb.UnimplementedLoginServiceServer
}

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49000")
	if err != nil {
		glog.Fatalf("[LoginServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	loginpb.RegisterLoginServiceServer(rpcServer, new(LoginRpcService))
	glog.Infoln("[LoginServer] Starting rpc server on :49000")
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[LoginServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *LoginRpcService) RegisterUser(ctx context.Context, req *loginpb.ReqRegisterUser) (*loginpb.RspRegisterUser, error) {
	var (
		exist int64
		res   *gorm.DB
	)
	// Check user if already exist by email.
	if res = LoginServerGetInstance().db.Model(&models.User{}).Where("email=?", req.Email).Count(&exist); res.Error != nil {
		return &loginpb.RspRegisterUser{
			ErrorCode: errorcode.UnknowError,
		}, res.Error
	} else if exist > 0 {
		return &loginpb.RspRegisterUser{
			ErrorCode: errorcode.UserAlreadyExist,
		}, nil
	}
	// Create the user record.
	user := &models.User{
		Email:    req.Email,
		Name:     req.Username,
		Password: req.Password,
	}
	res = LoginServerGetInstance().db.Create(user)
	if err := res.Error; err != nil {
		return &loginpb.RspRegisterUser{
			ErrorCode: errorcode.UnknowError,
		}, err
	}

	var user_id int64
	res = LoginServerGetInstance().db.Model(&models.User{}).Select("id").Where("email=?", req.Email).Find(&user_id)
	if err := res.Error; err != nil {
		return &loginpb.RspRegisterUser{
			ErrorCode: errorcode.UnknowError,
		}, err
	}
	return &loginpb.RspRegisterUser{
		ErrorCode: errorcode.Ok,
		UserId:    int32(user_id),
	}, nil
}

func (s *LoginRpcService) LoginUser(ctx context.Context, req *loginpb.ReqLoginUser) (*loginpb.RspLoginUser, error) {
	var user models.User
	res := LoginServerGetInstance().db.Select("id, password").Where("email=?", req.Email).Find(&user)
	if err := res.Error; err != nil {
		return &loginpb.RspLoginUser{
			ErrorCode: errorcode.UnknowError,
		}, err
	}
	// Check user if already exist by email.
	if res.RowsAffected == 0 {
		return &loginpb.RspLoginUser{
			ErrorCode: errorcode.UserNotExist,
		}, nil
	}
	// Check the password.
	if req.Password != user.Password {
		return &loginpb.RspLoginUser{
			ErrorCode: errorcode.LoginPasswordError,
		}, nil
	}
	// Delivery Jwt Token.
	ret, err := AuthClient().DeliverDoubleToken(ctx, &authpb.ReqDeliverDoubleToken{
		UserId: uint32(user.ID),
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
