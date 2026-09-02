// Small example
resource "opnsense_dnsmasq_tag" "servers" {
  tag = "servers"
}

// Full example
resource "opnsense_dnsmasq_tag" "workstations" {
  tag = "workstations"
}

resource "opnsense_dnsmasq_range" "workstations" {
  interface     = "lan"
  tag           = "workstations"
  start_address = "192.168.100.100"
  end_address   = "192.168.100.200"

  domain = "example.internal"

  description = "DHCP range for workstations"
}

resource "opnsense_dnsmasq_option" "workstations_ntp" {
  type   = "set"
  tag    = ["workstations"]
  option = 42
  value  = "192.168.100.1"

  description = "NTP server for workstations"
}
