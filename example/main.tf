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

resource "msgraph_app_redirect_uris" "uris" {
  app_object_id = "c52d8e63-4117-4b84-8856-879967d31606"
  #  redirect_uris {
  #    url  = "https://contoso.com"
  #    type = "Web"
  #  }
  redirect_uris {
    url  = "https://contoso2.com"
    type = "Web"
  }
  redirect_uris {
    url  = "https://contoso3.com"
    type = "Web"
  }

  tolerance_override = true
}

resource "msgraph_app_redirect_uris" "uris2" {
  app_object_id = "c52d8e63-4117-4b84-8856-879967d31606"
  #  redirect_uris {
  #    url  = "https://contoso.com"
  #    type = "Web"
  #  }
  redirect_uris {
    url  = "https://contoso2.com"
    type = "Web"
  }
  redirect_uris {
    url  = "https://contoso4.com"
    type = "Web"
  }

  tolerance_override = true
}
