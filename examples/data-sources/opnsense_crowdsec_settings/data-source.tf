data "opnsense_crowdsec_settings" "current" {}

output "crowdsec_agent_enabled" {
  value = data.opnsense_crowdsec_settings.current.agent_enabled
}
