---
layout: "msgraph"
subcategory: "Application"
page_title: "MsGraph: msgraph_app_redirect_uris"
description: |-
  Add redirect uris to Application
---

# Resource: msgraph_app_redirect_uris

Add redirect uris to Application

## Example Usage

```terraform
resource "azuread_application" "shared_app" {
   display_name = "Shared App"

  web {
    redirect_uris = [
      "https://*.contoso.com/",
    ]
  }
}

resource "msgraph_app_redirect_uris" "per_cluster_app" {
   app_object_id = azuread_application.shared_app.object_id
  
   redirect_uris {
     url = "https://*.my-first-cluster.contoso.com/"
     type = "Web"
   }
}
```

## Arguments Reference

* `app_object_id` - (Required, uuid) The object_id of target Application 
* `redirect_uris` - A `redirect_uris` block as documented blow, which configures redirect uri related settings for target application.
* `tolerance_override` - (Optional, bool) If some urls are already exist in target application, It may occur resource ownership conflict. If you want ignore this error, enable `tolerance_override` to true (default false)

---

`redirect_uris` block supports the following:

* `url` - (Required, string) A reply URL, You can find restrictions and limitations are this document. https://learn.microsoft.com/en-us/azure/active-directory/develop/reply-url
* `type` - (Required, enum) Type of Redirect Url, One of [Web, InstalledClient, Spa]

---

## Attributes Reference

* `id` - App redirect uris resource Id

## Import

Not support Terraform import
