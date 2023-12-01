package marqo

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/imroc/req/v3"
)

// options for the client
type Options func(*Client)

// WithLogger sets the logger for the client
func WithLogger(logger *slog.Logger) func(*Client) {
	return func(c *Client) {
		c.logger = logger
	}
}

// Client is the client for the marqo server
type Client struct {
	url       string
	logger    *slog.Logger
	reqClient *req.Client
}

// NewClient creates a new client
func NewClient(url string, opt ...Options) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}

	client := &Client{
		url: url,
	}

	for _, o := range opt {
		o(client)
	}

	// set req client
	if client.reqClient == nil {
		client.reqClient = req.NewClient()
		client.reqClient.BaseURL = url
	}

	// set default logger
	if client.logger == nil {
		client.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		}))
	}
	return client, nil
}
