package hamilton_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manicminer/hamilton/msgraph"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TraverseFunc func(groupId string, err error) error

func TraverseNestedGroups(client *msgraph.GroupsClient, ctx context.Context, groupId string, traverseFunc TraverseFunc) error {
	return traverseNestedGroups0(client, ctx, groupId, traverseFunc, map[string]bool{})
}

func traverseNestedGroups0(client *msgraph.GroupsClient, ctx context.Context, groupId string, traverseFunc TraverseFunc, visitCheck map[string]bool) error {
	visitCheck[groupId] = true

	resp, status, _, err := client.BaseClient.Get(ctx, msgraph.GetHttpRequestInput{
		ConsistencyFailureFunc: msgraph.RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: msgraph.Uri{
			Entity:      fmt.Sprintf("/groups/%s/members", groupId),
			Params:      url.Values{"$select": []string{"id"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		err = traverseFunc(groupId, err)
		return fmt.Errorf("GroupsClient.BaseClient.Get(): statusCode: %d, %+v", status, err)
	}
	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = traverseFunc(groupId, err)
		return fmt.Errorf("ioutil.ReadAll(): %v", err)
	}

	var data struct {
		Members []struct {
			Type string `json:"@odata.type"`
			Id   string `json:"id"`
		} `json:"value"`
	}
	if err = json.Unmarshal(rawBody, &data); err != nil {
		err = traverseFunc(groupId, err)
		return fmt.Errorf("json.Unmarshal(): %v", err)
	}

	err = traverseFunc(groupId, nil)
	if err != nil {
		return err
	}

	for _, v := range data.Members {
		if v.Type != "#microsoft.graph.group" {
			continue
		}
		if visitCheck[v.Id] {
			continue
		}
		err = traverseNestedGroups0(client, ctx, v.Id, traverseFunc, visitCheck)
		if err != nil {
			return err
		}
	}
	return nil
}
