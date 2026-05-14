resource "opnsense_kea_dhcpv4_reservation" "example" {
  subnet_id = opnsense_kea_dhcpv4_subnet.lan.id

  ip_address  = "10.8.2.102"
  mac_address = "00:25:96:12:34:55"

  description = "example host"
}

resource "opnsense_kea_dhcpv4_subnet" "lan" {
  subnet      = "10.8.0.0/16"
  description = "LAN"
}
