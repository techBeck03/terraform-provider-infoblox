---
page_title: "Fixed Address Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a fixed address from infoblox
---

# Data Source `infoblox_fixed_address`

Retrieves details for a fixed address from infoblox

## Example Usage

```terraform
data "infoblox_alias_record" "alias_record" {
  hostname = "example-alias.example.com"
}
```

```terraform
data "infoblox_alias_record" "alias_record" {
  ref = "record:alias/867530986753098675309867530986753098675309867530986753098675309:example-alias.example.com/default"
}
```

## Attributes Reference

The following attributes are exported.

- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of host fixed address object.
- `hostname` -  (MutuallyExclusiveGroup*/Computed, String) TThis field contains the name of this fixed address.
- `ip_address` -  (MutuallyExclusiveGroup*/Computed, String) The IPv4 Address of the fixed address.
- `cidr` - (Computed, String) The network to which this fixed address belongs, in IPv4 Address/CIDR format.
- `comment` - (Computed, String) Comment for the fixed address; maximum 256 characters.
- `network_view` -  (Computed, String) The name of the network view in which this fixed address resides.
- `mac` -  (Computed, String) The MAC address value for this fixed address.
- `match_client` -  (Computed, String) The match_client value for this fixed address.
- `disable` - (Computed, Bool) Determines whether a fixed address is disabled or not. When this is set to False, the fixed address is enabled.
- `option` - (Computed, Set of Objects) An array of DHCP option structs that lists the DHCP options associated with the object.  Attributes for each set item:
  - `name` - (Computed, String) Name of the DHCP option.
  - `code` - (Computed, Int) The code of the DHCP option.
  - `use_option` - (Computed, Bool) Only applies to special options that are displayed separately from other options and have a use flag.
  - `value` - (Computed, String) Value of the DHCP option.
  - `vendor_class` - (Computed, String) The name of the space this DHCP option is associated to.
- `query_params` - (Optional, Map) Additional query parameters used for fixed adress query (see infoblox documentation for full list)
- `extensible_attributes` - (Computed, Map) JSON string of extensible attributes associated with fixed address

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key