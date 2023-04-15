// 'A' record
resource "opnsense_unbound_host_override" "a_override" {
  enabled = true
  description = "A record override"

  hostname = "*"
  domain = "example.com"
  server = "192.168.1.1"
}

// Enabled alias with description
resource "opnsense_unbound_host_alias" "one_alias" {
  override = opnsense_unbound_host_override.a_override.id

  enabled = true
  hostname = "*"
  domain = "1.example.com"
  description = "Example 1"
}

// Disabled alias without description
resource "opnsense_unbound_host_alias" "two_alias" {
  override = opnsense_unbound_host_override.a_override.id

  enabled = false
  hostname = "*"
  domain = "2.example.com"
}
