// Small example
resource "opnsense_ipsec_connection" "example" {
  enabled                  = "1"
  proposals                = ["default"]
  unique                   = "no"
  aggressive               = "0"
  version                  = "2"
  mobike                   = "1"
  local_addresses          = ["192.168.1.1"]
  remote_addresses         = ["10.0.0.1"]
  local_port               = ""
  remote_port              = ""
  udp_encapsulation        = "0"
  reauthentication_time    = "3600"
  rekey_time               = "1800"
  ike_lifetime             = "3600"
  dpd_delay                = "10"
  dpd_timeout              = "60"
  send_certificate_request = "1"
  send_certificate         = "ifasked"
  keying_tries             = "1"
  description              = "Example IPsec Connection"
}

resource "opnsense_ipsec_child" "example" {
  enabled          = "1"
  ipsec_connection = opnsense_ipsec_connection.example.id
  proposals        = ["default"]
  sha256_96        = "0"
  start_action     = "trap|start"
  close_action     = "none"
  dpd_action       = "start"
  mode             = "tunnel"
  install_policies = "1"
  local_networks   = ["192.168.1.0/24"]
  remote_networks  = ["10.0.0.0/24"]
  request_id       = "100"
  rekey_time       = "1800"
  description      = "Example IPsec Child"
}
