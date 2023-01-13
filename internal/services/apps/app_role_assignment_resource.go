package apps

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pubg/terraform-provider-msgraph/internal/tf"
	"github.com/pubg/terraform-provider-msgraph/internal/validate"
)

func appRoleAssignmentResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: appRoleAssignmentResourceCreate,
		UpdateContext: appRoleAssignmentResourceUpdate,
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"principal_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"tolerance_duplicate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func appRoleAssignmentResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceCreateMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceUpdateMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceReadMsGraph(ctx, d, meta)
}

func appRoleAssignmentResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return appRoleAssignmentResourceDeleteMsGraph(ctx, d, meta)
}
