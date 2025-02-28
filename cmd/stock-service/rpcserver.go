package main

import (
	"context"
	"fmt"
	productpb "goshop/api/protobuf/product"
	stockpb "goshop/api/protobuf/stock"
	"goshop/configs"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type StockRpcService struct {
	stockpb.UnimplementedStockServiceServer
}

func rpcServerStart(cfg *configs.Config) bool {
	addr := fmt.Sprintf("%s:%s", cfg.StockCfg.Host, cfg.StockCfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("[StockServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	stockpb.RegisterStockServiceServer(rpcServer, new(StockRpcService))
	glog.Infof("[StockServer] Starting rpc server on [%s]\n", addr)
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[StockServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *StockRpcService) GetStock(ctx context.Context, req *stockpb.ReqGetStock) (*stockpb.RspGetStock, error) {
	count, err := models.NewStockQuery(db).GetStock(req.ProductId)
	if err != nil {
		return nil, err
	}
	return &stockpb.RspGetStock{Count: count}, nil
}

func (s *StockRpcService) AddStock(ctx context.Context, req *stockpb.ReqAddStock) (*stockpb.RspAddStock, error) {
	// check the product exist
	_, err := ProductClient().GetProduct(ctx, &productpb.ReqGetProduct{
		Id: req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	err = models.NewStockQuery(db).AddStock(req.ProductId, req.AddCount)
	if err != nil {
		return nil, err
	}
	return &stockpb.RspAddStock{ErrorCode: errorcode.Ok}, nil
}

func (s *StockRpcService) SubStock(ctx context.Context, req *stockpb.ReqSubStock) (*stockpb.RspSubStock, error) {
	// check the product exist
	_, err := ProductClient().GetProduct(ctx, &productpb.ReqGetProduct{
		Id: req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	err = models.NewStockQuery(db).SubStock(req.ProductId, req.SubCount)
	if err != nil {
		// FIXME: more info in return struct
		return nil, err
	}
	return &stockpb.RspSubStock{ErrorCode: errorcode.Ok}, nil
}
