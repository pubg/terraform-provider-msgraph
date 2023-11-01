---
layout: "msgraph"
page_title: "Provider: MsGraph"
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

Same as AzureAD Provider

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
    }
  }
}

# Access via ClientSecret 
data "vault_generic_secret" "azure_credential_root" {
  path = "secret/serviceprincipal/root"
}
        
provider "azuread" {
  alias = "tenant1"        
        
  tenant_id     = data.vault_generic_secret.azure_credential_root.data["tenant_id"]
  client_id     = data.vault_generic_secret.azure_credential_root.data["client_id"]
  client_secret = data.vault_generic_secret.azure_credential_root.data["client_secret"]
}

provider "msgraph" {
  alias = "tenant1"
          
  tenant_id     = data.vault_generic_secret.azure_credential_root.data["tenant_id"]
  client_id     = data.vault_generic_secret.azure_credential_root.data["client_id"]
  client_secret = data.vault_generic_secret.azure_credential_root.data["client_secret"]
}

# Access via AzureCLI
provider "azuread" {
  alias = "tenant2"
  use_cli = true
}

provider "msgraph" {
  alias = "tenant2"
  use_cli = true
}

# Enable clideside lock for avoiding resource ownership conflict (default: false)
provider "msgraph" {
  use_cli = true
  use_clientside_lock = true
}
...
```
