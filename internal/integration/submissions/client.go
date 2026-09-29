package submissionsingegration

import (
	"context"
	"fmt"

	grpcclient "github.com/markgredasov/viewcode-sandbox/internal/client/grpc"
	submissions_v1 "github.com/markgredasov/viewcode-submissions/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	address string
	conn    *grpcclient.Conn
	api     submissions_v1.SubmissionsServiceClient
}

func newClient(addr string) (*client, error) {
	cfg := grpcclient.NewConfig(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (c *client) CreateSubmission(ctx context.Context, in *submissions_v1.CreateRequest) (*submissions_v1.CreateResponse, error) {
	return c.api.Create(ctx, in)
}

func (c *client) SetRunning(ctx context.Context, in *submissions_v1.SetRunningRequest) (*submissions_v1.SetRunningResponse, error) {
	return c.api.SetRunning(ctx, in)
}

func (c *client) SetResult(ctx context.Context, in *submissions_v1.SetResultRequest) (*submissions_v1.SetResultResponse, error) {
	return c.api.SetResult(ctx, in)
}
