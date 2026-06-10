// A minimal OpenVPN server instance authenticating against the Local Database.
resource "opnsense_openvpn_instance" "server" {
  description = "example-server"
  vpn_id      = 2001

  role     = "server"
  dev_type = "tun"
  protocol = "udp"
  topology = "subnet"

  server                  = "10.10.20.0/24"
  port                    = 1194
  verify_client_cert      = "none"
  username_as_common_name = true
  auth_mode               = ["Local Database"]

  push_route = ["192.168.1.0/24"]
}

// An OpenVPN client instance dialing a remote server.
resource "opnsense_openvpn_instance" "client" {
  description = "example-client"
  vpn_id      = 2002

  role     = "client"
  dev_type = "tun"
  protocol = "udp"

  remote                = ["vpn.example.com:1194"]
  verify_client_cert    = "none"
  remote_cert_tls       = true
}
