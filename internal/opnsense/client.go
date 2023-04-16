package opnsense

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"
)

type Client struct {
	client *http.Client

	// Controllers
	Unbound *unbound

	// Mutexes
	routeMu *sync.Mutex

	opts Options
}

type controller interface {
	Client() *Client
	Mutex() *sync.Mutex
	Reconfigure(ctx context.Context) error
}

type Options struct {
	Uri           string
	APIKey        string
	APISecret     string
	AllowInsecure bool
}

func NewClient(options Options) *Client {
	client := &Client{
		client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: options.AllowInsecure},
		}},

		// Mutexes
		routeMu: &sync.Mutex{},
		opts:    options,
	}

	// Add controllers
	client.Unbound = newUnbound(client)

	return client
}

// Requests

func (c *Client) getAuth() string {
	auth := c.opts.APIKey + ":" + c.opts.APISecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, resp any) error {
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

	// Log request
	dReq, _ := httputil.DumpRequest(req, true)
	tflog.Info(ctx, fmt.Sprintf("\n%s\n%s\n", string(dReq), bodyBuf))

	// Do request
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Log response
	dRes, _ := httputil.DumpResponse(res, true)
	tflog.Info(ctx, fmt.Sprintf("\n%s\n", string(dRes)))

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
