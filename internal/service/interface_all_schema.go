package service

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/diagnostics"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InterfaceAllDataSourceModel struct {
	Interfaces types.List `tfsdk:"interfaces"`
}

func InterfaceAllDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "InterfacesAll can be used to get a list of all configurations of OPNsense interfaces. Allows for custom filtering.",

		Attributes: map[string]schema.Attribute{
			"interfaces": schema.ListNestedAttribute{
				MarkdownDescription: "A list of all interfaces present in OPNsense.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: InterfaceDataSourceSchema().Attributes,
				},
				Computed: true,
			},
		},
	}
}

func convertAllInterfaceConfigStructToSchema(d []diagnostics.Interface) (*InterfaceAllDataSourceModel, error) {
	var interfaces []InterfaceDataSourceModel
	for _, iface := range d {
		toSchema, err := convertInterfaceConfigStructToSchema(&iface)
		if err != nil {
			return nil, err
		}

		interfaces = append(interfaces, *toSchema)
	}

	// Create empty list first
	v, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: interfaceAttrTypes,
		},
		[]attr.Value{},
	)
	// Try to fill list
	if len(interfaces) > 0 {
		v, _ = types.ListValueFrom(
			context.Background(),
			types.ObjectType{
				AttrTypes: interfaceAttrTypes,
			},
			interfaces,
		)
	}

	model := &InterfaceAllDataSourceModel{
		Interfaces: v,
	}

	return model, nil
}
