// Configure a peer
resource "opnsense_wireguard_client" "example0" {
  enabled = false
  name = "example0"

  public_key = "/CPjuEdvHJulOIQ56TNyeNHkDJmRCMor4U9k68vMyac="
  psk        = "CJG05xgaLA8RiisoCAmp2U0v329LsIdK1GW4EMc9fmU="

  tunnel_address = [
    "192.168.1.1/32",
    "192.168.4.1/24",
  ]

  server_address = "10.10.10.10"
  server_port    = "1234"
}
