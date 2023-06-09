---
page_title: "opnsense_unbound_domain_override Resource - terraform-provider-opnsense"
subcategory: Unbound
description: |-
  Domain overrides can be used to forward queries for specific domains (and subsequent subdomains) to local or remote DNS servers.
---

# opnsense_unbound_domain_override (Resource)

Domain overrides can be used to forward queries for specific domains (and subsequent subdomains) to local or remote DNS servers.

## Example Usage

```terraform
// Enabled with description
resource "opnsense_unbound_domain_override" "one_override" {
  enabled = true
  description = "Example override"

  domain = "example.lan"
  server = "192.168.1.1"
}

// Disabled without description
resource "opnsense_unbound_domain_override" "two_override" {
  enabled = false

  domain = "example.arpa"
  server = "192.168.1.100"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) Domain to override (NOTE: this does not have to be a valid TLD!), e.g. `test` or `mycompany.localdomain` or `1.168.192.in-addr.arpa`.
- `server` (String) IP address of the authoritative DNS server for this domain, e.g. `192.168.100.100`. To use a nondefault port for communication, append an `@` with the port number.

### Optional

- `description` (String) Optional description here for your reference (not parsed).
- `enabled` (Boolean) Enable this domain override. Defaults to `true`.

### Read-Only

- `id` (String) UUID of the host override.

