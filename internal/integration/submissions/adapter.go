package submissionsingegration

import "fmt"

type Adapter struct {
	client *client
}

func NewAdapter(cfg Config) (*Adapter, error) {
	cl, err := newClient(cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("create adapter: %w", err)
	}

	return &Adapter{
		client: cl,
	}, nil
}
