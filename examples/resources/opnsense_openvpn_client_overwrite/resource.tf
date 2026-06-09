// Pin a specific tunnel address to a client connecting with this common name.
resource "opnsense_openvpn_client_overwrite" "alice" {
  common_name    = "alice"
  description    = "Alice — static tunnel IP"
  tunnel_network = "10.10.20.10/24"
  push_reset     = false

  dns_servers = ["10.10.20.1"]
  dns_domain  = ["example.lan"]
}
