// Enabled with description
resource "opnsense_unbound_domain_override" "one_override" {
  enabled = true
  description = "Example override"

  domain = "example.lan"
  server = "192.168.1.1"
}

// Disabled without description
resource "opnsense_unbound_domain_override" "two_override" {
  enabled = false

  domain = "example.arpa"
  server = "192.168.1.100"
}
