package mdnsrepeater

import (
	"encoding/json"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/mdnsrepeater"
)

// An empty blocklist comes back as a single selected "" option, which must
// not surface as an element, or every plan would show drift against [].
func TestConvertSettingsStructToSchemaEmptyBlocklist(t *testing.T) {
	body := `{"mdnsrepeater":{"enabled":"0","enablecarp":"0",` +
		`"interfaces":{"wan":{"value":"WAN","selected":1},"lan":{"value":"LAN","selected":0}},` +
		`"blocklist":{"":{"value":"","selected":1}}}}`

	var resp mdnsrepeater.SettingsMonad
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	model, err := convertSettingsStructToSchema(&resp.MDNSRepeater)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n := len(model.Blocklist.Elements()); n != 0 {
		t.Fatalf("expected empty blocklist, got %s", model.Blocklist)
	}
	if n := len(model.Interfaces.Elements()); n != 1 {
		t.Fatalf("expected 1 interface, got %s", model.Interfaces)
	}
}
