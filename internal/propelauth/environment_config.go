package propelauth

import (
	"encoding/json"
)

// GetEnvironmentConfig - Returns a project metadata
func (c *PropelAuthClient) GetEnvironmentConfig() (*EnvironmentConfigResponse, error) {
	res, err := c.get("config", nil)
	if err != nil {
		return nil, err
	}

	environmentConfig := EnvironmentConfigResponse{}
	err = json.Unmarshal(res.BodyBytes, &environmentConfig)
	if err != nil {
		return nil, err
	}

	return &environmentConfig, nil
}

// UpdateEnvironmentConfig - Updates the environment configuration ignoring the null values
func (c *PropelAuthClient) UpdateEnvironmentConfig(environmentConfig *EnvironmentConfigUpdate) (*EnvironmentConfigResponse, error) {
	body, err := json.Marshal(environmentConfig)
	if err != nil {
		return nil, err
	}

	_, err = c.patch("config", body)
	if err != nil {
		return nil, err
	}

	return c.GetEnvironmentConfig()
}
