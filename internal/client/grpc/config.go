package grpcclient

import "google.golang.org/grpc"

type Config struct {
	Address          string
	ExtraDialOptions []grpc.DialOption
}

func NewConfig(address string, opts ...grpc.DialOption) Config {
	return Config{
		Address:          address,
		ExtraDialOptions: opts,
	}
}
