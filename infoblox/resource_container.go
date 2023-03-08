package infoblox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	infoblox "github.com/techBeck03/infoblox-go-sdk"
)

func resourceContainer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerCreate,
		ReadContext:   resourceContainerRead,
		UpdateContext: resourceContainerUpdate,
		DeleteContext: resourceContainerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.Sequence(
			makeEACustomDiff("extensible_attributes"),
		),
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:             schema.TypeString,
				Description:      "The container network address in IPv4 Address/CIDR format.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"comment": {
				Type:             schema.TypeString,
				Description:      "Comment for the container; maximum 256 characters.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 256)),
			},
			"extensible_attributes": {
				Type:             schema.TypeMap,
				Description:      "Extensible attributes of A container (Values are JSON encoded).",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateEa,
				DiffSuppressFunc: eaSuppressDiff,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_view": {
				Type:        schema.TypeString,
				Description: "The name of the network view in which this network resides.",
				Default:     "default",
				Optional:    true,
				ForceNew:    true,
			},
			"ref": {
				Type:        schema.TypeString,
				Description: "Reference id of A container object.",
				Computed:    true,
			},
		},
	}
}

func convertContainerToResourceData(client *infoblox.Client, d *schema.ResourceData, container *infoblox.NetworkContainer) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("ref", container.Ref)
	d.Set("cidr", container.CIDR)
	d.Set("comment", container.Comment)
	d.Set("network_view", container.NetworkView)

	eas, err := client.ConvertEAsToJSONString(*container.ExtensibleAttributes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		d.Set("extensible_attributes", eas)
	}

	return diags
}

func convertResourceDataToContainer(client *infoblox.Client, d *schema.ResourceData) (*infoblox.NetworkContainer, error) {
	var container infoblox.NetworkContainer

	container.CIDR = d.Get("cidr").(string)
	container.Comment = d.Get("comment").(string)
	container.NetworkView = d.Get("network_view").(string)

	eaMap := d.Get("extensible_attributes").(map[string]interface{})
	if len(eaMap) > 0 {
		eas, err := createExtensibleAttributesFromJSON(eaMap)
		if err != nil {
			return &container, err
		}
		container.ExtensibleAttributes = &eas
	}

	if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
		if container.ExtensibleAttributes == nil {
			container.ExtensibleAttributes = &infoblox.ExtensibleAttribute{}
		}
		for k, v := range *client.OrchestratorEAs {
			(*container.ExtensibleAttributes)[k] = v
		}
	}

	return &container, nil
}

func resourceContainerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	container, err := client.GetContainerByRef(ref, nil)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	check := convertContainerToResourceData(client, d, &container)
	if check.HasError() {
		return check
	}

	d.SetId(container.Ref)

	return diags
}

func resourceContainerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics

	container, err := convertResourceDataToContainer(client, d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = client.CreateContainer(container)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(container.Ref)
	return resourceContainerRead(ctx, d, m)
}

func resourceContainerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	client := m.(*infoblox.Client)

	var container infoblox.NetworkContainer

	if d.HasChange("cidr") {
		container.CIDR = d.Get("cidr").(string)
	}
	if d.HasChange("comment") {
		container.Comment = d.Get("comment").(string)
	}
	if d.HasChange("view") {
		container.NetworkView = d.Get("network_view").(string)
	}
	if d.HasChange("extensible_attributes") {
		old, new := d.GetChange("extensible_attributes")
		oldKeys := Keys(old.(map[string]interface{}))
		oldEAs, err := createExtensibleAttributesFromJSON(old.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		newKeys := Keys(new.(map[string]interface{}))
		newEAs, err := createExtensibleAttributesFromJSON(new.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		removeEAs := sliceDiff(oldKeys, newKeys, false)
		if len(removeEAs) > 0 {
			container.ExtensibleAttributesRemove = &infoblox.ExtensibleAttribute{}
			for _, v := range removeEAs {
				(*container.ExtensibleAttributesRemove)[v] = oldEAs[v]
			}
		}
		for k, v := range newEAs {
			if !Contains(oldKeys, k) || (Contains(oldKeys, k) && v.Value != oldEAs[k].Value) {
				if container.ExtensibleAttributesAdd == nil {
					container.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
				}
				(*container.ExtensibleAttributesAdd)[k] = v
			}
		}
		if client.OrchestratorEAs != nil && len(*client.OrchestratorEAs) > 0 {
			if container.ExtensibleAttributesAdd == nil {
				container.ExtensibleAttributesAdd = &infoblox.ExtensibleAttribute{}
			}
			for k, v := range *client.OrchestratorEAs {
				(*container.ExtensibleAttributesAdd)[k] = v
			}
		}
	}
	changedcontainer, err := client.UpdateContainer(d.Id(), container)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(changedcontainer.Ref)
	return resourceContainerRead(ctx, d, m)
}

func resourceContainerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*infoblox.Client)

	var diags diag.Diagnostics
	ref := d.Id()

	err := client.DeleteContainer(ref)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
