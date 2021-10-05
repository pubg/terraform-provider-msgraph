package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"terraform-provider-msgraph/internal/clients"
)

type TestResource interface {
	Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error)
}

type TestResourceVerifyingRemoved interface {
	TestResource
	Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error)
}
