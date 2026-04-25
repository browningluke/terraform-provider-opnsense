// Small example
resource "opnsense_kea_subnet" "lan" {
  subnet = "10.8.0.0/16"
  description = "LAN"
}

// Full resource
resource "opnsense_kea_subnet" "example" {
  subnet = "10.8.0.0/16"

  next_server = "10.8.0.1"

  match_client_id = false

  auto_collect = false

  static_routes = [
    {
      destination_ip = "10.10.10.10"
      router_ip = "10.8.0.1"
    },
    {
      destination_ip = "10.10.10.11"
      router_ip = "10.8.50.1"
    }
  ]

  pools = [
    "10.8.2.1-10.8.2.100",
    "10.8.2.101-10.8.2.200",
    "10.8.3.0/24"
  ]

  routers = [
    "10.8.0.1",
    "10.8.50.2"
  ]

  dns_servers = [
    "10.8.0.160",
    "10.8.0.161"
  ]

  domain_name = "example.com"

  domain_search = [
    "search.example.com",
    "search2.example.com"
  ]

  ntp_servers = [
    "10.10.101.10",
    "10.10.101.11"
  ]

  time_servers = [
    "10.10.101.10",
    "10.10.101.11"
  ]

  tftp_server = "tftp.example.com"
  tftp_bootfile = "bootfile.txt"

  description = "EXAMPLE"
}
