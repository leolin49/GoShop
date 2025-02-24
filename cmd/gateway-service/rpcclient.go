package main

import (
	"github.com/golang/glog"
)

func rpcClientsStart() {
	runClient := func(clientName string, clientFunc func() bool) {
		if !clientFunc() {
			glog.Errorf("[GatewayServer] %s rpc client start failed\n", clientName)
		}
	}
	go runClient("login", LoginClientStart)
	go runClient("product", ProductClientStart)
	go runClient("cart", CartClientStart)
	go runClient("auth", AuthClientStart)
	go runClient("checkout", CheckoutClientStart)
}

func rpcClientsClose() {
	LoginClientClose()
	ProductClient()
	CartClientClose()
	AuthClientClose()
	CheckoutClient()
}
