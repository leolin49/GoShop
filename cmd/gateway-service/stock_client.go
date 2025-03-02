package main

import (
	"context"
	stockpb "goshop/api/protobuf/stock"
	"net/http"

	"github.com/gin-gonic/gin"
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
		glog.Errorln("[GatewayServer] consul service recover failed.")
		return err
	}
	stock_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[GatewayServer] new cart rpc client error: ", err.Error())
		return err
	}
	stock_client = stockpb.NewStockServiceClient(stock_conn)
	glog.Infoln("[Gatewayserver] connect [cart-service] server successful on: ", addr)
	return nil
}

func StockClient() stockpb.StockServiceClient {
	return stock_client
}

func StockClientClose() error {
	return stock_conn.Close()
}

func handleAddStock(c *gin.Context) {
	productId, err := getPostFormInt(c, "product_id")
	if err != nil {
		invalidParam(c)
		return
	}
	addCnt, err := getPostFormInt(c, "add_count")
	if err != nil {
		invalidParam(c)
		return
	}

	stocks := []*stockpb.Stock{}
	stocks = append(stocks, &stockpb.Stock{
		ProductId: uint32(productId),
		Count:     uint64(addCnt),
	})

	req := &stockpb.ReqAddStock{
		Stocks: stocks,
	}
	ret, err := StockClient().AddStock(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleSubStock(c *gin.Context) {
	productId, err := getPostFormInt(c, "product_id")
	if err != nil {
		invalidParam(c)
		return
	}
	subCnt, err := getPostFormInt(c, "sub_count")
	if err != nil {
		invalidParam(c)
		return
	}
	stocks := []*stockpb.Stock{}
	stocks = append(stocks, &stockpb.Stock{
		ProductId: uint32(productId),
		Count:     uint64(subCnt),
	})

	req := &stockpb.ReqSubStock{
		Stocks: stocks,
	}
	ret, err := StockClient().SubStock(context.Background(), req)
	if err != nil {
		rpcRequestError(c, err)
	}
	c.JSON(http.StatusOK, ret)
}

func handleGetStock(c *gin.Context) {
	productId, err := getPostFormInt(c, "product_id")
	if err != nil {
		invalidParam(c)
		return
	}
	req := &stockpb.ReqGetStock{
		ProductId: uint32(productId),
	}
	ret, err := StockClient().GetStock(context.Background(), req)
	if err != nil {
		rpcRequestError(c, err)
	}
	c.JSON(http.StatusOK, ret)
}

// --------- test ------------------- //
func handleFlashWarm(c *gin.Context) {
	ret, err := StockClient().FlashCacheWarmUp(context.Background(), &stockpb.ReqFlashCacheWarmUp{})
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleFlashBuy(c *gin.Context) {
	product_id, err := getPostFormInt(c, "product_id")
	if err != nil {
		invalidParam(c)
		return
	}
	count, err := getPostFormInt(c, "sub_count")
	if err != nil {
		invalidParam(c)
		return
	}
	ret, err := StockClient().FlashStock(context.Background(), &stockpb.ReqFlashStock{
		ProductId: uint32(product_id),
		SubCount:  uint64(count),
	})
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleFlashClear(c *gin.Context) {
	ret, err := StockClient().FlashCacheClear(context.Background(), &stockpb.ReqFlashCacheClear{})
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}
