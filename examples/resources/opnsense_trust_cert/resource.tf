resource "opnsense_trust_ca" "example" {
  description = "Example CA"
  action      = "internal"
  key_type    = "2048"
  digest      = "sha256"
  lifetime    = "3650"
  country     = "US"
  common_name = "example-ca.internal"
}

# Server certificate signed by the CA
resource "opnsense_trust_cert" "web" {
  description  = "Web Server Certificate"
  action       = "internal"
  caref        = opnsense_trust_ca.example.ref_id
  key_type     = "2048"
  digest       = "sha256"
  cert_type    = "server_cert"
  lifetime     = "397"
  country      = "US"
  organization = "Acme Corp"
  common_name  = "www.example.com"
  altnames_dns = "www.example.com\napi.example.com"
}
