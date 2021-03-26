---
page_title: "Fixed Address Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for a fixed address in infoblox
---

# Resource `infoblox_fixed_address`

Manages configuration details for a fixed address in infoblox

## Example Usage


### Specify IP
```terraform
resource "infoblox_fixed_address" "fixed-addr" {
  ip_address        = "172.19.4.251"
  name              = "HSRP-A"
  comment           = "example fixed address"
  match_client      = "RESERVED"
  restart_if_needed = true
  grid_ref          = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  extensible_attributes = {
    Owner = jsonencode({
      value = "leroyjenkins",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation",
      type  = "STRING"
    })
  }
}
```

### next_available_ip from CIDR
```terraform
resource "infoblox_fixed_address" "fixed-addr" {
  cidr              = "172.19.4.0/24"
  name              = "HSRP-A"
  comment           = "example fixed address"
  match_client      = "RESERVED"
  restart_if_needed = true
  grid_ref          = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  extensible_attributes = {
    Owner = jsonencode({
      value = "leroyjenkins",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation",
      type  = "STRING"
    })
  }
}
```

### next_available_ip from range
```terraform
resource "infoblox_fixed_address" "fixed-addr" {
  range_function_string = "172.19.4.2-172.19.4.10"
  name                  = "HSRP-A"
  comment               = "example fixed address"
  match_client          = "RESERVED"
  restart_if_needed     = true
  grid_ref              = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  extensible_attributes = {
    Owner = jsonencode({
      value = "leroyjenkins",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation",
      type  = "STRING"
    })
  }
}
```

### Specify IP and MAC
```terraform
resource "infoblox_fixed_address" "fixed-addr" {
  ip_address        = "172.19.4.251"
  name              = "HSRP-A"
  comment           = "example fixed address"
  mac               = "12:34:56:78:9A:BC"
  match_client      = "MAC_ADDRESS"
  restart_if_needed = true
  grid_ref          = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  extensible_attributes = {
    Owner = jsonencode({
      value = "leroyjenkins",
      type  = "STRING"
    })
    Location = jsonencode({
      value = "CollegeStation",
      type  = "STRING"
    })
  }
}
```

## Attributes Reference

The following attributes are exported.

- `cidr` - (AtLeastOneOfGroup*/Computed, String) The network to which this fixed address belongs, in IPv4 Address/CIDR format.
- `comment` - (Optional, String) Comment for the fixed address; maximum 256 characters.
- `disable` - (Optional, Bool) Determines whether a fixed address is disabled or not. When this is set to False, the fixed address is enabled.
- `extensible_attributes` - (Optional, Map) JSON string of extensible attributes associated with fixed address.
- `grid_ref` -  (Optional, String) Ref for grid needed for restarting services.
- `hostname` -  (Optional, String) This field contains the name of this fixed address.
- `ip_address` -  (AtLeastOneOfGroup*/Computed, String) The IPv4 Address of the fixed address.
- `mac` -  (Optional, String) The MAC address value for this fixed address.
- `match_client` -  (Optional, String) The match_client value for this fixed address.
- `member` - (Optional, Set of `1` Object) Grid member associated with fixed address (required to restart services).  Attributes for each set item:
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
- `range_function_string` -  (AtLeastOneOfGroup*, String) Range start and end string for next_available_ip function calls.
- `restart_if_needed` -  (Optional, Bool) Restart dhcp services if needed.

**_AtLeastOneOfGroup_**: At least one of the attritbutes in this group **MUST** be provided to determine the IP address

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of fixed address object.