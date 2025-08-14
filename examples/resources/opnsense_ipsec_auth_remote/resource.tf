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

resource "opnsense_ipsec_auth_local" "example" {
  enabled          = "1"
  ipsec_connection = opnsense_ipsec_connection.example.id
  round            = "0"
  authentication   = "psk"
  auth_id          = "auth-mail@tld.com"
  eap_id           = ""
  certificates     = []
  public_keys      = []
  description      = "Example IPsec Auth Local"
}
