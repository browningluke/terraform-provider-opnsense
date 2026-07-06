// Run firmware update check every day at 4am
resource "opnsense_cron_job" "firmware_check" {
  minutes = "0"
  hours   = "4"
  command = "firmware poll"
  description = "Daily firmware update check"
}

// Reload unbound every 6 hours, weekdays only
resource "opnsense_cron_job" "unbound_reload" {
  minutes  = "0"
  hours    = "*/6"
  weekdays = "1-5"
  command  = "unbound restart"
  description = "Periodic Unbound DNS reload"
}
