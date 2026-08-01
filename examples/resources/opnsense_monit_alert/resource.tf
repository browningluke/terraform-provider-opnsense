resource "opnsense_monit_alert" "admin" {
  recipient   = "admin@example.com"
  description = "Notify admin on service failures"
  events      = ["action", "timeout"]
}
