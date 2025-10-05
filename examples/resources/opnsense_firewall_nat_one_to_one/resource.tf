resource "opnsense_firewall_nat_one_to_one" "example_one" {
  external_net = "220.110.81.9/32"
  source = {
    net = "10.10.2.9/32"
  }
  description = "Example one"
}


resource "opnsense_firewall_nat_one_to_one" "example_two" {
  enabled = false
  log = true
  external_net = "220.110.81.9/32"
  type = "binat"
  
  source = {
    net = "10.10.2.9/32"
    invert = false
  }
  
  destination = {
    net = "any"
    invert = false
  }

  nat_reflection = "enable"
  categories = [ "8cb36e8e-1d72-480a-8268-bbdaf1ec6ed6" ]
  description = "Example two"
}

resource "opnsense_firewall_nat_one_to_one" "example_three" {
  enabled = true
  log = true
  external_net = "220.110.81.9/32"
  type = "nat"
  
  source = {
    net = "__lan_network" # aliases are only allowed in type nat rules
  }

  nat_reflection = "default"
  description = "Example tree"
}