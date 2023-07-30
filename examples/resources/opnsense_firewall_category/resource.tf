resource "opnsense_firewall_category" "example_one" {
  name  = "example"
  color = "ffaa00"
}

resource "opnsense_firewall_alias" "example_one" {
  name = "example"

  type = "geoip"
  content = [
    "FR",
    "CA",
  ]

  categories = [
    opnsense_firewall_category.example_one.id
  ]

  stats       = true
  description = "Example"
}
