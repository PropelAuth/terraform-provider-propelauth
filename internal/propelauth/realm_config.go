package propelauth

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetEnvironmentConfig - Get the realm's login/signup configuration.
func (c *PropelAuthClient) GetRealmConfig(environment string) (*RealmConfigResponse, error) {
	res, err := c.get("realm")
	if err != nil {
		return nil, err
	}

	realmConfigs := RealmConfigsResponse{}
	err = json.Unmarshal(res.BodyBytes, &realmConfigs)
	if err != nil {
		return nil, err
	}

	switch environment {
	case "Test":
		return &realmConfigs.Test, nil
	case "Staging":
		return realmConfigs.Staging, nil
	case "Prod":
		return realmConfigs.Prod, nil
	default:
		return nil, fmt.Errorf("invalid environment when fetching realm config: %s", environment)
	}
}

// UpdateRealmConfig - Updates the realms login/signup configuration ignoring any null values.
func (c *PropelAuthClient) UpdateRealmConfig(environment string, realmConfig RealmConfigUpdate) (*RealmConfigResponse, error) {
	body, err := json.Marshal(realmConfig)
	if err != nil {
		return nil, err
	}

	_, err = c.patch(fmt.Sprintf("realm/%s", strings.ToLower(environment)), body)
	if err != nil {
		return nil, err
	}

	return c.GetRealmConfig(environment)
}
