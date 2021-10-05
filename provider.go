package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-msgraph/internal/provider"
)

func Provider() *schema.Provider {
	return provider.Provider()
}
