package hamilton_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"io/ioutil"
	"net/http"
	"net/url"
)

func a() {
	b := msgraph.GroupsClient{}
	b.Get(nil, "", odata.Query{})
}

func ListNestedMemberGroups(client *msgraph.GroupsClient, ctx context.Context, id string, nestedGroups *map[string]bool) (*map[string]bool, int, error) {
	(*nestedGroups)[id] = true
	resp, status, _, err := client.BaseClient.Get(ctx, msgraph.GetHttpRequestInput{
		ConsistencyFailureFunc: msgraph.RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: msgraph.Uri{
			Entity:      fmt.Sprintf("/groups/%s/members", id),
			Params:      url.Values{"$select": []string{"id"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		Members []struct {
			Type string `json:"@odata.type"`
			Id   string `json:"id"`
		} `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	for _, v := range data.Members {
		if v.Type != "#microsoft.graph.group" {
			continue
		}
		if (*nestedGroups)[v.Id] {
			continue
		}
		ListNestedMemberGroups(client, ctx, v.Id, nestedGroups)
	}
	return nestedGroups, status, nil
}
