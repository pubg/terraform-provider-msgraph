---
layout: "msgraph"
subcategory: "EnterpriseApplication"
page_title: "MsGraph: msgraph_app_role_assignment"
description: |-
  Assign user or groups to EnterpriseApplication
---

# Resource: msgraph_app_role_assignment

Assign user or groups to EnterpriseApplication

## Example Usage

```terraform
resource "msgraph_app_role_assignment" "example" {
  # User or Group Id
   principal_id = "<uuid>"

   # Enterprise Application Id
   resource_id = "<uuid>"

   # Application Role Id
   app_role_id = "<uuid>"

  tolerance_duplicate = true
}
```

## Arguments Reference

* `app_role_id` - (Required) The Application Role Id 
* `principal_id` - (Required) The User or Group Id
* `resource_id` - (Required) The Enterprise Application Id
* `tolerance_duplicate` - (Optional) Allow create same `msgraph_app_role_assignment`. When this resource detect duplicated, then do nothing to real world.

## Attributes Reference

* `id` - App Role Assignment Resource Id

## Import

Not support Terraform import
