# This is a singleton resource. It must be imported before use:
# terraform import opnsense_monit_settings.settings monit_settings

resource "opnsense_monit_settings" "settings" {
  enabled    = true
  interval   = "120"
  startdelay = "120"
  mailserver = "localhost"
  port       = "25"
}
