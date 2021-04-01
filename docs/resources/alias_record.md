---
page_title: "Alias Record Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for an alias record in infoblox
---

# Resource `infoblox_alias_record`

Manages configuration details for an alias record in infoblox

## Example Usage

```terraform
resource "infoblox_alias_record" "alias" {
  name        = "alias.example.com"
  target_name = "realhost.example.com"
  target_type = "A"
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

- `comment` - (Optional, String) Comment for the record; maximum 256 characters.
- `disable` - (Optional, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `dns_name` -  (Computed, String) The name for an Alias record in punycode format.
- `dns_target_name` -  (Computed, String) Target name in punycode format.
- `extensible_attributes` - (Optional, Map) Extensible attributes of alias record (Values are JSON encoded).
- `name` -  (Required, String) The name for an Alias record in FQDN format.
- `target_name` - (Required, String) Target name in FQDN format.
- `target_type` - (Required, String) Target type.
- `view` - (Optional/Computed, String) The name of the DNS View in which the record resides. Example: “external”.
- `zone` - (Optional/Computed, String) The name of the zone in which the record resides.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of alias record object.