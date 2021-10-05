package approleassignment

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-msgraph/internal/tf"
	"terraform-provider-msgraph/internal/validate"
)

func appRoleAssignmentResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: appRoleAssignmentResourceCreateUpdate,
		UpdateContext: appRoleAssignmentResourceCreateUpdate,
		ReadContext:   appRoleAssignmentResourceRead,
		DeleteContext: appRoleAssignmentResourceDelete,

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			if _, err := uuid.ParseUUID(id); err != nil {
				return fmt.Errorf("specified ID (%q) is not valid: %s", id, err)
			}
			return nil
		}),

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_role_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"principal_display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"principal_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"principal_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"resource_display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func appRoleAssignmentResourceCreateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceCreateUpdateMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceReadMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceDeleteMsGraph(ctx, d, meta)
}
