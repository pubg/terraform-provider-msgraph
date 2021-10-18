package provider

import (
	"terraform-provider-msgraph/internal/services/approleassignment"
	"terraform-provider-msgraph/internal/services/groups"
)

func SupportedServices() []ServiceRegistration {
	return []ServiceRegistration{
		approleassignment.Registration{},
		groups.Registration{},
	}
}
