package quagga

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Singletons
		newGeneralResource,
		newBFDResource,
		newBGPResource,
		newOSPFResource,
		newOSPF6Resource,
		newRIPResource,
		newStaticResource,
		// BGP child resources
		newBGPASPathResource,
		newBGPCommunityListResource,
		newBGPNeighborResource,
		newBGPPeerGroupResource,
		newBGPPrefixListResource,
		newBGPRedistributionResource,
		newBGPRouteMapResource,
		// BFD child resources
		newBFDNeighborResource,
		// OSPF child resources
		newOSPFAreaResource,
		newOSPFInterfaceResource,
		newOSPFNeighborResource,
		newOSPFNetworkResource,
		newOSPFPrefixListResource,
		newOSPFRedistributionResource,
		newOSPFRouteMapResource,
		// OSPFv3 child resources
		newOSPF6InterfaceResource,
		newOSPF6NetworkResource,
		newOSPF6PrefixListResource,
		newOSPF6RedistributionResource,
		newOSPF6RouteMapResource,
		// Static child resources
		newStaticRouteResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// BGP child data sources
		newBGPASPathDataSource,
		newBGPCommunityListDataSource,
		newBGPNeighborDataSource,
		newBGPPeerGroupDataSource,
		newBGPPrefixListDataSource,
		newBGPRedistributionDataSource,
		newBGPRouteMapDataSource,
		// BFD child data sources
		newBFDNeighborDataSource,
		// OSPF child data sources
		newOSPFAreaDataSource,
		newOSPFInterfaceDataSource,
		newOSPFNeighborDataSource,
		newOSPFNetworkDataSource,
		newOSPFPrefixListDataSource,
		newOSPFRedistributionDataSource,
		newOSPFRouteMapDataSource,
		// OSPFv3 child data sources
		newOSPF6InterfaceDataSource,
		newOSPF6NetworkDataSource,
		newOSPF6PrefixListDataSource,
		newOSPF6RedistributionDataSource,
		newOSPF6RouteMapDataSource,
		// Static child data sources
		newStaticRouteDataSource,
	}
}
