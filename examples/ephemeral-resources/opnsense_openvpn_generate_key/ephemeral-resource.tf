// Generate a fresh tls-crypt key on every apply.
// The result is exposed only at apply time and is never written to Terraform
// state.
ephemeral "opnsense_openvpn_generate_key" "tls_crypt" {
  key_type = "tls-crypt"
}

// Ephemeral values can only be consumed by other ephemeral blocks, by
// provider configuration, or by `ephemeral = true` outputs / write-only
// attributes — never by a regular resource argument. The output below is
// the typical pattern when piping the key into a downstream provider such
// as a secrets manager.
output "tls_crypt_key" {
  value     = ephemeral.opnsense_openvpn_generate_key.tls_crypt.key
  ephemeral = true
  sensitive = true
}
