package propelauth

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetBeIntegrationInfo - Returns the BE integration info for the requested environment
func (c *PropelAuthClient) GetBeIntegrationInfo(environment string) (*BeIntegrationInfo, error) {
	res, err := c.get("be_integration", nil)
	if err != nil {
		return nil, err
	}

	beIntegration := BeIntegrationInfoResponse{}
	err = json.Unmarshal(res.BodyBytes, &beIntegration)
	if err != nil {
		return nil, err
	}

	switch environment {
		case "Test":
			return &beIntegration.Test, nil
		case "Staging":
			return &beIntegration.Staging, nil
		case "Prod":
			return &beIntegration.Prod, nil
		default:
			return nil, fmt.Errorf("invalid environment: %s", environment)
	}
}

// CreateBeApiKey - Creates a new BE API key and returns the result
func (c *PropelAuthClient) CreateBeApiKey(environment string, name string, isReadOnly bool) (*BeApiKey, error) {
	request := BeApiKeyCreateRequest{
		Name: name,
		IsReadOnly: isReadOnly,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	res, err := c.post(
		fmt.Sprintf("%v/be_integration/api_key", strings.ToLower(environment)),
		body,
	)
	if err != nil {
		return nil, err
	}

	beApiKey := BeApiKey{}
	err = json.Unmarshal(res.BodyBytes, &beApiKey)
	if err != nil {
		return nil, err
	}

	return &beApiKey, nil
}