package main

import (
	cartpb "goshop/api/protobuf/cart"
	"goshop/pkg/service"

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
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return false
	}
	cart_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new cart rpc client error: ", err.Error())
		return false
	}
	cart_client = cartpb.NewCartServiceClient(cart_conn)
	glog.Infoln("[Checkoutserver] connect [cart-service] server successful on: ", addr)
	return true
}

func CartClient() cartpb.CartServiceClient { return cart_client }

func CartClientClose() { cart_conn.Close() }
