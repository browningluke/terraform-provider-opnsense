// Network example
resource "opnsense_firewall_alias" "example_one" {
  name = "example_one"

  type = "network"
  content = [
    "10.8.0.1/24",
    "10.8.0.2/24"
  ]

  stats       = true
  description = "Example"
}

// With category
resource "opnsense_firewall_category" "example_one" {
  name  = "example"
  color = "ffaa00"
}

resource "opnsense_firewall_alias" "example_two" {
  name = "example_two"

  type = "geoip"
  content = [
    "FR",
    "CA",
  ]

  categories = [
    opnsense_firewall_category.example_one.id
  ]

  description = "Example two"
}
