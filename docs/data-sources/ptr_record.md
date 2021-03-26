---
page_title: "PTR Record Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a ptr record from infoblox
---

# Data Source `infoblox_ptr_record`

Retrieves details for a ptr record from infoblox

## Example Usage

```terraform
data "infoblox_ptr_record" "ptr_record" {
  pointer_domain_name = "example-hostname.example.com"
}
```

```terraform
data "infoblox_ptr_record" "ptr_record" {
  ref = "record:ptr/867530986753098675309867530986753098675309867530986753098675309:6.4.19.172.in-addr.arpa/default"
}
```

## Attributes Reference

The following attributes are exported.

- `comment` - (Computed, String) Comment for the fixed address; maximum 256 characters.
- `disable` - (Computed, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `dns_name` -  (MutuallyExclusiveGroup*/Computed, String) The name for a DNS PTR record in punycode format.
- `dns_pointer_domain_name` -  (MutuallyExclusiveGroup*/Computed, String) The domain name of the DNS PTR record in punycode format.
- `extensible_attributes` - (Computed, Map) Extensible attributes of ptr record (Values are JSON encoded).
- `ip_v4_address` -  (MutuallyExclusiveGroup*/Computed, String) The IPv4 Address of the record.
- `ip_v6_address` -  (MutuallyExclusiveGroup*/Computed, String) The IPv6 Address of the record.
- `name` -  (MutuallyExclusiveGroup*/Computed, String) The name of the DNS PTR record in FQDN format.
- `pointer_domain_name` -  (MutuallyExclusiveGroup*/Computed, String) The domain name of the DNS PTR record in FQDN format.
- `query_params` - (Optional, Map) Additional query parameters used for ptr record query (see infoblox documentation for full list)
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of ptr record object.
- `view` - (Optional/Computed, String) Name of the DNS View in which the record resides.
- `zone` - (Optional/Computed, String) The name of the zone in which the record resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key