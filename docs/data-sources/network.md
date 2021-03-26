---
page_title: "Network Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a network from infoblox
---

# Data Source `infoblox_network`

Retrieves details for a network from infoblox

## Example Usage

```terraform
data "infoblox_network" "network" {
  cidr = "172.19.4.0/24"
}
```

```terraform
data "infoblox_network" "network" {
  ref = "record:network/867530986753098675309867530986753098675309867530986753098675309:172.19.4.0/24/default"
}
```

## Attributes Reference

The following attributes are exported.

- `cidr` -  (MutuallyExclusiveGroup*/Computed, String) The network address in IPv4 Address/CIDR format.
- `comment` - (Computed, String) Comment for the record; maximum 256 characters.
- `disable_dhcp` - (Computed, Bool) Disable for DHCP.
- `extensible_attributes` - (Computed, Map) Extensible attributes of network (Values are JSON encoded).
- `member` - (Computed, Set of Objects) Grid members associated with network.  Attributes for each set item:
  - `struct` - (Computed, String) Struct type of member.
  - `ip_v4_address` - (Computed, String) IPv4 address.
  - `ip_v6_address` - (Computed, String) IPv6 address.
  - `hostname` - (Computed, String) Hostname of member.
- `network_view` -  (Computed, String) The name of the network view in which this fixed address resides.
- `option` - (Computed, Set of Objects) An array of DHCP option structs that lists the DHCP options associated with the object.  Attributes for each set item:
  - `name` - (Computed, String) Name of the DHCP option.
  - `code` - (Computed, Int) The code of the DHCP option.
  - `use_option` - (Computed, Bool) Only applies to special options that are displayed separately from other options and have a use flag.
  - `value` - (Computed, String) Value of the DHCP option.
  - `vendor_class` - (Computed, String) The name of the space this DHCP option is associated to.
- `query_params` - (Optional, Map) Additional query parameters used for network query (see infoblox documentation for full list)
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of network object.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key