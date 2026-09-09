package generic

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
)

func buildTerraformTypeName[OpnSenseType any, TerraformSchema any](providerTypeName string, definition Definition[OpnSenseType, TerraformSchema]) string {
	return fmt.Sprintf("%s_%s_%s", providerTypeName, definition.ServiceName, definition.TypeName)
}

func configureClient(providerData any) (opnsense.Client, error) {
	// Prevent panic if the provider has not been configured.
	if providerData == nil {
		return nil, nil
	}

	apiClient, ok := providerData.(*api.Client)
	if !ok {
		return nil, fmt.Errorf("expected *api.Client, got %T", providerData)
	}

	return opnsense.NewClient(apiClient), nil
}

func readResource[
	OpnSenseType any,
	TerraformSchema any,
](
	ctx context.Context,
	client opnsense.Client,
	definition Definition[OpnSenseType, TerraformSchema],
	data *TerraformSchema,
) (*TerraformSchema, error) {
	if definition.Read == nil {
		return nil, fmt.Errorf("Read operation is not defined for '%s_%s'", definition.ServiceName, definition.TypeName)
	}

	resourceStruct, err := definition.Read(ctx, client, definition.GetID(data).ValueString())
	if err != nil {
		return nil, err
	}

	model, err := definition.ToModel(resourceStruct)
	if err != nil {
		return nil, err
	}

	definition.SetID(model, definition.GetID(data))

	if definition.ReadHook != nil {
		if err := definition.ReadHook(ctx, client, model); err != nil {
			return nil, err
		}
	}

	return model, nil
}
