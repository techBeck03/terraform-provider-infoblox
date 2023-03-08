---
page_title: "Container Resource - terraform-provider-infoblox"
subcategory: ""
description: |-
  Manages configuration details for a network container from infoblox
---

# Resource `infoblox_container`

Manages configuration details for a network container from infoblox

## Example Usage

```terraform
resource "infoblox_container" "container" {
  cidr = "172.19.10.0/23"
  comment = "test network container"
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

## Arguments Reference

The following arguments are exported.

- `comment` - (Optional, String) Comment for the container; maximum 256 characters.
- `extensible_attributes` - (Optional, Map) Extensible attributes of container (Values are JSON encoded).
- `cidr` -  (Required, String) The network address in IPv4 Address/CIDR format.
- `network_view` - (Optional, String) The name of the network view in which this container resides. Default value is `default`

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `ref` -  (Computed, String) Reference id of network container object.
