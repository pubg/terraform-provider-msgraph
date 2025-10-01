package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/pubg/terraform-provider-msgraph/internal/helpers/hamilton_helper"

	"github.com/Azure/go-autorest/autorest"

	"github.com/pubg/terraform-provider-msgraph/internal/common"
)

// Client contains the handles to all the specific Azure AD resource classes' respective clients
type Client struct {
	Environment environments.Environment
	Claims      auth.Claims

	TerraformVersion string

	EnableMsGraphBeta bool // TODO: remove in v2.0

	StopContext context.Context

	ServicePrincipalClient struct {
		MsGraphClient *msgraph.ServicePrincipalsClient
	}
	GroupsClient *msgraph.GroupsClient
	AppClient    *msgraph.ApplicationsClient

	EnableResourceMutex bool
	ResourceMutex       *sync.Mutex
}

func (client *Client) build(ctx context.Context, o *common.ClientOptions) error { //nolint:unparam
	autorest.Count429AsRetry = false
	client.StopContext = ctx

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
		if client.Claims.TenantId == "" {
			return fmt.Errorf("parsing claims in access token: tid claim is empty")
		}
	}

	client.ServicePrincipalClient.MsGraphClient = msgraph.NewServicePrincipalsClient(client.Claims.TenantId)
	client.GroupsClient = msgraph.NewGroupsClient(client.Claims.TenantId)
	client.AppClient = msgraph.NewApplicationsClient(client.Claims.TenantId)

	o.ConfigureMsGraphClient(&client.ServicePrincipalClient.MsGraphClient.BaseClient)
	o.ConfigureMsGraphClient(&client.GroupsClient.BaseClient)
	o.ConfigureMsGraphClient(&client.AppClient.BaseClient)

	return nil
}
