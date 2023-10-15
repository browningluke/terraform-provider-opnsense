// Configure a prefix list
resource "opnsense_quagga_bgp_prefixlist" "example0" {
  enabled     = false

  description = "prefixlist0"
  name = "example0"

  number = 1234
  action = "permit"

  network = "10.10.0.0"
}

// Configure a route map
resource "opnsense_quagga_bgp_routemap" "example0" {
  enabled     = false
  description = "routemap0"

  name   = "example0"
  action = "deny"

  route_map_id = 100
  set = "local-preference 300"
}

// Configure a neighbor
resource "opnsense_quagga_bgp_neighbor" "example0" {
  enabled     = false

  description = "neighbor0"

  peer_ip   = "1.1.1.1"
  remote_as = 255

  md5_password            = "12345"
  weight                  = 1
  local_ip                = "2.2.2.2"
  update_source           = "wan"
  link_local_interface    = "wireguard"

  next_hop_self           = true
  next_hop_self_all       = true
  multi_hop               = true
  multi_protocol          = true
  rr_client               = true
  bfd                     = true

  keep_alive              = 100
  hold_down               = 10
  connect_timer           = 10

  default_route           = true
  as_override             = true
  disable_connected_check = true
  attribute_unchanged     = "as-path"

  prefix_list_in = opnsense_quagga_bgp_prefixlist.example0.id
  route_map_out = opnsense_quagga_bgp_routemap.example0.id
}
