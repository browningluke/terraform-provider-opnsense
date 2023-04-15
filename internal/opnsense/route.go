package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// Data structs

type Route struct {
	Disabled    string `json:"disabled"`
	Description string `json:"descr"`
	Gateway     string `json:"gateway"`
	Network     string `json:"network"`
}

// Response structs

type RouteGetResp struct {
	Route struct {
		Disabled    string `json:"disabled"`
		Description string `json:"descr"`
		Gateway     map[string]struct {
			Value    string `json:"value"`
			Selected int    `json:"selected"`
		} `json:"gateway"`
		Network string `json:"network"`
	} `json:"route"`
}

type routeAddResp struct {
	Result      string                 `json:"result"`
	UUID        string                 `json:"uuid"`
	Validations map[string]interface{} `json:"validations,omitempty"`
}

// CRUD operations

func (c *Client) AddRoute(ctx context.Context, route *Route) (string, error) {
	return c.makeRoute(ctx, route, "/routes/routes/addroute")
}

func (c *Client) GetRoute(ctx context.Context, id string) (*Route, error) {
	// Make get request to OPNsense
	respJson := &RouteGetResp{}
	err := c.doRequest(ctx, "GET",
		fmt.Sprintf("/routes/routes/getroute/%s", id), nil, respJson)

	// Handle errors
	if err != nil {
		// Handle unmarshal error (means ID is invalid, or was deleted upstream)
		if err.Error() == fmt.Sprintf("json: cannot unmarshal array into Go value of type %s",
			reflect.TypeOf(respJson).Elem().String()) {
			return nil, fmt.Errorf("unable to find resource. it may have been deleted upstream")
		}

		return nil, err
	}

	// Find selected gateway
	gateway := ""
	for k, v := range respJson.Route.Gateway {
		if v.Selected == 1 {
			gateway = k
		}
	}

	return &Route{
		Disabled:    respJson.Route.Disabled,
		Description: respJson.Route.Description,
		Gateway:     gateway,
		Network:     respJson.Route.Network,
	}, nil
}

func (c *Client) UpdateRoute(ctx context.Context, id string, route *Route) error {
	_, err := c.makeRoute(ctx, route, fmt.Sprintf("/routes/routes/setroute/%s", id))
	return err
}

func (c *Client) DeleteRoute(ctx context.Context, id string) error {
	// Since routes has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.routeMu.Lock()
	defer c.routeMu.Unlock()

	// Make delete request to OPNsense
	respJson := &unboundResp{}
	err := c.doRequest(ctx, "POST",
		fmt.Sprintf("/routes/routes/delroute/%s", id), nil, respJson)
	if err != nil {
		return err
	}

	// Validate that override was deleted
	if respJson.Result != "deleted" {
		return fmt.Errorf("route not deleted. result: %s", respJson.Result)
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
func (c *Client) makeRoute(ctx context.Context, route *Route, endpoint string) (string, error) {
	// Since routes has to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.routeMu.Lock()
	defer c.routeMu.Unlock()

	// Make request to OPNsense
	respJson := &routeAddResp{}
	err := c.doRequest(ctx, "POST", endpoint,
		map[string]*Route{
			"route": route,
		},
		respJson,
	)
	if err != nil {
		return "", err
	}

	// Validate result
	if respJson.Result != "saved" {
		return "", fmt.Errorf("route not changed. result: %s. errors: %s", respJson.Result, respJson.Validations)
	}

	// Reconfigure (i.e. restart) the unbound resolver
	err = c.reconfigureRoutes(ctx)
	if err != nil {
		return "", err
	}

	return respJson.UUID, nil
}

func (c *Client) reconfigureRoutes(ctx context.Context) error {
	// Send reconfigure request to OPNsense
	respJson := &struct {
		Status string `json:"status"`
	}{}
	err := c.doRequest(ctx, "POST", "/routes/routes/reconfigure", nil, respJson)
	if err != nil {
		return err
	}

	// Validate unbound restarted correctly
	if respJson.Status != "ok" {
		return fmt.Errorf("routes reconfigure failed. status: %s", respJson.Status)
	}

	return nil
}
