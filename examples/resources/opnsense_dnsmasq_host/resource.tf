// Small example
resource "opnsense_dnsmasq_host" "test_sm" {
  hostname          = "test"
  ip_addresses      = ["10.8.2.1", "::10"]
  harware_addresses = ["00:25:96:12:34:55"]
  client_id         = "01:02:f3"
}

// Full example
resource "opnsense_dnsmasq_host" "test_xl" {
  hostname          = "test"
  domain            = "example.com"
  is_local_domain   = true
  alias_records     = ["alias.example.com"]
  cname_records     = ["cname.example.com"]
  ip_addresses      = ["10.8.2.1", "::10"]
  harware_addresses = ["00:25:96:12:34:55"]
  client_id         = "01:02:f3"
  is_ignored        = false
  comment           = "This is a test host"
  description       = "Test host"
}
