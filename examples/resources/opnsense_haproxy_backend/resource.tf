resource "opnsense_haproxy_backend" "web" {
  name           = "web"
  linked_servers = [opnsense_haproxy_server.web.id]
}
