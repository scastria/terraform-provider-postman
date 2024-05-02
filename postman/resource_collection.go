package postman

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-postman/postman/client"
	"net/http"
	"net/url"
)

func resourceCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectionCreate,
		ReadContext:   resourceCollectionRead,
		UpdateContext: resourceCollectionUpdate,
		DeleteContext: resourceCollectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"collection_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fillCollection(c *client.CollectionContainer, d *schema.ResourceData) {
	c.Child.Info.WorkspaceId = d.Get("workspace_id").(string)
	c.Child.Info.Name = d.Get("name").(string)
	description, ok := d.GetOk("description")
	if ok {
		c.Child.Info.Description = description.(string)
	}
}

func fillCollectionCreate(c *client.CollectionCreateContainer, d *schema.ResourceData) {
	c.Child.Info.WorkspaceId = d.Get("workspace_id").(string)
	c.Child.Info.Name = d.Get("name").(string)
	c.Child.Info.Schema = client.CollectionSchema
	c.Child.Items = []string{}
	description, ok := d.GetOk("description")
	if ok {
		c.Child.Info.Description = description.(string)
	}
}

func fillResourceDataFromCollection(c *client.CollectionContainer, d *schema.ResourceData) {
	d.Set("workspace_id", c.Child.Info.WorkspaceId)
	d.Set("name", c.Child.Info.Name)
	d.Set("description", c.Child.Info.Description)
	d.Set("collection_id", c.Child.Info.Id)
	d.Set("uid", c.Child.Info.Uid)
}

func resourceCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newCollection := client.CollectionCreateContainer{}
	fillCollectionCreate(&newCollection, d)
	err := json.NewEncoder(&buf).Encode(newCollection)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CollectionPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	requestQuery := url.Values{}
	requestQuery[client.WorkspaceParam] = []string{newCollection.Child.Info.WorkspaceId}
	body, err := c.HttpRequest(ctx, http.MethodPost, requestPath, requestQuery, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.CollectionInfoContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal2 := client.CollectionContainer{}
	retVal2.Child.Info.Id = retVal.Child.CreateId
	retVal2.Child.Info.Uid = retVal.Child.Uid
	retVal2.Child.Info.Name = retVal.Child.Name
	retVal2.Child.Info.Description = newCollection.Child.Info.Description
	retVal2.Child.Info.WorkspaceId = newCollection.Child.Info.WorkspaceId
	d.SetId(retVal2.Child.Info.CollectionEncodeId())
	fillResourceDataFromCollection(&retVal2, d)
	return diags
}

func resourceCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	workspaceId, id := client.CollectionDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.CollectionContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.Child.Info.WorkspaceId = workspaceId
	fillResourceDataFromCollection(retVal, d)
	return diags
}

func resourceCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	_, id := client.CollectionDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCollection := client.CollectionContainer{}
	fillCollection(&upCollection, d)
	err := json.NewEncoder(&buf).Encode(upCollection)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.CollectionInfoContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal2 := client.CollectionContainer{}
	retVal2.Child.Info.Id = retVal.Child.CreateId
	uid, ok := d.GetOk("uid")
	if ok {
		retVal2.Child.Info.Uid = uid.(string)
	}
	retVal2.Child.Info.Name = retVal.Child.Name
	retVal2.Child.Info.Description = retVal.Child.Description
	retVal2.Child.Info.WorkspaceId = upCollection.Child.Info.WorkspaceId
	d.SetId(retVal2.Child.Info.CollectionEncodeId())
	fillResourceDataFromCollection(&retVal2, d)
	return diags
}

func resourceCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	_, id := client.CollectionDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
