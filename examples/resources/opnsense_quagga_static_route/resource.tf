resource "opnsense_quagga_static_route" "example" {
  network = "10.100.0.0/24"
  gateway = "192.168.1.1"
}
