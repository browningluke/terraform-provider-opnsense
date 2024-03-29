---
page_title: "opnsense_unbound_host_override Data Source - terraform-provider-opnsense"
subcategory: Unbound
description: |-
  Host overrides can be used to change DNS results from client queries or to add custom DNS records.
---

# opnsense_unbound_host_override (Data Source)

Host overrides can be used to change DNS results from client queries or to add custom DNS records.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) UUID of the resource.

### Read-Only

- `description` (String) Optional description here for your reference (not parsed).
- `domain` (String) Domain of the host, e.g. example.com
- `enabled` (Boolean) Whether this route is enabled.
- `hostname` (String) Name of the host, without the domain part. Use `*` to create a wildcard entry.
- `mx_host` (String) Host name of MX host, e.g. mail.example.com.
- `mx_priority` (Number) Priority of MX record, e.g. 10.
- `server` (String) IP address of the host, e.g. 192.168.100.100 or fd00:abcd::1.
- `type` (String) Type of resource record. Available values: `A`, `AAAA`, `MX`.

