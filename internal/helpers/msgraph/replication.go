package msgraph

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func WaitForCreationReplication(ctx context.Context, f func() (interface{}, int, error)) (interface{}, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context has no deadline")
	}
	timeout := time.Until(deadline)
	return (&resource.StateChangeConf{
		Pending:                   []string{"NotFound", "BadCast"},
		Target:                    []string{"Found"},
		Timeout:                   timeout,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 2,
		Refresh: func() (interface{}, string, error) {
			i, status, err := f()

			switch {
			case status >= 200 && status < 300:
				return i, "Found", nil
			case status == 404:
				return i, "NotFound", nil
			case i == nil:
				return nil, "BadCast", nil
			case err != nil:
				return i, "Error", fmt.Errorf("unable to retrieve object, received response with status %d: %v", status, err)
			}

			return i, "Error", fmt.Errorf("unrecognised response with status %d", status)
		},
	}).WaitForStateContext(ctx)
}
