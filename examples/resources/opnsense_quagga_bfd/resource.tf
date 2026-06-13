# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_bfd.bfd quagga_bfd

resource "opnsense_quagga_bfd" "bfd" {
  enabled = true
}
