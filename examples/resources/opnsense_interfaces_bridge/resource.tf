// Bridge two assigned interfaces, letting OPNsense generate the device name
resource "opnsense_interfaces_bridge" "bridge" {
  description = "Example bridge"
  members     = ["opt1", "opt2"]
}

// Bridge with spanning tree enabled
resource "opnsense_interfaces_bridge" "bridge_stp" {
  description    = "Example bridge with STP"
  members        = ["opt3", "opt4"]
  enable_stp     = true
  stp_proto      = "rstp"
  stp_interfaces = ["opt3", "opt4"]
}
