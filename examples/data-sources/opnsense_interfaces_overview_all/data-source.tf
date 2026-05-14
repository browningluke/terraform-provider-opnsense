// Get all interfaces
data "opnsense_interfaces_overview_all" "all" {}

// Filter for only enabled interfaces
output "enabled_interfaces" {
  value = [for i in data.opnsense_interfaces_overview_all.all.interfaces : i if i.enabled]
}

// Get the primary IPv4 address of each interface that has one
output "interface_ipv4_map" {
  value = { for i in data.opnsense_interfaces_overview_all.all.interfaces : i.identifier => i.addr4 if i.addr4 != "" }
}

// Filter for VLAN interfaces
output "vlan_interfaces" {
  value = [for i in data.opnsense_interfaces_overview_all.all.interfaces : i if i.vlan.tag != ""]
}
