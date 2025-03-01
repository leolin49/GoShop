package main

import (
	"github.com/golang/glog"
)

func rpcClientsStart() {

	clients := []struct {
		name      string
		startFunc func() error
	}{
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
		{"stock", StockClientClose},
	}

	for _, client := range clients {
		go func(clientName string, clientFunc func() error) {
			if err := clientFunc(); err != nil {
				glog.Errorf("[TimeServer] %s rpc client start failed: %s\n", clientName)
				return
			}
		}(client.name, client.closeFunc)
	}

}
