---
page_title: "Container Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a network container from infoblox
---

# Data Source `infoblox_container`

Retrieves details for a network container from infoblox

## Example Usage

```terraform
data "infoblox_container" "container" {
  cidr = "172.19.10.0/23"
}
```

```terraform
data "infoblox_container" "container" {
  ref = "networkcontainer/867530986753098675309867530986753098675309867530986753098675309:172.19.10.0/23/default"
}
```

## Attributes Reference

The following attributes are exported.

- `comment` - (Computed, String) Comment for the container; maximum 256 characters.
- `extensible_attributes` - (Computed, Map) Extensible attributes of container (Values are JSON encoded).
- `cidr` -  (MutuallyExclusiveGroup*/Computed, String) The network address in IPv4 Address/CIDR format.
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of container object.
- `network_view` - (Computed, String) The name of the network view in which this container resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key