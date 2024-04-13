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

func resourceWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkspaceCreate,
		ReadContext:   resourceWorkspaceRead,
		UpdateContext: resourceWorkspaceUpdate,
		DeleteContext: resourceWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"personal", "private", "public", "team", "partner"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func fillWorkspace(c *client.WorkspaceContainer, d *schema.ResourceData) {
	c.Child.Name = d.Get("name").(string)
	c.Child.Type = d.Get("type").(string)
	description, ok := d.GetOk("description")
	if ok {
		c.Child.Description = description.(string)
	}
}

func fillResourceDataFromWorkspace(c *client.WorkspaceContainer, d *schema.ResourceData) {
	d.Set("name", c.Child.Name)
	d.Set("type", c.Child.Type)
	d.Set("description", c.Child.Description)
}

func resourceWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newWorkspace := client.WorkspaceContainer{}
	fillWorkspace(&newWorkspace, d)
	err := json.NewEncoder(&buf).Encode(newWorkspace)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.WorkspacePath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.WorkspaceContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Child.Id)
	fillResourceDataFromWorkspace(retVal, d)
	return diags
}

func resourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.WorkspacePathGet, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.WorkspaceContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromWorkspace(retVal, d)
	return diags
}

func resourceWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upWorkspace := client.WorkspaceContainer{}
	fillWorkspace(&upWorkspace, d)
	err := json.NewEncoder(&buf).Encode(upWorkspace)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.WorkspacePathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.WorkspaceContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceDataFromWorkspace(retVal, d)
	return diags
}

func resourceWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.WorkspacePathGet, d.Id())
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
