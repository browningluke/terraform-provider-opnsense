package opnsense

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

func (c *Client) doRequest(method, endpoint string, body any, resp any) error {
	// Build request body
	var bodyBuf io.Reader = nil

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyBuf = bytes.NewBuffer(bodyBytes)
	}

	// Create request
	req, err := http.NewRequest(method, fmt.Sprintf("%s/api%s", c.opts.Uri, endpoint), bodyBuf)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.getAuth()))

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Do request
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for 200
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code non-200; status code %d", res.StatusCode)
	}

	// Unmarshal resp JSON data to struct
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return err
	}

	return nil
}
