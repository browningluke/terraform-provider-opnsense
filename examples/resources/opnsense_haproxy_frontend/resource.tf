# Plain HTTP frontend. Its rule list is set with opnsense_haproxy_frontend_actions.
# Binding port 80 on all addresses needs the web GUI's HTTP redirect disabled
# (System > Settings > Administration), which otherwise holds port 80.
resource "opnsense_haproxy_frontend" "http" {
  name            = "http-80"
  bind            = ["0.0.0.0:80"]
  default_backend = opnsense_haproxy_backend.web.id
}
