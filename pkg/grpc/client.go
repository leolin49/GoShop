package grpc

import (
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IGrpcClient interface {
	Start() error
	Close() error
	GetInstance() interface{}
}

type GrpcClient struct {
	ServiceName string
	Address     string
	Conn        *grpc.ClientConn
	Client      interface{}
	NewClient   func(conn *grpc.ClientConn) interface{}
}

func (c *GrpcClient) Start() (err error) {
	c.Conn, err = grpc.NewClient(c.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorf("New gRPC [%s] client [%s] to failed: %v\n", c.ServiceName, c.Address, err)
		return
	}
	c.Client = c.NewClient(c.Conn)
	glog.Infof("Connect to [%s] server on [%s] successful\n", c.ServiceName, c.Address)
	return
}

func (c *GrpcClient) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *GrpcClient) GetInstance() interface{} {
	return c.Client
}
