package main

import (
	"context"
	"errors"
	"fmt"
	cartpb "goshop/api/protobuf/cart"
	productpb "goshop/api/protobuf/product"
	"goshop/configs"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type CartRpcService struct {
	cartpb.UnimplementedCartServiceServer
}

func rpcServerStart(cfg *configs.Config) bool {
	addr := fmt.Sprintf("%s:%s", cfg.CartCfg.Host, cfg.CartCfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("[CartServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	cartpb.RegisterCartServiceServer(rpcServer, new(CartRpcService))
	glog.Infof("[CartServer] Starting rpc server on [%s]\n", addr)
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
	if req.Item == nil {
		return nil, errors.New("no any item need to add in cart")
	}
	_, err := ProductClient().GetProduct(ctx, &productpb.ReqGetProduct{
		Id: req.Item.ProductId,
	})
	if err != nil {
		return nil, err
	}
	cart := &models.Cart{
		UserId:     uint64(req.UserId),
		ProductId:  uint64(req.Item.ProductId),
		ProductCnt: uint64(req.Item.Quantity),
	}
	// Check the product if exist or not.
	err = models.NewCartQuery(db).AddProduct(cart)
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
