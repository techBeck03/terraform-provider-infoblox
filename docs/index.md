---
page_title: "Provider: Infoblox"
subcategory: ""
description: |-
  Terraform provider for interacting with Infoblox WAPI
---

# Infoblox Provider


The infoblox provider is intented to be used with Infoblox WAPI v2.11 and above.  Please open any issues or submit pull requests for any addtional features you'd like to see added.

Provider versioning will align with Infoblox's WAPI release versioning starting with release `2.11`

## Example Usage

To use this provider you will need the `hostname`, `username`, and `password` at a minimum.

```terraform
provider "guacamole" {
  url      = "https://guacamole.example.com"
  username = "guacadmin"
  password = "guacadmin"
  disable_tls_verification = true
  orchestrator_extensible_attributes = {
    Orchestrator = jsonencode({
      value = "Terraform",
      type  = "ENUM"
    })
  }
}
```

## Schema

- **hostname** (Required, String) Hostname or IP address of Grid master (defaults to environment variable `INFOBLOX_HOSTNAME`).
- **username** (Required, String) Username to authenticate to infoblox (defaults to environment variable `INFOBLOX_USERNAME`).
- **password** (Required, String) Password to authenticate to infoblox (defaults to environment variable `INFOBLOX_PASSWORD`).
- **port** (Required, String) Port on which to communicate with infoblox (defaults to environment variable `INFOBLOX_PORT` or `443` no value is set).
- **disable_tls_verification** (Optional, Bool) Whether to disable tls verification for ssl connections (defaults to environment variable `INFOBLOX_DISABLE_TLS` or `false` if no value is set).
- **wapi_version** (Optional, String) WAPI version (defaults to environment variable `INFOBLOX_VERSION` or `2.11` if no value is set).
- **orchestrator_extensible_attributes** (Optional, Map) Extensible attributes applied to all objects configured by provider. 

```