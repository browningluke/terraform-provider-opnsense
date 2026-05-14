// Read the current Unbound DNS resolver settings from OPNsense.
data "opnsense_unbound_settings" "settings" {}

output "unbound_enabled" {
  value = data.opnsense_unbound_settings.settings.general.enabled
}

output "unbound_port" {
  value = data.opnsense_unbound_settings.settings.general.port
}
