package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIpsecConnectionResourceConfig(
					"1",                                    // enabled
					[]string{"aes128-sha256-modp2048"},     // proposals
					"no",                                   // unique
					"0",                                    // aggressive
					"2",                                    // version
					"1",                                    // mobike
					[]string{"192.168.1.1", "192.168.2.1"}, // local_addresses
					[]string{"10.0.0.1"},                   // remote_addresses
					"",                                     // local_port (empty)
					"",                                     // remote_port (empty)
					"0",                                    // udp_encapsulation
					"3600",                                 // reauthentication_time
					"1800",                                 // rekey_time
					"3600",                                 // ike_lifetime
					"120",                                  // dpd_delay
					"540",                                  // dpd_timeout
					[]string{},                             // ip_pools
					"1",                                    // send_certificate_request
					"always",                               // send_certificate
					"3",                                    // keying_tries
					"Test IPsec Connection",                // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "proposals.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "proposals.0", "aes128-sha256-modp2048"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "unique", "no"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "aggressive", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "version", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "mobike", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.0", "192.168.1.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.1", "192.168.2.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_addresses.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_addresses.0", "10.0.0.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_port", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_port", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "udp_encapsulation", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "description", "Test IPsec Connection"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_connection.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_ipsec_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIpsecConnectionResourceConfig(
					"1", // enabled
					[]string{"aes256-sha256-modp2048", "aes128-sha256-modp2048"}, // proposals - updated
					"no",                                   // unique
					"0",                                    // aggressive
					"2",                                    // version
					"1",                                    // mobike - updated
					[]string{"192.168.1.1", "192.168.1.2"}, // local_addresses - updated
					[]string{"10.0.0.1"},                   // remote_addresses
					"",                                     // local_port (empty)
					"",                                     // remote_port (empty)
					"0",                                    // udp_encapsulation
					"7200",                                 // reauthentication_time - updated
					"3600",                                 // rekey_time - updated
					"7200",                                 // ike_lifetime - updated
					"30",                                   // dpd_delay - updated
					"120",                                  // dpd_timeout - updated
					[]string{},                             // ip_pools
					"1",                                    // send_certificate_request
					"never",                                // send_certificate - updated
					"3",                                    // keying_tries - updated
					"Updated Test IPsec Connection",        // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "proposals.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "proposals.1", "aes256-sha256-modp2048"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "proposals.0", "aes128-sha256-modp2048"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "mobike", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.0", "192.168.1.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.1", "192.168.1.2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "reauthentication_time", "7200"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "rekey_time", "3600"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "ike_lifetime", "7200"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "dpd_delay", "30"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "dpd_timeout", "120"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "send_certificate", "never"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "keying_tries", "3"),
					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "description", "Updated Test IPsec Connection"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// func TestAccIpsecConnectionResource_MinimalConfig(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccIpsecConnectionResourceConfigMinimal(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "enabled", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "proposals.#", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "unique", "no"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "version", "2"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.#", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_addresses.#", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "description", "Test IPsec Connection"),
// 					resource.TestCheckResourceAttrSet("opnsense_ipsec_connection.test", "id"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccIpsecConnectionResource_IKEv1(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccIpsecConnectionResourceConfigIKEv1(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "enabled", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "version", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "aggressive", "1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "mobike", "0"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccIpsecConnectionResource_MultipleAddresses(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccIpsecConnectionResourceConfigMultipleAddresses(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.#", "3"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_addresses.#", "2"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.0", "192.168.1.1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.1", "192.168.1.10"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "local_addresses.2", "10.10.10.1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_addresses.0", "203.0.113.1"),
// 					resource.TestCheckResourceAttr("opnsense_ipsec_connection.test", "remote_addresses.1", "203.0.113.10"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccIpsecConnectionResourceConfig(
	enabled string,
	proposals []string,
	unique string,
	aggressive string,
	version string,
	mobike string,
	localAddresses []string,
	remoteAddresses []string,
	localPort string,
	remotePort string,
	udpEncapsulation string,
	reauthenticationTime string,
	rekeyTime string,
	ikeLifetime string,
	dpdDelay string,
	dpdTimeout string,
	ipPools []string,
	sendCertificateRequest string,
	sendCertificate string,
	keyingTries string,
	description string,
) string {
	ipPoolsLine := ""
	if len(ipPools) > 0 {
		ipPoolsLine = fmt.Sprintf("  ip_pools               = %q\n", ipPools)
	}
	rval := fmt.Sprintf(`
resource "opnsense_ipsec_connection" "test" {
  enabled                = %[1]q
  proposals              = ["%[2]v"]
  unique                 = %[3]q
  aggressive             = %[4]q
  version                = %[5]q
  mobike                 = %[6]q
  local_addresses        = ["%[7]v"]
  remote_addresses       = ["%[8]v"]
  local_port             = %[9]q
  remote_port            = %[10]q
  udp_encapsulation      = %[11]q
  reauthentication_time  = %[12]q
  rekey_time             = %[13]q
  ike_lifetime           = %[14]q
  dpd_delay              = %[15]q
  dpd_timeout            = %[16]q
%[17]s  send_certificate_request = %[18]q
  send_certificate       = %[19]q
  keying_tries           = %[20]q
  description            = %[21]q
}
`, enabled, strings.Join(proposals, `", "`), unique, aggressive, version, mobike, strings.Join(localAddresses, `", "`), strings.Join(remoteAddresses, `", "`),
		localPort, remotePort, udpEncapsulation, reauthenticationTime, rekeyTime, ikeLifetime,
		dpdDelay, dpdTimeout, ipPoolsLine, sendCertificateRequest, sendCertificate, keyingTries, description)
	fmt.Println("Generated config:", rval)
	return rval
}

func testAccIpsecConnectionResourceConfigMinimal() string {
	return `
resource "opnsense_ipsec_connection" "test" {
  enabled                = "1"
  proposals              = ["aes128-sha256-modp2048"]
  unique                 = "no"
  aggressive             = "0"
  version                = "2"
  mobike                 = "1"
  local_addresses        = ["192.168.1.1"]
  remote_addresses       = ["10.0.0.1"]
  local_port             = ""
  remote_port            = ""
  udp_encapsulation      = "0"
  reauthentication_time  = "3600"
  rekey_time             = "1800"
  ike_lifetime           = "3600"
  dpd_delay              = "10"
  dpd_timeout            = "60"
  send_certificate_request = "1"
  send_certificate       = "ifasked"
  keying_tries           = "1"
  description            = "Test IPsec Connection"
}
`
}

func testAccIpsecConnectionResourceConfigIKEv1() string {
	return `
resource "opnsense_ipsec_connection" "test" {
  enabled                = "1"
  proposals              = ["aes128-sha1-modp1024"]
  unique                 = "no"
  aggressive             = "1"
  version                = "auto"
  mobike                 = "0"
  local_addresses        = ["192.168.1.1"]
  remote_addresses       = ["10.0.0.1"]
  local_port             = ""
  remote_port            = ""
  udp_encapsulation      = "0"
  reauthentication_time  = "3600"
  rekey_time             = "1800"
  ike_lifetime           = "3600"
  dpd_delay              = "10"
  dpd_timeout            = "60"
  send_certificate_request = "0"
  send_certificate       = "never"
  keying_tries           = "1"
  description            = "IKEv1 Test Connection"
}
`
}

func testAccIpsecConnectionResourceConfigMultipleAddresses() string {
	return `
resource "opnsense_ipsec_connection" "test" {
  enabled                = "1"
  proposals              = ["aes256-sha256-modp2048"]
  unique                 = "no"
  aggressive             = "0"
  version                = "2"
  mobike                 = "1"
  local_addresses        = ["192.168.1.1", "192.168.1.10", "10.10.10.1"]
  remote_addresses       = ["203.0.113.1", "203.0.113.10"]
  local_port             = ""
  remote_port            = ""
  udp_encapsulation      = "0"
  reauthentication_time  = "3600"
  rekey_time             = "1800"
  ike_lifetime           = "3600"
  dpd_delay              = "10"
  dpd_timeout            = "60"
  send_certificate_request = "1"
  send_certificate       = "ifasked"
  keying_tries           = "1"
  description            = "Multiple Addresses Test"
}
`
}
