package main

import (
	"context"
	"fmt"
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
	defer func() { // NOTE: delete the redis cache in the end.
		if err := rdb.Del(fmt.Sprintf("cart:%d", req.UserId)); err != nil {
			glog.Errorln("[CartServer] delete redis cache failed: ", err.Error())
		}
	}()
	cart := &models.Cart{
		UserId:     uint64(req.UserId),
		ProductId:  uint64(req.Item.ProductId),
		ProductCnt: uint64(req.Item.Quantity),
	}
	err := models.NewCartQuery(db).AddProduct(cart)
	if err != nil {
		return nil, err
	}
	return &cartpb.RspAddItem{
		ErrorCode: errorcode.Ok,
	}, nil
}

func (s *CartRpcService) CleanCart(ctx context.Context, req *cartpb.ReqCleanCart) (*cartpb.RspCleanCart, error) {
	defer func() {
		if err := rdb.Del(fmt.Sprintf("cart:%d", req.UserId)); err != nil {
			glog.Errorln("[CartServer] delete redis cache failed: ", err.Error())
		}
	}()
	err := models.NewCartQuery(db).CleanByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	return &cartpb.RspCleanCart{
		ErrorCode: errorcode.Ok,
	}, nil
}

func (s *CartRpcService) GetCart(ctx context.Context, req *cartpb.ReqGetCart) (*cartpb.RspGetCart, error) {
	// redis cache
	key := fmt.Sprintf("cart:%d", req.UserId)
	var ret cartpb.RspGetCart
	cache, err := rdb.GetProto(key, &ret)
	if err != nil {
		return nil, err
	}
	if cache {
		return &ret, nil
	}

	carts, err := models.NewCartQuery(db).GetByUserId(req.UserId)
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
	ret = cartpb.RspGetCart{
		ErrorCode: errorcode.Ok,
		Cart: &cartpb.Cart{
			UserId: req.UserId,
			Items:  items,
		},
	}

	// redis cache
	if err = rdb.SetProto(key, &ret); err != nil {
		glog.Errorln("[CartServer] cache write error: ", err.Error())
	}

	return &ret, nil
}
