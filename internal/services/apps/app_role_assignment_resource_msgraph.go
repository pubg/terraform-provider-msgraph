package apps

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/pubg/terraform-provider-msgraph/internal/helpers/hamilton_helper"

	"github.com/pubg/terraform-provider-msgraph/internal/clients"
	helpers "github.com/pubg/terraform-provider-msgraph/internal/helpers/msgraph"
	"github.com/pubg/terraform-provider-msgraph/internal/tf"
	"github.com/pubg/terraform-provider-msgraph/internal/utils"
)

var duplicatedPrefix = "duplicated_"

func appRoleAssignmentResourceCreateMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipalClient.MsGraphClient

	properties := msgraph.AppRoleAssignment{
		AppRoleId:   utils.String(d.Get("app_role_id").(string)),
		PrincipalId: utils.String(d.Get("principal_id").(string)),
		ResourceId:  utils.String(d.Get("resource_id").(string)),
	}

	appRoleAssignment, _, err := client.AssignAppRoleForResource(ctx, *properties.PrincipalId, *properties.ResourceId, *properties.AppRoleId)
	if err != nil {
		if d.Get("tolerance_duplicate").(bool) {
			if strings.Contains(err.Error(), "Permission being assigned already exists on the object") {
				d.SetId(duplicatedPrefix + utils.RandStringBytes(32))
				d.Set("app_role_id", properties.AppRoleId)
				d.Set("principal_id", properties.PrincipalId)
				d.Set("resource_id", properties.ResourceId)
				return nil
			}
		}
		return tf.ErrorDiagF(err, "Assigning role %+v", properties)
	}

	d.SetId(*appRoleAssignment.Id)
	_, err = helpers.WaitForCreationReplication(ctx, func() (interface{}, int, error) {
		pGroupRole, status, err := hamilton_helper.GetAppRoleAssignment(client, ctx, *properties.ResourceId, d.Id(), odata.Query{})
		if err != nil {
			return nil, status, err
		}
		return pGroupRole, status, nil
	})
	if err != nil {
		return tf.ErrorDiagF(err, "Waiting for AppRoleAssignment with object ID: %q", *appRoleAssignment.Id)
	}

	return appRoleAssignmentResourceReadMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceUpdateMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func appRoleAssignmentResourceReadMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("tolerance_duplicate").(bool) {
		if strings.HasPrefix(d.Id(), duplicatedPrefix) {
			return nil
		}
	}

	client := meta.(*clients.Client).ServicePrincipalClient.MsGraphClient

	resourceId := utils.String(d.Get("resource_id").(string))

	getAppRoleAssignment := func() (msgraph.AppRoleAssignment, error) {
		pGroupRole, _, err := hamilton_helper.GetAppRoleAssignment(client, ctx, *resourceId, d.Id(), odata.Query{})
		if err != nil {
			return msgraph.AppRoleAssignment{}, err
		}
		return *pGroupRole, nil
	}

	appRoleAssignment, err := getAppRoleAssignment()

	if !d.IsNewResource() && err != nil {
		log.Printf("[WARN] App Role Assignment (%q) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return tf.ErrorDiagF(err, "Could not find AppRoleAssignment: %q", d.Id())
	}

	tf.Set(d, "id", appRoleAssignment.Id)
	tf.Set(d, "app_role_id", appRoleAssignment.AppRoleId)
	tf.Set(d, "principal_id", appRoleAssignment.PrincipalId)
	tf.Set(d, "resource_id", appRoleAssignment.ResourceId)
	return nil
}

func appRoleAssignmentResourceDeleteMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("tolerance_duplicate").(bool) {
		if strings.HasPrefix(d.Id(), duplicatedPrefix) {
			return nil
		}
	}

	client := meta.(*clients.Client).ServicePrincipalClient.MsGraphClient

	resourceId := utils.String(d.Get("resource_id").(string))
	if _, err := client.RemoveAppRoleAssignment(ctx, *resourceId, d.Id()); err != nil {
		return tf.ErrorDiagF(err, "Deleting app role assignment with object ID: %q", d.Id())
	}

	return nil
}
