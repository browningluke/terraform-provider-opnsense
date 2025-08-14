package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecChildResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIpsecChildResourceConfig(
					"1",                                // enabled
					"connection-uuid-123",              // connection
					[]string{"aes128-sha256-modp2048"}, // proposals
					"0",                                // sha256_96
					"start",                            // start_action
					"none",                             // close_action
					"hold",                             // dpd_action
					"tunnel",                           // mode
					"1",                                // install_policies
					[]string{"192.168.1.0/24"},         // local_networks
					[]string{"10.0.0.0/24"},            // remote_networks
					"",                                 // request_id (empty)
					"0",                                // rekey_time
					"Test IPsec Child",                 // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "ipsec_connection", "connection-uuid-123"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.0", "aes128-sha256-modp2048"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "sha256_96", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "start_action", "start"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "close_action", "none"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "dpd_action", "hold"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "mode", "tunnel"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "install_policies", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "local_networks.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "local_networks.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "remote_networks.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "remote_networks.0", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "request_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "rekey_time", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "description", "Test IPsec Child"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_child.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_ipsec_child.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIpsecChildResourceConfig(
					"1",                   // enabled
					"connection-uuid-123", // connection
					[]string{"aes256-sha256-modp2048", "aes128-sha256-modp2048"}, // proposals - updated
					"1",         // sha256_96 - updated
					"route",     // start_action - updated
					"trap",      // close_action - updated
					"restart",   // dpd_action - updated
					"transport", // mode - updated
					"0",         // install_policies - updated
					[]string{"192.168.1.0/24", "192.168.2.0/24"}, // local_networks - updated
					[]string{"10.0.0.0/24", "10.1.0.0/24"},       // remote_networks - updated
					"custom-request-id",                          // request_id - updated
					"3600",                                       // rekey_time - updated
					"Updated Test IPsec Child",                   // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "sha256_96", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "start_action", "route"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "close_action", "trap"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "dpd_action", "restart"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "mode", "transport"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "install_policies", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "local_networks.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "remote_networks.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "request_id", "custom-request-id"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "rekey_time", "3600"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "description", "Updated Test IPsec Child"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIpsecChildResource_MinimalConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecChildResourceConfigMinimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "sha256_96", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "start_action", "start"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "close_action", "none"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "dpd_action", "hold"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "mode", "tunnel"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "install_policies", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "local_networks.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "remote_networks.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "rekey_time", "0"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_child.test", "id"),
				),
			},
		},
	})
}

func TestAccIpsecChildResource_MultipleProposals(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecChildResourceConfigMultipleProposals(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.#", "3"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "local_networks.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "remote_networks.#", "2"),
				),
			},
		},
	})
}

func testAccIpsecChildResourceConfig(
	enabled string,
	ipsec_connection string,
	proposals []string,
	sha256_96 string,
	startAction string,
	closeAction string,
	dpdAction string,
	mode string,
	installPolicies string,
	localNetworks []string,
	remoteNetworks []string,
	requestID string,
	rekeyTime string,
	description string,
) string {
	return fmt.Sprintf(`
resource "opnsense_ipsec_child" "test" {
  enabled          = %[1]q
  ipsec_connection = %[2]q
  proposals        = ["%[3]v"]
  sha256_96        = %[4]q
  start_action     = %[5]q
  close_action     = %[6]q
  dpd_action       = %[7]q
  mode             = %[8]q
  install_policies = %[9]q
  local_networks   = ["%[10]v"]
  remote_networks  = ["%[11]v"]
  request_id       = %[12]q
  rekey_time       = %[13]q
  description      = %[14]q
}
`, enabled, ipsec_connection, strings.Join(proposals, `", "`), sha256_96, startAction, closeAction,
		dpdAction, mode, installPolicies, strings.Join(localNetworks, `", "`),
		strings.Join(remoteNetworks, `", "`), requestID, rekeyTime, description)
}

func testAccIpsecChildResourceConfigMinimal() string {
	return `
resource "opnsense_ipsec_child" "test" {
  enabled          = "1"
  ipsec_connection = "connection-uuid-minimal"
  proposals        = ["aes128-sha256-modp2048"]
  local_networks   = ["192.168.1.0/24"]
  remote_networks  = ["10.0.0.0/24"]
}
`
}

func testAccIpsecChildResourceConfigMultipleProposals() string {
	return `
resource "opnsense_ipsec_child" "test" {
  enabled          = "1"
  ipsec_connection = "connection-uuid-multiple"
  proposals        = ["aes256-sha256-modp2048", "aes128-sha256-modp2048", "3des-sha1-modp1024"]
  local_networks   = ["192.168.1.0/24", "192.168.2.0/24"]
  remote_networks  = ["10.0.0.0/24", "10.1.0.0/24"]
  description      = "Multiple Proposals Test"
}
`
}
