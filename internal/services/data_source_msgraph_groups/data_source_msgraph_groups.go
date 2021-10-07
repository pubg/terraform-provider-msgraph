package data_source_msgraph_groups

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-msgraph/internal/clients"
	"terraform-provider-msgraph/internal/helpers/hamilton_helper"
)

func dataSourceMsgraphGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listup_nested_groups": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"group_ids": {
				Type:     schema.TypeList,
				Optional: true,
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
	listupNestedGroups := d.Get("listup_nested_groups").(bool)

	var groupIds []string

	if listupNestedGroups {
		err := hamilton_helper.TraverseNestedGroups(client, ctx, groupId, func(groupId string, err error) error {
			if err != nil {
				return err
			}

			groupIds = append(groupIds, groupId)
			return nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		groupIds = append(groupIds, groupId)
	}

	d.SetId(groupId)
	err := d.Set("group_ids", groupIds)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}
