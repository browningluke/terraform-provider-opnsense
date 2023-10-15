// Configure a prefix list
resource "opnsense_quagga_bgp_prefixlist" "example0" {
  enabled     = false

  description = "prefixlist0"
  name = "example0"

  number = 1234
  action = "permit"

  network = "10.10.0.0"
}
