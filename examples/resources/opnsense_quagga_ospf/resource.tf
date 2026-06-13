# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_ospf.ospf quagga_ospf

resource "opnsense_quagga_ospf" "ospf" {
  enabled   = true
  router_id = "10.0.0.1"
}
