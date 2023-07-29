---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opnsense_unbound_forward Data Source - terraform-provider-opnsense"
subcategory: ""
description: |-
  Query Forwarding section allows for entering arbitrary nameservers to forward queries to. Can forward queries normally, or over TLS.
---

# opnsense_unbound_forward (Data Source)

Query Forwarding section allows for entering arbitrary nameservers to forward queries to. Can forward queries normally, or over TLS.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) UUID of the resource.

### Read-Only

- `domain` (String) If a domain is entered here, queries for this specific domain will be forwarded to the specified server.
- `enabled` (Boolean) Whether this route is enabled.
- `server_ip` (String) IP address of DNS server to forward all requests.
- `server_port` (Number) Port of DNS server, for usual DNS use `53`, if you use DoT set it to `853`.
- `verify_cn` (String) The Common Name of the DNS server (e.g. `dns.example.com`). This field is required to verify its TLS certificate. DNS-over-TLS is susceptible to man-in-the-middle attacks unless certificates can be verified.

