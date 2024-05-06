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
			"num_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POSTMAN_NUM_RETRIES", 3),
			},
			"retry_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POSTMAN_RETRY_DELAY", 30),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"postman_workspace":       resourceWorkspace(),
			"postman_collection":      resourceCollection(),
			"postman_folder":          resourceFolder(),
			"postman_request":         resourceRequest(),
			"postman_collection_sort": resourceCollectionSort(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	numRetries := d.Get("num_retries").(int)
	retryDelay := d.Get("retry_delay").(int)
	var diags diag.Diagnostics
	c, err := client.NewClient(apiKey, numRetries, retryDelay)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
