package postman

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-postman/postman/client"
	"net/http"
)

func resourceGlobalEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalEnvironmentCreate,
		ReadContext:   resourceGlobalEnvironmentRead,
		UpdateContext: resourceGlobalEnvironmentUpdate,
		DeleteContext: resourceGlobalEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"var": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "default",
							ValidateFunc: validation.StringInSlice([]string{"default", "secret"}, false),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
	}
}

func fillGlobalEnvironment(c *client.GlobalEnvironment, d *schema.ResourceData) {
	c.WorkspaceId = d.Get("workspace_id").(string)
	variables, ok := d.GetOk("var")
	c.Variables = []client.EnvVariable{}
	if ok {
		for _, variable := range variables.([]interface{}) {
			variableMap := variable.(map[string]interface{})
			c.Variables = append(c.Variables, client.EnvVariable{
				Key:     variableMap["key"].(string),
				Value:   variableMap["value"].(string),
				Type:    variableMap["type"].(string),
				Enabled: variableMap["enabled"].(bool),
			})
		}
	}
}

func fillResourceDataFromGlobalEnvironment(c *client.GlobalEnvironment, d *schema.ResourceData) {
	d.Set("workspace_id", c.WorkspaceId)
	var variables []map[string]interface{}
	variables = nil
	if c.Variables != nil {
		for _, variable := range c.Variables {
			variableMap := map[string]interface{}{}
			variableMap["key"] = variable.Key
			variableMap["value"] = variable.Value
			variableMap["type"] = variable.Type
			variableMap["enabled"] = variable.Enabled
			variables = append(variables, variableMap)
		}
	}
	d.Set("var", variables)
}

func resourceGlobalEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newGlobalEnvironment := client.GlobalEnvironment{}
	fillGlobalEnvironment(&newGlobalEnvironment, d)
	err := json.NewEncoder(&buf).Encode(newGlobalEnvironment)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GlobalEnvironmentPath, newGlobalEnvironment.WorkspaceId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.GlobalEnvironment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.WorkspaceId = newGlobalEnvironment.WorkspaceId
	fillResourceDataFromGlobalEnvironment(retVal, d)
	d.SetId(retVal.WorkspaceId)
	return diags
}

func resourceGlobalEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.GlobalEnvironmentPath, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.GlobalEnvironment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.WorkspaceId = d.Id()
	fillResourceDataFromGlobalEnvironment(retVal, d)
	return diags
}

func resourceGlobalEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upGlobalEnvironment := client.GlobalEnvironment{}
	fillGlobalEnvironment(&upGlobalEnvironment, d)
	err := json.NewEncoder(&buf).Encode(upGlobalEnvironment)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GlobalEnvironmentPath, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.GlobalEnvironment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.WorkspaceId = d.Id()
	fillResourceDataFromGlobalEnvironment(retVal, d)
	return diags
}

func resourceGlobalEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	// Clear out all global variables
	upGlobalEnvironment := client.GlobalEnvironment{
		WorkspaceId: d.Id(),
	}
	err := json.NewEncoder(&buf).Encode(upGlobalEnvironment)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.GlobalEnvironmentPath, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
