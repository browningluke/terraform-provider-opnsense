package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// Data structs

type InterfacesVlan struct {
	Description string `json:"descr"`
	Tag         string `json:"tag"`
	Priority    string `json:"pcp"`
	Parent      string `json:"if"`
	Device      string `json:"vlanif"`
}

// Response structs

type interfacesVlanGetResp struct {
	Vlan struct {
		Description string `json:"descr"`
		Tag         string `json:"tag"`
		Priority    map[string]struct {
			Value    string `json:"value"`
			Selected int    `json:"selected"`
		} `json:"pcp"`
		Parent map[string]struct {
			Value    string `json:"value"`
			Selected int    `json:"selected"`
		} `json:"if"`
		Device string `json:"vlanif"`
	} `json:"vlan"`
}

// CRUD operations

func (c *Client) InterfacesAddVlan(ctx context.Context, vlan *InterfacesVlan) (string, error) {
	return c.interfacesMakeVlan(ctx, vlan, "/interfaces/vlan_settings/addItem")
}

func (c *Client) InterfacesGetVlan(ctx context.Context, id string) (*InterfacesVlan, error) {
	// Make get request to OPNsense
	respJson := &interfacesVlanGetResp{}
	err := c.doRequest(ctx, "GET",
		fmt.Sprintf("/interfaces/vlan_settings/getItem/%s", id), nil, respJson)

	// Handle errors
	if err != nil {
		// Handle unmarshal error (means ID is invalid, or was deleted upstream)
		if err.Error() == fmt.Sprintf("json: cannot unmarshal array into Go value of type %s",
			reflect.TypeOf(respJson).Elem().String()) {
			return nil, fmt.Errorf("unable to find resource. it may have been deleted upstream")
		}

		return nil, err
	}

	// Find selected pcp
	pcp := ""
	for k, v := range respJson.Vlan.Priority {
		if v.Selected == 1 {
			pcp = k
		}
	}

	// Find selected interface
	selectedInterface := ""
	for k, v := range respJson.Vlan.Parent {
		if v.Selected == 1 {
			selectedInterface = k
		}
	}

	return &InterfacesVlan{
		Description: respJson.Vlan.Description,
		Tag:         respJson.Vlan.Tag,
		Priority:    pcp,
		Parent:      selectedInterface,
		Device:      respJson.Vlan.Device,
	}, nil
}

func (c *Client) InterfacesUpdateVlan(ctx context.Context, id string, vlan *InterfacesVlan) error {
	_, err := c.interfacesMakeVlan(ctx, vlan, fmt.Sprintf("/interfaces/vlan_settings/setItem/%s", id))
	return err
}

func (c *Client) InterfacesDeleteVlan(ctx context.Context, id string) error {
	// Since VLANs have to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.routeMu.Lock()
	defer c.routeMu.Unlock()

	// Make delete request to OPNsense
	respJson := &unboundResp{}
	err := c.doRequest(ctx, "POST",
		fmt.Sprintf("/interfaces/vlan_settings/delItem/%s", id), nil, respJson)
	if err != nil {
		return err
	}

	// Validate that override was deleted
	if respJson.Result != "deleted" {
		return fmt.Errorf("vlan not deleted. result: %s", respJson.Result)
	}

	// Reconfigure the VLAN controller
	err = c.reconfigureUnbound(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func (c *Client) interfacesMakeVlan(ctx context.Context, vlan *InterfacesVlan, endpoint string) (string, error) {
	// Since VLANs have to be reconfigured after every change, locking the mutex prevents
	// the API from being written to while it's reconfiguring, which results in data loss.
	c.routeMu.Lock()
	defer c.routeMu.Unlock()

	// Make request to OPNsense
	respJson := &routeAddResp{}
	err := c.doRequest(ctx, "POST", endpoint,
		map[string]*InterfacesVlan{
			"vlan": vlan,
		},
		respJson,
	)
	if err != nil {
		return "", err
	}

	// Validate result
	if respJson.Result != "saved" {
		return "", fmt.Errorf("vlan not changed. result: %s. errors: %s", respJson.Result, respJson.Validations)
	}

	// Reconfigure the VLAN controller
	err = c.reconfigureInterfacesVlan(ctx)
	if err != nil {
		return "", err
	}

	return respJson.UUID, nil
}

func (c *Client) reconfigureInterfacesVlan(ctx context.Context) error {
	// Send reconfigure request to OPNsense
	respJson := &struct {
		Status string `json:"status"`
	}{}
	err := c.doRequest(ctx, "POST", "/interfaces/vlan_settings/reconfigure", nil, respJson)
	if err != nil {
		return err
	}

	// Validate unbound restarted correctly
	if respJson.Status != "ok" {
		return fmt.Errorf("vlan reconfigure failed. status: %s", respJson.Status)
	}

	return nil
}
