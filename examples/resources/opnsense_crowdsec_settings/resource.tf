// Import the singleton resource before managing it:
// terraform import opnsense_crowdsec_settings.settings crowdsec_settings

resource "opnsense_crowdsec_settings" "settings" {
  agent_enabled             = true
  lapi_enabled              = true
  firewall_bouncer_enabled  = true
  lapi_manual_configuration = false
  lapi_listen_address       = "127.0.0.1"
  lapi_listen_port          = "8080"
  rules_enabled             = true
  rules_log                 = false
  rules_tag                 = "crowdsec"
  crowdsec_firewall_verbose = false
}
