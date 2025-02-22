package main

import (
	"context"
	cartpb "goshop/api/protobuf/cart"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type CartRpcService struct {
	cartpb.UnimplementedCartServiceServer
}

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49200")
	if err != nil {
		glog.Fatalf("[CartServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	cartpb.RegisterCartServiceServer(rpcServer, new(CartRpcService))
	glog.Infoln("[CartServer] Starting rpc server on :49200")
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[CartServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *CartRpcService) AddItem(ctx context.Context, req *cartpb.ReqAddItem) (*cartpb.RspAddItem, error) {
	cart := &models.Cart{
		UserId:     uint64(req.UserId),
		ProductId:  uint64(req.Item.ProductId),
		ProductCnt: uint64(req.Item.Quantity),
	}
	err := models.NewCartQuery(Mysql()).AddProduct(cart)
	if err != nil {
		return nil, err
	}
	return &cartpb.RspAddItem{
		ErrorCode: errorcode.Ok,
	}, nil
}

func (s *CartRpcService) CleanCart(ctx context.Context, req *cartpb.ReqCleanCart) (*cartpb.RspCleanCart, error) {
	err := models.NewCartQuery(Mysql()).CleanByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	return &cartpb.RspCleanCart{
		ErrorCode: errorcode.Ok,
	}, nil
}

func (s *CartRpcService) GetCart(ctx context.Context, req *cartpb.ReqGetCart) (*cartpb.RspGetCart, error) {
	carts, err := models.NewCartQuery(Mysql()).GetByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	var items []*cartpb.CartItem
	for _, cart := range carts {
		items = append(items, &cartpb.CartItem{
			ProductId: uint32(cart.ProductId),
			Quantity:  int32(cart.ProductCnt),
		})
	}
	return &cartpb.RspGetCart{
		ErrorCode: errorcode.Ok,
		Cart: &cartpb.Cart{
			UserId: req.UserId,
			Items:  items,
		},
	}, nil
}
