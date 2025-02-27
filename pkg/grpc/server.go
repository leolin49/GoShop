package grpc

type GrpcServer interface {
	Start(address string) error
	Stop() error
}
