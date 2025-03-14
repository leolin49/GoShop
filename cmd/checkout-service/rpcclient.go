package main

import (
	"github.com/golang/glog"
)

func rpcClientsStart() {

	clients := []struct {
		name      string
		startFunc func() error
	}{
		{"product", ProductClientStart},
		{"cart", CartClientStart},
		{"order", OrderClientStart},
		{"pay", PayClientStart},
		{"stock", StockClientStart},
	}

	for _, client := range clients {
		go func(clientName string, clientFunc func() error) {
			if err := clientFunc(); err != nil {
				return
			}
		}(client.name, client.startFunc)
	}

}

func rpcClientsClose() {

	clients := []struct {
		name      string
		closeFunc func() error
	}{
		{"product", ProductClientClose},
		{"cart", CartClientClose},
		{"order", OrderClientClose},
		{"pay", PayClientClose},
		{"stock", StockClientClose},
	}

	for _, client := range clients {
		go func(clientName string, clientFunc func() error) {
			if err := clientFunc(); err != nil {
				glog.Errorf("[CheckoutServer] %s rpc client start failed: %s\n", clientName)
				return
			}
		}(client.name, client.closeFunc)
	}

}
