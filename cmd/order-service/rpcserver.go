package main

import (
	"context"
	"fmt"
	cartpb "goshop/api/protobuf/cart"
	orderpb "goshop/api/protobuf/order"
	"goshop/models"
	"net"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type OrderRpcService struct {
	orderpb.UnimplementedOrderServiceServer
}

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49600")
	if err != nil {
		glog.Fatalf("[OrderServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(rpcServer, new(OrderRpcService))
	glog.Infoln("[OrderServer] Starting rpc server on :49600")
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[OrderServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *OrderRpcService) PlaceOrder(ctx context.Context, req *orderpb.ReqPlaceOrder) (ret *orderpb.RspPlaceOrder, err error) {
	if len(req.OrderItems) == 0 {
		err = fmt.Errorf("order items is empty")
		return
	}

	err = Mysql().Transaction(func(tx *gorm.DB) error {
		orderId, _ := uuid.NewUUID()

		order := &models.Order {
			OrderId: orderId.String(),
			UserId: req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: models.Consignee{
				Email: req.Email,
			},
		}
		if req.Address != nil {
			address := req.Address
			order.Consignee.Country = address.Country
			order.Consignee.State = address.State
			order.Consignee.City = address.City
			order.Consignee.StreetAddress = address.StreetAddress
			order.Consignee.ZipCode = address.ZipCode
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		var items []*models.OrderItem
		for _, v := range req.OrderItems {
			items = append(items, &models.OrderItem{
				OrderIdRefer: order.OrderId,
				ProductId: v.Item.ProductId,
				Quantity: uint32(v.Item.Quantity),
				Cost: v.Cost,
			})
		}
		if err := tx.Create(&items).Error; err != nil {
			glog.Errorln("222222222222222222222")
			return err
		}

		ret = &orderpb.RspPlaceOrder{
			OrderResult: &orderpb.OrderResult{
				OrderId: orderId.String(),
			},	
		}

		return nil
	})

	return
}

func (s *OrderRpcService) ListOrder(ctx context.Context, req *orderpb.ReqListOrder) (*orderpb.RspListOrder, error) {
	list, err := models.NewOrderQuery(Mysql()).ListOrder(req.UserId)
	if err != nil {
		return nil, err
	}

	var orders []*orderpb.Order
	for _, v := range list {
		var items []*orderpb.OrderItem
		for _, oi := range v.OrderItem {
			items = append(items, &orderpb.OrderItem{
				Item: &cartpb.CartItem{
					ProductId: oi.ProductId,
					Quantity: int32(oi.Quantity),
				},
				Cost: oi.Cost,
			})
		}
		orders = append(orders, &orderpb.Order{
			OrderId: v.OrderId,
			UserId: v.UserId,
			UserCurrency: v.UserCurrency,
			Email: v.Consignee.Email,
			Address: &orderpb.Address{
				Country: v.Consignee.Country,
				State: v.Consignee.State,
				City: v.Consignee.City,
				StreetAddress: v.Consignee.StreetAddress,
				ZipCode: v.Consignee.ZipCode,
			},
			OrderItems: items,
		})	
	}
	return &orderpb.RspListOrder{
		Orders: orders,
	}, nil
}




