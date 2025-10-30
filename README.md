# Terraform Provider for OPNsense

This Terraform provider enables management of various configs and resources within OPNsense®.

> [!WARNING]
> This provider is under active development and makes no guarantee of stability. Breaking changes to resource and data source schemas will occur as needed until v1.0. **It is not recommended to use this provider in production environments.**


- [Example Usage](#example-usage)
- [Long Road to v1](#long-road-to-v1)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [Current API Coverage](#current-api-coverage)
  - [Core API](#core-api)
  - [Plugin API](#plugin-api)
- [License](#license)

## Example Usage

```hcl
# 1. Configure Terraform to use the provider
terraform {
  required_providers {
    opnsense = {
      source  = "browningluke/opnsense"
      version = "~> 0.16"
    }
  }
}

# 2. Configure the OPNsense provider with API credentials
provider "opnsense" {
  uri        = "https://opnsense.example.com"

  # Either reference the API credentials literally
  api_key    = "<api key>"
  api_secret = "<api password>"

  # Or specify them with environment variables
  # export OPNSENSE_API_KEY="<api key>"
  # export OPNSENSE_API_SECRET="<api key>"
}

# 3. Create resources - example: firewall rule
resource "opnsense_firewall_filter" "allow_https" {
  enabled     = true
  description = "Allow inbound HTTPS traffic"

  interface = {
    interface = ["wan"]
  }

  filter = {
    action    = "pass"
    direction = "in"
    protocol  = "TCP"

    source = {
      net = "any"
    }

    destination = {
      net  = "192.168.1.100"
      port = "https"
    }

    log = true
  }
}
```

## Long Road to v1

Version 1.0 will be released once the provider achieves feature-parity with the **Core** OPNsense API and all resources have comprehensive acceptance tests (see [Current API Coverage](#current-api-coverage)). Plugin resources will be added as requested (at a lower priority than requests for Core resources). There is no Plugin API converage requirement for v1.

v1 represents the first release where resource and data source schemas will be guaranteed to be stable, and breaking changes to these schemas will be forbidden. Any updates to these schemas will following appropriate SemVer conventions. Until v1.0 is reached, **schemas are subject to change as needed** to improve usability and align with best practices. Users should always check the release notes when upgrading between pre-v1.0 versions to understand any breaking changes that may affect their configurations.

## Documentation

- **[Terraform Registry Documentation](https://registry.terraform.io/providers/browningluke/opnsense/latest/docs)** - Full resource and data source reference
- **[Examples](./examples/)** - Working examples for all resources

## Contributing

Interested in contributing? Please see our [Contributing Guide](./CONTRIBUTING.md) for development setup, testing requirements, and guidelines.

## Current API Coverage

This provider is actively expanding to cover the OPNsense API. The tables below contain the current status of said coverage.

- ✅ = Fully implemented
- 🚧 = Missing acceptance tests
- ❌ = Not implemented

### Core API

| Module/Controller/Resource       | Resource | Data Source |
|----------------------------------|----------|-------------|
| `Auth/Group`                     | ❌        | ❌           |
| `Auth/Priv`                      | ❌        | ❌           |
| `Auth/User`                      | ❌        | ❌           |
| `Captiveportal/Settings`         | ❌        | ❌           |
| `Captiveportal/Service/Template` | ❌        | ❌           |
| `Captiveportal/Settings/Zone`    | ❌        | ❌           |
| `Core/Hasync`                    | ❌        | ❌           |
| `Core/Snapshots`                 | ❌        | ❌           |
| `Core/Tunables`                  | ❌        | ❌           |
| `Cron/Job`                       | ❌        | ❌           |
| `Dhcrelay/Settings`              | ❌        | ❌           |
| `Dhcrelay/Settings/Dest`         | ❌        | ❌           |
| `Dhcrelay/Settings/Relay`        | ❌        | ❌           |
| `Diagnostics/Interface`          | ❌        | 🚧          |
| `Diagnostics/Lvtemplate`         | ❌        | ❌           |
| `Diagnostics/Lvtemplate/Item`    | ❌        | ❌           |
| `Dnsmasq/Settings`               | ❌        | ❌           |
| `Dnsmasq/Settings/Boot`          | ❌        | ❌           |
| `Dnsmasq/Settings/Domain`        | ❌        | ❌           |
| `Dnsmasq/Settings/Host`          | ❌        | ❌           |
| `Dnsmasq/Settings/Option`        | ❌        | ❌           |
| `Dnsmasq/Settings/Range`         | ❌        | ❌           |
| `Dnsmasq/Settings/Tag`           | ❌        | ❌           |
| `Firewall/Alias`                 | ✅        | ✅           |
| `Firewall/Category`              | 🚧       | 🚧          |
| `Firewall/Filter`                | ✅        | ✅           |
| `Firewall/Group`                 | ❌        | ❌           |
| `Firewall/NPTv6`                 | ❌        | ❌           |
| `Firewall/Source NAT`            | 🚧       | 🚧          |
| `Firewall/One-to-One NAT`        | ✅        | ✅           |
| `Interfaces/Bridge`              | ❌        | ❌           |
| `Interfaces/Gif`                 | ❌        | ❌           |
| `Interfaces/Gre`                 | ❌        | ❌           |
| `Interfaces/Lagg`                | ❌        | ❌           |
| `Interfaces/Loopback`            | ❌        | ❌           |
| `Interfaces/Neighbor`            | ❌        | ❌           |
| `Interfaces/Overview`            |          | ❌           |
| `Interfaces/Vip`                 | ✅        | ✅           |
| `Interfaces/Vlan`                | ✅        | ✅           |
| `Interfaces/Vxlan`               | ❌        | ❌           |
| `Ipsec/Settings`                 | ❌        | ❌           |
| `Ipsec/Connections/Local`        | ✅        | ❌           |
| `Ipsec/Connections/Remote`       | ✅        | ❌           |
| `Ipsec/Connections/Child`        | ✅        | ❌           |
| `Ipsec/Connections/Connection`   | ✅        | ❌           |
| `Ipsec/KeyPairs`                 | ❌        | ❌           |
| `Ipsec/ManualSpd`                | ❌        | ❌           |
| `Ipsec/Pools`                    | ❌        | ❌           |
| `Ipsec/Psk`                      | ✅        | ❌           |
| `Ipsec/Vti`                      | ✅        | ❌           |
| `Kea/CtrlAgent`                  | ❌        | ❌           |
| `Kea/Dhcpv4/Peer`                | 🚧       | 🚧          |
| `Kea/Dhcpv4/Reservation`         | 🚧       | 🚧          |
| `Kea/Dhcpv4/Subnet`              | 🚧       | 🚧          |
| `Kea/Dhcpv6/PD Pool`             | ❌        | ❌           |
| `Kea/Dhcpv6/Peer`                | ❌        | ❌           |
| `Kea/Dhcpv6/Reservation`         | ❌        | ❌           |
| `Kea/Dhcpv6/Subnet`              | ❌        | ❌           |
| `Monit/Settings`                 | ❌        | ❌           |
| `Monit/Settings/Alert`           | ❌        | ❌           |
| `Monit/Settings/Service`         | ❌        | ❌           |
| `Monit/Settings/Test`            | ❌        | ❌           |
| `Openvpn/Client Overwrites`      | ❌        | ❌           |
| `Openvpn/Instances`              | ❌        | ❌           |
| `Openvpn/Instances/Static Key`   | ❌        | ❌           |
| `Openvpn/Instances/Generate Key` | ❌        |             |
| `Routes/Route`                   | 🚧       | 🚧          |
| `Routing/Gateway`                | ❌        | ❌           |
| `Syslog/Settings`                | ❌        | ❌           |
| `Syslog/Settings/Destination`    | ❌        | ❌           |
| `Trafficshaper/Pipe`             | ❌        | ❌           |
| `Trafficshaper/Queue`            | ❌        | ❌           |
| `Trafficshaper/Rule`             | ❌        | ❌           |
| `Trust/Settings`                 | ❌        | ❌           |
| `Trust/CA`                       | ❌        | ❌           |
| `Trust/Cert`                     | ❌        | ❌           |
| `Unbound/Settings`               | ❌        | ❌           |
| `Unbound/Settings/Forward`       | 🚧       | 🚧          |
| `Unbound/Settings/Host Alias`    | 🚧       | 🚧          |
| `Unbound/Settings/Host Override` | 🚧       | 🚧          |
| `Unbound/Settings/ACL`           | ❌        | ❌           |
| `Wireguard/Settings`             | ❌        | ❌           |
| `Wireguard/Client`               | 🚧       | 🚧          |
| `Wireguard/Server`               | 🚧       | 🚧          |
| `Wireguard/Generate Key Pair`    | ❌        | ❌           |
| `Wireguard/Generate PSK`         | ❌        | ❌           |

### Plugin API

The following is a non-exhaustive list of the plugin APIs OPNsense supports. The table shows those which are 'highest priority'. Please open a feature request to indicate interest for any plugin not listed here.

| Plugin/Controller/Resource     | Resource | Data Source |
|--------------------------------|----------|-------------|
| `Acmeclient/Settings`          | ❌        | ❌           |
| `Acmeclient/Account`           | ❌        | ❌           |
| `Acmeclient/Validation`        | ❌        | ❌           |
| `Acmeclient/Certificates`      | ❌        | ❌           |
| `Acmeclient/Action`            | ❌        | ❌           |
| `Haproxy/Maintenance`          | ❌        | ❌           |
| `Haproxy/Settings`             | ❌        | ❌           |
| `Haproxy/Settings/Acl`         | ❌        | ❌           |
| `Haproxy/Settings/Action`      | ❌        | ❌           |
| `Haproxy/Settings/Backend`     | ❌        | ❌           |
| `Haproxy/Settings/Cpu`         | ❌        | ❌           |
| `Haproxy/Settings/Errorfile`   | ❌        | ❌           |
| `Haproxy/Settings/Fcgi`        | ❌        | ❌           |
| `Haproxy/Settings/Frontend`    | ❌        | ❌           |
| `Haproxy/Settings/Group`       | ❌        | ❌           |
| `Haproxy/Settings/Healthcheck` | ❌        | ❌           |
| `Haproxy/Settings/Lua`         | ❌        | ❌           |
| `Haproxy/Settings/Mapfile`     | ❌        | ❌           |
| `Haproxy/Settings/Server`      | ❌        | ❌           |
| `Haproxy/Settings/User`        | ❌        | ❌           |
| `Quagga/General`               | ❌        | ❌           |
| `Quagga/Bfd`                   | ❌        | ❌           |
| `Quagga/Bfd/Neighbor`          | ❌        | ❌           |
| `Quagga/Bgp`                   | ❌        | ❌           |
| `Quagga/Bgp/AS Path`           | 🚧       | 🚧          |
| `Quagga/Bgp/Community List`    | 🚧       | 🚧          |
| `Quagga/Bgp/Neighbor`          | 🚧       | 🚧          |
| `Quagga/Bgp/Peer Group`        | ❌        | ❌           |
| `Quagga/Bgp/Prefix List`       | 🚧       | 🚧          |
| `Quagga/Bgp/Route Map`         | 🚧       | 🚧          |
| `Quagga/Ospf`                  | ❌        | ❌           |
| `Quagga/Ospf/Interface`        | ❌        | ❌           |
| `Quagga/Ospf/Neighbor`         | ❌        | ❌           |
| `Quagga/Ospf/Network`          | ❌        | ❌           |
| `Quagga/Ospf/Prefix List`      | ❌        | ❌           |
| `Quagga/Ospf/Redistribution`   | ❌        | ❌           |
| `Quagga/Ospf/Route Map`        | ❌        | ❌           |
| `Quagga/Ospf6`                 | ❌        | ❌           |
| `Quagga/Ospf6/Interface`       | ❌        | ❌           |
| `Quagga/Ospf6/Neighbor`        | ❌        | ❌           |
| `Quagga/Ospf6/Network`         | ❌        | ❌           |
| `Quagga/Ospf6/Prefix List`     | ❌        | ❌           |
| `Quagga/Ospf6/Redistribution`  | ❌        | ❌           |
| `Quagga/Rip`                   | ❌        | ❌           |
| `Quagga/Static`                | ❌        | ❌           |
| `Quagga/Static/Route`          | ❌        | ❌           |

The complete OPNsense API documentation can be found at: [docs.opnsense.org](https://docs.opnsense.org/development/api.html)

## License

This project is licensed under the Mozilla Public License v2.0 - see the [LICENSE](./LICENSE) file for details.
