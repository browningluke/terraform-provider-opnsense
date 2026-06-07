package firewall_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallNatPortForwardResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallNatPortForwardResourceConfig(
					true, false, 100, "wan", "inet", "tcp",
					"192.168.1.0/24", "1024", false,
					"10.10.10.22/32", "8080", false,
					"192.168.10.10", "80",
					"default", "Testing NAT port forward",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "log", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "sequence", "100"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "ip_protocol", "inet"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "source.net", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "source.port", "1024"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "source.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "destination.net", "10.10.10.22/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "destination.port", "8080"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "destination.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "target.ip", "192.168.10.10"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "target.port", "80"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "nat_reflection", "default"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "description", "Testing NAT port forward"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_nat_port_forward.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_nat_port_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallNatPortForwardResourceConfig(
					false, true, 200, "wan", "inet", "udp",
					"10.0.0.0/8", "1024-65535", true,
					"10.10.10.23/32", "53", true,
					"192.168.10.20", "53",
					"enable", "Updated NAT port forward",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "log", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "sequence", "200"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "source.net", "10.0.0.0/8"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "source.port", "1024-65535"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "source.invert", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "destination.net", "10.10.10.23/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "destination.port", "53"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "destination.invert", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "target.ip", "192.168.10.20"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "target.port", "53"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "nat_reflection", "enable"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_port_forward.test", "description", "Updated NAT port forward"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_nat_port_forward.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccFirewallNatPortForwardReflectionResource exercises all three
// nat_reflection values so the schema-to-API mapping (default -> "",
// enable -> "purenat", disable -> "disable") is covered end-to-end.
// An ImportStateVerify step is included to confirm the purenat <-> enable
// round-trip survives a state import.
func TestAccFirewallNatPortForwardReflectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatPortForwardResourceConfig(
					true, false, 1, "wan", "inet", "tcp",
					"", "", false,
					"10.10.10.30/32", "4443", false,
					"192.168.10.30", "443",
					"default", "NAT reflection default",
				),
				Check: resource.TestCheckResourceAttr(
					"opnsense_firewall_nat_port_forward.test", "nat_reflection", "default",
				),
			},
			{
				Config: testAccFirewallNatPortForwardResourceConfig(
					true, false, 1, "wan", "inet", "tcp",
					"", "", false,
					"10.10.10.30/32", "4443", false,
					"192.168.10.30", "443",
					"enable", "NAT reflection enable",
				),
				Check: resource.TestCheckResourceAttr(
					"opnsense_firewall_nat_port_forward.test", "nat_reflection", "enable",
				),
			},
			// ImportStateVerify confirms the purenat <-> enable mapping
			// survives an import (i.e. the API value round-trips correctly).
			{
				ResourceName:      "opnsense_firewall_nat_port_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallNatPortForwardResourceConfig(
					true, false, 1, "wan", "inet", "tcp",
					"", "", false,
					"10.10.10.30/32", "4443", false,
					"192.168.10.30", "443",
					"disable", "NAT reflection disable",
				),
				Check: resource.TestCheckResourceAttr(
					"opnsense_firewall_nat_port_forward.test", "nat_reflection", "disable",
				),
			},
		},
	})
}

// testAccFirewallNatPortForwardResourceConfig builds a NAT port forward HCL
// block. If sourceNet is empty the source block is omitted so the resource
// defaults apply (avoids the empty-string regex validator on source.port).
func testAccFirewallNatPortForwardResourceConfig(
	enabled, log bool,
	sequence int64,
	iface, ipProtocol, protocol,
	sourceNet, sourcePort string, sourceInvert bool,
	destNet, destPort string, destInvert bool,
	targetIP, targetPort,
	natReflection, description string,
) string {
	sourceBlock := ""
	if sourceNet != "" {
		sourceBlock = fmt.Sprintf(`
  source = {
    net    = %q
    port   = %q
    invert = %t
  }`, sourceNet, sourcePort, sourceInvert)
	}

	return fmt.Sprintf(`
resource "opnsense_firewall_nat_port_forward" "test" {
  enabled     = %[1]t
  log         = %[2]t
  sequence    = %[3]d
  interface   = %[4]q
  ip_protocol = %[5]q
  protocol    = %[6]q%[7]s
  destination = {
    net    = %[8]q
    port   = %[9]q
    invert = %[10]t
  }
  target = {
    ip   = %[11]q
    port = %[12]q
  }
  nat_reflection = %[13]q
  description    = %[14]q
}
`,
		enabled, log, sequence,
		iface, ipProtocol, protocol,
		sourceBlock,
		destNet, destPort, destInvert,
		targetIP, targetPort,
		natReflection, description,
	)
}
