// Allow queries from a specific subnet
resource "opnsense_unbound_acl" "lan_allow" {
  name     = "lan-allow"
  action   = "allow"
  networks = ["10.0.0.0/24", "10.0.1.0/24"]
}

// Refuse queries from a management network (returns REFUSED)
resource "opnsense_unbound_acl" "mgmt_refuse" {
  name        = "mgmt-refuse"
  action      = "refuse"
  networks    = ["172.16.0.0/12"]
  description = "Refuse external management traffic"
}

// Disabled deny entry
resource "opnsense_unbound_acl" "guest_deny" {
  enabled  = false
  name     = "guest-deny"
  action   = "deny"
  networks = ["192.168.100.0/24"]
}
