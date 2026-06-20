package kea

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Deprecated: use dhcpv4 variants
		func() resource.Resource { return &peerResource{} },
		func() resource.Resource { return &reservationResource{} },
		func() resource.Resource { return &subnetResource{} },
		// DHCPv4
		func() resource.Resource { return &dhcpv4PeerResource{} },
		func() resource.Resource { return &dhcpv4ReservationResource{} },
		func() resource.Resource { return &dhcpv4SubnetResource{} },
		// DHCPv6
		func() resource.Resource { return &dhcpv6PeerResource{} },
		func() resource.Resource { return &dhcpv6ReservationResource{} },
		func() resource.Resource { return &dhcpv6SubnetResource{} },
		func() resource.Resource { return &dhcpv6PdPoolResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Deprecated: use dhcpv4 variants
		func() datasource.DataSource { return &peerDataSource{} },
		func() datasource.DataSource { return &reservationDataSource{} },
		func() datasource.DataSource { return &subnetDataSource{} },
		// DHCPv4
		func() datasource.DataSource { return &dhcpv4PeerDataSource{} },
		func() datasource.DataSource { return &dhcpv4ReservationDataSource{} },
		func() datasource.DataSource { return &dhcpv4SubnetDataSource{} },
		// DHCPv6
		func() datasource.DataSource { return &dhcpv6PeerDataSource{} },
		func() datasource.DataSource { return &dhcpv6ReservationDataSource{} },
		func() datasource.DataSource { return &dhcpv6SubnetDataSource{} },
		func() datasource.DataSource { return &dhcpv6PdPoolDataSource{} },
	}
}
