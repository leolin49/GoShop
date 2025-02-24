package main

import (
	"context"
	checkoutpb "goshop/api/protobuf/checkout"
	paypb "goshop/api/protobuf/pay"
	"goshop/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	checkout_client checkoutpb.CheckoutServiceClient 
	checkout_conn *grpc.ClientConn
)

func CheckoutClientStart() bool {
	var err error
	// get address from consul register center.
	addr, err := service.ServiceRecover("checkout-service")
	if err != nil || addr == "" {
		glog.Errorln("[Gatewayserver] consul service recover failed.")
		return false
	}
	checkout_conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorln("[Gatewayserver] new cart rpc client error: ", err.Error())
		return false
	}
	checkout_client = checkoutpb.NewCheckoutServiceClient(checkout_conn) 
	glog.Infoln("[Gatewayserver] connect [checkout-service] server successful on: ", addr)
	return true
}

func CheckoutClient() checkoutpb.CheckoutServiceClient { return checkout_client }

func CheckoutClientClose() { checkout_conn.Close() }

func handleCheckout(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		invalidParam(c)	// TODO
		return
	}
	userId := user_id.(uint32)
	var (
		firstName = c.PostForm("first_name")
		lastName = c.PostForm("last_name")
		email = c.PostForm("email")
		// address
		streetAddress = c.PostForm("street")
		city = c.PostForm("city")
		state = c.PostForm("state")
		country = c.PostForm("country")
		zipCode, _ = getPostFormInt(c, "zip_code")
		// card
		cardNumber = c.PostForm("card_number")
		cardCvv, _ = getPostFormInt(c, "card_cvv")
		cardExpMonth, _ = getPostFormInt(c, "card_exp_month")
		cardExpYear, _ = getPostFormInt(c, "card_exp_year")
	)
	if firstName == "" || lastName == "" || email == "" {
		invalidParam(c)
		return
	}
	address := &checkoutpb.Address{
		StreetAddress: streetAddress,
		City: city,
		State: state,
		Country: country,
		ZipCode: int32(zipCode),
	}
	cardInfo := &paypb.CreditCardInfo{
		CreditCardNumber: cardNumber,
		CreditCardCvv: int32(cardCvv),
		CreditCardExpirationMonth: int32(cardExpMonth),
		CreditCardExpirationYear: int32(cardExpYear),
	}

	ret, err := CheckoutClient().Checkout(context.Background(), &checkoutpb.ReqCheckout{
		UserId: userId,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Address: address,
		CardInfo: cardInfo,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ret)
}

