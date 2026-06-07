resource "opnsense_kea_dhcpv6_peer" "example" {
  name = "example"
  role = "primary"
  url  = "http://[2001:db8::1]:8001/"
}
