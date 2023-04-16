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
}

func NewClient(options Options) *Client {
	client := &Client{
		client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: options.AllowInsecure},
		}},

		opts: options,
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
