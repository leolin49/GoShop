#!/bin/bash

logRoot="./log/"

go build -o gatewayserver 
if [ $? -eq 0 ];then
	echo "gatewayserver compile success!"
	if [ ! -d "$logRoot" ];then
		mkdir "$logRoot"
	fi
	./gatewayserver -log_dir="$logRoot" -alsologtostderr=true
else
	echo "gatewayserver compile success!"
fi
	
