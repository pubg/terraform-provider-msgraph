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

type TraverseFunc func(objectId string, oDataType odata.Type, err error) error

type traverseContext struct {
	client              *msgraph.GroupsClient
	ctx                 context.Context
	traverseFunc        TraverseFunc
	maxTraverseDepth    int
	objectVisitCheckMap map[string]bool
}

// maxTraverseDepth: 0 means traverse max depth
func TraverseNestedGroups(client *msgraph.GroupsClient, ctx context.Context, maxTraverseDepth int, groupId string, traverseFunc TraverseFunc) error {
	tCtx := &traverseContext{
		client:              client,
		ctx:                 ctx,
		traverseFunc:        traverseFunc,
		maxTraverseDepth:    maxTraverseDepth,
		objectVisitCheckMap: map[string]bool{},
	}
	return traverseNestedGroups0(tCtx, groupId, 1)
}

func traverseNestedGroups0(tCtx *traverseContext, groupId string, currentDepth int) error {
	tCtx.objectVisitCheckMap[groupId] = true

	resp, status, _, err := tCtx.client.BaseClient.Get(tCtx.ctx, msgraph.GetHttpRequestInput{
		ConsistencyFailureFunc: msgraph.RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: msgraph.Uri{
			Entity:      fmt.Sprintf("/groups/%s/members", groupId),
			Params:      url.Values{"$select": []string{"id"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		err = tCtx.traverseFunc(groupId, odata.TypeGroup, err)
		return fmt.Errorf("GroupsClient.BaseClient.Get(): statusCode: %d, %+v", status, err)
	}
	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = tCtx.traverseFunc(groupId, odata.TypeGroup, err)
		return fmt.Errorf("ioutil.ReadAll(): %v", err)
	}

	var data struct {
		Members []msgraph.DirectoryObject `json:"value"`
	}
	if err = json.Unmarshal(rawBody, &data); err != nil {
		err = tCtx.traverseFunc(groupId, odata.TypeGroup, err)
		return fmt.Errorf("json.Unmarshal(): %v", err)
	}

	err = tCtx.traverseFunc(groupId, odata.TypeGroup, err)
	if err != nil {
		return err
	}

	for _, v := range data.Members {
		if tCtx.objectVisitCheckMap[*v.ID] {
			continue
		}

		if *v.ODataType != odata.TypeGroup {
			tCtx.objectVisitCheckMap[*v.ID] = true
			err = tCtx.traverseFunc(*v.ID, *v.ODataType, err)
			if err != nil {
				return err
			}
			continue
		}
		if tCtx.maxTraverseDepth != 0 && tCtx.maxTraverseDepth < currentDepth+1 {
			continue
		}
		err = traverseNestedGroups0(tCtx, *v.ID, currentDepth+1)
		if err != nil {
			return err
		}
	}
	return nil
}
