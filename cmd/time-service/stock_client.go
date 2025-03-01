package main

import (
	stockpb "goshop/api/protobuf/stock"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	stock_client stockpb.StockServiceClient
	stock_conn   *grpc.ClientConn
)

func StockClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("stock-service")
	if err != nil || addr == "" {
		glog.Errorln("[TimeServer] consul service recover failed.")
		return err
	}
	stock_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[TimeServer] new cart rpc client error: ", err.Error())
		return err
	}
	stock_client = stockpb.NewStockServiceClient(stock_conn)
	glog.Infoln("[Checkoutserver] connect [stock-service] server successful on: ", addr)
	return nil
}

func StockClient() stockpb.StockServiceClient {
	return stock_client
}

func StockClientClose() error {
	return stock_conn.Close()
}
