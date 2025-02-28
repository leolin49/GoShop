package main

import (
	orderpb "goshop/api/protobuf/order"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	order_client orderpb.OrderServiceClient
	order_conn   *grpc.ClientConn
)

func OrderClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("order-service")
	if err != nil || addr == "" {
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return err 
	}
	order_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new cart rpc client error: ", err.Error())
		return err 
	}
	order_client = orderpb.NewOrderServiceClient(order_conn)
	glog.Infoln("[Checkoutserver] connect [order-service] server successful on: ", addr)
	return nil 
}

func OrderClient() orderpb.OrderServiceClient { return order_client }

func OrderClientClose() error { return order_conn.Close() }
