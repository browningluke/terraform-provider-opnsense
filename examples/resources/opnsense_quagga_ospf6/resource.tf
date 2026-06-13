# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_ospf6.ospf6 quagga_ospf6

resource "opnsense_quagga_ospf6" "ospf6" {
  enabled   = true
  router_id = "10.0.0.1"
}
