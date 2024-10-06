// Small example
resource "opnsense_kea_reservation" "test" {
  subnet_id = resource.opnsense_kea_subnet.lan.id

  ip_address = "10.8.2.102"
  mac_address = "00:25:96:12:34:55"

  description = "example host"
}

// LAN subnet example
resource "opnsense_kea_subnet" "lan" {
  subnet = "10.8.0.0/16"
  description = "LAN"
}