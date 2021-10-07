# terraform-provider-msgraph
Includes the following
1. Custom Terraform Provider `msgraph` for Azure AD app role assignments

## Caution
This provider using Microsoft Graph API `beta` version.

Microsoft don't recommend that you use them in your production apps. [MS Doc Here](https://docs.microsoft.com/ko-kr/graph/api/overview?view=graph-rest-beta)

But, This API is the only way to managing AzureAD apps with GitOps style.

## How to
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
2. Run `terraform init`

### Define Resources
```terraform
data "msgraph_groups" "example" {
  # Group Id
  group_id             = "<uuid>"
  listup_nested_groups = true
}

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

### Plan & Apply
### How to Build in local
```shell
# Build your OS dependent binary
make build

# Build all OS envs
goreleaser release --rm-dist --snapshot
```

### Debug in local PC
1. Check this document. https://www.terraform.io/docs/extend/debugging.html#starting-a-provider-in-debug-mode
2. Run Go Code with IDE Debug Mode with `--debug=true` option.
3. Set BreakPoints and debug this provider line by line.

### The End
