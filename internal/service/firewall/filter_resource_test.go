package firewall_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallFilterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallFilterResourceConfigMinimal("wan", "pass", "in", "any"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "sequence", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "interface.interface.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "interface.interface.*", "wan"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.action", "pass"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.direction", "in"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.protocol", "any"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.quick", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.log", "false"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_filter.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_filter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallFilterResourceConfigMinimalUpdated("wan", "block", "out", "TCP", "Updated minimal filter"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "description", "Updated minimal filter"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "interface.interface.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "interface.interface.*", "wan"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.action", "block"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.direction", "out"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.protocol", "TCP"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.log", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccFirewallFilterResource_HTTPS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallFilterResourceConfigHTTPS("192.168.1.100", "https", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "sequence", "100"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "description", "Allow HTTPS traffic to web server"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "categories.#", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "interface.interface.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "interface.interface.*", "wan"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.quick", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.action", "pass"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.direction", "in"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.protocol", "TCP"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.source.net", "any"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.net", "192.168.1.100"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.port", "https"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.log", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.type", "keep"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_filter.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_filter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallFilterResourceConfigHTTPS("192.168.1.200", "80", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.net", "192.168.1.200"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.port", "80"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.log", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccFirewallFilterResource_ICMP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallFilterResourceConfigICMP([]string{"echoreq", "echorep"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "description", "Allow ICMP echo request and reply"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "interface.interface.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "interface.interface.*", "wan"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.action", "pass"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.direction", "in"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.protocol", "ICMP"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.icmp_type.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.icmp_type.*", "echoreq"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.icmp_type.*", "echorep"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_filter.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_filter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallFilterResourceConfigICMP([]string{"unreach", "timex"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.icmp_type.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.icmp_type.*", "unreach"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.icmp_type.*", "timex"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccFirewallFilterResource_Comprehensive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallFilterResourceConfigComprehensive(
					"192.168.1.0/24",
					"1024-65535",
					"10.0.0.10",
					"https",
					3600,
					60000,
					120000,
					10000,
					"webtraffic",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "sequence", "200"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "no_xmlrpc_sync", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "description", "Comprehensive example with all attributes"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "categories.#", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "interface.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "interface.interface.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "interface.interface.*", "wan"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.quick", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.action", "pass"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.allow_options", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.direction", "in"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.ip_protocol", "inet"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.protocol", "TCP"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.source.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.source.net", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.source.port", "1024-65535"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.net", "10.0.0.10"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.port", "https"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.log", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.tcp_flags.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.tcp_flags.*", "syn"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.tcp_flags.*", "ack"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.tcp_flags_out_of.#", "4"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.tcp_flags_out_of.*", "syn"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.tcp_flags_out_of.*", "ack"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.tcp_flags_out_of.*", "fin"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_filter.test", "filter.tcp_flags_out_of.*", "rst"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.type", "keep"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.policy", ""),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.timeout", "3600"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.adaptive_timeouts.start", "60000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.adaptive_timeouts.end", "120000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.states", "10000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.source_nodes", "1000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.source_states", "500"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.source_connections", "100"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.new_connections.count", "10"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.new_connections.seconds", "10"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.overload_table", ""),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.no_pfsync", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "traffic_shaping.shaper", ""),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "traffic_shaping.reverse_shaper", ""),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "source_routing.gateway", ""),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "source_routing.disable_reply_to", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "source_routing.reply_to", ""),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "priority.match", "-1"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "priority.set", "3"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "priority.low_delay_set", "4"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "priority.match_tos", "lowdelay"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "internal_tagging.set_local", "webtraffic"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "internal_tagging.match_local", ""),
					resource.TestCheckResourceAttrSet("opnsense_firewall_filter.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_filter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallFilterResourceConfigComprehensive(
					"192.168.2.0/24",
					"2048-65535",
					"10.0.0.20",
					"http",
					7200,
					80000,
					150000,
					20000,
					"updated_tag",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.source.net", "192.168.2.0/24"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.source.port", "2048-65535"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.net", "10.0.0.20"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "filter.destination.port", "http"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.timeout", "7200"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.adaptive_timeouts.start", "80000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.adaptive_timeouts.end", "150000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "stateful_firewall.max.states", "20000"),
					resource.TestCheckResourceAttr("opnsense_firewall_filter.test", "internal_tagging.set_local", "updated_tag"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Helper functions to generate test configurations

func testAccFirewallFilterResourceConfigMinimal(iface, action, direction, protocol string) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_filter" "test" {
  interface = {
    interface = [%[1]q]
  }

  filter = {
    action    = %[2]q
    direction = %[3]q
    protocol  = %[4]q
  }
}
`, iface, action, direction, protocol)
}

func testAccFirewallFilterResourceConfigMinimalUpdated(iface, action, direction, protocol, description string) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_filter" "test" {
  description = %[5]q

  interface = {
    interface = [%[1]q]
  }

  filter = {
    action    = %[2]q
    direction = %[3]q
    protocol  = %[4]q
    log       = true
  }
}
`, iface, action, direction, protocol, description)
}

func testAccFirewallFilterResourceConfigHTTPS(destNet, destPort string, log bool) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_category" "test" {
  name = "Test Category"
}

resource "opnsense_firewall_filter" "test" {
  enabled     = true
  sequence    = 100
  description = "Allow HTTPS traffic to web server"

  categories = [
    opnsense_firewall_category.test.id,
  ]

  interface = {
    interface = ["wan"]
  }

  filter = {
    quick     = true
    action    = "pass"
    direction = "in"
    protocol  = "TCP"

    source = {
      net = "any"
    }

    destination = {
      net  = %[1]q
      port = %[2]q
    }

    log = %[3]t
  }

  stateful_firewall = {
    type = "keep"
  }
}
`, destNet, destPort, log)
}

func testAccFirewallFilterResourceConfigICMP(icmpTypes []string) string {
	icmpTypeList := ""
	for i, icmpType := range icmpTypes {
		if i > 0 {
			icmpTypeList += ", "
		}
		icmpTypeList += fmt.Sprintf("%q", icmpType)
	}

	return fmt.Sprintf(`
resource "opnsense_firewall_filter" "test" {
  description = "Allow ICMP echo request and reply"

  interface = {
    interface = ["wan"]
  }

  filter = {
    action    = "pass"
    direction = "in"
    protocol  = "ICMP"
	icmp_type = [%s]
  }
}
`, icmpTypeList)
}

func testAccFirewallFilterResourceConfigComprehensive(
	sourceNet, sourcePort, destNet, destPort string,
	timeout, adaptiveStart, adaptiveEnd, maxStates int,
	tag string,
) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_category" "test" {
  name = "Test Category"
}

resource "opnsense_firewall_filter" "test" {
  enabled        = true
  sequence       = 200
  no_xmlrpc_sync = false
  description    = "Comprehensive example with all attributes"

  categories = [
    opnsense_firewall_category.test.id,
  ]

  interface = {
    invert = false
    interface = [
      "wan",
    ]
  }

  filter = {
    quick         = true
    action        = "pass"
    allow_options = false
    direction     = "in"
    ip_protocol   = "inet"
    protocol      = "TCP"

    source = {
      invert = false
      net    = %[1]q
      port   = %[2]q
    }

    destination = {
      invert = false
      net    = %[3]q
      port   = %[4]q
    }

    log = true

    tcp_flags = [
      "syn",
      "ack",
    ]

    tcp_flags_out_of = [
      "syn",
      "ack",
      "fin",
      "rst",
    ]

    schedule = ""
  }

  stateful_firewall = {
    type    = "keep"
    policy  = ""
    timeout = %[5]d

    adaptive_timeouts = {
      start = %[6]d
      end   = %[7]d
    }

    max = {
      states             = %[8]d
      source_nodes       = 1000
      source_states      = 500
      source_connections = 100

      new_connections = {
        count   = 10
        seconds = 10
      }
    }

    overload_table = ""
    no_pfsync      = false
  }

  traffic_shaping = {
    shaper         = ""
    reverse_shaper = ""
  }

  source_routing = {
    gateway          = ""
    disable_reply_to = false
    reply_to         = ""
  }

  priority = {
    match         = -1
    set           = 3
    low_delay_set = 4
    match_tos     = "lowdelay"
  }

  internal_tagging = {
    set_local   = %[9]q
    match_local = ""
  }
}
`, sourceNet, sourcePort, destNet, destPort, timeout, adaptiveStart, adaptiveEnd, maxStates, tag)
}
