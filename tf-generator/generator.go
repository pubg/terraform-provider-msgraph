package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"terraform-provider-msgraph/hamilton/auth"
)

func replaceSpace(original string, replace string) string {
	return strings.Replace(original, " ", replace, -1)
}

func generateOutputFilepath(filepath string) string {
	prefix := strings.TrimSuffix(filepath, "json")
	return prefix + "tf"
}

func generateAppRolesTF(filepath string, authorizer auth.Authorizer, ctx context.Context) error {
	tfPath := generateOutputFilepath(filepath)
	jsonBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var appRoles *AppRoles
	err = json.Unmarshal(jsonBytes, &appRoles)
	if err != nil {
		return err
	}

	providerResourceName := strings.Replace(strings.TrimSuffix(filepath, ".json"), ".", "_", -1)
	appDisplayNameNoSpace := replaceSpace(appRoles.AppDisplayName, "_")
	appID, _, err := getAppInfo(appRoles.AppDisplayName, authorizer, ctx)

	if err != nil {
		return err
	}

	locals := fmt.Sprintf(LocalStart, providerResourceName)
	for _, appRole := range appRoles.AppRoles {
		appRoleNameNoSpace := replaceSpace(*appRole.DisplayName, "_")
		locals += fmt.Sprintf(
			RoleLocalTemplate,
			appDisplayNameNoSpace+"-"+appRoleNameNoSpace,
			appID,
			*appRole.AllowedMemberTypes,
			*appRole.DisplayName,
			*appRole.Description,
			*appRole.Enabled,
			*appRole.DisplayName)
	}
	locals += LocalEnd

	resource := fmt.Sprintf(RolesResourceTemplate, providerResourceName, providerResourceName)

	output := locals + resource
	bOutput := []byte(output)
	err = ioutil.WriteFile(tfPath, bOutput, 0644)
	if err != nil {
		return err
	}
	return nil
}

func generateAppRoleAssignmentsTF(filepath string, authorizer auth.Authorizer, ctx context.Context) error {
	tfPath := generateOutputFilepath(filepath)
	jsonBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var appAssignments *AppRoleAssignments
	err = json.Unmarshal(jsonBytes, &appAssignments)
	if err != nil {
		return err
	}

	providerResourceName := strings.Replace(strings.TrimSuffix(filepath, ".json"), ".", "_", -1)
	resourceDisplayName := *appAssignments.ResourceDisplayName
	resourceID, err := getServicePrincipalID(resourceDisplayName, authorizer, ctx)
	if err != nil {
		return err
	}
	_, appRoles, err := getAppInfo(resourceDisplayName, authorizer, ctx)
	if err != nil {
		return err
	}

	locals := fmt.Sprintf(LocalStart, providerResourceName)
	safeLocals := make(map[string]string)

	for _, appAssignment := range *appAssignments.AppRoleAssignments {
		assignedAppRoles := *appAssignment.AppRoles

		principalType := *appAssignment.PrincipalType
		err := checkPrincipalType(principalType)
		if err != nil {
			return err
		}

		relatedPrincipals := make(map[string]string)
		if principalType == "Group" {
			relatedPrincipals, err = getNestedGroups(*appAssignment.PrincipalDisplayName, authorizer, ctx)
		} else if principalType == "User" {
			relatedPrincipals, err = getUserInfo(*appAssignment.PrincipalDisplayName, authorizer, ctx)
		}
		if err != nil {
			return err
		}

		for _, appRoleDisplayName := range assignedAppRoles {
			for principalID, principalDisplayName := range relatedPrincipals {
				appRoleDisplayNameNoSpace := replaceSpace(appRoleDisplayName, "_")
				principalDisplayNameNoSpace := replaceSpace(principalDisplayName, "_")
				appRoleID := appRoles[appRoleDisplayName]["id"]
				safeKey := principalID + "_" + appRoleDisplayNameNoSpace
				safeLocals[safeKey] = fmt.Sprintf(
					AssignmentLocalTemplate,
					principalDisplayNameNoSpace+": "+appRoleDisplayName,
					safeKey,
					principalDisplayName,
					principalID,
					*appAssignment.PrincipalType,
					resourceDisplayName,
					resourceID,
					appRoleDisplayName,
					appRoleID)

			}
		}
	}

	for _, local := range safeLocals {
		locals += local
	}
	locals += LocalEnd

	resource := fmt.Sprintf(AssignmentsResourceTemplate, providerResourceName, providerResourceName)

	output := locals + resource
	bOutput := []byte(output)
	err = ioutil.WriteFile(tfPath, bOutput, 0644)
	if err != nil {
		return err
	}
	return nil
}
