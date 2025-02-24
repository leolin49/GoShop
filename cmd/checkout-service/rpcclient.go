package main

import (
	"github.com/golang/glog"
)

func rpcClientsStart() {
	runClient := func(clientName string, clientFunc func() bool) {
		if !clientFunc() {
			glog.Errorf("[CheckoutServer] %s rpc client start failed\n", clientName)
		}
	}
	go runClient("product", ProductClientStart)
	go runClient("cart", CartClientStart)
	go runClient("pay", PayClientStart)
	go runClient("order", OrderClientStart)
}

func rpcClientClose() {
	ProductClientClose()
	CartClientClose()
	PayClientClose()
	OrderClientClose()
}
