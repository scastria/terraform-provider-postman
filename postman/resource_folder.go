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
)

func resourceFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFolderCreate,
		ReadContext:   resourceFolderRead,
		UpdateContext: resourceFolderUpdate,
		DeleteContext: resourceFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"collection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_folder_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func fillFolder(c *client.FolderData, d *schema.ResourceData) {
	c.Name = d.Get("name").(string)
	description, ok := d.GetOk("description")
	if ok {
		c.Description = description.(string)
	}
	parentFolderId, ok := d.GetOk("parent_folder_id")
	if ok {
		c.ParentFolderId = parentFolderId.(string)
	}
	c.Events = []client.Event{}
	preRequestScript, ok := d.GetOk("pre_request_script")
	if ok {
		c.Events = append(c.Events, client.Event{
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
		c.Events = append(c.Events, client.Event{
			Listen: "test",
			Script: client.Script{
				Id:   "test",
				Type: client.TextJS,
				Exec: convertInterfaceArrayToStringArray(postResponseScript.([]interface{})),
			},
		})
	}
}

func fillResourceDataFromFolder(c *client.Folder, d *schema.ResourceData) {
	d.Set("collection_id", c.Data.CollectionId)
	d.Set("name", c.Data.Name)
	d.Set("description", c.Data.Description)
	d.Set("parent_folder_id", c.Data.ParentFolderId)
	d.Set("folder_id", c.Data.Id)
	preRequestScripts := []string{}
	postResponseScripts := []string{}
	if c.Data.Events != nil {
		for _, event := range c.Data.Events {
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

func resourceFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newFolder := client.FolderData{}
	fillFolder(&newFolder, d)
	err := json.NewEncoder(&buf).Encode(newFolder)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.FolderPath, d.Get("collection_id").(string))
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Folder{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.FolderEncodeId())
	fillResourceDataFromFolder(retVal, d)
	return diags
}

func resourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	collectionId, id := client.FolderDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.FolderPathGet, collectionId, id)
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Folder{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromFolder(retVal, d)
	return diags
}

func resourceFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	collectionId, id := client.FolderDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upFolder := client.FolderData{}
	fillFolder(&upFolder, d)
	err := json.NewEncoder(&buf).Encode(upFolder)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.FolderPathGet, collectionId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Folder{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.FolderEncodeId())
	fillResourceDataFromFolder(retVal, d)
	return diags
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	collectionId, id := client.FolderDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.FolderPathGet, collectionId, id)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
