package ipsec

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &authLocalResource{} },
		func() resource.Resource { return &authRemoteResource{} },
		func() resource.Resource { return &childResource{} },
		func() resource.Resource { return &connectionResource{} },
		func() resource.Resource { return &pskResource{} },
		func() resource.Resource { return &vtiResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
