resource "opnsense_quagga_ospf6_redistribution" "example" {
  redistribute = "connected"
  description  = "Redistribute connected routes into OSPFv3"
}
