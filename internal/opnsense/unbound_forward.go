package opnsense

import (
	"context"
	"fmt"
)

const (
	unboundForwardAddEndpoint    = "/unbound/settings/addDot"
	unboundForwardGetEndpoint    = "/unbound/settings/getDot"
	unboundForwardUpdateEndpoint = "/unbound/settings/setDot"
	unboundForwardDeleteEndpoint = "/unbound/settings/delDot"
)

// Data structs

type UnboundForward struct {
	Enabled  string      `json:"enabled"`
	Domain   string      `json:"domain"`
	Type     SelectedMap `json:"type"`
	Server   string      `json:"server"`
	Port     string      `json:"port"`
	VerifyCN string      `json:"verify"`
}

// CRUD operations

func (u *unbound) AddForward(ctx context.Context, forward *UnboundForward) (string, error) {
	return makeSetFunc(u, unboundForwardAddEndpoint, unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundForward{
			"dot": forward,
		},
	)
}

func (u *unbound) GetForward(ctx context.Context, id string) (*UnboundForward, error) {
	get, err := makeGetFunc(u.Client(), unboundForwardGetEndpoint,
		&struct {
			Dot UnboundForward `json:"dot"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Dot, nil
}

func (u *unbound) UpdateForward(ctx context.Context, id string, forward *UnboundForward) error {
	_, err := makeSetFunc(u, fmt.Sprintf("%s/%s", unboundForwardUpdateEndpoint, id),
		unboundReconfigureEndpoint)(ctx,
		map[string]*UnboundForward{
			"dot": forward,
		},
	)
	return err
}

func (u *unbound) DeleteForward(ctx context.Context, id string) error {
	return makeDeleteFunc(u, unboundForwardDeleteEndpoint, unboundReconfigureEndpoint)(ctx, id)
}
