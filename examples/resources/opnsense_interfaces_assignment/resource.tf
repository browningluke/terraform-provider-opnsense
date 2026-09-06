# Descriptions may only contain letters, digits and underscores — OPNsense
# rejects spaces and other punctuation (model mask: /^[a-z_0-9]{0,255}$/i).

# Import the existing wan assignment before managing it:
#   terraform import opnsense_interfaces_assignment.wan wan
resource "opnsense_interfaces_assignment" "wan" {
  device      = "vtnet0"
  description = "WAN"
}

# Assigning a device to a new logical name creates an optional interface
# (opt1, opt2, ...); it does not create "wan" or "lan".
resource "opnsense_interfaces_assignment" "opt_example" {
  device      = "vtnet2"
  description = "opt_example"
}

# Assigning a VLAN
#
# The VLAN's *device* is what gets assigned, not the VLAN resource's id (which
# is a UUID). Set `device` explicitly on the VLAN so it can be referenced
# directly:
resource "opnsense_interfaces_vlan" "servers" {
  description = "Servers VLAN"
  tag         = 20
  priority    = 0
  parent      = opnsense_interfaces_assignment.wan.device
  device      = "vlan0.20"
}

resource "opnsense_interfaces_assignment" "servers" {
  device      = opnsense_interfaces_vlan.servers.device
  description = "servers"
}

# If you leave `device` unset on the VLAN, OPNsense generates the device name
# (vlan01, vlan02, ...) and it is not reflected back into the VLAN resource's
# state. In that case look it up by tag via the overview data source:
resource "opnsense_interfaces_vlan" "guests" {
  description = "Guests VLAN"
  tag         = 30
  priority    = 0
  parent      = opnsense_interfaces_assignment.wan.device
}

data "opnsense_interfaces_overview_all" "all" {
  depends_on = [opnsense_interfaces_vlan.guests]
}

locals {
  vlan_devices = {
    for i in data.opnsense_interfaces_overview_all.all.interfaces :
    i.vlan.tag => i.device
    if i.vlan.tag != ""
  }
}

resource "opnsense_interfaces_assignment" "guests" {
  device      = local.vlan_devices[tostring(opnsense_interfaces_vlan.guests.tag)]
  description = "guests"
}

# OPNsense's REST API exposes only the assignment, description and lock flag.
# Addressing (IPv4/IPv6, MTU, ...) still has to be configured by hand. A
# postcondition on the overview data source can block the apply until that is
# done:
data "opnsense_interfaces_overview" "servers" {
  device = opnsense_interfaces_vlan.servers.device

  depends_on = [opnsense_interfaces_assignment.servers]

  lifecycle {
    postcondition {
      condition = try(self.ipv4[0].ipaddr == "192.168.20.1", false)

      error_message = <<-EOT
        Manual configuration required.

        Configure OPNsense interface ${opnsense_interfaces_assignment.servers.id}:
          IPv4: "192.168.20.1/24"

        Run 'terraform apply' again after completing the manual configuration.
      EOT
    }
  }
}
