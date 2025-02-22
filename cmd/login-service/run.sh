#!/bin/bash

logRoot="./log/"

go build -o loginserver
if [ $? -eq 0 ];then
	echo "loginserver compile success!"
	if [ ! -d "$logRoot" ];then
		mkdir "$logRoot"
	fi
	./loginserver -log_dir="$logRoot" -alsologtostderr=true
else
	echo "loginserver compile success!"
fi
	
