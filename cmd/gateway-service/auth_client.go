package main

import (
	"context"
	authpb "goshop/api/protobuf/auth"
	errorcode "goshop/pkg/error"
	"net/http"

	"github.com/gin-gonic/gin"
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
		glog.Errorln("[GatewayServer] consul service recover failed.")
		return err
	}
	authConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[GatewayServer] new login rpc client error: ", err.Error())
		return err
	}
	authClient = authpb.NewAuthServiceClient(authConn)
	glog.Infoln("[GatewayServer] connect [auth-service] server successful on: ", addr)
	return nil
}

func AuthClient() authpb.AuthServiceClient {
	return authClient
}

func AuthClientClose() error {
	return authConn.Close()
}

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
