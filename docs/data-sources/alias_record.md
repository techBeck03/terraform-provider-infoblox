---
page_title: "Alias Record Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for an alias record from infoblox
---

# Data Source `infoblox_alias_record`

Retrieves details for an alias record from infoblox

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

- `comment` - (Computed, String) Comment for the record; maximum 256 characters.
- `disable` - (Computed, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `dns_name` -  (Computed, String) The name for an Alias record in punycode format.
- `dns_target_name` -  (Computed, String) Target name in punycode format.
- `extensible_attributes` - (Computed, Map) Extensible attributes of alias record (Values are JSON encoded).
- `name` -  (MutuallyExclusiveGroup*/Computed, String) The name for an Alias record in FQDN format.
- `query_params` - (Optional, Map) Additional query parameters used for alias record query (see infoblox documentation for full list)
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of alias record object.
- `target_name` - (Computed, String) Target name in FQDN format.
- `target_type` - (Computed, String) Target type.
- `view` - (Optional/Computed, String) The name of the DNS View in which the record resides. Example: “external”.
- `zone` - (Optional/Computed, String) The name of the zone in which the record resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key