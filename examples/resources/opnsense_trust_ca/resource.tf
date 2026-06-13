# Create an internal root CA
resource "opnsense_trust_ca" "root" {
  description  = "Internal Root CA"
  action       = "internal"
  key_type     = "4096"
  digest       = "sha256"
  lifetime     = "3650"
  country      = "US"
  organization = "Acme Corp"
  common_name  = "Acme Internal Root CA"
}

# Create an intermediate CA signed by the root CA
resource "opnsense_trust_ca" "intermediate" {
  description  = "Internal Intermediate CA"
  action       = "internal"
  key_type     = "2048"
  digest       = "sha256"
  lifetime     = "1825"
  caref        = opnsense_trust_ca.root.ref_id
  country      = "US"
  organization = "Acme Corp"
  common_name  = "Acme Intermediate CA"
}
