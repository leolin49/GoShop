package main

import (
	paypb "goshop/api/protobuf/pay"
	"goshop/pkg/service"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	pay_client paypb.PayServiceClient
	pay_conn   *grpc.ClientConn
)

func PayClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := service.ServiceRecover("pay-service")
	if err != nil || addr == "" {
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return false
	}
	pay_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new pay rpc client error: ", err.Error())
		return false
	}
	pay_client = paypb.NewPayServiceClient(pay_conn)
	glog.Infoln("[Checkoutserver] connect [pay-service] server successful on: ", addr)
	return true
}

func PayClient() paypb.PayServiceClient { return pay_client }

func PayClientClose() { pay_conn.Close() }

