package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	helpers "terraform-provider-msgraph/internal/helpers/msgraph"
	"terraform-provider-msgraph/tf-generator/hamilton_helper"
)

const filter string = "displayName eq '%s'"

func getAppInfo(appDisplayName string, authorizer auth.Authorizer, ctx context.Context) (string, map[string]map[string]interface{}, error) {
	client := msgraph.NewApplicationsClient(tenantID)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf(filter, appDisplayName)})
	if err != nil {
		return "", nil, err
	}

	switch {
	case result == nil || len(*result) == 0:
		return "", nil, errors.New(fmt.Sprintf("No applications found matching filter: %q", filter))
	case len(*result) > 1:
		return "", nil, errors.New(fmt.Sprintf("Found multiple applications matching filter: %q", filter))
	}

	app := &(*result)[0]
	appID := *app.ID

	return appID, helpers.ApplicationFlattenAppRoles(app.AppRoles), nil
}

func getUserInfo(userDisplayName string, authorizer auth.Authorizer, ctx context.Context) (map[string]string, error) {
	client := msgraph.NewUsersClient(tenantID)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf(filter, userDisplayName)})
	if err != nil {
		return nil, err
	}

	switch {
	case result == nil || len(*result) == 0:
		return nil, errors.New(fmt.Sprintf("No users found matching filter: %q", filter))
	case len(*result) > 1:
		return nil, errors.New(fmt.Sprintf("Found multiple users matching filter: %q", filter))
	}

	user := &(*result)[0]
	userID := *user.ID

	userInfo := map[string]string{
		userID: userDisplayName,
	}

	return userInfo, nil
}

func getNestedGroups(groupDisplayName string, authorizer auth.Authorizer, ctx context.Context) (map[string]string, error) {
	client := msgraph.NewGroupsClient(tenantID)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf(filter, groupDisplayName)})
	if err != nil {
		return nil, err
	}

	switch {
	case result == nil || len(*result) == 0:
		return nil, errors.New(fmt.Sprintf("No groups found matching filter: %q", filter))
	case len(*result) > 1:
		return nil, errors.New(fmt.Sprintf("Found multiple groups matching filter: %q", filter))
	}

	group := &(*result)[0]
	groupID := *group.ID

	groupInfo := make(map[string]string)
	nestedGroups := make(map[string]bool)
	_, _, err = hamilton_helper.ListNestedMemberGroups(client, ctx, groupID, &nestedGroups)
	if err != nil {
		return nil, err
	}

	for key := range nestedGroups {
		group, _, err := client.Get(ctx, key, odata.Query{})
		if err != nil {
			return nil, err
		}
		displayName := *group.DisplayName
		groupInfo[key] = displayName
	}

	return groupInfo, nil
}

func getServicePrincipalID(appDisplayName string, authorizer auth.Authorizer, ctx context.Context) (string, error) {
	client := msgraph.NewServicePrincipalsClient(tenantID)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf(filter, appDisplayName)})
	if err != nil {
		return "", err
	}

	switch {
	case result == nil || len(*result) == 0:
		return "", errors.New(fmt.Sprintf("No service principlas found matching filter: %q", filter))
	case len(*result) > 1:
		return "", errors.New(fmt.Sprintf("Found multiple service principals matching filter: %q", filter))
	}

	servicePrincipal := &(*result)[0]
	servicePrincipalID := *servicePrincipal.ID

	return servicePrincipalID, nil
}
