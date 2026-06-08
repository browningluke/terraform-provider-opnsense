resource "opnsense_firewall_nat_port_forward" "wan_https_k3s_ingress" {
  enabled     = true
  sequence    = 100
  interface   = ["wan"]
  ip_protocol = "inet"
  protocol    = "tcp"

  source = {
    net = "any"
  }

  destination = {
    net  = "wanip"
    port = "443"
  }

  target = {
    ip   = "10.1.1.20"
    port = "443"
  }

  description = "WAN HTTPS to k3s Traefik ingress VIP"
}
