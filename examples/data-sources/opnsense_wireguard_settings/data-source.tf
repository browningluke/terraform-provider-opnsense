data "opnsense_wireguard_settings" "current" {}

output "wireguard_enabled" {
  value = data.opnsense_wireguard_settings.current.enabled
}
