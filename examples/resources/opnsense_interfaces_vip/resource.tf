
resource "opnsense_interfaces_vip" "vip" {
  description = "Example vip"
  interface = "wan"
  mode = "proxyarp"
  network = "10.10.1.32/32"
}