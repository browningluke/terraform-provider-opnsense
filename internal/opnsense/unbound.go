package opnsense

import (
	"context"
	"fmt"
)

// Response structs

type unboundResp struct {
	Result string `json:"result"`
}

type unboundAddResp struct {
	Result      string            `json:"result"`
	UUID        string            `json:"uuid"`
	Validations map[string]string `json:"validations,omitempty"`
}

// Helper functions

func (c *Client) reconfigureUnbound(ctx context.Context) error {
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
		return fmt.Errorf("unbound reconfigure failed. status: %s", respJson.Status)
	}

	return nil
}
