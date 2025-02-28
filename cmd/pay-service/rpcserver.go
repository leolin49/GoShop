package main

import (
	"context"
	"fmt"
	paypb "goshop/api/protobuf/pay"
	"goshop/configs"
	"goshop/models"
	"net"
	"strconv"
	"time"

	creditcard "github.com/durango/go-credit-card"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type PayRpcService struct {
	paypb.UnimplementedPayServiceServer
}

func rpcServerStart(cfg *configs.Config) bool {
	addr := fmt.Sprintf("%s:%s", cfg.PayCfg.Host, cfg.PayCfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("[PayServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	paypb.RegisterPayServiceServer(rpcServer, new(PayRpcService))
	glog.Infof("[PayServer] Starting rpc server on [%s]\n", addr)
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[PayServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *PayRpcService) Charge(ctx context.Context, req *paypb.ReqCharge) (*paypb.RspCharge, error) {
	card := creditcard.Card{
		Number: req.CardInfo.CreditCardNumber,
		Cvv:    strconv.Itoa(int(req.CardInfo.CreditCardCvv)),
		Month:  strconv.Itoa(int(req.CardInfo.CreditCardExpirationMonth)),
		Year:   strconv.Itoa(int(req.CardInfo.CreditCardExpirationYear)),
	}

	err := card.Validate(true)
	if err != nil {
		return nil, err
	}
	transactionId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	err = models.NewPaymentLogQuery(db).CreatePaymentLog(&models.PaymentLog{
		UserId:        req.UserId,
		OrderId:       req.OrderId,
		TransactionId: transactionId.String(),
		Amount:        req.Amount,
		PayAt:         time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &paypb.RspCharge{
		TransactionId: transactionId.String(),
	}, nil
}
