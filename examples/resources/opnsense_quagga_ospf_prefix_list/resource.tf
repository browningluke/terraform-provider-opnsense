resource "opnsense_quagga_ospf_prefix_list" "example" {
  name            = "OSPF_PL"
  action          = "permit"
  sequence_number = "10"
  network         = "10.0.0.0/8"
}
