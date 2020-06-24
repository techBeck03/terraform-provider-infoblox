package infoblox

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,

		Schema: map[string]*schema.Schema{
			"network_view_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Network view name available in NIOS Server.",
			},
			"network_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of your network block.",
			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network block in cidr format.",
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of your tenant in cloud.",
			},
			"reserve_ip": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The no of IP's you want to reserve.",
			},
			"gateway": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "gateway ip address of your network block.By default first IPv4 address is set as gateway address.",
				Computed:    true,
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning network block Creation", resourceNetworkIDString(d))

	networkViewName := d.Get("network_view_name").(string)
	cidr := d.Get("cidr").(string)
	networkName := d.Get("network_name").(string)
	reserveIP := d.Get("reserve_ip").(int)
	gateway := d.Get("gateway").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	ZeroMacAddr := "00:00:00:00:00:00"
	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)
	ea := make(ibclient.EA)

	nwname, err := objMgr.CreateNetwork(networkViewName, cidr, networkName)
	if err != nil {
		return fmt.Errorf("Creation of network block failed in network view (%s) : %s", networkViewName, err)
	}

	// Check whether gateway or ip address already allocated
	gatewayIP, err := objMgr.GetFixedAddress(networkViewName, cidr, gateway, "")
	if err == nil && gatewayIP != nil {
		fmt.Printf("Gateway already created")
	} else if gatewayIP == nil {
		gatewayIP, err = objMgr.AllocateIP(networkViewName, cidr, gateway, ZeroMacAddr, "", ea)
		if err != nil {
			return fmt.Errorf("Gateway Creation failed in network block(%s) error: %s", cidr, err)
		}
	}

	for i := 1; i <= reserveIP; i++ {
		_, err = objMgr.AllocateIP(networkViewName, cidr, gateway, ZeroMacAddr, "", ea)
		if err != nil {
			return fmt.Errorf("Reservation in network block failed in network view(%s):%s", networkViewName, err)
		}
	}

	_, err = objMgr.UpdateNetwork(nwname.Ref, networkName, tenantID, gatewayIP.IPAddress)
	if err != nil {
		return fmt.Errorf("Getting Network block from network view (%s) failed : %s", networkViewName, err)
	}

	d.Set("gateway", gatewayIP.IPAddress)
	d.SetId(nwname.Ref)

	log.Printf("[DEBUG] %s: Creation on network block complete", resourceNetworkIDString(d))
	return resourceNetworkRead(d, m)
}
func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Reading the required network block", resourceNetworkIDString(d))

	networkViewName := d.Get("network_view_name").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetNetworkwithref(d.Id())
	if err != nil {
		return fmt.Errorf("Getting Network block from network view (%s) failed : %s", networkViewName, err)
	}
	d.SetId(obj.Ref)
	d.Set("network_name", obj.Ea["Network Name"])
	d.Set("tenant_id", obj.Ea["Tenant ID"])
	d.Set("gateway", obj.Ea["Gateway"])
	log.Printf("[DEBUG] %s: Completed reading network block", resourceNetworkIDString(d))
	return nil
}
func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("cidr") || d.HasChange("network_view_name") || d.HasChange("gateway") || d.HasChange("reserveIP") {
		return fmt.Errorf("Only network name and tenant_id updates are allowed")
	}
	networkViewName := d.Get("network_view_name").(string)
	networkName := d.Get("network_name").(string)
	tenantID := d.Get("tenant_id").(string)
	gateway := d.Get("gateway").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.UpdateNetwork(d.Id(), networkName, tenantID, gateway)
	if err != nil {
		return fmt.Errorf("Getting Network block from network view (%s) failed : %s", networkViewName, err)
	}
	d.Set("network_name", obj.Ea["Network Name"])
	d.Set("tenant_id", obj.Ea["Tenant ID"])
	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Deletion of network block", resourceNetworkIDString(d))

	networkViewName := d.Get("network_view_name").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	_, err := objMgr.DeleteNetwork(d.Id(), d.Get("network_view_name").(string))
	if err != nil {
		return fmt.Errorf("Deletion of Network block failed from network view(%s): %s", networkViewName, err)
	}
	d.SetId("")

	log.Printf("[DEBUG] %s: Deletion of network block complete", resourceNetworkIDString(d))
	return nil
}

type resourceNetworkIDStringInterface interface {
	Id() string
}

func resourceNetworkIDString(d resourceNetworkIDStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("infoblox_ip_allocation (ID = %s)", id)
}
