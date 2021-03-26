---
page_title: "Grid Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a grid from infoblox
---

# Data Source `infoblox_grid`

Retrieves details for a grid from infoblox

## Example Usage

```terraform
data "infoblox_grid" "grid" {
  name = "Infoblox"
}
```

## Attributes Reference

The following attributes are exported.

- `dns_resolvers` -  (Computed, List) List of DNS resolvers.
- `dns_search_domains` - (Optional, List) List of DNS search domains.
- `name` -  (MutuallyExclusiveGroup*/Computed, String) Name of grid.
- `query_params` - (Optional, Map) Additional query parameters used for grid member query (see infoblox documentation for full list)
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of grid object.
- `service_status` - (Computed, String) Service status of grid.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key