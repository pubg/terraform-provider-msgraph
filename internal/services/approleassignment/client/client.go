package client

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/manicminer/hamilton/msgraph"
	"terraform-provider-msgraph/internal/common"
)

type Client struct {
	AadClient *graphrbac.ServicePrincipalsClient
	MsClient  *msgraph.ServicePrincipalsClient
}

func NewClient(o *common.ClientOptions) *Client {
	aadClient := graphrbac.NewServicePrincipalsClientWithBaseURI(o.AadGraphEndpoint, o.TenantID)
	msClient := msgraph.NewServicePrincipalsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient, &aadClient.Client)

	return &Client{
		AadClient: &aadClient,
		MsClient:  msClient,
	}
}
