resource "opnsense_quagga_ospf_interface" "example" {
  interface_name = "wan"
  area           = "0.0.0.0"
}
