resource "opnsense_haproxy_server" "web" {
  name    = "web01"
  address = "192.0.2.20"
  port    = 80
}
