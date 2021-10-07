package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
	"terraform-provider-msgraph/internal/provider"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: ProviderFunc}

	if debugMode {
		// TODO: update this string with the full name of your provider as used in your configs
		err := plugin.Debug(context.Background(), "pubg/msgraph", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}

func ProviderFunc() *schema.Provider {
	return provider.Provider()
}
