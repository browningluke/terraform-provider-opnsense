resource "opnsense_quagga_bfd_neighbor" "example" {
  address     = "192.0.2.1"
  description = "BFD neighbor to router"
}
