package msgraph

import "github.com/manicminer/hamilton/msgraph"

func ApplicationFlattenAppRoles(in *[]msgraph.AppRole) map[string]map[string]interface{} {
	if in == nil {
		return map[string]map[string]interface{}{}
	}

	appRoles := make(map[string]map[string]interface{}, 0)
	for _, role := range *in {
		roleId := ""
		if role.ID != nil {
			roleId = *role.ID
		}
		allowedMemberTypes := make([]interface{}, 0)
		if v := role.AllowedMemberTypes; v != nil {
			for _, m := range *v {
				allowedMemberTypes = append(allowedMemberTypes, m)
			}
		}
		description := ""
		if role.Description != nil {
			description = *role.Description
		}
		displayName := ""
		if role.DisplayName != nil {
			displayName = *role.DisplayName
		}
		enabled := false
		if role.IsEnabled != nil && *role.IsEnabled {
			enabled = true
		}
		value := ""
		if role.Value != nil {
			value = *role.Value
		}
		appRoles[displayName] = map[string]interface{}{
			"id":                   roleId,
			"allowed_member_types": allowedMemberTypes,
			"description":          description,
			"display_name":         displayName,
			"enabled":              enabled,
			"value":                value,
		}
	}

	return appRoles
}
