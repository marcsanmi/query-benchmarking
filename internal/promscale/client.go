package promscale

import (
	"fmt"
	"net/http"
	"net/url"
)

// Config is the configuration for the Promscale client.
type Config struct {
	URL string `default:"http://127.0.0.1:9201" required:"true"`
}

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, config Config) (Client, error) {
	baseURL, err := url.Parse(config.URL)
	if err != nil {
		return Client{}, fmt.Errorf("invalid promscale url: %v", err)
	}

	return Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}
