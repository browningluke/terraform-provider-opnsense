package opnsense

import (
	"net/http"
)

type Client struct {
	client *http.Client

	opts Options
}

type Options struct {
	Uri       string
	APIKey    string
	APISecret string
}

func NewClient(options Options) *Client {
	return &Client{
		client: &http.Client{},
		opts:   options,
	}
}
