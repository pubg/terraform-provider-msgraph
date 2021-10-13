---
layout: "msgraph"
subcategory: "Groups"
page_title: "MsGraph: msgraph_groups"
description: |-
  Get nested groups.
---

# Data Source: msgraph_groups

The data source can get nested groups of top group.

You can assign role to all groups belong to big organization or division.

## Example Usage

```terraform
data "msgraph_groups" "my_groups" {
  group_id             = "4729d0a8-2cea-446b-95fb-43c7e8973816"
  listup_nested_groups = true
}

resource "msgraph_app_role_assignment" "my_assign" {
  for_each = toset(data.msgraph_groups.my_groups.group_ids)

  principal_id = each.key
  resource_id = azuread_service_principal.my_app.object_id
  app_role_id = azuread_application_app_role.my_role.role_id
}
```

## Arguments Reference

* `group_id` - (Required) The Group's UUID.
* `listup_nested_groups` - (Required) Bool flag of search nested groups. 

## Attributes Reference

* `group_ids` - Type: String List, list of nested or single group ids.

