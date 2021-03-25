---
page_title: "Cname Record Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a cname record from infoblox
---

# Data Source `infoblox_cname_record`

Retrieves details for a cname record from infoblox

## Example Usage

```terraform
data "infoblox_fixed_address" "fixed_address" {
  alias = "example-alias.example.com"
}
```

```terraform
data "infoblox_cname_record" "cname_record" {
  ref = "record:cname/867530986753098675309867530986753098675309867530986753098675309:example-alias.example.com/default"
}
```

## Attributes Reference

The following attributes are exported.

- `ref` -  (MutuallyExclusiveGroup*/Computed, String) reference string.
- `alias` -  (MutuallyExclusiveGroup*/Computed, String) The name for a CNAME record in FQDN format.
- `canonical` - (Computed, String) Canonical name in FQDN format.
- `dns_name` -  (Computed, String) The name for the CNAME record in punycode format.
- `dns_canonical` -  (Computed,String) Canonical name in punycode format.
- `comment` - (Computed, String) Comment for the record; maximum 256 characters.
- `disable` - (Computed, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `view` - (Optional/Computed, String) The name of the DNS view in which the record resides.
- `zone` - (Optional/Computed, String) TThe name of the zone in which the record resides.
- `query_params` - (Optional, Map) Additional query parameters used for cname record query (see infoblox documentation for full list)
- `extensible_attributes` - (Computed, Map) Extensible attributes of cname record (Values are JSON encoded).

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key