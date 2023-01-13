package provider

import (
	"github.com/pubg/terraform-provider-msgraph/internal/services/apps"
	"github.com/pubg/terraform-provider-msgraph/internal/services/groups"
)

func SupportedServices() []ServiceRegistration {
	return []ServiceRegistration{
		apps.Registration{},
		groups.Registration{},
	}
}
