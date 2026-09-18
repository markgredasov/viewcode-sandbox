package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
)

type Conn struct {
	config Config
	conn   *grpc.ClientConn
}

func NewConn(cfg Config) (*Conn, error) {
	opts := cfg.ExtraDialOptions

	conn, err := grpc.NewClient(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("create gRPC connection to: %s", cfg.Address)
	}

	return &Conn{
		config: cfg,
		conn:   conn,
	}, nil
}

func (c *Conn) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
