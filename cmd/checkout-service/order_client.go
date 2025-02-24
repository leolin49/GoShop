package main

import (
	orderpb "goshop/api/protobuf/order"
	"goshop/pkg/service"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	order_client orderpb.OrderServiceClient
	order_conn   *grpc.ClientConn
)

func OrderClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := service.ServiceRecover("order-service")
	if err != nil || addr == "" {
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return false
	}
	order_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new cart rpc client error: ", err.Error())
		return false
	}
	order_client = orderpb.NewOrderServiceClient(order_conn)
	glog.Infoln("[Checkoutserver] connect [order-service] server successful on: ", addr)
	return true
}

func OrderClient() orderpb.OrderServiceClient { return order_client }

func OrderClientClose() { order_conn.Close() }
