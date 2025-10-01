package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"

	"github.com/pubg/terraform-provider-msgraph/internal/common"
)

type ClientBuilder struct {
	AuthConfig           *auth.Config
	EnableMsGraph        bool
	PartnerID            string
	TerraformVersion     string
	EnableClientsideLock bool
}

// Build is a helper method which returns a fully instantiated *Client based on the auth Config's current settings.
func (b *ClientBuilder) Build(ctx context.Context) (*Client, error) {
	// client declarations:
	client := Client{
		TerraformVersion: b.TerraformVersion,
	}

	if b.EnableClientsideLock {
		client.EnableResourceMutex = true
		client.ResourceMutex = &sync.Mutex{}
	}

	if b.AuthConfig != nil {
		client.Environment = b.AuthConfig.Environment
	}

	// MS Graph
	if b.AuthConfig == nil {
		return nil, fmt.Errorf("building client: AuthConfig is nil")
	}

	client.EnableMsGraphBeta = true
	msGraphAuthorizer, err := b.AuthConfig.NewAuthorizer(ctx, environments.MsGraphGlobal)
	if err != nil {
		return nil, err
	}

	// Obtain the tenant ID from Azure CLI
	if cli, ok := msGraphAuthorizer.(*auth.AzureCliAuthorizer); ok {
		if cli.TenantID == "" {
			return nil, fmt.Errorf("azure-cli could not determine tenant ID to use")
		}
	}

	o := &common.ClientOptions{
		Environment:       client.Environment,
		PartnerID:         b.PartnerID,
		TerraformVersion:  client.TerraformVersion,
		MsGraphAuthorizer: msGraphAuthorizer,
	}
	if err := client.build(ctx, o); err != nil {
		return nil, fmt.Errorf("building client: %+v", err)
	}

	return &client, nil
}
