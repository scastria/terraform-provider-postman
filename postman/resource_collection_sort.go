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
	"sort"
)

func resourceCollectionSort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectionSortCreate,
		ReadContext:   resourceCollectionSortRead,
		UpdateContext: resourceCollectionSortUpdate,
		DeleteContext: resourceCollectionSortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"collection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"order": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ASC",
				ValidateFunc: validation.StringInSlice([]string{"ASC", "DESC", "NONE"}, false),
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCollectionSortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	d.SetId(d.Get("collection_id").(string))
	clearId, err := applySort(ctx, d, c)
	if clearId {
		d.SetId("")
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCollectionSortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CollectionPathGet, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.CollectionSortContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	isASC := sort.SliceIsSorted(retVal.Child.Items, func(i, j int) bool {
		return retVal.Child.Items[i]["name"].(string) < retVal.Child.Items[j]["name"].(string)
	})
	isDESC := false
	if !isASC {
		isDESC = sort.SliceIsSorted(retVal.Child.Items, func(i, j int) bool {
			return retVal.Child.Items[i]["name"].(string) > retVal.Child.Items[j]["name"].(string)
		})
	}
	order := "NONE"
	if isASC {
		order = "ASC"
	} else if isDESC {
		order = "DESC"
	}
	d.Set("order", order)
	d.Set("collection_id", retVal.Child.Info.Id)
	return diags
}

func applySort(ctx context.Context, d *schema.ResourceData, c *client.Client) (bool, error) {
	// Read current
	requestPath := fmt.Sprintf(client.CollectionPathGet, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return true, nil
		}
		return true, err
	}
	retVal := &client.CollectionSortContainer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return true, err
	}
	// Sort items
	order := d.Get("order").(string)
	if order == "NONE" {
		return false, nil
	}
	sort.Slice(retVal.Child.Items, func(i, j int) bool {
		if order == "ASC" {
			return retVal.Child.Items[i]["name"].(string) < retVal.Child.Items[j]["name"].(string)
		} else {
			return retVal.Child.Items[i]["name"].(string) > retVal.Child.Items[j]["name"].(string)
		}
	})
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(retVal)
	if err != nil {
		return false, err
	}
	requestPath = fmt.Sprintf(client.CollectionPathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return false, err
	}
	return false, nil
}

func resourceCollectionSortUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	clearId, err := applySort(ctx, d, c)
	if clearId {
		d.SetId("")
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCollectionSortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}
