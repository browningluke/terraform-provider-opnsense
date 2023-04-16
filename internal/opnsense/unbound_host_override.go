package opnsense

import (
	"context"
	"fmt"
)

const (
	unboundHostOverrideAddEndpoint    = "/unbound/settings/addHostOverride"
	unboundHostOverrideGetEndpoint    = "/unbound/settings/getHostOverride"
	unboundHostOverrideUpdateEndpoint = "/unbound/settings/setHostOverride"
	unboundHostOverrideDeleteEndpoint = "/unbound/settings/delHostOverride"
)

// Data structs

type UnboundHostOverride struct {
	Enabled     string      `json:"enabled"`
	Hostname    string      `json:"hostname"`
	Domain      string      `json:"domain"`
	Type        SelectedMap `json:"rr"`
	Server      string      `json:"server"`
	MXPriority  string      `json:"mxprio"`
	MXDomain    string      `json:"mx"`
	Description string      `json:"description"`
}

// CRUD operations

func (u *unbound) AddHostOverride(ctx context.Context, host *UnboundHostOverride) (string, error) {
	return makeSetFunc(u, unboundHostOverrideAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundHostOverride{
			"host": host,
		},
	)
}

func (u *unbound) GetHostOverride(ctx context.Context, id string) (*UnboundHostOverride, error) {
	get, err := makeGetFunc(u.Client(), unboundHostOverrideGetEndpoint,
		&struct {
			Host UnboundHostOverride `json:"host"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Host, nil
}

func (u *unbound) UpdateHostOverride(ctx context.Context, id string, host *UnboundHostOverride) error {
	_, err := makeSetFunc(u, fmt.Sprintf("%s/%s", unboundHostOverrideUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundHostOverride{
			"host": host,
		},
	)
	return err
}

func (u *unbound) DeleteHostOverride(ctx context.Context, id string) error {
	return makeDeleteFunc(u, unboundHostOverrideDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
