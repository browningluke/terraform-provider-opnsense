// Query Forward
resource "opnsense_unbound_forward" "query" {
  domain = "example.lan"
  server_ip = "192.168.1.2"
  server_port = 853
}

// DoT forward
resource "opnsense_unbound_forward" "dot" {
  enabled = false
  type = "dot"

  domain = "example.dev"
  server_ip = "192.168.1.1"
  server_port = 53
  verify_cn = "example.dev"
}
