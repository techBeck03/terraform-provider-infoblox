---
page_title: "Host Record Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for a host record in infoblox
---

# Resource `infoblox_host_record`

Manages configuration details for a host record in infoblox

## Example Usage

### Specify IP address
```terraform
resource "infoblox_host_record" "static" {
  hostname   = "realhost.example.com"
  comment    = "example host record"
  enable_dns = true
  ip_v4_address {
    ip_address             = "172.19.4.31"
    use_for_ea_inheritance = true
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

### next_available_ip from range
```terraform
resource "infoblox_host_record" "from-range" {
  hostname   = "realhost.example.com"
  comment    = "example host record"
  enable_dns = true
  ip_v4_address {
    range_function_string  = "172.19.4.2-172.19.4.10"
    use_for_ea_inheritance = true
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

### next_available_ip from network
```terraform
resource "infoblox_host_record" "from-network" {
  hostname   = "realhost.example.com"
  comment    = "example host record"
  enable_dns = true
  ip_v4_address {
    network                = "172.19.4.0/24"
    use_for_ea_inheritance = true
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


## Argument Reference

The following attributes are exported.

- `comment` - (Optional, String) Comment for the fixed address; maximum 256 characters.
- `enable_dns` - (Optional, Bool) When false, the host does not have parent zone information.
- `extensible_attributes` - (Optional, Map) Extensible attributes of host record (Values are JSON encoded).
- `hostname` -  (Required, String) The host name in FQDN format.
- `ip_v4_address` - (Optional/Computed, Set of Objects) IPv4 addresses associated with host record.  Attributes for each set item:
  - `configure_for_dhcp` - (Optional, Bool) Set this to True to enable the DHCP configuration for this host address.
  - `hostname` - (Computed, String) Hostname associated with IP address.
  - `ip_address` - (MutuallyExclusiveGroup*/Computed, String) IP address.
  - `mac_address` - (Optional, String) MAC address associated with IP address.
  - `network` - (MutuallyExclusiveGroup*, String) The network to which this fixed address belongs, in IPv4 Address/CIDR format.
  - `range_function_string` -  (MutuallyExclusiveGroup*, String) Range start and end string for next_available_ip function calls.
  - `ref` - (Computed, String) Reference id of address object.
  - `use_for_ea_inheritance` - (Optional, Bool) Set this to True when using this host address for EA inheritance.
- `network_view` -  (Optional, String) The name of the network view in which this fixed address resides.
- `view` - (Optional, String) The name of the DNS view in which the record resides.
- `zone` - (Computed, String) The name of the zone in which the record resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided to determine the IP address.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of host record object.