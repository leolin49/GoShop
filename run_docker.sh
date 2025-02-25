#!/bin/bash

start_consul() {
    echo "Stopping and removing existing Consul container..."
    docker stop consul1
    docker rm consul1
    sleep 1

    echo "Starting Consul container..."
    docker run --name consul1 \
        -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600 \
        hashicorp/consul agent -server -bootstrap-expect=1 -ui -bind=0.0.0.0 -client=0.0.0.0
}

start_mysql() {
    echo "Stopping and removing existing MySQL container..."
    docker stop mysql
    docker rm mysql
    sleep 1

    echo "Starting MySQL container..."
    docker run --name mysql \
        --restart=always \
        --privileged=true \
        -p 3306:3306 \
        -v ~/mysql/logs:/logs \
        -v ~/mysql/data:/var/lib/mysql \
        -v ~/mysql/conf.d:/etc/mysql/conf.d \
        -e MYSQL_ROOT_PASSWORD='123' \
        -d mysql
}

start_nginx() {
    echo "Stopping and removing existing Nginx container..."
    docker stop nginx
    docker rm nginx
    sleep 1

    echo "Starting Nginx container..."
    docker run --name nginx \
        -p 80:80 \
        -v ~/docker-conf/nginx/conf/nginx.conf:/etc/nginx/nginx.conf \
        -v ~/docker-conf/nginx/conf/conf.d:/etc/nginx/conf.d \
        -v ~/docker-conf/nginx/log:/var/log/nginx \
        -v ~/docker-conf/nginx/html:/usr/share/nginx/html \
        -d nginx:latest
}

start_rabbitmq() {
    echo "Stopping and removing existing RabbitMQ container..."
    docker stop rabbitmq 
    docker rm rabbitmq 
    sleep 1

    echo "Starting RabbitMQ container..."
	docker run -id --name=rabbitmq \
	-p 5671:5671 \
	-p 5672:5672 \
	-p 4369:4369 \
	-p 15671:15671 \
	-p 15672:15672 \
	-p 25672:25672 \
	-e RABBITMQ_DEFAULT_USER=root \
	-e RABBITMQ_DEFAULT_PASS=goshop \
	rabbitmq:management
}

start_redis() {
    echo "Stopping and removing existing Redis container..."
    docker stop redis 
    docker rm redis 
    sleep 1
	
    echo "Starting Redis container..."
	docker run -itd --name redis \
		-p 6379:6379 \
		redis
}

case "$1" in
    consul)
        start_consul
        ;;
    mysql)
        start_mysql
        ;;
    nginx)
        start_nginx
        ;;
	rabbitmq)
		start_rabbitmq
		;;
	redis)
		start_redis
		;;
    all)
        start_consul
        start_mysql
        start_nginx
        ;;
    *)
        echo "Usage: $0 {consul|mysql|nginx|rabbitmq|redis|all}"
        exit 1
        ;;
esac
