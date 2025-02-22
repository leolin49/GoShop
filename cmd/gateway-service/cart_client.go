package main

import (
	"context"
	cartpb "goshop/api/protobuf/cart"
	"goshop/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cart_client cartpb.CartServiceClient
	cart_conn   *grpc.ClientConn
)

func CartClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := service.ServiceRecover("cart-service")
	if err != nil || addr == "" {
		glog.Errorln("[Gatewayserver] consul service recover failed.")
		return false
	}
	cart_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Gatewayserver] new cart rpc client error: ", err.Error())
		return false
	}
	cart_client = cartpb.NewCartServiceClient(cart_conn)
	glog.Infoln("[Gatewayserver] connect [cart-service] server successful on: ", addr)
	return true
}

func CartClient() cartpb.CartServiceClient { return cart_client }

func CartClientClose() { login_conn.Close() }

func handleAddCart(c *gin.Context) {
	user_id, err := getPostFormInt(c, "user_id")
	if err != nil {
		invalidParam(c)
		return
	}
	product_id, err := getPostFormInt(c, "product_id")
	if err != nil {
		invalidParam(c)
		return
	}
	product_cnt, err := getPostFormInt(c, "product_cnt")
	if err != nil {
		invalidParam(c)
		return
	}

	req := &cartpb.ReqAddItem{
		UserId: uint32(user_id),
		Item: &cartpb.CartItem{
			ProductId: uint32(product_id),
			Quantity:  int32(product_cnt),
		},
	}
	ret, err := CartClient().AddItem(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleCleanCart(c *gin.Context) {
	user_id, err := getPostFormInt(c, "user_id")
	if err != nil {
		invalidParam(c)
		return
	}
	req := &cartpb.ReqCleanCart{
		UserId: uint32(user_id),
	}
	ret, err := CartClient().CleanCart(context.Background(), req)
	if err != nil {
		rpcRequestError(c)
	}
	c.JSON(http.StatusOK, ret)
}

func handleGetCart(c *gin.Context) {
	user_id, err := getPostFormInt(c, "user_id")
	if err != nil {
		invalidParam(c)
		return
	}
	req := &cartpb.ReqGetCart{
		UserId: uint32(user_id),
	}
	ret, err := CartClient().GetCart(context.Background(), req)
	if err != nil {
		rpcRequestError(c)
	}
	c.JSON(http.StatusOK, ret)
}
