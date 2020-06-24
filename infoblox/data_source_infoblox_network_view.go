package infoblox

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func dataSourceInfobloxNetworkView() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkViewRead,
		Schema: map[string]*schema.Schema{
			"network_view_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Desired name of the view shown in NIOS appliance.",
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of your tenant in cloud.",
			},
			"extensible_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Optional map of extensible attributes",
			},
		},
	}
}

func dataSourceNetworkViewRead(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning to get network view ", resourceNetworkViewIDString(d))

	tenantID := d.Get("tenant_id").(string)
	Connector := m.(*ibclient.Connector)
	objMgr := ibclient.NewObjectManager(Connector, "Terraform", tenantID)

	obj, err := objMgr.GetNetworkView(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to get Network View : %s", err)
	}
	for key, value := range obj.Ea {
		convertedValue := ""
		switch value.(type) {
		case bool, ibclient.Bool:
			convertedValue = fmt.Sprintf("%t", value)
		case int:
			convertedValue = fmt.Sprintf("%d", value)
		default:
			convertedValue = value.(string)
		}
		obj.Ea[key] = convertedValue
	}
	d.Set("tenant_id", obj.Ea["Tenant ID"])
	d.Set("extensible_attributes", obj.Ea)
	d.SetId(obj.Name)
	log.Printf("[DEBUG] %s: got Network View", resourceNetworkViewIDString(d))

	return nil
}
