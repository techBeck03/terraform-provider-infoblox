---
page_title: "PTR Record Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for a ptr record in infoblox
---

# Resource `infoblox_ptr_record`

Manages configuration details for a ptr record in infoblox

## Example Usage

```terraform
resource "infoblox_ptr_record" "ptr" {
  pointer_domain_name = "realhost.example.com"
  ip_v4_address       = "172.19.4.6"
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
- `disable` - (Optional, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `dns_name` -  (Computed, String) The name for a DNS PTR record in punycode format.
- `dns_pointer_domain_name` -  (Computed, String) The domain name of the DNS PTR record in punycode format.
- `extensible_attributes` - (Optional, Map) Extensible attributes of ptr record (Values are JSON encoded).
- `ip_v4_address` -  (MutuallyExclusiveGroup1*, String) The IPv4 Address of the record.
- `ip_v6_address` -  (MutuallyExclusiveGroup1*, String) The IPv6 Address of the record.
- `name` -  (MutuallyExclusiveGroup2*, String) The name of the DNS PTR record in FQDN format.
- `pointer_domain_name` -  (MutuallyExclusiveGroup2*, String) The domain name of the DNS PTR record in FQDN format.
- `view` - (Optional/Computed, String) Name of the DNS View in which the record resides.
- `zone` - (Computed, String) The name of the zone in which the record resides.

**_MutuallyExclusiveGroup1_**: One and only one of the attritbutes in this group **MUST** be provided to determine IP address

**_MutuallyExclusiveGroup2_**: One and only one of the attritbutes in this group **MUST** be provided to determine mapping zone

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of ptr record object.