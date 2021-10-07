package hamilton_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"io"
	"net/http"
)

func GetAppRoleAssignment(c *msgraph.ServicePrincipalsClient, ctx context.Context, resourceId string, assignmentId string, query odata.Query) (*msgraph.AppRoleAssignment, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, msgraph.GetHttpRequestInput{
		ConsistencyFailureFunc: msgraph.RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: msgraph.Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/appRoleAssignedTo/%s", resourceId, assignmentId),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ServicePrincipalsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var assignment msgraph.AppRoleAssignment

	if err := json.Unmarshal(respBody, &assignment); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &assignment, status, nil
}
