# Send ACME HTTP-01 challenges for *.web.example.com to the web backend.
resource "opnsense_haproxy_action" "web_acme" {
  name        = "web-acme-challenge"
  linked_acls = [opnsense_haproxy_acl.host.id, opnsense_haproxy_acl.acme_path.id]
  operator    = "and"
  type        = "use_backend"
  use_backend = opnsense_haproxy_backend.web.id
}

# Redirect everything except ACME challenges to HTTPS, same host and path.
# HAProxy runs http-request rules before use_backend rules, so without the
# "unless" condition the challenges would be redirected too.
resource "opnsense_haproxy_action" "https_redirect" {
  name                = "https-redirect"
  test_type           = "unless"
  linked_acls         = [opnsense_haproxy_acl.acme_path.id]
  type                = "http-request"
  http_request_action = "redirect"
  http_request_option = "scheme https code 301"
}
