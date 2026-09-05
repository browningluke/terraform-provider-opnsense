package trust

import (
	"encoding/base64"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/trust"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The OPNsense "existing" action parses the certificate and private key
// from crt_payload/prv_payload (raw PEM). The persisted crt/prv fields are
// ignored on write. See opnsense/core Trust/Ca.xml and
// CaController::setBaseHook.
func TestConvertCaSchemaToStructSendsPayloadFields(t *testing.T) {
	crtPEM := "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"
	prvPEM := "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----\n"

	ca, err := convertCaSchemaToStruct(&caResourceModel{
		Action:     types.StringValue("existing"),
		Crt:        types.StringValue(crtPEM),
		Prv:        types.StringValue(prvPEM),
		CaRef:      types.StringValue(""),
		KeyType:    types.StringValue("2048"),
		Digest:     types.StringValue("sha256"),
		Country:    types.StringValue(""),
		CommonName: types.StringValue("example-ca.internal"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ca.CrtPayload != crtPEM {
		t.Fatalf("expected CrtPayload to be raw PEM %q, got %q", crtPEM, ca.CrtPayload)
	}
	if ca.PrvPayload != prvPEM {
		t.Fatalf("expected PrvPayload to be raw PEM %q, got %q", prvPEM, ca.PrvPayload)
	}
}

func TestConvertCaSchemaToStructOmitsEmptyPayloads(t *testing.T) {
	ca, err := convertCaSchemaToStruct(&caResourceModel{
		Action:     types.StringValue("internal"),
		Crt:        types.StringValue(""),
		Prv:        types.StringValue(""),
		CaRef:      types.StringValue(""),
		KeyType:    types.StringValue("2048"),
		Digest:     types.StringValue("sha256"),
		Country:    types.StringValue(""),
		CommonName: types.StringValue("example-ca.internal"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ca.CrtPayload != "" || ca.PrvPayload != "" {
		t.Fatalf("expected empty payload fields, got crt_payload=%q prv_payload=%q", ca.CrtPayload, ca.PrvPayload)
	}
}

// convertCaStructToSchema must decode the base64 that OPNsense persists in
// crt/prv into raw PEM, so that re-applying the same raw-PEM configuration
// produces no diff (state must match config after a round trip).
func TestConvertCaStructToSchemaDecodesBase64(t *testing.T) {
	crtPEM := "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"
	prvPEM := "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----\n"

	model, err := convertCaStructToSchema(&trust.Ca{
		Crt: base64.StdEncoding.EncodeToString([]byte(crtPEM)),
		Prv: base64.StdEncoding.EncodeToString([]byte(prvPEM)),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if model.Crt.ValueString() != crtPEM {
		t.Fatalf("expected decoded crt %q, got %q", crtPEM, model.Crt.ValueString())
	}
	if model.Prv.ValueString() != prvPEM {
		t.Fatalf("expected decoded prv %q, got %q", prvPEM, model.Prv.ValueString())
	}
}

func TestConvertCaStructToSchemaEmptyFieldsRoundtrip(t *testing.T) {
	model, err := convertCaStructToSchema(&trust.Ca{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if model.Crt.ValueString() != "" || model.Prv.ValueString() != "" {
		t.Fatalf("expected empty crt/prv, got crt=%q prv=%q", model.Crt.ValueString(), model.Prv.ValueString())
	}
}

func TestConvertCaStructToSchemaInvalidBase64Errors(t *testing.T) {
	if _, err := convertCaStructToSchema(&trust.Ca{Prv: "not-valid-base64!!"}); err == nil {
		t.Fatal("expected error decoding invalid base64 prv, got nil")
	}
}
