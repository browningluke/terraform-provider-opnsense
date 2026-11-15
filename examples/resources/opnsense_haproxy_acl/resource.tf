resource "opnsense_haproxy_acl" "host" {
  name       = "host-web"
  expression = "hdr_end"
  value      = ".web.example.com"
}

resource "opnsense_haproxy_acl" "acme_path" {
  name       = "path-acme-challenge"
  expression = "path_beg"
  value      = "/.well-known/acme-challenge/"
}
