package dnsmasq

import (
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var genericGroupResource = &generic.GenericResource[dnsmasq.Host, hostResourceModel]{
	Definition: hostDefinition,
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = genericGroupResource
var _ resource.ResourceWithConfigure = genericGroupResource

func NewHostResource() resource.Resource {
	return genericGroupResource
}
