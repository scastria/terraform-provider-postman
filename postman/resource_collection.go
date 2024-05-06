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
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"pre_request_script": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_response_script": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func fillCollectionUpdate(c *client.CollectionUpdateContainer, d *schema.ResourceData) {
	c.Child.Info.WorkspaceId = d.Get("workspace_id").(string)
	c.Child.Info.Name = d.Get("name").(string)
	description, ok := d.GetOk("description")
	if ok {
		c.Child.Info.Description = description.(string)
	}
	variables, ok := d.GetOk("var")
	c.Child.Variables = []client.Variable{}
	if ok {
		for _, variable := range variables.([]interface{}) {
			variableMap := variable.(map[string]interface{})
			c.Child.Variables = append(c.Child.Variables, client.Variable{
				Key:      variableMap["key"].(string),
				Value:    variableMap["value"].(string),
				Type:     "string",
				Disabled: variableMap["disabled"].(bool),
			})
		}
	}
	c.Child.Events = []client.Event{}
	preRequestScript, ok := d.GetOk("pre_request_script")
	if ok {
		c.Child.Events = append(c.Child.Events, client.Event{
			Listen: "prerequest",
			Script: client.Script{
				Id:   "prerequest",
				Type: client.TextJS,
				Exec: convertInterfaceArrayToStringArray(preRequestScript.([]interface{})),
			},
		})
	}
	postResponseScript, ok := d.GetOk("post_response_script")
	if ok {
		c.Child.Events = append(c.Child.Events, client.Event{
			Listen: "test",
			Script: client.Script{
				Id:   "test",
				Type: client.TextJS,
				Exec: convertInterfaceArrayToStringArray(postResponseScript.([]interface{})),
			},
		})
	}
}

func fillCollectionCreate(c *client.CollectionCreateContainer, d *schema.ResourceData) {
	c.Child.Info.WorkspaceId = d.Get("workspace_id").(string)
	c.Child.Info.Name = d.Get("name").(string)
	c.Child.Info.Schema = client.CollectionSchema
	c.Child.Items = []interface{}{}
	description, ok := d.GetOk("description")
	if ok {
		c.Child.Info.Description = description.(string)
	}
	variables, ok := d.GetOk("var")
	c.Child.Variables = []client.Variable{}
	if ok {
		for _, variable := range variables.([]interface{}) {
			variableMap := variable.(map[string]interface{})
			c.Child.Variables = append(c.Child.Variables, client.Variable{
				Key:      variableMap["key"].(string),
				Value:    variableMap["value"].(string),
				Type:     "string",
				Disabled: variableMap["disabled"].(bool),
			})
		}
	}
	c.Child.Events = []client.Event{}
	preRequestScript, ok := d.GetOk("pre_request_script")
	if ok {
		c.Child.Events = append(c.Child.Events, client.Event{
			Listen: "prerequest",
			Script: client.Script{
				Id:   "prerequest",
				Type: client.TextJS,
				Exec: convertInterfaceArrayToStringArray(preRequestScript.([]interface{})),
			},
		})
	}
	postResponseScript, ok := d.GetOk("post_response_script")
	if ok {
		c.Child.Events = append(c.Child.Events, client.Event{
			Listen: "test",
			Script: client.Script{
				Id:   "test",
				Type: client.TextJS,
				Exec: convertInterfaceArrayToStringArray(postResponseScript.([]interface{})),
			},
		})
	}
}

func fillResourceDataFromCollection(c *client.CollectionContainer, d *schema.ResourceData) {
	d.Set("workspace_id", c.Child.Info.WorkspaceId)
	d.Set("name", c.Child.Info.Name)
	d.Set("description", c.Child.Info.Description)
	d.Set("collection_id", c.Child.Info.Id)
	d.Set("uid", c.Child.Info.Uid)
	var variables []map[string]interface{}
	variables = nil
	if c.Child.Variables != nil {
		for _, variable := range c.Child.Variables {
			variableMap := map[string]interface{}{}
			variableMap["key"] = variable.Key
			variableMap["value"] = variable.Value
			variableMap["disabled"] = variable.Disabled
			variables = append(variables, variableMap)
		}
	}
	d.Set("var", variables)
	preRequestScripts := []string{}
	postResponseScripts := []string{}
	if c.Child.Events != nil {
		for _, event := range c.Child.Events {
			if event.Listen == "prerequest" {
				preRequestScripts = event.Script.Exec
			} else if event.Listen == "test" {
				postResponseScripts = event.Script.Exec
			}
		}
	}
	d.Set("pre_request_script", preRequestScripts)
	d.Set("post_response_script", postResponseScripts)
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
	// Must re-read the collection to get the full response
	requestPath = fmt.Sprintf(client.CollectionPathGet, retVal.Child.CreateId)
	body, err = c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	retVal2 := &client.CollectionContainer{}
	err = json.NewDecoder(body).Decode(retVal2)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal2.Child.Info.WorkspaceId = newCollection.Child.Info.WorkspaceId
	fillResourceDataFromCollection(retVal2, d)
	d.SetId(retVal2.Child.Info.CollectionEncodeId())
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
	workspaceId, id := client.CollectionDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCollection := client.CollectionUpdateContainer{}
	fillCollectionUpdate(&upCollection, d)
	err := json.NewEncoder(&buf).Encode(upCollection)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	// Must re-read the collection to get the full response
	requestPath = fmt.Sprintf(client.CollectionPathGet, id)
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.CollectionContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.Child.Info.WorkspaceId = workspaceId
	fillResourceDataFromCollection(retVal, d)
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
