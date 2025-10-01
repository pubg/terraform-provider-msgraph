package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	"github.com/pubg/terraform-provider-msgraph/version"
)

type ClientOptions struct {
	Environment environments.Environment

	PartnerID        string
	TerraformVersion string

	MsGraphAuthorizer auth.Authorizer // TODO: rename in v2.0
}

func (o ClientOptions) ConfigureMsGraphClient(c *msgraph.Client) {
	if o.MsGraphAuthorizer != nil {
		c.Authorizer = o.MsGraphAuthorizer
		c.Endpoint = o.Environment.MsGraph.Endpoint
		c.UserAgent = o.UserAgent(c.UserAgent)
	}
}

func (o ClientOptions) UserAgent(sdkUserAgent string) (userAgent string) {
	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", o.TerraformVersion, meta.SDKVersionString())
	providerUserAgent := fmt.Sprintf("%s terraform-provider-azuread/%s", tfUserAgent, version.ProviderVersion)
	userAgent = strings.TrimSpace(fmt.Sprintf("%s %s", sdkUserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, azureAgent)
	}

	if o.PartnerID != "" {
		userAgent = fmt.Sprintf("%s pid-%s", userAgent, o.PartnerID)
	}

	return
}
