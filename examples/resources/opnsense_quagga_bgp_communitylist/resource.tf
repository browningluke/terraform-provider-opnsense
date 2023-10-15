// Configure a community list
resource "opnsense_quagga_bgp_communitylist" "example0" {
  enabled     = false
  description = "communitylist0"

  number     = 100
  seq_number = 99
  action     = "deny"

  community = "example.*"
}
