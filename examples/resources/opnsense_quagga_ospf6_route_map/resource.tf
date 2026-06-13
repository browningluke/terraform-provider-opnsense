resource "opnsense_quagga_ospf6_route_map" "example" {
  name         = "OSPF6_RM"
  action       = "permit"
  route_map_id = "10"
  set          = "metric 100"
}
