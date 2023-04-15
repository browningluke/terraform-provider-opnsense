// OPNsense generates a device name
resource "opnsense_interfaces_vlan" "vlan" {
  description = "Example vlan"
  tag = 10
  priority = 0
  parent = "vtnet0"
}

// Manually configure a device name
resource "opnsense_interfaces_vlan" "vlan04" {
  description = "Example vlan 4"
  tag = 50
  priority = 5
  parent = "vtnet0"
  device = "vlan04"
}
