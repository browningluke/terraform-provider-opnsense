resource "opnsense_auth_user" "alice" {
  name               = "alice"
  full_name          = "Alice Example"
  email              = "alice@example.com"
  shell              = "/bin/sh"
  scramble_password  = true
  groups             = ["admins"]
  privileges         = ["page-all"]
  authorized_keys = [
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIexample alice@example.com",
  ]
}
