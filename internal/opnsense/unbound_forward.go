package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// Data structs

type UnboundForward struct {
	Enabled  string `json:"enabled"`
	Domain   string `json:"domain"`
	Type     string `json:"type"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	VerifyCN string `json:"verify"`
}

// Response structs

type unboundForwardGetResp struct {
	Dot struct {
		Enabled string `json:"enabled"`
		Domain  string `json:"domain"`
		Type    map[string]struct {
			Value    string `json:"value"`
			Selected int    `json:"selected"`
		} `json:"type"`
		Server   string `json:"server"`
		Port     string `json:"port"`
		VerifyCN string `json:"verify"`
	} `json:"dot"`
}

// CRUD operations

func (c *Client) UnboundAddForward(ctx context.Context, forward *UnboundForward) (string, error) {
	return c.unboundMakeForward(ctx, forward, "/unbound/settings/addDot")
}

func (c *Client) UnboundGetForward(ctx context.Context, id string) (*UnboundForward, error) {
	// Make get request to OPNsense
	respJson := &unboundForwardGetResp{}
	err := c.doRequest(ctx, "GET",
		fmt.Sprintf("/unbound/settings/getDot/%s", id), nil, respJson)

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
	forwardType := ""
	for k, v := range respJson.Dot.Type {
		if v.Selected == 1 {
			forwardType = k
		}
	}

	return &UnboundForward{
		Enabled:  respJson.Dot.Enabled,
		Domain:   respJson.Dot.Domain,
		Type:     forwardType,
		Server:   respJson.Dot.Server,
		Port:     respJson.Dot.Port,
		VerifyCN: respJson.Dot.VerifyCN,
	}, nil
}

func (c *Client) UnboundUpdateForward(ctx context.Context, id string, forward *UnboundForward) error {
	_, err := c.unboundMakeForward(ctx, forward, fmt.Sprintf("/unbound/settings/setDot/%s", id))
	return err
}

func (c *Client) UnboundDeleteForward(ctx context.Context, id string) error {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make delete request to OPNsense
	respJson := &unboundResp{}
	err := c.doRequest(ctx, "POST",
		fmt.Sprintf("/unbound/settings/delDot/%s", id), nil, respJson)
	if err != nil {
		return err
	}

	// Validate that override was deleted
	if respJson.Result != "deleted" {
		return fmt.Errorf("forward not deleted. result: %s", respJson.Result)
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
func (c *Client) unboundMakeForward(ctx context.Context, forward *UnboundForward, endpoint string) (string, error) {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make request to OPNsense
	respJson := &unboundAddResp{}
	err := c.doRequest(ctx, "POST", endpoint,
		map[string]*UnboundForward{
			"dot": forward,
		},
		respJson,
	)
	if err != nil {
		return "", err
	}

	// Validate result
	if respJson.Result != "saved" {
		return "", fmt.Errorf("forward not changed. result: %s", respJson.Result)
	}

	// Reconfigure (i.e. restart) the unbound resolver
	err = c.reconfigureUnbound(ctx)
	if err != nil {
		return "", err
	}

	return respJson.UUID, nil
}
