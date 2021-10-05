package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"terraform-provider-msgraph/hamilton/auth"
	"terraform-provider-msgraph/hamilton/environments"
)

var (
	tenantID     = os.Getenv("TENANT_ID")
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

const (
	defaultFilePath = "./cyk.prod.roles.json"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "filepath", defaultFilePath, "File Path")
	flag.Parse()

	err := checkFilepathFormat(filepath)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	authConfig := &auth.Config{
		Environment:         environments.Global,
		TenantID:            tenantID,
		EnableAzureCliToken: true,
		ClientID:            clientID,
		ClientSecret:        clientSecret,
	}

	authorizer, err := authConfig.NewAuthorizer(ctx, auth.MsGraph)
	if err != nil {
		log.Fatal(err)
	}

	resourceType := strings.Split(filepath, ".")[2]
	if resourceType == "roles" {
		err = generateAppRolesTF(filepath, authorizer, ctx)
	} else {
		err = generateAppRoleAssignmentsTF(filepath, authorizer, ctx)
	}
	if err != nil {
		log.Fatal(err)
	}
}
