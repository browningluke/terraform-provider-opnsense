package validators

import "context"

// descriptionValidator holds a static description string and satisfies the
// Description / MarkdownDescription pair required by all Terraform validator
// interfaces.  Embed it in any validator whose description is a fixed string.
type descriptionValidator struct {
	desc string
}

func (d descriptionValidator) Description(_ context.Context) string {
	return d.desc
}

func (d descriptionValidator) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
