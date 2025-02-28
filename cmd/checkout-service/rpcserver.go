package main

import (
	"context"
	"fmt"
	cartpb "goshop/api/protobuf/cart"
	checkoutpb "goshop/api/protobuf/checkout"
	orderpb "goshop/api/protobuf/order"
	paypb "goshop/api/protobuf/pay"
	productpb "goshop/api/protobuf/product"
	stockpb "goshop/api/protobuf/stock"
	"goshop/configs"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type CheckoutRpcService struct {
	checkoutpb.UnimplementedCheckoutServiceServer
}

func rpcServerStart(cfg *configs.Config) bool {
	addr := fmt.Sprintf("%s:%s", cfg.CheckoutCfg.Host, cfg.CheckoutCfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("[CheckoutServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	checkoutpb.RegisterCheckoutServiceServer(rpcServer, new(CheckoutRpcService))
	glog.Infof("[CheckoutServer] Starting rpc server on [%s]\n", addr)
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
		glog.Errorln("[CheckoutServer] get cart error:", err.Error())
		return nil, err
	}
	if cartRet == nil || cartRet.Cart == nil || len(cartRet.Cart.Items) == 0 {
		glog.Errorln("[CheckoutServer] get cart error: cart is empty")
		return nil, fmt.Errorf("cart is empty")
	}

	var (
		total      float32
		orderId    string
		stockItems []*stockpb.Stock
		orderItems []*orderpb.OrderItem
	)

	for _, item := range cartRet.Cart.Items {
		productRet, err := ProductClient().GetProduct(ctx, &productpb.ReqGetProduct{
			Id: item.ProductId,
		})
		if err != nil {
			glog.Errorln("[CheckoutServer] get product error:", err.Error())
			return nil, err
		}
		if productRet.Product == nil {
			glog.Warningf("[CheckoutServer] product [%d] is not existed!\n", item.ProductId)
			continue
		}

		price := productRet.Product.Price
		cost := price * float32(item.Quantity)

		total += cost

		orderItems = append(orderItems, &orderpb.OrderItem{
			Item: &cartpb.CartItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			},
			Cost: cost,
		})

		stockItems = append(stockItems, &stockpb.Stock{
			ProductId: item.ProductId,
			Count:     uint64(item.Quantity),
		})
	}

	// create the order.
	orderRet, err := OrderClient().PlaceOrder(ctx, &orderpb.ReqPlaceOrder{
		UserId: req.UserId,
		Email:  req.Email,
		Address: &orderpb.Address{
			Country:       req.Address.Country,
			State:         req.Address.State,
			City:          req.Address.City,
			StreetAddress: req.Address.StreetAddress,
			ZipCode:       req.Address.ZipCode,
		},
		OrderItems: orderItems,
	})
	if err != nil {
		glog.Errorln("[CheckoutServer] place order error:", err.Error())
		return nil, err
	}
	if orderRet != nil && orderRet.OrderResult != nil {
		orderId = orderRet.OrderResult.OrderId
	}

	// clean cart.
	_, err = CartClient().CleanCart(ctx, &cartpb.ReqCleanCart{
		UserId: req.UserId,
	})
	if err != nil {
		glog.Errorln("[CheckoutServer] clean cart error:", err.Error())
		return nil, err
	}

	// create the transaction.
	payReq := &paypb.ReqCharge{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CardInfo: &paypb.CreditCardInfo{
			CreditCardNumber:          req.CardInfo.CreditCardNumber,
			CreditCardCvv:             req.CardInfo.CreditCardCvv,
			CreditCardExpirationMonth: req.CardInfo.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CardInfo.CreditCardExpirationYear,
		},
	}
	payRet, err := PayClient().Charge(ctx, payReq)
	if err != nil {
		glog.Errorln("[CheckoutServer] pay charge error:", err.Error())
		return nil, err
	}

	// sub the stock of product
	_, err = StockClient().SubStock(ctx, &stockpb.ReqSubStock{
		Stocks: stockItems,
	})
	if err != nil {
		glog.Errorln("[CheckoutServer] stock sub error:", err.Error())
		return nil, err
	}

	glog.Infof("[Checkoutserver] %v checkout success\n", req.UserId)

	return &checkoutpb.RspCheckout{
		OrderId:       orderId,
		TransactionId: payRet.TransactionId,
	}, nil
}
