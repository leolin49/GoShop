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

// func Helper(ctx context.Context, rdb *redis.Rdb, prefix string, req interface{}, handlerFunc func()) (interface{}, error) {
// 	key := fmt.Sprintf("%s:%d", prefix, req.UserId)
// 	ret := reflect.New(reflect.TypeOf(req))
// 	cache, err := rdb.GetProto(key, ret.Interface().(proto.Message))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if cache {
// 		return ret, nil
// 	}
// 	return ret, nil
// }

func (s *CartRpcService) GetCart(ctx context.Context, req *cartpb.ReqGetCart) (*cartpb.RspGetCart, error) {
	// redis cache
	rdb := CartServerGetInstance().rdb
	key := fmt.Sprintf("card:%d", req.UserId)
	var ret cartpb.RspGetCart
	cache, err := rdb.GetProto(key, &ret)
	if err != nil {
		return nil, err
	}
	if cache {
		return &ret, nil
	}

	db := CartServerGetInstance().db
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

// NOTE: grpc cache middleware, don't use !
// func CacheMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	rdb, ok := ctx.Value("rdb").(*redis.Rdb)
// 	if !ok {
// 		return nil, errors.New("grpc cache middleware: no redis client")
// 	}

// 	key := fmt.Sprintf("%s:%v", info.FullMethod, req)
// 	data, err := rdb.Get(key)
// 	if err == nil {
// 		ret := &cartpb.RspGetCart{}
// 		if err = util.Deserialize([]byte(data), ret); err == nil {
// 			return ret, nil
// 		} 	
// 	}
// 	ret, err := handler(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cacheData, err := util.Serialize(ret.(proto.Message))
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = rdb.Set(key, string(cacheData))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ret, nil
// }

