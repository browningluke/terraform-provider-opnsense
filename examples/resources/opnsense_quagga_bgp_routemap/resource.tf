// Configure an AS Path
resource "opnsense_quagga_bgp_aspath" "example0" {
  enabled     = false
  description = "aspath0"

  number = 123
  action = "permit"

  as = "_2$"
}

// Configure a prefix list
resource "opnsense_quagga_bgp_prefixlist" "example0" {
  enabled     = false

  description = "prefixlist0"
  name = "example0"

  number = 1234
  action = "permit"

  network = "10.10.0.0"
}
// Configure a community list
resource "opnsense_quagga_bgp_communitylist" "example0" {
  enabled     = false
  description = "communitylist0"

  number     = 100
  seq_number = 99
  action     = "deny"

  community = "example.*"
}

// Configure a route map
resource "opnsense_quagga_bgp_routemap" "example0" {
  enabled     = false
  description = "routemap0"

  name   = "example0"
  action = "deny"

  route_map_id = 100

  aspaths = [
    opnsense_quagga_bgp_aspath.example0.id
  ]

  prefix_lists = [
    opnsense_quagga_bgp_prefixlist.example0.id
  ]

  community_lists = [
    opnsense_quagga_bgp_communitylist.example0.id
  ]

  set = "local-preference 300"
}
