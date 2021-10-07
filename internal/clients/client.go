package clients

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
	"terraform-provider-msgraph/internal/helpers/hamilton_helper"

	"github.com/Azure/go-autorest/autorest"

	"terraform-provider-msgraph/internal/common"
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

	ServicePrincipalClient struct {
		AadClient     *graphrbac.ServicePrincipalsClient
		MsGraphClient *msgraph.ServicePrincipalsClient
	}
	GroupsClient *msgraph.GroupsClient
}

func (client *Client) build(ctx context.Context, o *common.ClientOptions) error { //nolint:unparam
	autorest.Count429AsRetry = false
	client.StopContext = ctx

	spAadClient := graphrbac.NewServicePrincipalsClientWithBaseURI(o.AadGraphEndpoint, o.TenantID)
	client.ServicePrincipalClient.AadClient = &spAadClient
	client.ServicePrincipalClient.MsGraphClient = msgraph.NewServicePrincipalsClient(o.TenantID)
	client.GroupsClient = msgraph.NewGroupsClient(o.TenantID)

	o.ConfigureAadClient(&client.ServicePrincipalClient.AadClient.Client)
	o.ConfigureMsGraphClient(&client.ServicePrincipalClient.MsGraphClient.BaseClient)
	o.ConfigureMsGraphClient(&client.GroupsClient.BaseClient)

	if client.EnableMsGraphBeta {
		// Acquire an access token upfront so we can decode and populate the JWT claims
		token, err := o.MsGraphAuthorizer.Token()
		if err != nil {
			return fmt.Errorf("unable to obtain access token: %v", err)
		}
		client.Claims, err = hamilton_helper.ParseClaims(token)
		if err != nil {
			return fmt.Errorf("unable to parse claims in access token: %v", err)
		}
		if client.Claims.ObjectId == "" {
			return fmt.Errorf("parsing claims in access token: oid claim is empty")
		}
	}
	return nil
}
