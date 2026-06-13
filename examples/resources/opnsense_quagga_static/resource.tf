# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_static.static quagga_static

resource "opnsense_quagga_static" "static" {
  enabled = true
}
