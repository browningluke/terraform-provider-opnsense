resource "opnsense_quagga_ospf_network" "example" {
  ip_addr  = "10.0.0.0"
  area     = "0.0.0.0"
  net_mask = "24"
}
