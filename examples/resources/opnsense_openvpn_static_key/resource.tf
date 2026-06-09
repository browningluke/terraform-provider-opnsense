resource "opnsense_openvpn_static_key" "example" {
  description = "tls-crypt key for example server"
  mode        = "crypt"
  key         = <<-EOT
    #
    # 2048 bit OpenVPN static key
    #
    -----BEGIN OpenVPN Static key V1-----
    # ... key bytes (16 lines) ...
    -----END OpenVPN Static key V1-----
  EOT
}
