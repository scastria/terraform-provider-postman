package postman

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-postman/postman/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POSTMAN_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"postman_workspace":  resourceWorkspace(),
			"postman_collection": resourceCollection(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	var diags diag.Diagnostics
	c, err := client.NewClient(apiKey)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
