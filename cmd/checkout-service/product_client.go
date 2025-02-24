package main

import (
	productpb "goshop/api/protobuf/product"
	"goshop/pkg/service"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	product_client productpb.ProductServiceClient
	product_conn   *grpc.ClientConn
)

func ProductClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := service.ServiceRecover("product-service")
	if err != nil || addr == "" {
		glog.Errorln("[Checkoutserver] consul service recover failed.")
		return false
	}
	product_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Checkoutserver] new product rpc client error: ", err.Error())
		return false
	}
	product_client = productpb.NewProductServiceClient(product_conn)
	glog.Infoln("[Checkoutserver] connect [product-service] server successful on: ", addr)
	return true
}

func ProductClient() productpb.ProductServiceClient { return product_client }

func ProductClientClose() { product_conn.Close() }
