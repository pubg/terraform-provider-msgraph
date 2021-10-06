---
layout: "msgraph"
subcategory: "ServicePrincipal"
page_title: "MsGraph: msgraph_app_role_assignment"
description: |-
  Get AWS CloudTrail Service Account ID for storing trail data in S3.
---

# Resource: msgraph_app_role_assignment

Assign Subscription's role to ServicePrincipal

## Example Usage

```terraform
resource "msgraph_app_role_assignment" "example" {
  for_each = toset(data.msgraph_groups.example.group_ids)

  # User or Group Id
   principal_id = each.key

   # Enterprise Application Id
   resource_id = "<uuid>"

   # Application Role Id
   app_role_id = "<uuid>"
}
```

## Arguments Reference

## Attributes Reference

## Import

Not support Terraform import
