# This is a singleton resource. It must be imported before use:
# terraform import opnsense_mdnsrepeater_settings.settings mdnsrepeater_settings

resource "opnsense_mdnsrepeater_settings" "settings" {
  enabled    = true
  interfaces = ["lan", "opt1"]
  blocklist  = ["192.168.1.0/24"]
}
