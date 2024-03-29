---
page_title: "opnsense_quagga_bgp_communitylist Data Source - terraform-provider-opnsense"
subcategory: Quagga
description: |-
  Configure community lists for BGP.
---

# opnsense_quagga_bgp_communitylist (Data Source)

Configure community lists for BGP.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) UUID of the resource.

### Read-Only

- `action` (String) Set permit for match or deny to negate the rule.
- `community` (String) The community you want to match. You can also regex and it is not validated so please be careful.
- `description` (String) An optional description for this prefix list.
- `enabled` (Boolean) Enable this community list.
- `number` (Number) Set the number of your Community-List. 1-99 are standard lists while 100-500 are expanded lists.
- `seq_number` (Number) The ACL sequence number (10-99).

