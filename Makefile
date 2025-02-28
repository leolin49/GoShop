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
	   stock-service \
       gateway-service

.PHONY: all
all: build

.PHONY: build
build: $(addprefix $(BIN)/,$(SERVICES))

$(BIN)/%: $(CMD)/%
	@mkdir -p $(BIN)
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
run: $(SERVICES)
	@echo "Listed services are running:"
	@ps -a | grep service

.PHONY: stop
stop:
	@echo "Stopping all services..."
	$(foreach service,$(SERVICES), pkill -f $(service) || true;)
	@echo "All services has been stopped."

.PHONY: rebuild
rebuild: clean build 

.PHONY: gateway
gateway: $(BIN)/gateway-service

.PHONY: login
login: $(BIN)/login-service

.PHONY: product
product: $(BIN)/product-service

.PHONY: cart 
cart: $(BIN)/cart-service

.PHONY: auth 
auth: $(BIN)/auth-service

.PHONY: checkout 
checkout: $(BIN)/checkout-service

.PHONY: pay 
pay: $(BIN)/pay-service

.PHONY: order 
order: $(BIN)/order-service

.PHONY: stock 
stock: $(BIN)/stock-service

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "Code formatted successfully."

.PHONY: lint
lint:
	@golangci-lint run ./...


