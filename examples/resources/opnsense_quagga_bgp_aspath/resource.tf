// Configure an AS Path
resource "opnsense_quagga_bgp_aspath" "example0" {
  enabled     = false
  description = "aspath0"

  number = 123
  action = "permit"

  as = "_2$"
}
