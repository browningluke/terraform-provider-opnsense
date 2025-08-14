// Small example
resource "opnsense_ipsec_vti" "example" {
  enabled     = "1"
  description = "Example IPsec VTI"

  request_id = "100"

  local_ip          = "1.2.3.4"
  remote_ip         = "5.6.7.8"
  tunnel_local_ip   = "100.64.101.100"
  tunnel_remote_ip  = "100.64.102.100"
  tunnel_local_ip2  = ""
  tunnel_remote_ip2 = ""
}
