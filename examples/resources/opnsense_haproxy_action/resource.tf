# Send ACME HTTP-01 challenges for *.web.example.com to the web backend.
resource "opnsense_haproxy_action" "web_acme" {
  name        = "web-acme-challenge"
  linked_acls = [opnsense_haproxy_acl.host.id, opnsense_haproxy_acl.acme_path.id]
  operator    = "and"
  use_backend = opnsense_haproxy_backend.web.id
}
