package ipsec_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecChildResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccChildResourceConfig(
					"1",                                // enabled
					"connection-uuid-123",              // connection
					[]string{"aes128-sha256-modp2048"}, // proposals
					"0",                                // sha256_96
					"start",                            // start_action
					"none",                             // close_action
					"clear",                            // dpd_action
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
					resource.TestCheckResourceAttrPair("opnsense_ipsec_child.test", "ipsec_connection", "opnsense_ipsec_connection.parent", "id"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.0", "aes128-sha256-modp2048"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "sha256_96", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "start_action", "start"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "close_action", "none"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "dpd_action", "clear"),
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
				Config: testAccChildResourceConfig(
					"1",                   // enabled
					"connection-uuid-123", // connection
					[]string{"aes256-sha256-modp2048", "aes128-sha256-modp2048"}, // proposals - updated
					"1",         // sha256_96 - updated
					"route",     // start_action - updated
					"trap",      // close_action - updated
					"clear",     // dpd_action - updated
					"transport", // mode - updated
					"0",         // install_policies - updated
					[]string{"192.168.1.0/24", "192.168.2.0/24"}, // local_networks - updated
					[]string{"10.0.0.0/24", "10.1.0.0/24"},       // remote_networks - updated
					"55",                                         // request_id - updated
					"3600",                                       // rekey_time - updated
					"Updated Test IPsec Child",                   // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "proposals.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "sha256_96", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "start_action", "route"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "close_action", "trap"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "dpd_action", "clear"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "mode", "transport"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "install_policies", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "local_networks.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "remote_networks.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "request_id", "55"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "rekey_time", "3600"),
					resource.TestCheckResourceAttr("opnsense_ipsec_child.test", "description", "Updated Test IPsec Child"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccChildResourceConfig(
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
resource "opnsense_ipsec_connection" "parent" {
  enabled                  = "1"
  proposals                = ["aes128-sha256-modp2048"]
  unique                   = "no"
  aggressive               = "0"
  version                  = "2"
  mobike                   = "1"
  local_addresses          = ["192.168.1.1"]
  remote_addresses         = ["10.0.0.1"]
  local_port               = ""
  remote_port              = ""
  udp_encapsulation        = "0"
  reauthentication_time    = "3600"
  rekey_time               = "1800"
  ike_lifetime             = "3600"
  dpd_delay                = "10"
  dpd_timeout              = "60"
  send_certificate_request = "1"
  send_certificate         = "ifasked"
  keying_tries             = "1"
  description              = "Test IPsec Connection for Child"
}

resource "opnsense_ipsec_child" "test" {
  enabled          = %[1]q
  ipsec_connection = opnsense_ipsec_connection.parent.id
  proposals        = ["%[2]v"]
  sha256_96        = %[3]q
  start_action     = %[4]q
  close_action     = %[5]q
  dpd_action       = %[6]q
  mode             = %[7]q
  install_policies = %[8]q
  local_networks   = ["%[9]v"]
  remote_networks  = ["%[10]v"]
  request_id       = %[11]q
  rekey_time       = %[12]q
  description      = %[13]q
}
`, enabled, strings.Join(proposals, `", "`), sha256_96, startAction, closeAction,
		dpdAction, mode, installPolicies, strings.Join(localNetworks, `", "`),
		strings.Join(remoteNetworks, `", "`), requestID, rekeyTime, description)
}
