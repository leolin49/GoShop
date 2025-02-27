package main

import (
	"context"
	loginpb "goshop/api/protobuf/login"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	login_client loginpb.LoginServiceClient
	login_conn   *grpc.ClientConn
)

func LoginClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("login-service")
	if err != nil || addr == "" {
		glog.Errorln("[Gatewayserver] consul service recover failed.")
		return false
	}
	login_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Gatewayserver] new login rpc client error: ", err.Error())
		return false
	}
	login_client = loginpb.NewLoginServiceClient(login_conn)
	glog.Infoln("[Gatewayserver] connect [login-service] server successful on: ", addr)
	return true
}

func LoginClient() loginpb.LoginServiceClient {
	return login_client
}

func LoginClientClose() error {
	return login_conn.Close()
}

func handleRegister(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirm_password := c.PostForm("confirm_password")
	if email == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing email or username",
		})
		return
	}
	if password != confirm_password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "two password is not same",
		})
		return
	}
	req := &loginpb.ReqRegisterUser{
		Email:    email,
		Username: username,
		Password: password,
	}
	ret, err := LoginClient().RegisterUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing username or password",
		})
		return
	}
	req := &loginpb.ReqLoginUser{
		Email:    email,
		Password: password,
	}
	ret, err := LoginClient().LoginUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ret)
}
