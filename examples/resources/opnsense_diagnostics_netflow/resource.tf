// opnsense_diagnostics_netflow is a SINGLETON resource.
//
// There should only ever be ONE instance of this resource in your Terraform
// configuration. It manages the global Netflow diagnostics configuration for
// your OPNsense appliance.
//
// IMPORTANT: This resource cannot be created via `terraform apply`. You must
// import it first, then manage it going forward.
//
// Import using the fixed ID "netflow_settings":
//
//   terraform import opnsense_diagnostics_netflow.netflow netflow_settings
//
// After importing, running `terraform destroy` will only remove the resource
// from Terraform state — it will NOT reset or delete the upstream configuration.

// Import block (Terraform v1.5+)
import {
  to = opnsense_diagnostics_netflow.netflow
  id = "netflow_settings"
}

resource "opnsense_diagnostics_netflow" "netflow" {
  enabled                = true
  interfaces             = ["wan"]
  egress_only_interfaces = []
  version                = "v9"
  targets                = ["10.0.0.10:2055"]
  active_timeout         = 1800
  inactive_timeout       = 15
}
