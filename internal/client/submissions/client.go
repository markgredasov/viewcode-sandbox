package submissionsclient

import (
	"context"
	"fmt"

	submissions_v1 "github.com/markgredasov/viewcode-submissions/pkg/api/v1"
	"google.golang.org/grpc"
)

type client struct {
	conn *grpc.ClientConn
	api  submissions_v1.SubmissionsServiceClient
}

func newClient(addr string, opts ...grpc.DialOption) (*client, error) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("create gRPC connection to %s: %w", addr, err)
	}
	return &client{
		conn: conn,
		api:  submissions_v1.NewSubmissionsServiceClient(conn),
	}, nil
}

func (c *client) close() error { return c.conn.Close() }

func (c *client) CreateSubmission(ctx context.Context, in *submissions_v1.CreateRequest) (*submissions_v1.CreateResponse, error) {
	return c.api.Create(ctx, in)
}

func (c *client) SetRunning(ctx context.Context, in *submissions_v1.SetRunningRequest) (*submissions_v1.SetRunningResponse, error) {
	return c.api.SetRunning(ctx, in)
}

func (c *client) SetResult(ctx context.Context, in *submissions_v1.SetResultRequest) (*submissions_v1.SetResultResponse, error) {
	return c.api.SetResult(ctx, in)
}
