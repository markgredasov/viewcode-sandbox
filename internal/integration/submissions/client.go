package submissionsingegration

import (
	"fmt"

	grpcclient "github.com/markgredasov/viewcode-sandbox/internal/client/grpc"
	submissions_v1 "github.com/markgredasov/viewcode-submissions/pkg/api/v1"
)

type client struct {
	address string
	conn    *grpcclient.Conn
	api     submissions_v1.SubmissionsServiceClient
}

func newClient(addr string) (*client, error) {
	cfg := grpcclient.NewConfig(addr)
	conn, err := grpcclient.NewConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s: %w", addr, err)
	}

	return &client{
		address: addr,
		conn:    conn,
		api:     submissions_v1.NewSubmissionsServiceClient(conn.Conn()),
	}, nil
}
