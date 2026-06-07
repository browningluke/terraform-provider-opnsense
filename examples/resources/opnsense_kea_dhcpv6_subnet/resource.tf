# Minimal example
resource "opnsense_kea_dhcpv6_subnet" "lan" {
  subnet      = "2001:db8::/64"
  description = "LAN IPv6"
}

# Full resource
resource "opnsense_kea_dhcpv6_subnet" "example" {
  subnet = "2001:db8::/64"

  pools = [
    "2001:db8::100 - 2001:db8::200",
    "2001:db8::300 - 2001:db8::400"
  ]

  allocator    = "random"
  pd_allocator = "random"

  interface = "em0"

  dns_servers = [
    "2001:db8::53",
    "2001:db8::54"
  ]

  domain_search = [
    "example.com",
    "search.example.com"
  ]

  description = "EXAMPLE IPv6 Subnet"
}
