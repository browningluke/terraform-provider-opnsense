// Generate random 256-bit base64 public key
resource "random_id" "pubkey" {
  byte_length = 32
}

// Generate random 256-bit base64 private key
resource "random_id" "privkey" {
  byte_length = 32
}

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

// Configure the server
resource "opnsense_wireguard_server" "example0" {
  name = "example0"

  private_key = random_id.privkey.b64_std
  public_key  = random_id.pubkey.b64_std

  dns = [
    "1.1.1.1",
    "8.8.8.8"
  ]

  tunnel_address = [
    "192.168.1.100/32",
    "10.10.0.0/24"
  ]

  peers = [
    opnsense_wireguard_client.example0.id
  ]
}
