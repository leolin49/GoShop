# build
FROM golang:1.23.7-alpine AS builder
LABEL MAINTAINER="linyf49@qq.com"

ENV APP_HOME=/goshop
RUN mkdir $APP_HOME
WORKDIR $APP_HOME 

COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache make

RUN make build

# run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /goshop/bin/ ./bin/

COPY run.sh .

RUN chmod +x run.sh

CMD echo "OvO... Hello, welcome to goshop container ^_^"

