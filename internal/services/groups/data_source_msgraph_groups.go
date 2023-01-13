package groups

import (
	"context"
	"fmt"

	"github.com/pubg/terraform-provider-msgraph/internal/clients"
	"github.com/pubg/terraform-provider-msgraph/internal/helpers/hamilton_helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/odata"
)

func dataSourceMsgraphGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"max_traverse_depth": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: readGroups,
	}
}

func readGroups(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).GroupsClient

	groupId := d.Get("group_id").(string)
	traverseDepth := d.Get("max_traverse_depth").(int)

	var groupIds []string
	var userIds []string

	err := hamilton_helper.TraverseNestedGroups(client, ctx, traverseDepth, groupId, func(objectId string, oDataType odata.Type, err error) error {
		if err != nil {
			return err
		}
		if oDataType == odata.TypeGroup {
			groupIds = append(groupIds, objectId)
		} else if oDataType == odata.TypeUser {
			userIds = append(userIds, objectId)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s#%d", groupId, traverseDepth))

	err = d.Set("group_ids", groupIds)
	if err != nil {
		diag.FromErr(err)
	}
	err = d.Set("user_ids", userIds)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}
