package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsmasqBootResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDnsmasqBootResourceConfig(
					"wan",
					"pxelinux.0",
					"bootserver",
					"192.168.1.10",
					"Test boot configuration",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "file_name", "pxelinux.0"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "server_name", "bootserver"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "server_address", "192.168.1.10"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "description", "Test boot configuration"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_boot.test", "id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "opnsense_dnsmasq_boot.test",
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Update and Read testing
			{
				Config: testAccDnsmasqBootResourceConfig(
					"wan",
					"grubx64.efi",
					"updated-bootserver",
					"192.168.2.20",
					"Updated boot configuration",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "file_name", "grubx64.efi"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "server_name", "updated-bootserver"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "server_address", "192.168.2.20"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_boot.test", "description", "Updated boot configuration"),
				),
			},
		},
	})

}

func testAccDnsmasqBootResourceConfig(
	iface string,
	fileName string,
	serverName string,
	serverAddress string,
	description string,
) string {
	return fmt.Sprintf(`

resource "opnsense_dnsmasq_boot" "test" {
interface      = "%s"
file_name      = "%s"
server_name    = "%s"
server_address = "%s"
description    = "%s"
}
`,
		iface,
		fileName,
		serverName,
		serverAddress,
		description,
	)
}
