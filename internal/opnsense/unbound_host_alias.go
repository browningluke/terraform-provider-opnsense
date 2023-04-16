package opnsense

import (
	"context"
	"fmt"
)

const (
	unboundHostAliasAddEndpoint    = "/unbound/settings/addHostAlias"
	unboundHostAliasGetEndpoint    = "/unbound/settings/getHostAlias"
	unboundHostAliasUpdateEndpoint = "/unbound/settings/setHostAlias"
	unboundHostAliasDeleteEndpoint = "/unbound/settings/delHostAlias"
)

// Data structs

type UnboundHostAlias struct {
	Enabled     string      `json:"enabled"`
	Host        SelectedMap `json:"host"`
	Hostname    string      `json:"hostname"`
	Domain      string      `json:"domain"`
	Description string      `json:"description"`
}

// CRUD operations

func (u *unbound) AddHostAlias(ctx context.Context, alias *UnboundHostAlias) (string, error) {
	return makeSetFunc(u, unboundHostAliasAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundHostAlias{
			"alias": alias,
		},
	)
}

func (u *unbound) GetHostAlias(ctx context.Context, id string) (*UnboundHostAlias, error) {
	get, err := makeGetFunc(u.Client(), unboundHostAliasGetEndpoint,
		&struct {
			Alias UnboundHostAlias `json:"alias"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Alias, nil
}

func (u *unbound) UpdateHostAlias(ctx context.Context, id string, alias *UnboundHostAlias) error {
	_, err := makeSetFunc(u, fmt.Sprintf("%s/%s", unboundHostAliasUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundHostAlias{
			"alias": alias,
		},
	)
	return err
}

func (u *unbound) DeleteHostAlias(ctx context.Context, id string) error {
	return makeDeleteFunc(u, unboundHostAliasDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
