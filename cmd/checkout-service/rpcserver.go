package main

import (
	"context"
	"fmt"
	cartpb "goshop/api/protobuf/cart"
	checkoutpb "goshop/api/protobuf/checkout"
	paypb "goshop/api/protobuf/pay"
	productpb "goshop/api/protobuf/product"
	"goshop/models"
	"net"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type CheckoutRpcService struct {
	checkoutpb.UnimplementedCheckoutServiceServer
}

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49500")
	if err != nil {
		glog.Fatalf("[CheckoutServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	checkoutpb.RegisterCheckoutServiceServer(rpcServer, new(CheckoutRpcService))
	glog.Infoln("[CheckoutServer] Starting rpc server on :49500")
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[CheckoutServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *CheckoutRpcService) Checkout(ctx context.Context, req *checkoutpb.ReqCheckout) (*checkoutpb.RspCheckout, error) {
	cartRet, err := CartClient().GetCart(ctx, &cartpb.ReqGetCart{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	if cartRet == nil || cartRet.Cart == nil || len(cartRet.Cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	var total float32
	for _, item := range cartRet.Cart.Items {
		productRet, err := ProductClient().GetProduct(ctx, &productpb.ReqGetProduct{
			Id: item.ProductId,	
		})
		if err != nil {
			return nil, err
		}
		if productRet.Product == nil {
			glog.Warningf("[CheckoutServer] product [%d] is not existed!\n", item.ProductId)
			continue
		}

		price := productRet.Product.Price
		cost := price * float32(item.Quantity)
		
		total += cost
	}

	// create the order.
	var orderId string

	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	orderId = u.String()

	payReq := &paypb.ReqCharge{
		UserId: req.UserId,
		OrderId: orderId,
		Amount: total,
		CardInfo: &paypb.CreditCardInfo{
			CreditCardNumber: req.CardInfo,
		},
	}
}

