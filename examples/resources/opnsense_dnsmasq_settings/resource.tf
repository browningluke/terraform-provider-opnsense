// Small example
// ```console
// % terraform import opnsense_dnsmasq_settings.small_settings dnsmasq_settings
// ```
// import {
//   to = opnsense_dnsmasq_settings.small_settings
//   id = "dnsmasq_settings"
// }
resource "opnsense_dnsmasq_settings" "small_settings" {
  enabled = true

  interface = [
    "lan"
  ]

  dhcp {
    authoritative = true
    fqdn          = true
    default_domain = "example.internal"
  }
}

// Full example
// ```console
// % terraform import opnsense_dnsmasq_settings.full_settings dnsmasq_settings
// ```
// import {
//   to = opnsense_dnsmasq_settings.full_settings
//   id = "dnsmasq_settings"
// }
resource "opnsense_dnsmasq_settings" "full_settings" {
  enabled = true

  interface = [
    "lan"
  ]

  strict_interface_binding = false

  dns {
    dnssec               = true
    cache_size           = 10000
    max_concurrent_queries = 150
    local_entry_ttl      = 60
  }

  dns_query_forwarding {
    query_sequentially          = false
    require_domain              = true
    do_not_forward_private_reverse = true
  }

  dhcp {
    authoritative        = true
    fqdn                 = true
    default_domain       = "example.internal"
    local_domain         = true
    max_leases           = 1000
    register_firewall_rules = true
    router_advertisements = true
    host_ping             = true
    log_dhcp              = false
    log_quiet             = false
  }
}


// Use Unbound together with dnsmasq_settings
// ```console
// % terraform import opnsense_dnsmasq_settings.local_dns dnsmasq_settings
// ```
// import {
//   to = opnsense_dnsmasq_settings.local_dns
//   id = "dnsmasq_settings"
// }
resource "opnsense_dnsmasq_settings" "local_dns" {
  enabled = true

  interface = [
    "lan"
  ]

  dns_query_forwarding {
    do_not_forward_system_dns   = true
    do_not_forward_private_reverse = true
  }

  dhcp {
    authoritative  = true
    fqdn           = true
    default_domain = "lan.example.internal"
    local_domain   = true
  }
}
