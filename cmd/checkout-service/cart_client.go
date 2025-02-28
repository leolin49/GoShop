package main

import (
	cartpb "goshop/api/protobuf/cart"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cart_client cartpb.CartServiceClient
	cart_conn   *grpc.ClientConn
)

func CartClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("cart-service")
	if err != nil || addr == "" {
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return err
	}
	cart_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new cart rpc client error: ", err.Error())
		return err
	}
	cart_client = cartpb.NewCartServiceClient(cart_conn)
	glog.Infoln("[Checkoutserver] connect [cart-service] server successful on: ", addr)
	return nil
}

func CartClient() cartpb.CartServiceClient {
	return cart_client
}

func CartClientClose() error {
	return cart_conn.Close()
}
