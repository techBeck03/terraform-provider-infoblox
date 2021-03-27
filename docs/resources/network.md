---
page_title: "Network Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for a network in infoblox
---

# Resource `infoblox_network`

Manages configuration details for a network in infoblox

## Example Usage

```terraform
resource "infoblox_network" "net" {
  cidr       = "172.19.4.0/24"
  comment    = "example network"
  network_view      = "default"
  restart_if_needed = true
  grid_ref          = data.infoblox_grid.grid.ref
  member {
    hostname = data.infoblox_grid_member.member.hostname
  }
  option {
    code  = 3
    name  = "routers"
    value = "172.19.4.1"
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

## Argument Reference

The following attributes are exported.

- `cidr` -  (Required, String) The network address in IPv4 Address/CIDR format.
- `comment` - (Optional, String) Comment for the record; maximum 256 characters.
- `disable_dhcp` - (Optional, Bool) Disable for DHCP.
- `extensible_attributes` - (Optional, Map) Extensible attributes of network (Values are JSON encoded).
- `grid_ref` -  (Optional, String) Ref for grid needed for restarting services.
- `member` - (Optional, Set of `1` Object) Grid member associated with network (required to restart services).  Attributes for each set item:
  - `struct` - (Optional, String) Struct type of member (default = `dhcpmember`).
  - `ip_v4_address` - (Optional, String) IPv4 address.
  - `ip_v6_address` - (Optional, String) IPv6 address.
  - `hostname` - (Required, String) Hostname of member.
- `network_view` -  (Computed, String) The name of the network view in which this fixed address resides.
- `option` - (Optional, Set of Objects) An array of DHCP option structs that lists the DHCP options associated with the object.  Attributes for each set item:
  - `name` - (Required, String) Name of the DHCP option.
  - `code` - (Optional, Int) The code of the DHCP option.
  - `use_option` - (Optional, Bool) Only applies to special options that are displayed separately from other options and have a use flag (Default = `true`).
  - `value` - (Required, String) Value of the DHCP option.
  - `vendor_class` - (Optional, String) The name of the space this DHCP option is associated to.
- `restart_if_needed` -  (Optional, Bool) Restart dhcp services if needed.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of network object.