package opnsense

import (
	"encoding/base64"
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

// Requests

func (c *Client) getAuth() string {
	auth := c.opts.APIKey + ":" + c.opts.APISecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
