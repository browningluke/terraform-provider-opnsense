resource "opnsense_kea_dhcpv6_reservation" "example" {
  subnet_id = opnsense_kea_dhcpv6_subnet.lan.id

  ip_address = "2001:db8::102"
  duid       = "00:03:00:01:00:25:96:12:34:55"

  hostname = "myhost.example.com"

  domain_search = [
    "example.com"
  ]

  description = "example IPv6 host"
}

resource "opnsense_kea_dhcpv6_subnet" "lan" {
  subnet      = "2001:db8::/64"
  description = "LAN IPv6"
}
