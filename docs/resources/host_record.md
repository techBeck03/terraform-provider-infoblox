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
  hostname   = "realhost.example.com"
  comment    = "example host record"
  enable_dns = true
  ip_v4_address {
    ip_address = "172.19.4.31"
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
```

### next_available_ip from range
```terraform
  hostname   = "realhost.example.com"
  comment    = "example host record"
  enable_dns = true
  range_function_string = "172.19.4.2-172.19.4.10"
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
```

### next_available_ip from network
```terraform
  hostname   = "realhost.example.com"
  comment    = "example host record"
  enable_dns = true
  network    = "172.19.4.0/24"
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
```


## Argument Reference

The following attributes are exported.

- `comment` - (Optional, String) Comment for the fixed address; maximum 256 characters.
- `enable_dns` - (Optional, Bool) When false, the host does not have parent zone information.
- `extensible_attributes` - (Optional, Map) Extensible attributes of host record (Values are JSON encoded).
- `hostname` -  (Required, String) The host name in FQDN format.
- `ip_v4_address` - (MutuallyExclusiveGroup*/Computed, Set of Objects) IPv4 addresses associated with host record.  Attributes for each set item:
  - `ref` - (Computed, String) Reference id of address object.
  - `ip_address` - (Required, String) IP address.
  - `hostname` - (Computed, String) Hostname associated with IP address.
  - `network` - (Optional, String) Network associated with IP address.
  - `mac_address` - (Optional, String) MAC address associated with IP address.
  - `configure_for_dhcp` - (Optional, Bool) Set this to True to enable the DHCP configuration for this host address.
- `network` - (AtLeastOneOfGroup*/Computed, String) The network to which this fixed address belongs, in IPv4 Address/CIDR format.
- `network_view` -  (Optional, String) The name of the network view in which this fixed address resides.
- `range_function_string` -  (AtLeastOneOfGroup*, String) Range start and end string for next_available_ip function calls.
- `view` - (Optional, String) The name of the DNS view in which the record resides.
- `zone` - (Computed, String) The name of the zone in which the record resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided to determine the IP address.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of host record object.