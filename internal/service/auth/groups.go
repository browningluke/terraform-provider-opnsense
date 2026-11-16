package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/browningluke/opnsense-go/pkg/api"
)

type groupSearchResult struct {
	Rows []struct {
		UUID string `json:"uuid"`
		GID  string `json:"gid"`
		Name string `json:"name"`
	} `json:"rows"`
}

type groupIndex struct {
	byName map[string]string
	byGID  map[string]string
}

func loadGroupIndex(ctx context.Context, client *api.Client) (*groupIndex, error) {
	result := &groupSearchResult{}
	_, err := api.Call(client, ctx, api.RPCOpts{
		BaseEndpoint: "/auth/group/search",
		Method:       "GET",
	}, result)
	if err != nil {
		return nil, fmt.Errorf("search auth groups: %w", err)
	}

	idx := &groupIndex{
		byName: make(map[string]string, len(result.Rows)),
		byGID:  make(map[string]string, len(result.Rows)),
	}
	for _, row := range result.Rows {
		idx.byName[row.Name] = row.GID
		idx.byGID[row.GID] = row.Name
	}
	return idx, nil
}

// resolveGroupMembershipsToGIDs accepts group names only. Read-back always
// stores names, so numeric GIDs in config would drift after apply.
func resolveGroupMembershipsToGIDs(ctx context.Context, client *api.Client, values []string) ([]string, error) {
	if len(values) == 0 {
		return []string{}, nil
	}

	idx, err := loadGroupIndex(ctx, client)
	if err != nil {
		return nil, err
	}

	gids := make([]string, 0, len(values))
	for _, v := range values {
		if _, err := strconv.Atoi(v); err == nil {
			return nil, fmt.Errorf("auth group %q looks numeric; use the group name (e.g. admins), not the GID", v)
		}
		gid, ok := idx.byName[v]
		if !ok {
			return nil, fmt.Errorf("unknown auth group name %q", v)
		}
		gids = append(gids, gid)
	}
	return gids, nil
}

func groupGIDsToNames(ctx context.Context, client *api.Client, gids []string) ([]string, error) {
	if len(gids) == 0 {
		return []string{}, nil
	}

	idx, err := loadGroupIndex(ctx, client)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(gids))
	for _, gid := range gids {
		if name, ok := idx.byGID[gid]; ok {
			names = append(names, name)
			continue
		}
		names = append(names, gid)
	}
	return names, nil
}
