// PXE lan
resource "opnsense_dnsmasq_boot" "pxe_lan" {
  interface = "lan"
  file_name = "pxelinux.0"
}

// PXE Boot
resource "opnsense_dnsmasq_boot" "pxe_install" {
  interface = "lan"

  tag = [
    "pxe",
    "install"
  ]

  file_name      = "ipxe/boot.ipxe"
  server_name    = "pxe-server"
  server_address = "192.168.110.10"

  description = "PXE boot for installation clients"
}

// PXE Provisioning
resource "opnsense_dnsmasq_boot" "provisioning" {
  interface = "kubernetes"

  tag = [
    "provisioning"
  ]

  file_name      = "ipxe/k8s.ipxe"
  server_address = "192.168.111.10"

  description = "iPXE boot for Kubernetes provisioning"
}

// PXE Boot Server
resource "opnsense_dnsmasq_boot" "pxe_server" {
  interface    = "lan"
  file_name    = "pxelinux.0"
  server_name  = "pxe.example.internal"
  description  = "PXE boot from internal boot server"
}