resource "opnsense_interfaces_vlan" "vlan" {
  description = "Example vlan"
  tag         = 10
  priority    = 0
  parent      = "vtnet0"
}
