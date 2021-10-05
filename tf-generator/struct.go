package main

type AppRoles struct {
	AppDisplayName string    `json:"app_display_name"`
	AppRoles       []AppRole `json:"app_roles"`
}

type AppRole struct {
	DisplayName        *string   `json:"display_name"`
	Description        *string   `json:"description"`
	AllowedMemberTypes *[]string `json:"allowed_member_types"`
	Enabled            *bool     `json:"enabled"`
}

type AppRoleAssignments struct {
	AppRoleAssignments  *[]AppRoleAssignment `json:"app_role_assignments,omitempty"`
	ResourceDisplayName *string              `json:"name"`
}

type AppRoleAssignment struct {
	PrincipalDisplayName *string   `json:"display_name"`
	PrincipalType        *string   `json:"type"`
	AppRoles             *[]string `json:"roles"`
}
