package main

import (
	"context"
	productpb "goshop/api/protobuf/product"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		glog.Errorln("[Gatewayserver] consul service recover failed.")
		return err
	}
	product_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Gatewayserver] new product rpc client error: ", err.Error())
		return err
	}
	product_client = productpb.NewProductServiceClient(product_conn)
	glog.Infoln("[Gatewayserver] connect [product-service] server successful on: ", addr)
	return nil
}

func ProductClient() productpb.ProductServiceClient {
	return product_client
}

func ProductClientClose() error {
	return product_conn.Close()
}

func handleAddProduct(c *gin.Context) {
	productName := c.PostForm("product_name")
	description := c.PostForm("description")
	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		invalidParam(c)
		return
	}
	categories := c.PostForm("categories")
	categoriesName := strings.Split(categories, ",")
	pd := &productpb.Product{
		Name:        productName,
		Description: description,
		Price:       float32(price),
		Categories:  categoriesName,
	}
	req := &productpb.ReqAddProduct{
		Product: pd,
	}
	ret, err := ProductClient().AddProduct(context.Background(), req)
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error_code": ret.ErrorCode,
	})
}

func handleListProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.PostForm("page"))
	if err != nil {
		invalidParam(c)
		return
	}
	page_size, err := strconv.Atoi(c.PostForm("page_size"))
	if err != nil {
		invalidParam(c)
		return
	}
	category_name := c.PostForm("category")
	req := &productpb.ReqListProducts{
		Page:         int32(page),
		PageSize:     int64(page_size),
		CategoryName: category_name,
	}
	response, err := ProductClient().ListProducts(context.Background(), req)
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": response.Products,
	})
}

func handleSearchProducts(c *gin.Context) {
	query := c.PostForm("query")
	ret, err := ProductClient().SearchProducts(context.Background(), &productpb.ReqSearchProducts{
		Query: query,
	})
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": ret.Results,
	})
}

func handleGetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("product_id"))
	ret, err := ProductClient().GetProduct(context.Background(), &productpb.ReqGetProduct{
		Id: uint32(id),
	})
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product": ret.Product,
	})
}
