package postman

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-postman/postman/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"pat": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KONNECT_PAT", nil),
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("KONNECT_REGION", "us"),
				ValidateFunc: validation.StringInSlice([]string{"us", "eu", "au"}, false),
			},
			//"default_tags": {
			//	Type:     schema.TypeSet,
			//	Optional: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//},
		},
		ResourcesMap:         map[string]*schema.Resource{},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	pat := d.Get("pat").(string)
	region := d.Get("region").(string)
	//defaultTags := []string{}
	//defaultTagsSet, ok := d.GetOk("default_tags")
	//if ok {
	//	defaultTags = convertSetToArray(defaultTagsSet.(*schema.Set))
	//}

	var diags diag.Diagnostics
	//c, err := client.NewClient(pat, region, defaultTags)
	c, err := client.NewClient(pat, region)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
