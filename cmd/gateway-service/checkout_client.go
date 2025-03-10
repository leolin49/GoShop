package main

import (
	"context"
	checkoutpb "goshop/api/protobuf/checkout"
	paypb "goshop/api/protobuf/pay"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	checkout_client checkoutpb.CheckoutServiceClient
	checkout_conn   *grpc.ClientConn
)

func CheckoutClientStart() error {
	var err error
	// get address from consul register center.
	addr, err := consul.ServiceRecover("checkout-service")
	if err != nil || addr == "" {
		glog.Errorln("[Gatewayserver] consul service recover failed.")
		return err
	}

	checkout_conn, err = grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		glog.Errorln("[Gatewayserver] new cart rpc client error: ", err.Error())
		return err
	}
	checkout_client = checkoutpb.NewCheckoutServiceClient(checkout_conn)
	glog.Infoln("[Gatewayserver] connect [checkout-service] server successful on: ", addr)
	return nil
}

func CheckoutClient() checkoutpb.CheckoutServiceClient {
	return checkout_client
}

func CheckoutClientClose() error {
	return checkout_conn.Close()
}

func handleFlashCheckout(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		invalidParam(c) // TODO
		return
	}
	userId := user_id.(uint32)
	var (
		firstName       = c.PostForm("first_name")
		lastName        = c.PostForm("last_name")
		email           = c.PostForm("email")
		streetAddress   = c.PostForm("street")
		city            = c.PostForm("city")
		state           = c.PostForm("state")
		country         = c.PostForm("country")
		zipCode, _      = getPostFormInt(c, "zip_code")
		cardNumber      = c.PostForm("card_number")
		cardCvv, _      = getPostFormInt(c, "card_cvv")
		cardExpMonth, _ = getPostFormInt(c, "card_exp_month")
		cardExpYear, _  = getPostFormInt(c, "card_exp_year")
		// diff
		product_id, _  = getPostFormInt(c, "product_id")
		flash_count, _ = getPostFormInt(c, "flash_count")
	)
	if firstName == "" || lastName == "" || email == "" {
		invalidParam(c)
		return
	}
	address := &checkoutpb.Address{
		StreetAddress: streetAddress,
		City:          city,
		State:         state,
		Country:       country,
		ZipCode:       int32(zipCode),
	}
	cardInfo := &paypb.CreditCardInfo{
		CreditCardNumber:          cardNumber,
		CreditCardCvv:             int32(cardCvv),
		CreditCardExpirationMonth: int32(cardExpMonth),
		CreditCardExpirationYear:  int32(cardExpYear),
	}

	req := &checkoutpb.ReqFlashCheckout{
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Address:   address,
		CardInfo:  cardInfo,
		ProductId: uint32(product_id),
		Count:     uint64(flash_count),
	}

	ret, err := CheckoutClient().FlashCheckout(context.Background(), req)
	if err != nil {
		rpcRequestError(c, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleCheckout(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		invalidParam(c) // TODO
		return
	}
	userId := user_id.(uint32)
	var (
		firstName       = c.PostForm("first_name")
		lastName        = c.PostForm("last_name")
		email           = c.PostForm("email")
		streetAddress   = c.PostForm("street")
		city            = c.PostForm("city")
		state           = c.PostForm("state")
		country         = c.PostForm("country")
		zipCode, _      = getPostFormInt(c, "zip_code")
		cardNumber      = c.PostForm("card_number")
		cardCvv, _      = getPostFormInt(c, "card_cvv")
		cardExpMonth, _ = getPostFormInt(c, "card_exp_month")
		cardExpYear, _  = getPostFormInt(c, "card_exp_year")
	)
	if firstName == "" || lastName == "" || email == "" {
		invalidParam(c)
		return
	}
	address := &checkoutpb.Address{
		StreetAddress: streetAddress,
		City:          city,
		State:         state,
		Country:       country,
		ZipCode:       int32(zipCode),
	}
	cardInfo := &paypb.CreditCardInfo{
		CreditCardNumber:          cardNumber,
		CreditCardCvv:             int32(cardCvv),
		CreditCardExpirationMonth: int32(cardExpMonth),
		CreditCardExpirationYear:  int32(cardExpYear),
	}

	req := &checkoutpb.ReqCheckout{
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Address:   address,
		CardInfo:  cardInfo,
	}

	go GatewayServerGetInstance().MQClient.PublishProtoMsgSimple(req)

	c.JSON(http.StatusOK, gin.H{
		"info": "checkout request finish",
	})
}
