#!/bin/bash

BIN="$(pwd)/bin"
SERVICES=("login" "product" "cart" "auth" "pay" "checkout" "order" "stock" "gateway")

main() {
	case "$1" in
		start)
			if [ -z "$2" ]; then 
				run_all
			elif [[ " ${SERVICES[@]} " =~ " $2 " ]]; then
				run_service "$2"
			else
				echo "Unknow service: $2"
				exit 1
			fi
			;;
		stop)
			if [ -z "$2" ]; then
				stop_all
			elif [[ " ${SERVICES[@]} " =~ " $2 " ]]; then
				stop_service "$2"
			else
				echo "Unknow service: $2"
				exit 1
			fi
			;;
		restart)
			if [ -z "$2" ]; then
				stop_all
				sleep 2
				run_all
			elif [[ " ${SERVICES[@]} " =~ " $2 " ]]; then
				stop_service "$2"
				sleep 1
				run_service "$2"
			else
				echo "Unknow service: $2"
				exit 1
			fi
			;;
		*)
			echo "Usage: $0 {start|stop|restart} [service]"
			echo "Example:"
			echo "  $0 start            # Run all the services."
			echo "  $0 start login      # Run login-service."
			echo "  $0 stop             # Stop all the services."
			echo "  $0 stop cart        # Stop cart-service"
			echo "  $0 restart order    # Restart order-service"
			exit 1
			;;
	esac
}

run_gateway() {
	local service_name="gateway"
	for node in node1 node2 node3; do
		echo "running ${service_name}-service ${node}"
		mkdir -p "${BIN}/${service_name}-${node}"
		nohup "${BIN}/${service_name}-service" -node="${node}" -log_dir="${BIN}/${service_name}-${node}" > /dev/null 2>&1 &
		sleep 0.5
	done
	pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_login() {
    local service_name="login"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_product() {
    local service_name="product"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_auth() {
    local service_name="auth"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_cart() {
    local service_name="cart"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_pay() {
    local service_name="pay"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_checkout() {
    local service_name="checkout"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_order() {
    local service_name="order"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_stock() {
    local service_name="stock"
    echo "running ${service_name}"
    mkdir -p "${BIN}/${service_name}" 
    nohup "${BIN}/${service_name}-service" -log_dir="${BIN}/${service_name}" > /dev/null 2>&1 &
    sleep 1
    pgrep -f "${BIN}/${service_name}-service" > /dev/null && echo "${service_name} start success" || echo "${service_name} start failed"
}

run_all() {
	for service in "${SERVICES[@]}"; do
		run_service "${service}"
	done
	ps -a | grep "${BIN}/.*-service"
}

run_service() {
	service_name="$1"
	case "${service_name}" in
		gateway)
			run_gateway
			;;
		login)
			run_login
			;;
		product)
			run_product
			;;
		cart)
			run_cart
			;;
		auth)
			run_auth
			;;
		pay)
			run_pay
			;;
		order)
			run_order
			;;
		checkout)
			run_checkout
			;;
		stock)
			run_stock
			;;
		*)
			echo "Unknow service: $2"
			exit 1
			;;
	esac
}

stop_service() {
	service_name="$1"
	echo "Stopping ${service_name}-service..."
	pkill -9 -f "${BIN}/${service_name}-service" 
	echo "${service_name}-service already stopped."
}

stop_all() { 
	echo "Stopping all services..."
	for service in "${SERVICES[@]}"; do
		stop_service "${service}"
	done
	echo "All services already stopped."
}

main "$@"
