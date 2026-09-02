// 
resource "opnsense_dnsmasq_range" "lan" {
  interface    = "lan"
  start_address = "192.168.100.100"
  end_address   = "192.168.100.200"
}

//
resource "opnsense_dnsmasq_range" "lan" {
  interface     = "lan"
  tag           = "lan-dhcp"
  start_address = "192.168.100.100"
  end_address   = "192.168.100.200"

  lease_time = 86400

  domain_type = "interface"
  domain      = "example.internal"

  description = "DHCP range for LAN clients"
}

// 
resource "opnsense_dnsmasq_range" "static" {
  interface     = "servers"
  start_address = "192.168.110.1"
  end_address   = "192.168.110.254"

  mode = [
    "static"
  ]

  domain = "example.internal"

  description = "Static DHCP range for infrastructure hosts"
}

//
resource "opnsense_dnsmasq_range" "lan_ipv6" {
  interface     = "lan"
  constructor   = "lan"
  start_address = "::"
  prefix_length = 64

  ra_mode = [
    "slaac",
    "ra-stateless",
    "ra-names"
  ]

  ra_priority = "high"

  ra_interval        = 10
  ra_router_lifetime = 30

  description = "IPv6 Router Advertisements for LAN"
}