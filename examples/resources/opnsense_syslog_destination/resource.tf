// Forward firewall logs to a remote CrowdSec log collector
resource "opnsense_syslog_destination" "crowdsec" {
  hostname    = "192.0.2.10"
  port        = "514"
  transport   = "udp4"
  facility    = ["local0"]
  description = "CrowdSec log collector"
}
