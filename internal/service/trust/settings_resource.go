package trust

import (
	"github.com/browningluke/opnsense-go/pkg/trust"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var genericSettingsResource = &generic.GenericSingletonResource[trust.TrustSettings, settingsResourceModel]{
	Definition: settingsDefinition,
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = genericSettingsResource
var _ resource.ResourceWithConfigure = genericSettingsResource
var _ resource.ResourceWithImportState = genericSettingsResource

func newSettingsResource() resource.Resource {
	return genericSettingsResource
}
