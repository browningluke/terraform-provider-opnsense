resource "opnsense_monit_service" "unbound" {
  name        = "unbound"
  description = "DNS resolver"
  type        = "process"
  pidfile     = "/var/run/unbound.pid"
  start       = "/usr/local/etc/rc.d/unbound start"
  stop        = "/usr/local/etc/rc.d/unbound stop"
}
