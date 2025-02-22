#!/bin/bash

logRoot="./log/"

go build -o productserver 
if [ $? -eq 0 ];then
	echo "productserver compile success!"
	if [ ! -d "$logRoot" ];then
		mkdir "$logRoot"
	fi
	./productserver -log_dir="$logRoot" -alsologtostderr=true
else
	echo "productserver compile success!"
fi
	
