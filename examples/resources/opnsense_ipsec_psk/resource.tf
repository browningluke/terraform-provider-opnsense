// Small example
resource "opnsense_ipsec_psk" "example" {
  identity_local  = "psk-mail@tld.com"
  identity_remote = "1.2.3.4"
  type            = "PSK"
  pre_shared_key  = "someSuperSecretKey"
  description     = "Example IPsec PSK"
}
