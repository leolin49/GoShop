package main

import "github.com/golang/glog"

func rpcClientsStart() {
	runClient := func(clientName string, clientFunc func() bool) {
		if !clientFunc() {
			glog.Errorf("[GatewayServer] %s rpc client start failed\n", clientName)
		}
	}
	go runClient("auth", AuthClientStart)
}
