// 'A' record
resource "opnsense_unbound_host_override" "a_override" {
  enabled = true
  description = "A record override"

  hostname = "*"
  domain = "example.com"
  server = "192.168.1.1"
}

// 'AAAA' record
resource "opnsense_unbound_host_override" "aaaa_override" {
  enabled = true

  type = "AAAA"
  hostname = "*"
  domain = "example.com"
  server = "fd00:abcd::1"
}

// 'MX' record
resource "opnsense_unbound_host_override" "mx_override" {
  enabled = false
  description = "MX record override"

  type = "MX"
  hostname = "*"
  domain = "example.com"

  mx_priority = 10
  mx_host = "mail.example.dev"
}
