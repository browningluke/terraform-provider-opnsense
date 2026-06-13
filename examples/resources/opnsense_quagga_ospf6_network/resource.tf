resource "opnsense_quagga_ospf6_network" "example" {
  ip_addr  = "2001:db8::"
  net_mask = "32"
  area     = "0.0.0.0"
}
