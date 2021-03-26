---
page_title: "Host Record Data Source - terraform-provider-infoblox"
subcategory: ""
description: |-
  Retrieves details for a host record from infoblox
---

# Data Source `infoblox_host_record`

Retrieves details for a host record from infoblox

## Example Usage

```terraform
data "infoblox_host_record" "host_record" {
  hostname = "example-host.example.com"
}
```

```terraform
data "infoblox_host_record" "host_record" {
  ref = "record:host/867530986753098675309867530986753098675309867530986753098675309:example-host.example.com/default"
}
```

## Attributes Reference

The following attributes are exported.

- `comment` - (Computed, String) Comment for the fixed address; maximum 256 characters.
- `enable_dns` - (Computed, Bool) When false, the host does not have parent zone information.
- `extensible_attributes` - (Computed, Map) Extensible attributes of host record (Values are JSON encoded).
- `hostname` -  (MutuallyExclusiveGroup*/Computed, String) The host name in FQDN format.
- `ip_v4_address` - (Computed, Set of Objects) IPv4 addresses associated with host record.  Attributes for each set item:
  - `ref` - (Computed, String) Reference id of address object.
  - `ip_address` - (Computed, String) IP address.
  - `hostname` - (Computed, String) Hostname associated with IP address.
  - `network` - (Computed, String) Network associated with IP address.
  - `mac_address` - (Computed, String) MAC address associated with IP address.
  - `configure_for_dhcp` - (Computed, Bool) Set this to True to enable the DHCP configuration for this host address.
- `network_view` -  (Computed, String) The name of the network view in which this fixed address resides.
- `query_params` - (Optional, Map) Additional query parameters used for host record query (see infoblox documentation for full list)
- `ref` -  (MutuallyExclusiveGroup*/Computed, String) Reference id of host record object.
- `view` - (Optional/Computed, String) The name of the DNS view in which the record resides.
- `zone` - (Optional/Computed, String) The name of the zone in which the record resides.

**_MutuallyExclusiveGroup_**: One and only one of the attritbutes in this group **MUST** be provided as a primary search key