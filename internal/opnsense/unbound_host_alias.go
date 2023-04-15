package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// Data structs

type UnboundHostAlias struct {
	Enabled     string `json:"enabled"`
	Host        string `json:"host"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
}

// Response structs

type unboundHostAliasGetResp struct {
	Alias struct {
		Enabled     string `json:"enabled"`
		Hostname    string `json:"hostname"`
		Domain      string `json:"domain"`
		Description string `json:"description"`
	} `json:"alias"`
}

// CRUD operations

func (c *Client) UnboundAddHostAlias(ctx context.Context, alias *UnboundHostAlias) (string, error) {
	return c.unboundMakeHostAlias(ctx, alias, "/unbound/settings/addHostAlias")
}

func (c *Client) UnboundGetHostAlias(ctx context.Context, id string) (*UnboundHostAlias, error) {
	// Make get request to OPNsense
	respJson := &unboundHostAliasGetResp{}
	err := c.doRequest(ctx, "GET",
		fmt.Sprintf("/unbound/settings/getHostAlias/%s", id), nil, respJson)

	// Handle errors
	if err != nil {
		// Handle unmarshal error (means ID is invalid, or was deleted upstream)
		if err.Error() == fmt.Sprintf("json: cannot unmarshal array into Go value of type %s",
			reflect.TypeOf(respJson).Elem().String()) {
			return nil, fmt.Errorf("unable to find resource. it may have been deleted upstream")
		}

		return nil, err
	}

	return &UnboundHostAlias{
		Enabled:     respJson.Alias.Enabled,
		Hostname:    respJson.Alias.Hostname,
		Domain:      respJson.Alias.Domain,
		Description: respJson.Alias.Description,
	}, nil
}

func (c *Client) UnboundUpdateHostAlias(ctx context.Context, id string, alias *UnboundHostAlias) error {
	_, err := c.unboundMakeHostAlias(ctx, alias, fmt.Sprintf("/unbound/settings/setHostAlias/%s", id))
	return err
}

func (c *Client) UnboundDeleteHostAlias(ctx context.Context, id string) error {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make delete request to OPNsense
	respJson := &unboundResp{}
	err := c.doRequest(ctx, "POST",
		fmt.Sprintf("/unbound/settings/delHostAlias/%s", id), nil, respJson)
	if err != nil {
		return err
	}

	// Validate that override was deleted
	if respJson.Result != "deleted" {
		return fmt.Errorf("alias not deleted. result: %s", respJson.Result)
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
func (c *Client) unboundMakeHostAlias(ctx context.Context, alias *UnboundHostAlias, endpoint string) (string, error) {
	// Since unbound has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.unboundMu.Lock()
	defer c.unboundMu.Unlock()

	// Make 'add' request to OPNsense
	respJson := &unboundAddResp{}
	err := c.doRequest(ctx, "POST", endpoint,
		map[string]*UnboundHostAlias{
			"alias": alias,
		},
		respJson,
	)
	if err != nil {
		return "", err
	}

	// Validate result
	if respJson.Result != "saved" {
		return "", fmt.Errorf("alias not changed. result: %s", respJson.Result)
	}

	// Reconfigure (i.e. restart) the unbound resolver
	err = c.reconfigureUnbound(ctx)
	if err != nil {
		return "", err
	}

	return respJson.UUID, nil
}
