---
page_title: "Cname Record Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for an cname record in infoblox
---

# Resource `infoblox_cname_record`

Manages configuration details for an cname record in infoblox

## Example Usage

```terraform
resource "infoblox_cname_record" "record" {
  alias     = "alias.example.com"
  comment   = "test cname record"
  canonical = "realhost.example.com"
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

- `alias` -  (Required, String) The name for a CNAME record in FQDN format.
- `canonical` - (Requried, String) Canonical name in FQDN format.
- `comment` - (Optional, String) Comment for the record; maximum 256 characters.
- `disable` - (Optional, Bool) Determines if the record is disabled or not. False means that the record is enabled.
- `dns_name` -  (Computed, String) The name for the CNAME record in punycode format.
- `dns_canonical` -  (Computed,String) Canonical name in punycode format.
- `extensible_attributes` - (Optional, Map) Extensible attributes of cname record (Values are JSON encoded).
- `view` - (Optional/Computed, String) The name of the DNS view in which the record resides.
- `zone` - (Optional/Computed, String) The name of the zone in which the record resides.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of cname record object.