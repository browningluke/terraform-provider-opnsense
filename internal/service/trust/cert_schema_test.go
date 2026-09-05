package trust

import (
	"encoding/base64"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/trust"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The OPNsense import/import_csr/sign_csr/reissue/manual actions parse the
// certificate, CSR and private key from crt_payload/csr_payload/prv_payload
// (raw PEM). The persisted crt/csr/prv fields are ignored on write. See
// opnsense/core Trust/Cert.xml and CertController::setBaseHook.
func TestConvertCertSchemaToStructSendsPayloadFields(t *testing.T) {
	crtPEM := "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"
	prvPEM := "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----\n"
	csrPEM := "-----BEGIN CERTIFICATE REQUEST-----\nMII...\n-----END CERTIFICATE REQUEST-----\n"

	cert, err := convertCertSchemaToStruct(&certResourceModel{
		Action:             types.StringValue("import"),
		CaRef:              types.StringValue(""),
		Crt:                types.StringValue(crtPEM),
		Csr:                types.StringValue(csrPEM),
		Prv:                types.StringValue(prvPEM),
		KeyType:            types.StringValue("2048"),
		Digest:             types.StringValue("sha256"),
		CertType:           types.StringValue("usr_cert"),
		PrivateKeyLocation: types.StringValue("firewall"),
		Country:            types.StringValue(""),
		CommonName:         types.StringValue("example.com"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cert.CrtPayload != crtPEM {
		t.Fatalf("expected CrtPayload to be raw PEM %q, got %q", crtPEM, cert.CrtPayload)
	}
	if cert.CsrPayload != csrPEM {
		t.Fatalf("expected CsrPayload to be raw PEM %q, got %q", csrPEM, cert.CsrPayload)
	}
	if cert.PrvPayload != prvPEM {
		t.Fatalf("expected PrvPayload to be raw PEM %q, got %q", prvPEM, cert.PrvPayload)
	}
}

// convertCertSchemaToStruct must not send crt_payload/csr_payload/prv_payload
// for actions that don't need them (e.g. "internal"), since the fields carry
// `omitempty` and an unexpectedly-present empty payload could confuse other
// action branches server-side.
func TestConvertCertSchemaToStructOmitsEmptyPayloads(t *testing.T) {
	cert, err := convertCertSchemaToStruct(&certResourceModel{
		Action:             types.StringValue("internal"),
		CaRef:              types.StringValue(""),
		Crt:                types.StringValue(""),
		Csr:                types.StringValue(""),
		Prv:                types.StringValue(""),
		KeyType:            types.StringValue("2048"),
		Digest:             types.StringValue("sha256"),
		CertType:           types.StringValue("usr_cert"),
		PrivateKeyLocation: types.StringValue("firewall"),
		Country:            types.StringValue(""),
		CommonName:         types.StringValue("example.com"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.CrtPayload != "" || cert.CsrPayload != "" || cert.PrvPayload != "" {
		t.Fatalf("expected empty payload fields, got crt_payload=%q csr_payload=%q prv_payload=%q",
			cert.CrtPayload, cert.CsrPayload, cert.PrvPayload)
	}
}

// convertCertStructToSchema must decode the base64 that OPNsense persists in
// crt/csr/prv into raw PEM, so that re-applying the same raw-PEM
// configuration produces no diff (state must match config after a round
// trip).
func TestConvertCertStructToSchemaDecodesBase64(t *testing.T) {
	crtPEM := "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"
	prvPEM := "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----\n"
	csrPEM := "-----BEGIN CERTIFICATE REQUEST-----\nMII...\n-----END CERTIFICATE REQUEST-----\n"

	model, err := convertCertStructToSchema(&trust.Cert{
		Crt: base64.StdEncoding.EncodeToString([]byte(crtPEM)),
		Csr: base64.StdEncoding.EncodeToString([]byte(csrPEM)),
		Prv: base64.StdEncoding.EncodeToString([]byte(prvPEM)),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if model.Crt.ValueString() != crtPEM {
		t.Fatalf("expected decoded crt %q, got %q", crtPEM, model.Crt.ValueString())
	}
	if model.Csr.ValueString() != csrPEM {
		t.Fatalf("expected decoded csr %q, got %q", csrPEM, model.Csr.ValueString())
	}
	if model.Prv.ValueString() != prvPEM {
		t.Fatalf("expected decoded prv %q, got %q", prvPEM, model.Prv.ValueString())
	}
}

func TestConvertCertStructToSchemaEmptyFieldsRoundtrip(t *testing.T) {
	model, err := convertCertStructToSchema(&trust.Cert{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if model.Crt.ValueString() != "" || model.Csr.ValueString() != "" || model.Prv.ValueString() != "" {
		t.Fatalf("expected empty crt/csr/prv, got crt=%q csr=%q prv=%q",
			model.Crt.ValueString(), model.Csr.ValueString(), model.Prv.ValueString())
	}
}

func TestConvertCertStructToSchemaInvalidBase64Errors(t *testing.T) {
	if _, err := convertCertStructToSchema(&trust.Cert{Crt: "not-valid-base64!!"}); err == nil {
		t.Fatal("expected error decoding invalid base64 crt, got nil")
	}
}
