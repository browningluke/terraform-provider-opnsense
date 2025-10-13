package ipsec_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecVtiResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVtiResourceConfig(
					"1",              // enabled
					"1234",           // request_id
					"2.3.4.5",        // local_ip
					"5.6.7.8",        // remote_ip
					"1.2.3.4",        // tunnel_local_ip
					"4.3.2.1",        // tunnel_remote_ip
					"7.8.9.10",       // tunnel_local_ip2
					"8.7.6.5",        // tunnel_remote_ip2
					"Test IPsec VTI", // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "request_id", "1234"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "local_ip", "2.3.4.5"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "remote_ip", "5.6.7.8"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip", "4.3.2.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip2", "7.8.9.10"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip2", "8.7.6.5"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "description", "Test IPsec VTI"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_vti.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_ipsec_vti.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccVtiResourceConfig(
					"1",                      // enabled
					"5678",                   // request_id - updated
					"10.20.30.40",            // local_ip - updated
					"40.30.20.10",            // remote_ip - updated
					"192.168.100.1",          // tunnel_local_ip - updated
					"192.168.100.2",          // tunnel_remote_ip - updated
					"172.16.1.1",             // tunnel_local_ip2 - updated
					"172.16.1.2",             // tunnel_remote_ip2 - updated
					"Updated Test IPsec VTI", // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "request_id", "5678"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "local_ip", "10.20.30.40"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "remote_ip", "40.30.20.10"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip", "192.168.100.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip", "192.168.100.2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip2", "172.16.1.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip2", "172.16.1.2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "description", "Updated Test IPsec VTI"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIpsecVtiResource_MinimalConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVtiResourceConfigMinimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "local_ip", "192.168.1.10"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "remote_ip", "203.0.113.10"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip", "10.0.1.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip", "10.0.1.2"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_vti.test", "id"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_vti.test", "request_id"),
				),
			},
		},
	})
}

func TestAccIpsecVtiResource_WithOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVtiResourceConfigWithOptionals(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "enabled", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "request_id", "9999"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "local_ip", "172.16.100.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "remote_ip", "172.16.200.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip", "10.100.1.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip", "10.100.1.2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_local_ip2", "10.200.1.1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "tunnel_remote_ip2", "10.200.1.2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_vti.test", "description", "VTI with all optional fields"),
				),
			},
		},
	})
}

func testAccVtiResourceConfig(
	enabled string,
	requestId string,
	localIP string,
	remoteIP string,
	tunnelLocalIP string,
	tunnelRemoteIP string,
	tunnelLocalIP2 string,
	tunnelRemoteIP2 string,
	description string,
) string {
	return fmt.Sprintf(`
resource "opnsense_ipsec_vti" "test" {
  enabled           = %[1]q
  request_id        = %[2]q
  local_ip          = %[3]q
  remote_ip         = %[4]q
  tunnel_local_ip   = %[5]q
  tunnel_remote_ip  = %[6]q
  tunnel_local_ip2  = %[7]q
  tunnel_remote_ip2 = %[8]q
  description       = %[9]q
}
`, enabled, requestId, localIP, remoteIP, tunnelLocalIP, tunnelRemoteIP, tunnelLocalIP2, tunnelRemoteIP2, description)
}

func testAccVtiResourceConfigMinimal() string {
	return `
resource "opnsense_ipsec_vti" "test" {
  local_ip         = "192.168.1.10"
  remote_ip        = "203.0.113.10"
  tunnel_local_ip  = "10.0.1.1"
  tunnel_remote_ip = "10.0.1.2"
  request_id       = "1234"
}
`
}

func testAccVtiResourceConfigWithOptionals() string {
	return `
resource "opnsense_ipsec_vti" "test" {
  enabled           = "0"
  request_id        = "9999"
  local_ip          = "172.16.100.1"
  remote_ip         = "172.16.200.1"
  tunnel_local_ip   = "10.100.1.1"
  tunnel_remote_ip  = "10.100.1.2"
  tunnel_local_ip2  = "10.200.1.1"
  tunnel_remote_ip2 = "10.200.1.2"
  description       = "VTI with all optional fields"
}
`
}
