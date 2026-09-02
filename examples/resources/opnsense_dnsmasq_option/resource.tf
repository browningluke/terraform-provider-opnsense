// DNS server
resource "opnsense_dnsmasq_option" "dns_server" {
  type   = "set"
  option = 6
  value  = "192.168.100.1"
}

// NTP server
resource "opnsense_dnsmasq_option" "ntp_server" {
  type      = "set"
  interface = "lan"

  tag = [
    "workstations"
  ]

  option = 42
  value  = "192.168.100.1"

  force       = true
  description = "NTP server for LAN workstations"
}

//
resource "opnsense_dnsmasq_option" "vendor_class" {
  type     = "match"
  option   = 60
  set_tag  = "pxe-client"
  value    = "PXEClient"
  description = "Identify PXE clients by DHCP vendor class"
}

//
resource "opnsense_dnsmasq_option" "dns_server_ipv6" {
  type    = "set"
  option6 = 23
  value   = "fd00:abcd::1"

  description = "IPv6 DNS server for DHCPv6 clients"
}

//

