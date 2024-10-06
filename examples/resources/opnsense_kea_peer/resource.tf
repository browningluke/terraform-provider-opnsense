// Small example
resource "opnsense_kea_peer" "example" {
  name = "example"
  role = "primary"
  url = "http://192.0.2.1:8001/"
}