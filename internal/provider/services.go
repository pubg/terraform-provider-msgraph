package provider

import (
	"github.com/pubg/terraform-provider-msgraph/internal/services/approleassignment"
	"github.com/pubg/terraform-provider-msgraph/internal/services/groups"
)

func SupportedServices() []ServiceRegistration {
	return []ServiceRegistration{
		approleassignment.Registration{},
		groups.Registration{},
	}
}
