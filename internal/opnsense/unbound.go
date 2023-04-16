package opnsense

import (
	"sync"
)

const unboundReconfigureEndpoint = "/unbound/service/reconfigure"

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
