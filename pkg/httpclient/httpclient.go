package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type Config struct {
	Timeout time.Duration `json:"timeout"`
}

type clientImpl struct {
	*http.Client
}

func New(cfg *Config) *clientImpl {
	return &clientImpl{
		Client: &http.Client{
			Timeout: cfg.Timeout * time.Second,
		},
	}
}

func (c *clientImpl) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}
	return c.Do(req)
}
