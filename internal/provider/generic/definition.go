package generic

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/opnsense"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Definition[OpnSenseType any, TerraformSchema any] struct {
	ServiceName string
	TypeName    string

	ResourceSchema   schema.Schema
	DataSourceSchema dschema.Schema

	GetID func(*TerraformSchema) types.String
	SetID func(*TerraformSchema, types.String)

	FromModel func(*TerraformSchema) (*OpnSenseType, error)
	ToModel   func(*OpnSenseType) (*TerraformSchema, error)

	Create           func(context.Context, opnsense.Client, *OpnSenseType) (string, error)
	Read             func(context.Context, opnsense.Client, string) (*OpnSenseType, error)
	Update           func(context.Context, opnsense.Client, string, *OpnSenseType) error
	Delete           func(context.Context, opnsense.Client, string) error
	ConfigValidators func(context.Context) []resource.ConfigValidator
	UpgradeState     func(context.Context) map[int64]resource.StateUpgrader
	ImportState      func(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse)

	CreateHook func(context.Context, opnsense.Client, *TerraformSchema, string) error
	ReadHook   func(context.Context, opnsense.Client, *TerraformSchema) error
	UpdateHook func(context.Context, opnsense.Client, *TerraformSchema) error
	DeleteHook func(context.Context, opnsense.Client, *TerraformSchema) error
}
