package opnsense

import (
	"context"
	"fmt"
)

const (
	interfacesVlanReconfigureEndpoint = "/interfaces/vlan_settings/reconfigure"
	interfacesVlanAddEndpoint         = "/interfaces/vlan_settings/addItem"
	interfacesVlanGetEndpoint         = "/interfaces/vlan_settings/getItem"
	interfacesVlanUpdateEndpoint      = "/interfaces/vlan_settings/setItem"
	interfacesVlanDeleteEndpoint      = "/interfaces/vlan_settings/delItem"
)

// Data structs

type InterfacesVlan struct {
	Description string      `json:"descr"`
	Tag         string      `json:"tag"`
	Priority    SelectedMap `json:"pcp"`
	Parent      SelectedMap `json:"if"`
	Device      string      `json:"vlanif"`
}

// CRUD operations

func (i *interfaces) AddVlan(ctx context.Context, vlan *InterfacesVlan) (string, error) {
	return makeSetFunc(i, interfacesVlanAddEndpoint, interfacesVlanReconfigureEndpoint)(ctx,
		map[string]*InterfacesVlan{
			"vlan": vlan,
		},
	)
}

func (i *interfaces) GetVlan(ctx context.Context, id string) (*InterfacesVlan, error) {
	get, err := makeGetFunc(i.Client(), interfacesVlanGetEndpoint,
		&struct {
			Vlan InterfacesVlan `json:"vlan"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Vlan, nil
}

func (i *interfaces) UpdateVlan(ctx context.Context, id string, vlan *InterfacesVlan) error {
	_, err := makeSetFunc(i, fmt.Sprintf("%s/%s", interfacesVlanUpdateEndpoint, id),
		interfacesVlanReconfigureEndpoint)(ctx,
		map[string]*InterfacesVlan{
			"vlan": vlan,
		},
	)
	return err
}

func (i *interfaces) DeleteVlan(ctx context.Context, id string) error {
	return makeDeleteFunc(i, interfacesVlanDeleteEndpoint, interfacesVlanReconfigureEndpoint)(ctx, id)
}
