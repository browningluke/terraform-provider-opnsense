resource "opnsense_quagga_bgp_peer_group" "example" {
  name      = "PEER_GROUP_EXAMPLE"
  remote_as = "65100"
  family    = "ipv4"
}
