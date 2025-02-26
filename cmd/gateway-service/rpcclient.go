package main

import (
	"github.com/golang/glog"
)

// TODO

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
	closeClient := func(clientName string, clientFunc func() error) {
		if err := clientFunc(); err != nil {
			glog.Errorf("[GatewayServer] %s rpc client start failed: %s\n", clientName, clientName)
		}
	}
	go closeClient("login", LoginClientClose)
	go closeClient("product", ProductClientClose)
	go closeClient("cart", CartClientClose)
	go closeClient("auth", AuthClientClose)
	go closeClient("checkout", CheckoutClientClose)
}
