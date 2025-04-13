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
	"net/url"
	"strings"
)

func resourceRequest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRequestCreate,
		ReadContext:   resourceRequestRead,
		UpdateContext: resourceRequestUpdate,
		DeleteContext: resourceRequestDelete,
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
			"folder_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GET",
				ValidateFunc: validation.StringInSlice([]string{
					"GET",
					"PUT",
					"POST",
					"PATCH",
					"DELETE",
					"COPY",
					"HEAD",
					"OPTIONS",
					"LINK",
					"UNLINK",
					"PURGE",
					"LOCK",
					"UNLOCK",
					"PROPFIND",
					"VIEW",
				}, false),
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_param": {
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
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"header": {
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
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
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

func fillRequest(c *client.RequestData, d *schema.ResourceData) {
	c.Name = d.Get("name").(string)
	baseURL, ok := d.GetOk("base_url")
	if ok {
		c.BaseURL = baseURL.(string)
	}
	description, ok := d.GetOk("description")
	if ok {
		c.Description = description.(string)
	}
	body, ok := d.GetOk("body")
	if ok {
		c.Body = body.(string)
		if c.Body != "" {
			c.DataMode = "raw"
			c.Options = client.DataOptions{
				Raw: client.RawDataOptions{
					Language: "json",
				},
			}
		}
	}
	folderId, ok := d.GetOk("folder_id")
	if ok {
		c.FolderId = folderId.(string)
	}
	method, ok := d.GetOk("method")
	if ok {
		c.Method = method.(string)
	}
	queryParams, ok := d.GetOk("query_param")
	if ok {
		c.QueryParams = []client.QueryParamHeader{}
		for _, queryParam := range queryParams.([]interface{}) {
			queryParamMap := queryParam.(map[string]interface{})
			c.QueryParams = append(c.QueryParams, client.QueryParamHeader{
				Key:         queryParamMap["key"].(string),
				Value:       queryParamMap["value"].(string),
				Description: queryParamMap["description"].(string),
				Enabled:     queryParamMap["enabled"].(bool),
			})
		}
	}
	headerValues, ok := d.GetOk("header")
	if ok {
		c.Headers = []client.QueryParamHeader{}
		for _, header := range headerValues.([]interface{}) {
			headerMap := header.(map[string]interface{})
			c.Headers = append(c.Headers, client.QueryParamHeader{
				Key:         headerMap["key"].(string),
				Value:       headerMap["value"].(string),
				Description: headerMap["description"].(string),
				Enabled:     headerMap["enabled"].(bool),
			})
		}
	}
	requestQuery := url.Values{}
	for _, queryParam := range c.QueryParams {
		if !queryParam.Enabled {
			continue
		}
		requestQuery.Add(queryParam.Key, queryParam.Value)
	}
	encodedQuery := requestQuery.Encode()
	// Manually append query params to prevent URL parsing from escaping variables in the host and path
	if len(encodedQuery) > 0 {
		c.URL = c.BaseURL + "?" + encodedQuery
	} else {
		c.URL = c.BaseURL
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

func fillResourceDataFromRequest(c *client.Request, d *schema.ResourceData) {
	d.Set("collection_id", c.Data.CollectionId)
	d.Set("name", c.Data.Name)
	d.Set("description", c.Data.Description)
	d.Set("body", c.Data.Body)
	d.Set("folder_id", c.Data.FolderId)
	d.Set("method", c.Data.Method)
	// Manually strip query params to prevent URL parsing from escaping variables in the host and path
	d.Set("base_url", strings.Split(c.Data.URL, "?")[0])
	d.Set("url", c.Data.URL)
	d.Set("request_id", c.Data.Id)
	var queryParams []map[string]interface{}
	queryParams = nil
	if c.Data.QueryParams != nil {
		for _, queryParam := range c.Data.QueryParams {
			queryParamMap := map[string]interface{}{}
			queryParamMap["key"] = queryParam.Key
			queryParamMap["value"] = queryParam.Value
			queryParamMap["description"] = queryParam.Description
			queryParamMap["enabled"] = queryParam.Enabled
			queryParams = append(queryParams, queryParamMap)
		}
	}
	d.Set("query_param", queryParams)
	var headerValues []map[string]interface{}
	headerValues = nil
	if c.Data.Headers != nil {
		for _, header := range c.Data.Headers {
			headerMap := map[string]interface{}{}
			headerMap["key"] = header.Key
			headerMap["value"] = header.Value
			headerMap["description"] = header.Description
			headerMap["enabled"] = header.Enabled
			headerValues = append(headerValues, headerMap)
		}
	}
	d.Set("header", headerValues)
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

func resourceRequestCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newRequest := client.RequestData{}
	fillRequest(&newRequest, d)
	err := json.NewEncoder(&buf).Encode(newRequest)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RequestPath, d.Get("collection_id").(string))
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Request{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.RequestEncodeId())
	fillResourceDataFromRequest(retVal, d)
	return diags
}

func resourceRequestRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	collectionId, id := client.RequestDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RequestPathGet, collectionId, id)
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Request{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromRequest(retVal, d)
	return diags
}

func resourceRequestUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	collectionId, id := client.RequestDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upRequest := client.RequestData{}
	fillRequest(&upRequest, d)
	err := json.NewEncoder(&buf).Encode(upRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RequestPathGet, collectionId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Request{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.RequestEncodeId())
	fillResourceDataFromRequest(retVal, d)
	return diags
}

func resourceRequestDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	collectionId, id := client.RequestDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RequestPathGet, collectionId, id)
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
