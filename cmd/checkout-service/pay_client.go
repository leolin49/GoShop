package main

import (
	paypb "goshop/api/protobuf/pay"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	pay_client paypb.PayServiceClient
	pay_conn   *grpc.ClientConn
)

func PayClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("pay-service")
	if err != nil || addr == "" {
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return err
	}
	pay_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new pay rpc client error: ", err.Error())
		return err
	}
	pay_client = paypb.NewPayServiceClient(pay_conn)
	glog.Infoln("[Checkoutserver] connect [pay-service] server successful on: ", addr)
	return nil 
}

func PayClient() paypb.PayServiceClient { return pay_client }

func PayClientClose() error { return pay_conn.Close() }
