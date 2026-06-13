package trust

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/trust"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type certResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	RefId              types.String `tfsdk:"ref_id"`
	Description        types.String `tfsdk:"description"`
	CaRef              types.String `tfsdk:"caref"`
	Crt                types.String `tfsdk:"crt"`
	Csr                types.String `tfsdk:"csr"`
	Prv                types.String `tfsdk:"prv"`
	Action             types.String `tfsdk:"action"`
	KeyType            types.String `tfsdk:"key_type"`
	Digest             types.String `tfsdk:"digest"`
	CertType           types.String `tfsdk:"cert_type"`
	Lifetime           types.String `tfsdk:"lifetime"`
	PrivateKeyLocation types.String `tfsdk:"private_key_location"`
	Country            types.String `tfsdk:"country"`
	State              types.String `tfsdk:"state"`
	City               types.String `tfsdk:"city"`
	Organization       types.String `tfsdk:"organization"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
	Email              types.String `tfsdk:"email"`
	CommonName         types.String `tfsdk:"common_name"`
	OcspUri            types.String `tfsdk:"ocsp_uri"`
	AltnamesDns        types.String `tfsdk:"altnames_dns"`
	AltnamesIp         types.String `tfsdk:"altnames_ip"`
	AltnamesUri        types.String `tfsdk:"altnames_uri"`
	AltnamesEmail      types.String `tfsdk:"altnames_email"`
	Rfc3280Purpose     types.String `tfsdk:"rfc3280_purpose"`
	InUse              types.String `tfsdk:"in_use"`
	IsUser             types.String `tfsdk:"is_user"`
	CrtPayload         types.String `tfsdk:"crt_payload"`
	CsrPayload         types.String `tfsdk:"csr_payload"`
	PrvPayload         types.String `tfsdk:"prv_payload"`
	Name               types.String `tfsdk:"name"`
	ValidFrom          types.String `tfsdk:"valid_from"`
	ValidTo            types.String `tfsdk:"valid_to"`
}

func certResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages an end-entity certificate in the OPNsense Trust store.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ref_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Short hex reference ID used by other OPNsense subsystems to reference this certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this certificate.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"caref": schema.StringAttribute{
				MarkdownDescription: "Reference ID (`ref_id`) of the signing CA. Required for `internal` and `sign_csr` actions.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"crt": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded PEM certificate body. Required when `action` is `import`. Computed when `action` is `internal`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"csr": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded PEM Certificate Signing Request. Required when `action` is `sign_csr` or `import_csr`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"prv": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded PEM private key. Required when `action` is `import`. Computed when `action` is `internal`.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Certificate action: `internal` (generate signed by CA), `external` (CSR only), `import` (import existing), `sign_csr`, `import_csr`, `reissue`, `manual`. Defaults to `internal`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("internal"),
			},
			"key_type": schema.StringAttribute{
				MarkdownDescription: "Key type and size: `512`, `1024`, `2048` (default), `3072`, `4096`, `8192`, `prime256v1`, `secp384r1`, `secp521r1`. Defaults to `2048`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("2048"),
			},
			"digest": schema.StringAttribute{
				MarkdownDescription: "Digest algorithm: `sha224`, `sha256` (default), `sha384`, `sha512`. Defaults to `sha256`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("sha256"),
			},
			"cert_type": schema.StringAttribute{
				MarkdownDescription: "Certificate type: `server_cert`, `usr_cert` (default), `combined_server_client`. Defaults to `usr_cert`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("usr_cert"),
			},
			"lifetime": schema.StringAttribute{
				MarkdownDescription: "Certificate validity period in days. Defaults to `397` (maximum accepted by modern browsers for server certificates).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("397"),
			},
			"private_key_location": schema.StringAttribute{
				MarkdownDescription: "Where to store the private key: `firewall` (default, store on the firewall) or `external` (CSR only, no key stored). Defaults to `firewall`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("firewall"),
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "ISO 3166-1 alpha-2 country code (e.g. `US`, `CA`, `DE`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"state": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"city": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"organization": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"organizational_unit": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"common_name": schema.StringAttribute{
				MarkdownDescription: "Common name (CN) for the certificate.",
				Required:            true,
			},
			"ocsp_uri": schema.StringAttribute{
				MarkdownDescription: "OCSP responder URI.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"altnames_dns": schema.StringAttribute{
				MarkdownDescription: "DNS Subject Alternative Names, newline or CRLF separated (e.g. `www.example.com\\napi.example.com`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"altnames_ip": schema.StringAttribute{
				MarkdownDescription: "IP Subject Alternative Names, newline separated.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"altnames_uri": schema.StringAttribute{
				MarkdownDescription: "URI Subject Alternative Names, newline separated.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"altnames_email": schema.StringAttribute{
				MarkdownDescription: "Email Subject Alternative Names, newline separated.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"rfc3280_purpose": schema.StringAttribute{
				MarkdownDescription: "Extended Key Usage OID string (e.g. `id-kp-serverAuth`, `id-kp-clientAuth`). Computed by OPNsense based on `cert_type` when not explicitly set.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"in_use": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this certificate is currently referenced by another OPNsense subsystem (`0`/`1`, read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_user": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this certificate is linked to a user account (`0`/`1`, read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"crt_payload": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Decoded PEM certificate body (read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"csr_payload": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Decoded PEM CSR (read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prv_payload": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Decoded PEM private key (read-only, sensitive).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Distinguished name string of the certificate (read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"valid_from": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp of certificate start date (read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"valid_to": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unix timestamp of certificate expiry date (read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func certDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves an end-entity certificate from the OPNsense Trust store by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the certificate.",
			},
			"ref_id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Short hex reference ID used by other OPNsense subsystems.",
			},
			"description": dschema.StringAttribute{
				Computed: true,
			},
			"caref": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reference ID of the signing CA.",
			},
			"crt": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Base64-encoded PEM certificate body.",
			},
			"csr": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Base64-encoded PEM CSR.",
			},
			"action": dschema.StringAttribute{
				Computed: true,
			},
			"key_type": dschema.StringAttribute{
				Computed: true,
			},
			"digest": dschema.StringAttribute{
				Computed: true,
			},
			"cert_type": dschema.StringAttribute{
				Computed: true,
			},
			"lifetime": dschema.StringAttribute{
				Computed: true,
			},
			"private_key_location": dschema.StringAttribute{
				Computed: true,
			},
			"country": dschema.StringAttribute{
				Computed: true,
			},
			"state": dschema.StringAttribute{
				Computed: true,
			},
			"city": dschema.StringAttribute{
				Computed: true,
			},
			"organization": dschema.StringAttribute{
				Computed: true,
			},
			"organizational_unit": dschema.StringAttribute{
				Computed: true,
			},
			"email": dschema.StringAttribute{
				Computed: true,
			},
			"common_name": dschema.StringAttribute{
				Computed: true,
			},
			"ocsp_uri": dschema.StringAttribute{
				Computed: true,
			},
			"altnames_dns": dschema.StringAttribute{
				Computed: true,
			},
			"altnames_ip": dschema.StringAttribute{
				Computed: true,
			},
			"altnames_uri": dschema.StringAttribute{
				Computed: true,
			},
			"altnames_email": dschema.StringAttribute{
				Computed: true,
			},
			"rfc3280_purpose": dschema.StringAttribute{
				Computed: true,
			},
			"in_use": dschema.StringAttribute{
				Computed: true,
			},
			"is_user": dschema.StringAttribute{
				Computed: true,
			},
			"crt_payload": dschema.StringAttribute{
				Computed: true,
			},
			"csr_payload": dschema.StringAttribute{
				Computed: true,
			},
			"prv_payload": dschema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"name": dschema.StringAttribute{
				Computed: true,
			},
			"valid_from": dschema.StringAttribute{
				Computed: true,
			},
			"valid_to": dschema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func convertCertSchemaToStruct(d *certResourceModel) (*trust.Cert, error) {
	return &trust.Cert{
		Description:        d.Description.ValueString(),
		CaRef:              api.SelectedMap(d.CaRef.ValueString()),
		Crt:                d.Crt.ValueString(),
		Csr:                d.Csr.ValueString(),
		Prv:                d.Prv.ValueString(),
		Action:             api.SelectedMap(d.Action.ValueString()),
		KeyType:            api.SelectedMap(d.KeyType.ValueString()),
		Digest:             api.SelectedMap(d.Digest.ValueString()),
		CertType:           api.SelectedMap(d.CertType.ValueString()),
		Lifetime:           d.Lifetime.ValueString(),
		PrivateKeyLocation: api.SelectedMap(d.PrivateKeyLocation.ValueString()),
		Country:            api.SelectedMap(d.Country.ValueString()),
		State:              d.State.ValueString(),
		City:               d.City.ValueString(),
		Organization:       d.Organization.ValueString(),
		OrganizationalUnit: d.OrganizationalUnit.ValueString(),
		Email:              d.Email.ValueString(),
		CommonName:         d.CommonName.ValueString(),
		OcspUri:            d.OcspUri.ValueString(),
		AltnamesDns:        d.AltnamesDns.ValueString(),
		AltnamesIp:         d.AltnamesIp.ValueString(),
		AltnamesUri:        d.AltnamesUri.ValueString(),
		AltnamesEmail:      d.AltnamesEmail.ValueString(),
		Rfc3280Purpose:     d.Rfc3280Purpose.ValueString(),
	}, nil
}

func convertCertStructToSchema(d *trust.Cert) (*certResourceModel, error) {
	return &certResourceModel{
		RefId:              tools.StringOrNull(d.RefId),
		Description:        types.StringValue(d.Description),
		CaRef:              types.StringValue(d.CaRef.String()),
		Crt:                types.StringValue(d.Crt),
		Csr:                types.StringValue(d.Csr),
		Prv:                types.StringValue(d.Prv),
		Action:             types.StringValue(d.Action.String()),
		KeyType:            types.StringValue(d.KeyType.String()),
		Digest:             types.StringValue(d.Digest.String()),
		CertType:           types.StringValue(d.CertType.String()),
		Lifetime:           types.StringValue(d.Lifetime),
		PrivateKeyLocation: types.StringValue(d.PrivateKeyLocation.String()),
		Country:            types.StringValue(d.Country.String()),
		State:              types.StringValue(d.State),
		City:               types.StringValue(d.City),
		Organization:       types.StringValue(d.Organization),
		OrganizationalUnit: types.StringValue(d.OrganizationalUnit),
		Email:              types.StringValue(d.Email),
		CommonName:         types.StringValue(d.CommonName),
		OcspUri:            types.StringValue(d.OcspUri),
		AltnamesDns:        types.StringValue(d.AltnamesDns),
		AltnamesIp:         types.StringValue(d.AltnamesIp),
		AltnamesUri:        types.StringValue(d.AltnamesUri),
		AltnamesEmail:      types.StringValue(d.AltnamesEmail),
		Rfc3280Purpose:     types.StringValue(d.Rfc3280Purpose),
		InUse:              tools.StringOrNull(d.InUse),
		IsUser:             tools.StringOrNull(d.IsUser),
		CrtPayload:         tools.StringOrNull(d.CrtPayload),
		CsrPayload:         tools.StringOrNull(d.CsrPayload),
		PrvPayload:         tools.StringOrNull(d.PrvPayload),
		Name:               tools.StringOrNull(d.Name),
		ValidFrom:          tools.StringOrNull(d.ValidFrom),
		ValidTo:            tools.StringOrNull(d.ValidTo),
	}, nil
}
