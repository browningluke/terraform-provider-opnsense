package opnsense

import (
	"context"
	"fmt"
	"sync"
)

// Unbound controller
type unbound struct {
	client *Client
	mu     *sync.Mutex
}

func newUnbound(c *Client) *unbound {
	return &unbound{
		client: c,
		mu:     &sync.Mutex{},
	}
}

func (u *unbound) Client() *Client {
	return u.client
}

func (u *unbound) Mutex() *sync.Mutex {
	return u.mu
}

func (u *unbound) Reconfigure(ctx context.Context) error {
	// Send reconfigure request to OPNsense
	respJson := &struct {
		Status string `json:"status"`
	}{}
	err := u.client.doRequest(ctx, "POST", "/unbound/service/reconfigure", nil, respJson)
	if err != nil {
		return err
	}

	// Validate unbound restarted correctly
	if respJson.Status != "ok" {
		return fmt.Errorf("unbound reconfigure failed. status: %s", respJson.Status)
	}

	return nil
}
