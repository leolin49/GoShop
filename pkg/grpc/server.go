package grpc

type IGrpcServer interface {
	Start(address string) error
	Stop() error
}
