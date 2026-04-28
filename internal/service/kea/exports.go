package kea

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Deprecated: use dhcpv4 variants
		newPeerResource,
		newReservationResource,
		newSubnetResource,
		// DHCPv4
		newDhcpv4PeerResource,
		newDhcpv4ReservationResource,
		newDhcpv4SubnetResource,
		// DHCPv6
		newDhcpv6PeerResource,
		newDhcpv6ReservationResource,
		newDhcpv6SubnetResource,
		newDhcpv6PdPoolResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Deprecated: use dhcpv4 variants
		newPeerDataSource,
		newReservationDataSource,
		newSubnetDataSource,
		// DHCPv4
		newDhcpv4PeerDataSource,
		newDhcpv4ReservationDataSource,
		newDhcpv4SubnetDataSource,
		// DHCPv6
		newDhcpv6PeerDataSource,
		newDhcpv6ReservationDataSource,
		newDhcpv6SubnetDataSource,
		newDhcpv6PdPoolDataSource,
	}
}
