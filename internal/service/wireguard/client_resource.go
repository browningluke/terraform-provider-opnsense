package wireguard

import (
	"github.com/browningluke/opnsense-go/pkg/wireguard"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var genericGroupResource = &generic.GenericResource[wireguard.Client, clientResourceModel]{
	Definition: clientDefinition,
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = genericGroupResource
var _ resource.ResourceWithConfigure = genericGroupResource
var _ resource.ResourceWithImportState = genericGroupResource

func newClientResource() resource.Resource {
	return genericGroupResource
}
