// Enabled with description
resource "opnsense_route" "one_route" {
  description = "Example route"
  gateway = "LAN_DHCP"
  network = "10.9.0.0/24"
}

// Disabled without description
resource "opnsense_route" "two_route" {
  enabled = false

  gateway = "LAN"
  network = "10.10.0.0/24"
}




