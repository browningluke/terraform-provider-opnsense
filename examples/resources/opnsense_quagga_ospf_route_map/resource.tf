resource "opnsense_quagga_ospf_route_map" "example" {
  name         = "OSPF_RM"
  action       = "permit"
  route_map_id = "10"
  set          = "metric 100"
}
