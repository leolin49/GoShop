package main

import (
	productpb "goshop/api/protobuf/product"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	product_client productpb.ProductServiceClient
	product_conn   *grpc.ClientConn
)

func ProductClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("product-service")
	if err != nil || addr == "" {
		glog.Errorln("[CartServer] consul service recover failed.")
		return err
	}
	product_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[CartServer] new product rpc client error: ", err.Error())
		return err
	}
	product_client = productpb.NewProductServiceClient(product_conn)
	glog.Infoln("[CartServer] connect [product-service] server successful on: ", addr)
	return nil
}

func ProductClient() productpb.ProductServiceClient {
	return product_client
}

func ProductClientClose() error {
	return product_conn.Close()
}
