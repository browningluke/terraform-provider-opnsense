package opnsense

import (
	"context"
	"fmt"
	"reflect"
)

// makeSetFunc creates a func that creates/updates the resource
func makeSetFunc(c controller, endpoint string, reconfigureEndpoint string) func(ctx context.Context, data any) (string, error) {
	return func(ctx context.Context, data any) (string, error) {
		// Since the OPNsense controller has to be reconfigured after every change, locking the mutex prevents
		// the API from being written to while it's reconfiguring, which results in data loss.
		c.Mutex().Lock()
		defer c.Mutex().Unlock()

		// Make request to OPNsense
		respJson := &addResp{}
		err := c.Client().doRequest(ctx, "POST", endpoint, data, respJson)
		if err != nil {
			return "", err
		}

		// Validate result
		if respJson.Result != "saved" {
			return "", fmt.Errorf("resource not changed. result: %s. errors: %s", respJson.Result, respJson.Validations)
		}

		// Reconfigure (i.e. restart) the OPNsense service
		err = c.Client().ReconfigureService(ctx, reconfigureEndpoint)
		if err != nil {
			return "", err
		}

		return respJson.UUID, nil
	}
}

// makeGetFunc creates a func that reads the resource
func makeGetFunc[K any](c *Client, endpoint string, data *K) func(ctx context.Context, id string) (*K, error) {
	return func(ctx context.Context, id string) (*K, error) {
		err := c.doRequest(ctx, "GET",
			fmt.Sprintf("%s/%s", endpoint, id), nil, data)

		// Handle errors
		if err != nil {
			// Handle unmarshal error (means ID is invalid, or was deleted upstream)
			if err.Error() == fmt.Sprintf("json: cannot unmarshal array into Go value of type %s",
				reflect.TypeOf(data).Elem().String()) {
				return nil, fmt.Errorf("unable to find resource. it may have been deleted upstream")
			}

			return nil, err
		}

		return data, nil
	}
}

// makeDeleteFunc creates a func that deletes the resource
func makeDeleteFunc(c controller, endpoint, reconfigureEndpoint string) func(ctx context.Context, id string) error {
	return func(ctx context.Context, id string) error {
		// Since the OPNsense controller has to be reconfigured after every change, locking the mutex prevents
		// the API from being written to while it's reconfiguring, which results in data loss.
		c.Mutex().Lock()
		defer c.Mutex().Unlock()

		respJson := &deleteResp{}
		err := c.Client().doRequest(ctx, "POST", fmt.Sprintf("%s/%s", endpoint, id), nil, respJson)
		if err != nil {
			return err
		}

		// Validate that override was deleted
		if respJson.Result != "deleted" {
			return fmt.Errorf("resource not deleted. result: %s", respJson.Result)
		}

		// Reconfigure (i.e. restart) the OPNsense service
		err = c.Client().ReconfigureService(ctx, reconfigureEndpoint)
		if err != nil {
			return err
		}

		return nil
	}
}
