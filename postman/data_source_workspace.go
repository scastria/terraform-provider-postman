package postman

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-postman/postman/client"
	"net/http"
	"net/url"
	"strings"
)

func dataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"search_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"personal", "private", "public", "team", "partner"}, false),
			},
		},
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestQuery := url.Values{}
	filterType, ok := d.GetOk("type")
	if ok {
		requestQuery["type"] = []string{filterType.(string)}
	}
	requestPath := fmt.Sprintf(client.WorkspacePath)
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.WorkspaceCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Check for a quick exit
	if len(retVals.Data) == 0 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No workspace exists with that filter criteria"))
	}
	//Do manual searching
	filteredList := []client.Workspace{}
	searchName, ok := d.GetOk("search_name")
	if ok {
		searchNameLower := strings.ToLower(searchName.(string))
		for _, w := range retVals.Data {
			if strings.Contains(strings.ToLower(w.Name), searchNameLower) {
				filteredList = append(filteredList, w)
			}
		}
		retVals.Data = filteredList
		filteredList = []client.Workspace{}
	}
	name, ok := d.GetOk("name")
	if ok {
		nameLower := strings.ToLower(name.(string))
		for _, w := range retVals.Data {
			if strings.ToLower(w.Name) == nameLower {
				filteredList = append(filteredList, w)
			}
		}
		retVals.Data = filteredList
		filteredList = []client.Workspace{}
	}
	numWorkspaces := len(retVals.Data)
	if numWorkspaces > 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single workspace"))
	} else if numWorkspaces != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No workspace exists with that filter criteria"))
	}
	retVal := retVals.Data[0]
	d.Set("name", retVal.Name)
	d.Set("type", retVal.Type)
	d.SetId(retVal.Id)
	return diags
}
