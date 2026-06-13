# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_rip.rip quagga_rip

resource "opnsense_quagga_rip" "rip" {
  enabled = true
  version = "2"
}
