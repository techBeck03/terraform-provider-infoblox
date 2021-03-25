---
page_title: "A Record Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for an A record from infoblox
---

# Data Source `infoblox_a_record`

Retrieves details for an A record from infoblox

## Example Usage

```terraform
data "infoblox_a_record" "arecord" {
  hostname = "example-hostname.example.com"
}
```

```terraform
data "infoblox_a_record" "arecord" {
  ref = "record:a/867530986753098675309867530986753098675309867530986753098675309:example-hostname.example.com/default"
}
```

## Attributes Reference

The following attributes are exported.

- `ref` -  (MutuallyExclusiveGroup*/Computed, String) reference string.
- `hostname` -  (MutuallyExclusiveGroup*/Computed, String) Name for A record in FQDN format.
- `ip_address` - (MutuallyExclusiveGroup*/Computed, String) The IPv4 Address of the record.
- `dns_name` -  (Computed, String) The name for an A record in punycode format.
- `comment` - (Computed, String) Comment for the record; maximum 256 characters.
- `disable` - (Computed, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `view` - (Optional/Computed, String) The name of the DNS view in which the record resides. Example: “external”.
- `zone` - (Optional/Computed, String) The name of the zone in which the record resides. If a view is not specified when searching by zone, the default view is used.
- `query_params` - (Optional, Map) Additional query parameters used for A record query (see infoblox documentation for full list)
- `extensible_attributes` - (Computed, Map) Extensible attributes of A record (Values are JSON encoded).

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key