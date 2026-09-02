// Set internal dns server
resource "opnsense_dnsmasq_domain" "internal" {
  sequence = 10
  domain   = "example.internal"
  ip       = "192.168.110.53"
}

resource "opnsense_dnsmasq_domain" "internal_dns" {
  sequence = 10
  domain   = "corp.example.com"

  srcip = "192.168.110.1"
  port  = 5353
  ip    = "192.168.110.53"

  description = "Resolve corporate DNS through the internal DNS server"
}

// Firewall alias usage
resource "opnsense_dnsmasq_domain" "updates" {
  sequence       = 100
  domain         = "updates.example.com"
  ip             = "192.168.110.53"
  firewall_alias = "allowed_updates"

  description = "Resolve update servers and populate the firewall alias"
}

// Block DNS usage
resource "opnsense_dnsmasq_domain" "blocked_domain" {
  sequence = 900
  domain   = "blocked.example.com"

  description = "Prevent DNS lookups for this domain"
}
