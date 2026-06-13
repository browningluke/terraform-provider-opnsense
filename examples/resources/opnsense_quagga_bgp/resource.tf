# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_bgp.bgp quagga_bgp

resource "opnsense_quagga_bgp" "bgp" {
  enabled   = true
  as_number = "65001"
  router_id = "10.0.0.1"
}
