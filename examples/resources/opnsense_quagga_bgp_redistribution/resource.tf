resource "opnsense_quagga_bgp_redistribution" "example" {
  redistribute = "connected"
  description  = "Redistribute connected routes into BGP"
}
