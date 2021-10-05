package main

const LocalStart = `locals {
  %s = {`

const LocalEnd = `  }
}
`

const RoleLocalTemplate = `
    "%s" = {
      application_object_id = %q
      allowed_member_types  = %q
      display_name          = %q
      description           = %q
      enabled               = %t
      value                 = %q
    }
`

const RolesResourceTemplate = `
resource "azuread_application_app_role" %q {
  for_each = {
    for app_role in local.%s : app_role.display_name => {
      application_object_id = app_role.application_object_id
      allowed_member_types  = app_role.allowed_member_types
      description           = app_role.description
      display_name          = app_role.display_name
      enabled               = app_role.enabled
      value                 = app_role.value
    }
  }

  application_object_id = each.value.application_object_id
  allowed_member_types  = each.value.allowed_member_types
  description           = each.value.description
  display_name          = each.value.display_name
  enabled               = each.value.enabled
  value                 = each.value.value
}
`

const AssignmentLocalTemplate = `
    "%s" = {
      principal_role         = %q
      principal_display_name = %q
      principal_id           = %q
      principal_type         = %q
      resource_display_name  = %q
      resource_id            = %q
      app_role_display_name  = %q
      app_role_id            = %q
    }
`

const AssignmentsResourceTemplate = `
resource "msgraph_app_role_assignment" %q {
  for_each = {
    for assignment in local.%s : assignment.principal_role => {
      principal_display_name = assignment.principal_display_name
      principal_id           = assignment.principal_id
      principal_type         = assignment.principal_type
      resource_display_name  = assignment.resource_display_name
      resource_id            = assignment.resource_id
      app_role_id            = assignment.app_role_id
    }
  }

  principal_display_name = each.value.principal_display_name
  principal_id           = each.value.principal_id
  principal_type         = each.value.principal_type
  resource_display_name  = each.value.resource_display_name
  resource_id            = each.value.resource_id
  app_role_id            = each.value.app_role_id
}
`
