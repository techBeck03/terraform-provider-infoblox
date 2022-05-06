---
page_title: "Range Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a range from infoblox
---

# Data Source `infoblox_range`

Retrieves details for a range from infoblox

## Example Usage

```terraform
data "infoblox_range" "range" {
  cidr          = "172.19.4.0/24"
  start_address = "172.19.4.2"
  end_address   = "172.19.4.10"
}
```

```terraform
data "infoblox_range" "range" {
  ref = "range/867530986753098675309867530986753098675309867530986753098675309:172.19.4.2/172.19.4.10/default/default"
}
```

## Attributes Reference

The following attributes are exported.

- `address_list` -  (Computed) The list of IP Addresses associated with this range
- `cidr` -  (MutuallyExclusiveGroup*/Computed, String) The network to which this range belongs, in IPv4 Address/CIDR format.
- `comment` - (Computed, String) Comment for the range; maximum 256 characters.
- `disable_dhcp` - (Computed, Bool) Disable for DHCP.
- `end_address` -  (Computed, String) The IPv4 Address end address of the range.
- `extensible_attributes` - (Computed, Map) Extensible attributes of ptr record (Values are JSON encoded).
- `member` - (Computed, Set of Objects) Grid members associated with range.  Attributes for each set item:
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
- `query_params` - (Optional, Map) Additional query parameters used for ptr record query (see infoblox documentation for full list)
- `range_function_string` -  (Computed, String) String representation of start and end addresses to be used with function calls. record in FQDN format.
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of range object.
- `start_address` -  (Computed, String) The IPv4 Address starting address of the range.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key