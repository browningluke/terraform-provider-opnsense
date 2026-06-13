# This is a singleton resource. Import it before managing:
# terraform import opnsense_quagga_general.general quagga_general

resource "opnsense_quagga_general" "general" {
  enabled      = true
  profile      = "traditional"
  enable_carp  = false
  syslog_level = "notifications"
  fw_rules     = true
}
