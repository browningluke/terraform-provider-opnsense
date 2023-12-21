// Get all interface configs
data "opnsense_interface_all" "all" {}

// Filter for a specific MAC address
output "specific_mac" {
  value = [for i in data.opnsense_interface_all.all.interfaces : i if i.macaddr == "5a:dc:1f:24:12:c6"]
}

// Filter for specific device (advisable to use opnsense_interface instead)
 output "wireguard" {
   value = [for i in data.opnsense_interface_all.all.interfaces : i if i.device == "wg1"]
 }
