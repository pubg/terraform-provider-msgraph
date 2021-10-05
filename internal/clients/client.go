package clients

import (
	"context"
	"fmt"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"terraform-provider-msgraph/internal/helpers/hamilton_jwt"

	"github.com/Azure/go-autorest/autorest"

	"terraform-provider-msgraph/internal/common"
	approleassignment "terraform-provider-msgraph/internal/services/approleassignment/client"
)

// Client contains the handles to all the specific Azure AD resource classes' respective clients
type Client struct {
	Environment environments.Environment
	TenantID    string
	ClientID    string
	ObjectID    string
	Claims      auth.Claims

	TerraformVersion string

	AuthenticatedAsAServicePrincipal bool
	EnableMsGraphBeta                bool // TODO: remove in v2.0

	StopContext context.Context

	AppRoleAssignment *approleassignment.Client
}

func (client *Client) build(ctx context.Context, o *common.ClientOptions) error { //nolint:unparam
	autorest.Count429AsRetry = false
	client.StopContext = ctx

	client.AppRoleAssignment = approleassignment.NewClient(o)

	if client.EnableMsGraphBeta {
		// Acquire an access token upfront so we can decode and populate the JWT claims
		token, err := o.MsGraphAuthorizer.Token()
		if err != nil {
			return fmt.Errorf("unable to obtain access token: %v", err)
		}
		client.Claims, err = hamilton_jwt.ParseClaims(token)
		if err != nil {
			return fmt.Errorf("unable to parse claims in access token: %v", err)
		}
		if client.Claims.ObjectId == "" {
			return fmt.Errorf("parsing claims in access token: oid claim is empty")
		}
	}

	return nil
}
