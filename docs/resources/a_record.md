---
page_title: "A Record Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for an A record in infoblox
---

# Resource `infoblox_a_record`

Manages configuration details for an A record in infoblox

## Example Usage

```terraform
resource "infoblox_a_record" "owen" {
  ip_address = "172.19.4.6"
  comment    = "test a record"
  hostname   = "realhost.example.com"
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
- `dns_name` -  (Computed, String) The name for an A record in punycode format.
- `extensible_attributes` - (Optional, Map) Extensible attributes of A record (Values are JSON encoded).
- `hostname` -  (Required, String) Name for A record in FQDN format.
- `ip_address` - (Required, String) The IPv4 Address of the record.
- `view` - (Optional, String) The name of the DNS view in which the record resides. Example: “external”.
- `zone` - (Computed, String) The name of the zone in which the record resides. If a view is not specified when searching by zone, the default view is used.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of A record object.