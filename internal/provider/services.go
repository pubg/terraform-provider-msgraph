package provider

import (
	"terraform-provider-msgraph/internal/services/approleassignment"
	"terraform-provider-msgraph/internal/services/data_source_msgraph_groups"
)

func SupportedServices() []ServiceRegistration {
	return []ServiceRegistration{
		approleassignment.Registration{},
		data_source_msgraph_groups.Registration{},
	}
}
