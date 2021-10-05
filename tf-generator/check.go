package main

import (
	"errors"
	"fmt"
	"strings"
)

func checkFilepathFormat(filepath string) error {
	if strings.HasSuffix(filepath, "json") == false {
		return errors.New("wrong file format: should be `json`")
	}

	filepaths := strings.Split(filepath, "/")
	filename := filepaths[len(filepaths)-1]
	array := strings.Split(filename, ".")
	if len(array) != 4 {
		return errors.New("wrong filename format: should be `{app}.{env}.{resource}.json`")
	}

	resourceType := array[2]
	if resourceType != "roles" && resourceType != "assignments" {
		return errors.New("unsupported resource type: should be either `roles` or `assignments`")
	}

	return nil
}

func checkPrincipalType(principalType string) error {
	switch principalType {
	case "User", "Group":
		return nil
	default:
		return errors.New(fmt.Sprintf("unsupported principal type %q: should be either `User` or `Group`", principalType))
	}
}
