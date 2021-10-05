package approleassignment

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"log"

	"terraform-provider-msgraph/internal/clients"
	helpers "terraform-provider-msgraph/internal/helpers/msgraph"
	"terraform-provider-msgraph/internal/tf"
	"terraform-provider-msgraph/internal/utils"
)

const appRoleAssignmentResourceName = "msgraph_app_role_assignment"

func appRoleAssignmentResourceCreateUpdateMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AppRoleAssignment.MsClient

	properties := msgraph.AppRoleAssignment{
		AppRoleId:           utils.String(d.Get("app_role_id").(string)),
		PrincipalId:         utils.String(d.Get("principal_id").(string)),
		PrincipalType:       utils.String(d.Get("principal_type").(string)),
		ResourceDisplayName: utils.String(d.Get("resource_display_name").(string)),
		ResourceId:          utils.String(d.Get("resource_id").(string)),
	}

	appRoleAssignment, _, err := client.AssignAppRoleForResource(ctx, *properties.PrincipalId, *properties.ResourceId, *properties.AppRoleId)
	if err != nil {
		return tf.ErrorDiagF(err, "Assigning role %s to %s %s", *properties.AppRoleId, *properties.PrincipalType, *properties.PrincipalDisplayName)
	}

	d.SetId(*appRoleAssignment.Id)

	_, err = helpers.WaitForCreationReplication(ctx, func() (interface{}, int, error) {
		pGroupRoles, status, err := client.ListAppRoleAssignments(ctx, *properties.ResourceId, odata.Query{})
		if err != nil {
			return nil, status, err
		}
		for _, value := range *pGroupRoles {
			if *value.Id == *appRoleAssignment.Id {
				return appRoleAssignment, status, nil
			}
		}
		return nil, 404, errors.New("App role assignment not yet created")
	})

	if err != nil {
		return tf.ErrorDiagF(err, "Waiting for AppRoleAssignment with object ID: %q", *appRoleAssignment.Id)
	}

	return appRoleAssignmentResourceReadMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceReadMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AppRoleAssignment.MsClient

	resourceId := utils.String(d.Get("resource_id").(string))

	getAppRoleAssignment := func() (msgraph.AppRoleAssignment, error) {
		pGroupRoles, _, err := client.ListAppRoleAssignments(ctx, *resourceId, odata.Query{})
		if err != nil {
			return msgraph.AppRoleAssignment{}, err
		}

		for _, value := range *pGroupRoles {
			if *value.Id == d.Id() {
				return value, nil
			}
		}

		return msgraph.AppRoleAssignment{}, helpers.ErrNoSuchAssignment()
	}

	appRoleAssignment, err := getAppRoleAssignment()

	if !d.IsNewResource() && errors.Is(err, helpers.ErrNoSuchAssignment()) {
		log.Printf("[WARN] App Role Assignment (%q) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return tf.ErrorDiagF(err, "Could not find AppRoleAssignment: %q", d.Id())
	}

	tf.Set(d, "id", appRoleAssignment.Id)
	tf.Set(d, "app_role_id", appRoleAssignment.AppRoleId)
	tf.Set(d, "principal_display_name", appRoleAssignment.PrincipalDisplayName)
	tf.Set(d, "principal_type", appRoleAssignment.PrincipalType)
	tf.Set(d, "principal_id", appRoleAssignment.PrincipalId)
	tf.Set(d, "resource_id", appRoleAssignment.ResourceId)
	tf.Set(d, "resource_display_name", appRoleAssignment.ResourceDisplayName)

	return nil
}

func appRoleAssignmentResourceDeleteMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AppRoleAssignment.MsClient

	resourceId := utils.String(d.Get("resource_id").(string))
	if _, err := client.RemoveAppRoleAssignment(ctx, *resourceId, d.Id()); err != nil {
		return tf.ErrorDiagF(err, "Deleting app role assignment with object ID: %q", d.Id())
	}

	return nil
}
