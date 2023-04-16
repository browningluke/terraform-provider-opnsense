package opnsense

import (
	"context"
	"fmt"
	"sync"
)

const (
	routesRouteReconfigureEndpoint = "/routes/routes/reconfigure"
	routesRouteAddEndpoint         = "/routes/routes/addroute"
	routesRouteGetEndpoint         = "/routes/routes/getroute"
	routesRouteUpdateEndpoint      = "/routes/routes/setroute"
	routesRouteDeleteEndpoint      = "/routes/routes/delroute"
)

// Routes controller
type routes struct {
	client *Client
	mu     *sync.Mutex
}

func newRoutes(c *Client) *routes {
	return &routes{
		client: c,
		mu:     &sync.Mutex{},
	}
}

func (r *routes) Client() *Client {
	return r.client
}

func (r *routes) Mutex() *sync.Mutex {
	return r.mu
}

// Data structs

type Route struct {
	Disabled    string      `json:"disabled"`
	Description string      `json:"descr"`
	Gateway     SelectedMap `json:"gateway"`
	Network     string      `json:"network"`
}

// CRUD operations

func (r *routes) AddRoute(ctx context.Context, route *Route) (string, error) {
	return makeSetFunc(r, routesRouteAddEndpoint, routesRouteReconfigureEndpoint)(ctx,
		map[string]*Route{
			"route": route,
		},
	)
}

func (r *routes) GetRoute(ctx context.Context, id string) (*Route, error) {
	get, err := makeGetFunc(r.Client(), routesRouteGetEndpoint,
		&struct {
			Route Route `json:"route"`
		}{},
	)(ctx, id)
	if err != nil {
		return nil, err
	}
	return &get.Route, nil
}

func (r *routes) UpdateRoute(ctx context.Context, id string, route *Route) error {
	_, err := makeSetFunc(r, fmt.Sprintf("%s/%s", routesRouteUpdateEndpoint, id),
		routesRouteReconfigureEndpoint)(ctx,
		map[string]*Route{
			"route": route,
		},
	)
	return err
}

func (r *routes) DeleteRoute(ctx context.Context, id string) error {
	return makeDeleteFunc(r, routesRouteDeleteEndpoint, routesRouteReconfigureEndpoint)(ctx, id)
}
