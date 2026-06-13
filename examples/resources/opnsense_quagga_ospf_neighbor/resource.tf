resource "opnsense_quagga_ospf_neighbor" "example" {
  address     = "192.0.2.1"
  description = "OSPF neighbor"
}
