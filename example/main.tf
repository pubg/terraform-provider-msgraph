terraform {
  required_providers {
    msgraph = {
      source = "pubg/msgraph"
    }
  }
}

provider "msgraph" {
  use_cli = true
}

data "msgraph_groups" "groups" {
  group_id = "d2ec52ab-f8ec-463a-ba99-7561719b984b"
}

output "group_ids" {
  value = data.msgraph_groups.groups.group_ids
}

output "user_ids" {
  value = data.msgraph_groups.groups.user_ids
}
