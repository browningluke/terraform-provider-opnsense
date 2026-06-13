resource "opnsense_quagga_ospf6_prefix_list" "example" {
  name            = "OSPF6_PL"
  action          = "permit"
  sequence_number = "10"
  network         = "2001:db8::/32"
}
