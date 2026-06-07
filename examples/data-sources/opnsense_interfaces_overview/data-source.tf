// Look up the live state of the WAN interface
data "opnsense_interfaces_overview" "wan" {
  device = "vtnet0"
}

// Use the primary IPv4 address in another resource
output "wan_ip" {
  value = data.opnsense_interfaces_overview.wan.addr4
}
