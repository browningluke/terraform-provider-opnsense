// opnsense_unbound_settings is a SINGLETON resource.
//
// There should only ever be ONE instance of this resource in your Terraform
// configuration. It manages the global Unbound DNS resolver settings for your
// OPNsense appliance.
//
// IMPORTANT: This resource cannot be created via `terraform apply`. You must
// import it first, then manage it going forward.
//
// Import using the fixed ID "unbound_settings":
//
//   terraform import opnsense_unbound_settings.settings unbound_settings
//
// After importing, running `terraform destroy` will only remove the resource
// from Terraform state — it will NOT reset or delete the upstream configuration.

// Import block (Terraform v1.5+)
import {
  to = opnsense_unbound_settings.settings
  id = "unbound_settings"
}

// Configure Unbound DNS resolver settings.
// All attributes are optional — omit any block to keep the upstream default.
resource "opnsense_unbound_settings" "settings" {
  general = {
    enabled = true
    port    = 53
  }

  advanced = {
    hide_identity = true
    hide_version  = true

    logging = {
      verbosity_level = 1
    }
  }

  acls = {
    default_action = "allow"
  }

  forwarding = {
    enabled = false
  }
}
