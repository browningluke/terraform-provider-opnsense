package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// Data structs

type UnboundDomainOverride struct {
	Enabled     string `json:"enabled"`
	Domain      string `json:"domain"`
	Server      string `json:"server"`
	Description string `json:"description"`
}

// Response structs

type unboundDomainOverrideGetResp struct {
	Domain UnboundDomainOverride `json:"domain"`
}

// CRUD operations

func (c *Client) UnboundAddDomainOverride(ctx context.Context, domain *UnboundDomainOverride) (string, error) {
	return c.unboundMakeDomainOverride(ctx, domain, "/unbound/settings/addDomainOverride")
}

func (c *Client) UnboundGetDomainOverride(ctx context.Context, id string) (*UnboundDomainOverride, error) {
	// Make get request to OPNsense
	respJson := &unboundDomainOverrideGetResp{}
	err := c.doRequest(ctx, "GET",
		fmt.Sprintf("/unbound/settings/getDomainOverride/%s", id), nil, respJson)

	// Handle errors
	if err != nil {
		// Handle unmarshal error (means ID is invalid, or was deleted upstream)
		if err.Error() == fmt.Sprintf("json: cannot unmarshal array into Go value of type %s",
			reflect.TypeOf(respJson).Elem().String()) {
			return nil, fmt.Errorf("unable to find resource. it may have been deleted upstream")
		}

		return nil, err
	}

	return &respJson.Domain, nil
}

func (c *Client) UnboundUpdateDomainOverride(ctx context.Context, id string, domain *UnboundDomainOverride) error {
	_, err := c.unboundMakeDomainOverride(ctx, domain, fmt.Sprintf("/unbound/settings/setDomainOverride/%s", id))
	return err
}

func (c *Client) UnboundDeleteDomainOverride(ctx context.Context, id string) error {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make delete request to OPNsense
	respJson := &unboundResp{}
	err := c.doRequest(ctx, "POST",
		fmt.Sprintf("/unbound/settings/delDomainOverride/%s", id), nil, respJson)
	if err != nil {
		return err
	}

	// Validate that override was deleted
	if respJson.Result != "deleted" {
		return fmt.Errorf("override not deleted. result: %s", respJson.Result)
	}

	// Reconfigure (i.e. restart) the unbound resolver
	err = c.reconfigureUnbound(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

// unboundMakeHostOverride creates/updates a host override, depending on the endpoint parameter
func (c *Client) unboundMakeDomainOverride(ctx context.Context, domain *UnboundDomainOverride, endpoint string) (string, error) {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make request to OPNsense
	respJson := &unboundAddResp{}
	err := c.doRequest(ctx, "POST", endpoint,
		map[string]*UnboundDomainOverride{
			"domain": domain,
		},
		respJson,
	)
	if err != nil {
		return "", err
	}

	// Validate result
	if respJson.Result != "saved" {
		return "", fmt.Errorf("override not changed. result: %s", respJson.Result)
	}

	// Reconfigure (i.e. restart) the unbound resolver
	err = c.reconfigureUnbound(ctx)
	if err != nil {
		return "", err
	}

	return respJson.UUID, nil
}
