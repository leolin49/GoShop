package main

import (
	"context"
	paypb "goshop/api/protobuf/pay"
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

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49400")
	if err != nil {
		glog.Fatalf("[PayServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	paypb.RegisterPayServiceServer(rpcServer, new(PayRpcService))
	glog.Infoln("[PayServer] Starting rpc server on :49400")
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
		Cvv: strconv.Itoa(int(req.CardInfo.CreditCardCvv)),
		Month: strconv.Itoa(int(req.CardInfo.CreditCardExpirationMonth)),
		Year: strconv.Itoa(int(req.CardInfo.CreditCardExpirationYear)),
	}

	err := card.Validate(true)
	if err != nil {
		return nil, err
	}
	transactionId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	err = models.NewPaymentLogQuery(Mysql()).CreatePaymentLog(&models.PaymentLog{
		UserId: req.UserId,
		OrderId: req.OrderId,
		TransactionId: transactionId.String(),
		Amount: req.Amount,
		PayAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &paypb.RspCharge{
		TransactionId: transactionId.String(),
	}, nil
}

