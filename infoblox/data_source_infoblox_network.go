package infoblox

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/techBeck03/infoblox-go-client"
)

func dataSourceInfobloxNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,
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

func dataSourceNetworkRead(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning to get network ", resourceNetworkIDString(d))

	connector := m.(*ibclient.Connector)

	cidr := d.Get("cidr").(string)
	networkViewName := d.Get("network_view_name").(string)
	tenantID := d.Get("tenant_id").(string)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetNetwork(networkViewName, cidr, nil)
	if err != nil {
		return fmt.Errorf("Getting Network block from network (%s) failed : %s", cidr, err)
	}

	if obj == nil {
		return fmt.Errorf("API returns a nil/empty id on network (%s) failed", cidr)
	}
	d.SetId(obj.Ref)
	if obj.Ea["Network Name"] != nil {
		d.Set("network_name", obj.Ea["Network Name"])
	}
	d.Set("gateway", obj.Ea["Gateway"])
	return nil
}
