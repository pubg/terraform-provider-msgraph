---
layout: "msgraph"
page_title: "Provider: MsGraph"
description: |-
  Use the Amazon Web Services (AWS) provider to interact with the many resources supported by AWS. You must configure the provider with the proper credentials before you can use it.
---

# terraform-provider-msgraph
Includes the following
1. Extends `hashicorp/azuread` provider
2. Use msgraph 2.0 beta API

## Caution
This provider using Microsoft Graph API `beta` version.

Microsoft don't recommend that you use them in your production apps. [MS Doc Here](https://docs.microsoft.com/ko-kr/graph/api/overview?view=graph-rest-beta)

But, This API is the only way to managing AzureAD apps with GitOps style.

### Set up Custom Terraform Provider
1. Add the providers. For example,
```terraform
terraform {
  ...
  required_providers {
    ...

    azuread = {
      source  = "hashicorp/azuread"
      version = "1.6.0"
    }

    msgraph = {
      source = "pubg/msgraph"
      version = "0.0.3"
    }
  }
}

provider "azuread" {
  use_microsoft_graph = true
}

provider "msgraph" {
  use_microsoft_graph = true
}
```
