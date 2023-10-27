resource "opnsense_firewall_filter" "example_one" {
  enabled = false

  sequence = 1
  action   = "block"
  quick    = false

  interface = [
    "lan",
    "lo0",
  ]

  direction = "in"
  ip_protocol = "inet"
  protocol    = "UDP"

  source = {
    net    = "any"
    invert = true
  }

  destination = {
    net    = "examplealias"
    port   = "https"
  }

  log = false
  description = "example rule"
}

resource "opnsense_firewall_filter" "example_two" {
  action = "pass"
  interface = [
    "wan",
  ]

  direction = "in"
  protocol  = "TCP"

  source = {
    net = "wan" # This is equiv. to WAN Net
  }

  destination = {
    net  = "10.8.0.1"
    port = "443"
  }

  description = "example rule"
}

resource "opnsense_firewall_filter" "example_three" {
  action = "pass"
  interface = [
    "wan",
  ]

  direction = "out"
  protocol  = "TCP"

  source = {
    net = "192.168.0.0/16"
  }

  destination = {
    net  = "wanip" # This is equiv. to WAN Address
    port = "80-443"
  }

  description = "example rule"
  log         = true
}