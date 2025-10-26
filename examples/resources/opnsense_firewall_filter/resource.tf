# Define a category for use in filter examples
resource "opnsense_firewall_category" "example" {
  name = "Example Category"
}

# Example 1: Minimal firewall filter with only required attributes
resource "opnsense_firewall_filter" "minimal" {
  interface = {
    interface = ["lan"]
  }

  filter = {
    action    = "pass"
    direction = "in"
    protocol  = "any"
  }
}

# Example 2: Typical use case: Allow HTTPS traffic from specific source to web server
resource "opnsense_firewall_filter" "allow_https" {
  enabled     = true
  sequence    = 100
  description = "Allow HTTPS traffic to web server"

  categories = [
    opnsense_firewall_category.example.id,
  ]

  interface = {
    interface = ["wan"]
  }

  filter = {
    quick     = true
    action    = "pass"
    direction = "in"
    protocol  = "TCP"

    source = {
      net = "any"
    }

    destination = {
      net  = "192.168.1.100"
      port = "https"
    }

    log = true
  }

  stateful_firewall = {
    type = "keep"
  }
}

# Example 3: Comprehensive example with all attributes explicitly configured
resource "opnsense_firewall_filter" "comprehensive" {
  enabled        = true
  sequence       = 200
  no_xmlrpc_sync = false
  description    = "Comprehensive example with all attributes"

  categories = [
    opnsense_firewall_category.example.id,
  ]

  interface = {
    invert = false
    interface = [
      "lan",
      "opt1",
    ]
  }

  filter = {
    quick         = true
    action        = "pass"
    allow_options = false
    direction     = "in"
    ip_protocol   = "inet"
    protocol      = "TCP"

    icmp_type = [
      "echoreq",
      "echorep",
    ]

    source = {
      invert = false
      net    = "192.168.1.0/24"
      port   = "1024-65535"
    }

    destination = {
      invert = false
      net    = "10.0.0.10"
      port   = "https"
    }

    log = true

    tcp_flags = [
      "syn",
      "ack",
    ]

    tcp_flags_out_of = [
      "syn",
      "ack",
      "fin",
      "rst",
    ]

    schedule = ""
  }

  stateful_firewall = {
    type    = "keep"
    policy  = ""
    timeout = 3600

    adaptive_timeouts = {
      start = 60000
      end   = 120000
    }

    max = {
      states             = 10000
      source_nodes       = 1000
      source_states      = 500
      source_connections = 100

      new_connections = {
        count   = 10
        seconds = 10
      }
    }

    overload_table = ""
    no_pfsync      = false
  }

  traffic_shaping = {
    shaper         = ""
    reverse_shaper = ""
  }

  source_routing = {
    gateway          = ""
    disable_reply_to = false
    reply_to         = ""
  }

  priority = {
    match         = -1
    set           = 3
    low_delay_set = 4
    match_tos     = "lowdelay"
  }

  internal_tagging = {
    set_local   = "webtraffic"
    match_local = ""
  }
}
