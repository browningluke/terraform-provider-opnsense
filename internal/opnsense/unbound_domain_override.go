package opnsense

import (
	"context"
	"fmt"
)

const (
	unboundDomainOverrideAddEndpoint    = "/unbound/settings/addDomainOverride"
	unboundDomainOverrideGetEndpoint    = "/unbound/settings/getDomainOverride"
	unboundDomainOverrideUpdateEndpoint = "/unbound/settings/setDomainOverride"
	unboundDomainOverrideDeleteEndpoint = "/unbound/settings/delDomainOverride"
)

// Data structs

type UnboundDomainOverride struct {
	Enabled     string `json:"enabled"`
	Domain      string `json:"domain"`
	Server      string `json:"server"`
	Description string `json:"description"`
}

// CRUD operations

func (u *unbound) AddDomainOverride(ctx context.Context, domain *UnboundDomainOverride) (string, error) {
	return makeSetFunc(u, unboundDomainOverrideAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundDomainOverride{
			"domain": domain,
		},
	)
}

func (u *unbound) GetDomainOverride(ctx context.Context, id string) (*UnboundDomainOverride, error) {
	get, err := makeGetFunc(u.Client(), unboundDomainOverrideGetEndpoint,
		&struct {
			Domain UnboundDomainOverride `json:"domain"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Domain, nil
}

func (u *unbound) UpdateDomainOverride(ctx context.Context, id string, domain *UnboundDomainOverride) error {
	_, err := makeSetFunc(u, fmt.Sprintf("%s/%s", unboundDomainOverrideUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundDomainOverride{
			"domain": domain,
		},
	)
	return err
}

func (u *unbound) DeleteDomainOverride(ctx context.Context, id string) error {
	return makeDeleteFunc(u, unboundDomainOverrideDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
