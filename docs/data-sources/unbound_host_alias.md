---
page_title: "opnsense_unbound_host_alias Data Source - terraform-provider-opnsense"
subcategory: Unbound
description: |-
  Host aliases can be used to create alternative names for a Host.
---

# opnsense_unbound_host_alias (Data Source)

Host aliases can be used to create alternative names for a Host.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) UUID of the resource.

### Read-Only

- `description` (String) Optional description here for your reference (not parsed).
- `domain` (String) Domain of the host, e.g. example.com
- `enabled` (Boolean) Whether this route is enabled.
- `hostname` (String) Name of the host, without the domain part.
- `override` (String) The associated host override to apply this alias on.

