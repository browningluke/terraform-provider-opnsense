package opnsense

import (
	"sync"
)

// Interfaces controller
type interfaces struct {
	client *Client
	mu     *sync.Mutex
}

func newInterfaces(c *Client) *interfaces {
	return &interfaces{
		client: c,
		mu:     &sync.Mutex{},
	}
}

func (i *interfaces) Client() *Client {
	return i.client
}

func (i *interfaces) Mutex() *sync.Mutex {
	return i.mu
}
