package main

import (
	"context"
	authpb "goshop/api/protobuf/auth"
	errorcode "goshop/pkg/error"
	"goshop/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
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

func handleRefreshToken(c *gin.Context) {
	refreshToken := c.PostForm("refreshToken")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing refreshToken",
		})
		return
	}
	ret, err := AuthClient().VerifyToken(context.Background(), &authpb.ReqVerifyToken{
		Token:    refreshToken,
		IsAccess: false,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ret2, err := AuthClient().DeliverDoubleToken(context.Background(), &authpb.ReqDeliverDoubleToken{
		UserId: ret.UserId,
	})
	if err != nil || ret.ErrorCode != errorcode.Ok {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, ret2)
}
