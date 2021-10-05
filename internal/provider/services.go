package provider

import (
	"terraform-provider-msgraph/internal/services/approleassignment"
)

func SupportedServices() []ServiceRegistration {
	return []ServiceRegistration{
		approleassignment.Registration{},
	}
}
