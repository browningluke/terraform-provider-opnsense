package haproxy_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// haproxyPreCheck skips unless OPNSENSE_TEST_HAPROXY=1: the os-haproxy plugin is
// not part of the CI image.
func haproxyPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	acctest.AccPreCheck(t)
	if os.Getenv("OPNSENSE_TEST_HAPROXY") != "1" {
		t.Skip("OPNSENSE_TEST_HAPROXY=1 required: needs the os-haproxy plugin on the test instance")
	}
}

// testFrontend creates a frontend through the API, standing in for one that is
// managed outside Terraform.
func testFrontend(t *testing.T) string {
	t.Helper()
	c := &haproxy.Controller{Api: api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: os.Getenv("OPNSENSE_ALLOW_INSECURE") == "true",
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})}
	ctx := context.Background()
	id, err := c.AddFrontend(ctx, &haproxy.Frontend{Enabled: "1", Name: "tf-acc-test-frontend", Bind: api.SelectedMapList{"127.0.0.1:18080"}, Mode: "http"})
	if err != nil {
		t.Fatalf("creating test frontend: %s", err)
	}
	t.Cleanup(func() {
		if err := c.DeleteFrontend(ctx, id); err != nil {
			t.Logf("deleting test frontend %s: %s", id, err)
		}
	})
	return id
}

func TestAccHAProxyRuleChain(t *testing.T) {
	// Runs before resource.Test because the configs below need the frontend UUID.
	haproxyPreCheck(t)
	frontendID := testFrontend(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHAProxyConfig(frontendID, 80, "web", "other"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_haproxy_server.test", "address", "192.0.2.20"),
					resource.TestCheckResourceAttr("opnsense_haproxy_server.test", "port", "80"),
					resource.TestCheckResourceAttr("opnsense_haproxy_server.test", "mode", "active"),
					resource.TestCheckResourceAttr("opnsense_haproxy_backend.test", "linked_servers.#", "1"),
					resource.TestCheckResourceAttr("opnsense_haproxy_backend.test", "algorithm", "source"),
					resource.TestCheckResourceAttr("opnsense_haproxy_acl.host", "expression", "hdr_end"),
					resource.TestCheckResourceAttr("opnsense_haproxy_acl.host", "value", ".web.example.com"),
					resource.TestCheckResourceAttr("opnsense_haproxy_acl.path", "value", "/.well-known/acme-challenge/"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.web", "linked_acls.#", "2"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.web", "operator", "and"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_action.web", "use_backend", "opnsense_haproxy_backend.test", "id"),
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend_actions.test", "action_ids.#", "2"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.test", "action_ids.0", "opnsense_haproxy_action.web", "id"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.test", "action_ids.1", "opnsense_haproxy_action.other", "id"),
				),
			},
			{ResourceName: "opnsense_haproxy_server.test", ImportState: true, ImportStateVerify: true},
			{ResourceName: "opnsense_haproxy_backend.test", ImportState: true, ImportStateVerify: true},
			{ResourceName: "opnsense_haproxy_acl.host", ImportState: true, ImportStateVerify: true},
			{ResourceName: "opnsense_haproxy_action.web", ImportState: true, ImportStateVerify: true},
			{ResourceName: "opnsense_haproxy_frontend_actions.test", ImportState: true, ImportStateVerify: true},
			{
				// Reorder the rules and change the server.
				Config: testAccHAProxyConfig(frontendID, 8080, "other", "web"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_haproxy_server.test", "port", "8080"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.test", "action_ids.0", "opnsense_haproxy_action.other", "id"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.test", "action_ids.1", "opnsense_haproxy_action.web", "id"),
				),
			},
		},
	})
}

func testAccHAProxyConfig(frontendID string, port int, first, second string) string {
	return fmt.Sprintf(`
resource "opnsense_haproxy_server" "test" {
  name    = "tf-acc-test-server"
  address = "192.0.2.20"
  port    = %[2]d
}

resource "opnsense_haproxy_backend" "test" {
  name                 = "tf-acc-test-backend"
  linked_servers       = [opnsense_haproxy_server.test.id]
  health_check_enabled = false
}

resource "opnsense_haproxy_acl" "host" {
  name       = "tf-acc-test-host"
  expression = "hdr_end"
  value      = ".web.example.com"
}

resource "opnsense_haproxy_acl" "path" {
  name       = "tf-acc-test-path"
  expression = "path_beg"
  value      = "/.well-known/acme-challenge/"
}

resource "opnsense_haproxy_action" "web" {
  name        = "tf-acc-test-web"
  linked_acls = [opnsense_haproxy_acl.host.id, opnsense_haproxy_acl.path.id]
  type        = "use_backend"
  use_backend = opnsense_haproxy_backend.test.id
}

resource "opnsense_haproxy_action" "other" {
  name        = "tf-acc-test-other"
  linked_acls = [opnsense_haproxy_acl.path.id]
  type        = "use_backend"
  use_backend = opnsense_haproxy_backend.test.id
}

resource "opnsense_haproxy_frontend_actions" "test" {
  frontend_id = %[1]q
  action_ids  = [opnsense_haproxy_action.%[3]s.id, opnsense_haproxy_action.%[4]s.id]
}
`, frontendID, port, first, second)
}

// TestAccHAProxyFrontend points opnsense_haproxy_frontend and
// opnsense_haproxy_frontend_actions at the same frontend. Every step ends with
// a refresh and an empty-plan check, so a resource that resets the other's
// field fails here.
func TestAccHAProxyFrontend(t *testing.T) {
	haproxyPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHAProxyFrontendConfig("first", `["127.0.0.1:18081"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend.test", "name", "tf-acc-test-fe"),
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend.test", "mode", "http"),
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend.test", "enabled", "true"),
					resource.TestCheckTypeSetElemAttr("opnsense_haproxy_frontend.test", "bind.*", "127.0.0.1:18081"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend.test", "default_backend", "opnsense_haproxy_backend.fe", "id"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.fe", "frontend_id", "opnsense_haproxy_frontend.test", "id"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.fe", "action_ids.0", "opnsense_haproxy_action.fe", "id"),
					resource.TestCheckResourceAttrPair("opnsense_haproxy_frontend_actions.fe", "action_ids.1", "opnsense_haproxy_action.redirect", "id"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.redirect", "type", "http-request"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.redirect", "http_request_action", "redirect"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.redirect", "http_request_option", "scheme https code 301"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.redirect", "linked_acls.#", "0"),
					resource.TestCheckResourceAttr("opnsense_haproxy_action.redirect", "use_backend", ""),
				),
			},
			{ResourceName: "opnsense_haproxy_frontend.test", ImportState: true, ImportStateVerify: true},
			{ResourceName: "opnsense_haproxy_action.redirect", ImportState: true, ImportStateVerify: true},
			{
				// Changing the frontend must keep the rule list.
				Config: testAccHAProxyFrontendConfig("second", `["127.0.0.1:18081", "127.0.0.1:18082"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend.test", "description", "second"),
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend.test", "bind.#", "2"),
					resource.TestCheckResourceAttr("opnsense_haproxy_frontend_actions.fe", "action_ids.#", "2"),
				),
			},
		},
	})
}

func testAccHAProxyFrontendConfig(description, bind string) string {
	return fmt.Sprintf(`
resource "opnsense_haproxy_server" "fe" {
  name    = "tf-acc-test-fe-server"
  address = "192.0.2.21"
  port    = 80
}

resource "opnsense_haproxy_backend" "fe" {
  name                 = "tf-acc-test-fe-backend"
  linked_servers       = [opnsense_haproxy_server.fe.id]
  health_check_enabled = false
}

resource "opnsense_haproxy_acl" "fe" {
  name       = "tf-acc-test-fe-path"
  expression = "path_beg"
  value      = "/.well-known/acme-challenge/"
}

resource "opnsense_haproxy_action" "fe" {
  name        = "tf-acc-test-fe-rule"
  linked_acls = [opnsense_haproxy_acl.fe.id]
  type        = "use_backend"
  use_backend = opnsense_haproxy_backend.fe.id
}

# No conditions: applies to every request.
resource "opnsense_haproxy_action" "redirect" {
  name                = "tf-acc-test-fe-https"
  type                = "http-request"
  http_request_action = "redirect"
  http_request_option = "scheme https code 301"
}

resource "opnsense_haproxy_frontend" "test" {
  name            = "tf-acc-test-fe"
  description     = %[1]q
  bind            = %[2]s
  default_backend = opnsense_haproxy_backend.fe.id
}

resource "opnsense_haproxy_frontend_actions" "fe" {
  frontend_id = opnsense_haproxy_frontend.test.id
  action_ids  = [opnsense_haproxy_action.fe.id, opnsense_haproxy_action.redirect.id]
}
`, description, bind)
}

func TestAccHAProxyActionValidation(t *testing.T) {
	haproxyPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// use_backend rule without a backend.
				Config: `
resource "opnsense_haproxy_action" "bad" {
  name = "tf-acc-test-bad"
  type = "use_backend"
}`,
				ExpectError: regexp.MustCompile(`use_backend`),
			},
			{
				// redirect arguments on a use_backend rule.
				Config: `
resource "opnsense_haproxy_action" "bad" {
  name                = "tf-acc-test-bad"
  type                = "use_backend"
  use_backend         = "6b1a0c6e-0000-4000-8000-00000000000a"
  http_request_option = "scheme https code 301"
}`,
				ExpectError: regexp.MustCompile(`http_request_option`),
			},
		},
	})
}
