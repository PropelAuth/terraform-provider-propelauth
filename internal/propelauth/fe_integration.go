package propelauth

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetTestFeIntegrationInfo - Returns the FE integration info for the test environment.
func (c *PropelAuthClient) GetTestFeIntegrationInfo() (*TestFeIntegrationInfo, error) {
	res, err := c.get("fe_integration")
	if err != nil {
		return nil, err
	}

	feIntegration := FeIntegrationInfoResponse{}
	err = json.Unmarshal(res.BodyBytes, &feIntegration)
	if err != nil {
		return nil, err
	}

	return &feIntegration.Test, nil
}

type FeIntegrationUpdate struct {
	ApplicationUrl string
	LoginRedirectPath string
	LogoutRedirectPath string
	AdditionalFeLocations []AdditionalFeLocation
}

// UpdateTestFeIntegration - Updates the FE integration info for the test environment.
func (c *PropelAuthClient) UpdateTestFeIntegration(update FeIntegrationUpdate) (*TestFeIntegrationInfo, error) {
	request := feIntegrationUpdateRequest{
		LoginRedirectPath: update.LoginRedirectPath,
		LogoutRedirectPath: update.LogoutRedirectPath,
		TestEnvFeIntegrationApplicationUrl: &testEnvFeIntegrationApplicationUrl{
			ApplicationUrl: update.ApplicationUrl,
			Type: "SchemeAndDomain",
		},
		AdditionalFeLocations: AdditionalFeLocations{
			AdditionalFeLocations: update.AdditionalFeLocations,
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, err = c.put("fe_integration/test", body)
	if err != nil {
		return nil, err
	}

	return c.GetTestFeIntegrationInfo()
}

// UpdateLiveFeIntegration - Updates the FE integration info for a live staging or prod environment.
func (c *PropelAuthClient) UpdateLiveFeIntegration(environment string, update FeIntegrationUpdate) (*FeIntegrationInfoForEnv, error) {
	request := feIntegrationUpdateRequest{
		ApplicationHostnameWithScheme: update.ApplicationUrl,
		LoginRedirectPath: update.LoginRedirectPath,
		LogoutRedirectPath: update.LogoutRedirectPath,
		AdditionalFeLocations: AdditionalFeLocations{
			AdditionalFeLocations: update.AdditionalFeLocations,
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, err = c.put(fmt.Sprintf("fe_integration/%s", strings.ToLower(environment)), body)
	if err != nil {
		return nil, err
	}

	return c.GetLiveFeIntegrationInfo(environment)
}


// GetLiveFeIntegrationInfo - Returns the FE integration info for a live staging or prod environment.
func (c *PropelAuthClient) GetLiveFeIntegrationInfo(environment string) (*FeIntegrationInfoForEnv, error) {
	res, err := c.get("fe_integration")
	if err != nil {
		return nil, err
	}

	feIntegration := FeIntegrationInfoResponse{}
	err = json.Unmarshal(res.BodyBytes, &feIntegration)
	if err != nil {
		return nil, err
	}

	switch environment {
		case "Staging":
			return &feIntegration.Staging, nil
		case "Prod":
			return &feIntegration.Prod, nil
		default:
			return nil, fmt.Errorf("invalid environment when fetching FE integration info: %s", environment)
	}
}