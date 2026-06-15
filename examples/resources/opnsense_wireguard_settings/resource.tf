// Import the singleton resource before managing it:
// terraform import opnsense_wireguard_settings.settings wireguard_settings

resource "opnsense_wireguard_settings" "settings" {
  enabled = true
}
