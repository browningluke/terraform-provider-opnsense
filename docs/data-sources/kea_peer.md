---
page_title: "opnsense_kea_peer Data Source - terraform-provider-opnsense"
subcategory: Kea
description: |-
  Configure HA Peers for Kea.
---

# opnsense_kea_peer (Data Source)

Configure HA Peers for Kea.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) UUID of the peer.

### Read-Only

- `name` (String) Peer name, there should be one entry matching this machine's "This server name".
- `role` (String) Peer's role.
- `url` (String) URL of the server instance.
