package propelauth

import (
	"encoding/json"
	"fmt"
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