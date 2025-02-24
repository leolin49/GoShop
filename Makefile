GO ?= go
PROTOC ?= protoc
PROTOC_GEN_GO ?= protoc-gen-go
PROTOC_GEN_GO_GRPC ?= protoc-gen-go-grpc

ROOT := $(shell pwd)
API := $(ROOT)/api/protobuf
BIN := $(ROOT)/bin
CMD := $(ROOT)/cmd

SERVICES := \
	login-service \
	product-service \
	cart-service \
	auth-service \
	pay-service \
	checkout-service \
	order-service \
	gateway-service

.PHONY: all
all: build

.PHONY: build
build: $(addprefix $(BIN)/,$(SERVICES))

$(BIN)/%: $(CMD)/%
	@echo "Building $*..."
	@$(GO) build -o $@ $<	
	@echo "Building $* successfully."

.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	@$(PROTOC) \
		--proto_path=$(API) \
		--go_out=.. \
		--go-grpc_out=.. \
		$(API)/*.proto
	@echo "Generating protobuf files successfully." 

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN)/*
	@echo "Cleaning up successfully."

.PHONY: run
run: $(addprefix $(BIN)/,$(SERVICES))
	@echo "Running all services..."
	$(foreach service,$(SERVICES), \
		nohup $(BIN)/$(service) -log_dir=$(BIN) > /dev/null 2>&1 & \
		pgrep -f $(service) && echo "$(service) is running" || echo "$(service) failed to start"; \
		sleep 2; \
	)
	@echo "Listed services are running:"
	@ps -a | grep service

.PHONY: stop
stop:
	@echo "Stopping all services..."
	$(foreach service,$(SERVICES), pkill -f $(service) || true;)
	@echo "All services has been stopped."

.PHONY: restart
restart: stop run

.PHONY: gateway
gateway: $(BIN)/gateway-service
	@echo "Running gateway-service..."
	@nohup $(BIN)/gateway-service -log_dir=$(BIN) > /dev/null 2>&1 &
	@pgrep -f gateway-service && echo "gateway-service is running" || echo "gateway-service failed to start"

.PHONY: login
login: $(BIN)/login-service
	@echo "Running login-service..."
	@nohup $(BIN)/login-service -log_dir=$(BIN) > /dev/null 2>&1 &
	@pgrep -f login-service && echo "login-service is running" || echo "login-service failed to start"

.PHONY: product
product: $(BIN)/product-service
	@echo "Running product-service..."
	@nohup $(BIN)/product-service -log_dir=$(BIN) > /dev/null 2>&1 &
	@pgrep -f product-service && echo "product-service is running" || echo "product-service failed to start"

.PHONY: cart 
cart: $(BIN)/cart-service
	@echo "Running cart-service..."
	@nohup $(BIN)/cart-service -log_dir=$(BIN) > /dev/null 2>&1 &
	@pgrep -f cart-service && echo "cart-service is running" || echo "cart-service failed to start"

.PHONY: auth 
auth: $(BIN)/auth-service
	@echo "Running auth-service..."
	@nohup $(BIN)/auth-service -log_dir=$(BIN) > /dev/null 2>&1 &
	@pgrep -f auth-service && echo "auth-service is running" || echo "auth-service failed to start"

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "Code formatted successfully."

.PHONY: lint
lint:
	@golangci-lint run ./...


