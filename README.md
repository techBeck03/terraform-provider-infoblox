# Terraform Provider for Infoblox

This is a fork of [terraform-providers/terraform-provider-infoblox](https://github.com/terraform-providers/terraform-provider-infoblox) to support string based extensible attributes and additional data resources.

## Usage Examples

### Create a network with Comment
```
# Retrieve network view for new network
data "infoblox_network_view" "nv"{
  network_view_name="default"
}
resource "infoblox_network" "net-1"{
  network_view_name=data.infoblox_network_view.nv.network_view_name
  network_name="Network 1"
  cidr="1.1.1.0/24"
  tenant_id="infra"
}
```

### Create a host record with EAs
```
# Retrieve network view for new network
data "infoblox_network_view" "nv"{
  network_view_name="default"
}

# Retrieve network for host record
data "infoblox_network" "nw"{
  network_view_name=data.infoblox_network_view.nv.network_view_name
  cidr="1.1.1.0/24"
  tenant_id="infra"
}

# Create allocation
resource "infoblox_ip_allocation" "allocation"{
  network_view_name=data.infoblox_network_view.nv.network_view_name
  vm_name="server-1"
  zone="example.com"
  dns_view="default"
  enable_dns=true
  cidr=data.infoblox_network.nw.cidr
  tenant_id="infra"
  extensible_attributes = {
    Environment = "Dev"
    Deployment = "1234567"
    Owner = "ACME User"
  }
}

# Assuming virtual machine is created elsewhere
resource "infoblox_ip_association" "associate"{
  vm_name=vsphere_virtual_machine.vm.name
  cidr=infoblox_ip_allocation.allocation.cidr
  mac_addr=vsphere_virtual_machine.vm.network_interface.0.mac_address
  ip_addr=vsphere_virtual_machine.vm.default_ip_address
  vm_id=vsphere_virtual_machine.vm.id
  tenant_id="infra"
  zone="example.com"
  dns_view="default"
}
```


## Disclaimer
This is a community provider so use at your own risk
