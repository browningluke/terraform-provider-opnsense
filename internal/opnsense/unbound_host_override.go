package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// Data structs

type UnboundHostOverride struct {
	Enabled     string `json:"enabled"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Type        string `json:"rr"`
	Server      string `json:"server"`
	MXPriority  string `json:"mxprio"`
	MXDomain    string `json:"mx"`
	Description string `json:"description"`
}

// Response structs

type unboundHostOverrideGetResp struct {
	Host struct {
		Enabled  string `json:"enabled"`
		Hostname string `json:"hostname"`
		Domain   string `json:"domain"`
		Type     map[string]struct {
			Value    string `json:"value"`
			Selected int    `json:"selected"`
		} `json:"rr"`
		MXPriority  string `json:"mxprio"`
		MXDomain    string `json:"mx"`
		Server      string `json:"server"`
		Description string `json:"description"`
	} `json:"host"`
}

// CRUD operations

func (c *Client) UnboundAddHostOverride(ctx context.Context, host *UnboundHostOverride) (string, error) {
	return c.unboundMakeHostOverride(ctx, host, "/unbound/settings/addHostOverride")
}

func (c *Client) UnboundGetHostOverride(ctx context.Context, id string) (*UnboundHostOverride, error) {
	// Make get request to OPNsense
	respJson := &unboundHostOverrideGetResp{}
	err := c.doRequest(ctx, "GET",
		fmt.Sprintf("/unbound/settings/getHostOverride/%s", id), nil, respJson)

	// Handle errors
	if err != nil {
		// Handle unmarshal error (means ID is invalid, or was deleted upstream)
		if err.Error() == fmt.Sprintf("json: cannot unmarshal array into Go value of type %s",
			reflect.TypeOf(respJson).Elem().String()) {
			return nil, fmt.Errorf("unable to find resource. it may have been deleted upstream")
		}

		return nil, err
	}

	// Find selected type
	rr := ""
	for k, v := range respJson.Host.Type {
		if v.Selected == 1 {
			rr = k
		}
	}

	return &UnboundHostOverride{
		Enabled:     respJson.Host.Enabled,
		Hostname:    respJson.Host.Hostname,
		Domain:      respJson.Host.Domain,
		Type:        rr,
		Server:      respJson.Host.Server,
		MXPriority:  respJson.Host.MXPriority,
		MXDomain:    respJson.Host.MXDomain,
		Description: respJson.Host.Description,
	}, nil
}

func (c *Client) UnboundUpdateHostOverride(ctx context.Context, id string, host *UnboundHostOverride) error {
	_, err := c.unboundMakeHostOverride(ctx, host, fmt.Sprintf("/unbound/settings/setHostOverride/%s", id))
	return err
}

func (c *Client) UnboundDeleteHostOverride(ctx context.Context, id string) error {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make delete request to OPNsense
	respJson := &unboundResp{}
	err := c.doRequest(ctx, "POST",
		fmt.Sprintf("/unbound/settings/delHostOverride/%s", id), nil, respJson)
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
func (c *Client) unboundMakeHostOverride(ctx context.Context, host *UnboundHostOverride, endpoint string) (string, error) {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make request to OPNsense
	respJson := &unboundAddResp{}
	err := c.doRequest(ctx, "POST", endpoint,
		map[string]*UnboundHostOverride{
			"host": host,
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
