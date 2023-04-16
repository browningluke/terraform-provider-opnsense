package opnsense

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

const (
	clientMaxBackoff = 30
	clientMinBackoff = 1
	clientMaxRetries = 4
)

type Client struct {
	client *retryablehttp.Client

	// Controllers
	Routes     *routes
	Interfaces *interfaces
	Unbound    *unbound

	opts Options
}

type controller interface {
	Client() *Client
	Mutex() *sync.Mutex
}

type Options struct {
	Uri           string
	APIKey        string
	APISecret     string
	AllowInsecure bool

	// Retries
	MaxBackoff int64
	MinBackoff int64
	MaxRetries int64
}

func NewClient(options Options) *Client {
	client := &Client{
		client: retryablehttp.NewClient(),
		opts:   options,
	}

	// Configure HTTP client
	client.client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: options.AllowInsecure},
	}

	//   Set defaults for retries
	client.client.RetryWaitMax = clientMaxBackoff
	client.client.RetryWaitMin = clientMinBackoff
	client.client.RetryMax = clientMaxRetries

	//   Override defaults for retries, if set
	if options.MaxBackoff != 0 {
		client.client.RetryWaitMax = time.Duration(options.MaxBackoff) * time.Second
	}
	if options.MinBackoff != 0 {
		client.client.RetryWaitMin = time.Duration(options.MinBackoff) * time.Second
	}
	if options.MaxRetries != 0 {
		client.client.RetryMax = int(options.MaxRetries)
	}

	// Add controllers
	client.Routes = newRoutes(client)
	client.Interfaces = newInterfaces(client)
	client.Unbound = newUnbound(client)

	return client
}

// Requests

func (c *Client) getAuth() string {
	auth := c.opts.APIKey + ":" + c.opts.APISecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, resp any) error {
	// Create IO readers
	var bodyReader io.Reader
	if body != nil {
		// Marshal body into bytes
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		// Convert body bytes into reader
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := retryablehttp.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/api%s", c.opts.Uri, endpoint), bodyReader)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.getAuth()))
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Log request
	dReq, _ := httputil.DumpRequest(req.Request, true)
	tflog.Info(ctx, fmt.Sprintf("\n%s\n", string(dReq)))

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

// ReconfigureService defined at the endpoint.
func (c *Client) ReconfigureService(ctx context.Context, endpoint string) error {
	// Send reconfigure request to OPNsense
	respJson := &struct {
		Status string `json:"status"`
	}{}
	err := c.doRequest(ctx, "POST", "/unbound/service/reconfigure", nil, respJson)
	if err != nil {
		return err
	}

	// Validate unbound restarted correctly
	if respJson.Status != "ok" {
		return fmt.Errorf("reconfigure failed. status: %s", respJson.Status)
	}

	return nil
}
