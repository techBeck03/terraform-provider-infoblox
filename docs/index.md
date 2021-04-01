---
page_title: "Provider: Infoblox"
subcategory: ""
description: |-
  Terraform provider for interacting with Infoblox WAPI
---

# Infoblox Provider


The infoblox provider is intented to be used with Infoblox WAPI v2.11 and above.  Please open any issues or submit pull requests for any addtional features you'd like to see added.

## Example Usage

To use this provider you will need the `hostname`, `username`, and `password` at a minimum.

```terraform
provider "infoblox" {
  url      = "https://infoblox.example.com"
  username = "admin"
  password = "password"
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

# Extensible Attributes

Extensible attributes are supported for all resource types within this provider and can be defined in the `extensible_attributes` argument.  Because of the varying value types and structure of extensible attributes within Infoblox, `extensible_attributes` are defined as a map of JSON encoded strings as shown in the example below:

```hcl
extensible_attributes = {
  Owner = jsonencode({
    value = "leeroyjenkins",
    type  = "STRING"
  })
  Location = jsonencode({
    value = "CollegeStation",
    type  = "STRING"
  })
  Orchestrator = jsonencode({
    value = "Terraform",
    type  = "ENUM"
  })
}
```

The `type` property tells the underlying go sdk how to type cast the value when sending to infoblox.  Supported values for `type` are:

- `STRING`
- `ENUM`
- `EMAIL`
- `URL`
- `DATE`
- `INTEGER`

## Extensible Attribute Inheritance

Each `extensible_attribute` also supports optional inheritance operations/actions such as `descendents_action` and `inheritance_operation`.  These values are not stored in state as they are one-time actions.  Subsequent terraform applies would always view these arguments as a new change since they are not stored in state.  Below are examples of how to use each:

### `inheritance_operation` (applied to a child infoblox object)

```hcl
extensible_attributes = {
  Owner = jsonencode({
    type  = "STRING",
    inheritance_operation = "INHERIT"
  })
  Location = jsonencode({
    value = "CollegeStation",
    type  = "STRING",
  })
}
```

### `descendents_action` (applied to a parent infoblox object)

```hcl
extensible_attributes = {
  Owner = jsonencode({
    value = "leeroyjenkins",
    type  = "STRING",
    descendants_action = jsonencode({
      option_with_ea    = "CONVERT"
      option_without_ea = "INHERIT"
    })
  })
  Location = jsonencode({
    value = "CollegeStation",
    type  = "STRING",
  })
  Gateway = jsonencode({
    value = "172.19.4.1",
    type  = "STRING",
  })
}
```

