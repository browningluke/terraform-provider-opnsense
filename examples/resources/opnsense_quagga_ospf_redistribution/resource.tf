resource "opnsense_quagga_ospf_redistribution" "example" {
  redistribute = "connected"
  description  = "Redistribute connected routes into OSPF"
}
