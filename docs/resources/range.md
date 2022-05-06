---
page_title: "Range Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for a range in infoblox
---

# Resource `infoblox_range`

Manages configuration details for a range in infoblox

## Example Usage

### Next Available Sequential Range

The configuration below will find the next 10 available sequential IP addresses in the specified network and put them in a new range.

```terraform
resource "infoblox_range" "sequential_range" {
  cidr             = "172.19.4.0/24"
  comment          = "test range"
  sequential_count = 10
  disable_dhcp     = false
  restart_if_needed = true
  grid_ref          = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  extensible_attributes = {
    Owner = jsonencode({
      value = "leeroyjenkins",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation",
      type  = "STRING"
    })
  }
}
```

### Specified Range
```terraform
resource "infoblox_range" "specified_range" {
  cidr             = infoblox_network.net-1.cidr
  comment          = "test range"
  disable_dhcp     = false
  start_address    = "172.19.4.2"
  end_address      = "172.19.4.11"
  restart_if_needed = true
  grid_ref          = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  extensible_attributes = {
    Owner = jsonencode({
      value = "robbeck",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "Austin",
      type  = "STRING"
    })
  }
}
```

## Argument Reference

The following attributes are exported.

- `address_list` -  (Computed) The list of IP Addresses associated with this range
- `cidr` -  (Required, String) The network to which this range belongs, in IPv4 Address/CIDR format.
- `comment` - (Optional, String) Comment for the range; maximum 256 characters.
- `disable_dhcp` - (Optional, Bool) Disable for DHCP.
- `end_address` -  (MutuallyExclusiveGroup*/Computed, String) The IPv4 Address end address of the range.
- `extensible_attributes` - (Optional, Map) Extensible attributes of ptr record (Values are JSON encoded).
- `grid_ref` -  (Optional, String) Ref for grid needed for restarting services.
- `member` - (Optional, Set of `1` Object) Grid member associated with range (required to restart services).  Attributes for each set item:
  - `struct` - (Optional, String) Struct type of member (default = `dhcpmember`).
  - `ip_v4_address` - (Optional, String) IPv4 address.
  - `ip_v6_address` - (Optional, String) IPv6 address.
  - `hostname` - (Required, String) Hostname of member.
- `network_view` -  (Optional, String) The name of the network view in which this fixed address resides.
- `option` - (Optional, Set of Objects) An array of DHCP option structs that lists the DHCP options associated with the object.  Attributes for each set item:
  - `name` - (Required, String) Name of the DHCP option.
  - `code` - (Optional, Int) The code of the DHCP option.
  - `use_option` - (Optional, Bool) Only applies to special options that are displayed separately from other options and have a use flag (Default = `true`).
  - `value` - (Required, String) Value of the DHCP option.
  - `vendor_class` - (Optional, String) The name of the space this DHCP option is associated to.
- `range_function_string` -  (Computed, String) String representation of start and end addresses to be used with function calls. record in FQDN format.
- `restart_if_needed` -  (Optional, Bool) Restart dhcp services if needed.
- `sequential_count` - (MutuallyExclusiveGroup*/Computed, Int) Sequential count of addresses.
- `start_address` -  (MutuallyExclusiveGroup*/Computed, String) The IPv4 Address starting address of the range.

**_MutuallyExclusiveGroup_**: Either `sequential_count` OR `start_address` AND `end_address` must be specified

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of range object.