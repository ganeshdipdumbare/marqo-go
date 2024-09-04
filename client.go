package marqo

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/imroc/req/v3"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

// Options for the client
type Options func(*Client)

// WithLogger sets the logger for the client
func WithLogger(logger *slog.Logger) func(*Client) {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithMarqoCloudAuth sets the API key for authentication if you are using MarqoCloud
func WithMarqoCloudAuth(apiKey string) func(*Client) {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// Client is the client for the Marqo server
type Client struct {
	url       string
	logger    *slog.Logger
	reqClient *req.Client
	apiKey    string // Field to hold the API key for use with MarqoCloud
}

// NewClient creates a new client for the Marqo server.
//
// This method initializes a new client with the specified endpoint URL and optional parameters.
//
// Parameters:
//
//	url (string): The endpoint URL of your Marqo instance.
//	opt (...Options): Optional parameters for the client.
//
// Returns:
//
//	*Client: A new Marqo client instance.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the url parameter.
// 2. Initializes a new Client instance.
// 3. Applies the optional parameters to the client.
// 4. Sets the reqClient if not already set.
// 5. Returns the new client instance if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	client, err := marqo.NewClient("http://localhost:8882")
//	if err != nil {
//	    log.Fatalf("Failed to create client: %v", err)
//	}
//	fmt.Printf("Client: %+v\n", client)
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

		// Add the API key header if provided
		if client.apiKey != "" {
			client.reqClient.SetCommonHeader("x-api-key", client.apiKey)
		}
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