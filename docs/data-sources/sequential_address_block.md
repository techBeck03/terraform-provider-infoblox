---
page_title: "Sequential Address Block Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves available sequential address block from infoblox network
---

# Data Source `infoblox_sequential_address_block`

Retrieves available sequential address block from infoblox network

## Example Usage

```terraform
data "infoblox_sequential_address_block" "block-1" {
  cidr          = "172.19.4.0/24"
  address_count = 10
}
```

## Attributes Reference

The following attributes are exported.

- `address_count` - (Required, Int) Number of IPs to allocate.
- `addresses` - (Computed, List of Objects) List of IP addresses for sequential block.  Attributes for each set item:
  - `ref` - (Computed, String) Reference id of address object.
  - `ip_address` - (Computed, String) IPv4 address.
  - `hostnames` - (Computed, List) List of hostnames associated with IP address.
  - `mac_address` - (Computed, String) MAC address associated with IP address.
  - `network_view` - (Computed, String) Network view associated with IP address.
  - `cidr` - (Computed, String) CIDR associated with IP address.
  - `usage` - (Computed, String) Usage associated with IP address.
  - `types` - (Computed, String) Types associated with IP address.
  - `objects` - (Computed, String) Objects associated with IP address.
  - `status` - (Computed, String) Status associated with IP address.
- `cidr` -  (Required, String) Network for address block in IPv4 Address/CIDR format.
- `network_view` - (Optional/Computed, String) The name of the network view in which the address block resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key