---
page_title: "Grid Member Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a grid member from infoblox
---

# Data Source `infoblox_grid_member`

Retrieves details for a grid member from infoblox

## Example Usage

```terraform
data "infoblox_grid_member" "grid_member" {
  hostname = "infoblox.example.com"
}
```

## Attributes Reference

The following attributes are exported.

- `config_address_type` - (Computed, String) Configured IP address type.
- `hostname` -  (MutuallyExclusiveGroup*/Computed, String) Hostname of member in FQDN format.
- `query_params` - (Optional, Map) Additional query parameters used for grid member query (see infoblox documentation for full list)
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of member.
- `service_type_configuration` -  (Computed, String) Service type configuration.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key