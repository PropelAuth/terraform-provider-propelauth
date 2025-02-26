package propelauth

import (
	"encoding/json"
)

// GetApiKeyAlert - Returns the configuration of api key alerts if any.
func (c *PropelAuthClient) GetApiKeyAlert() (*ApiKeyAlert, error) {
	res, err := c.get("end_user_api_key_alerts")
	if err != nil {
		return nil, err
	}

	apiKeyAlert := ApiKeyAlert{}
	err = json.Unmarshal(res.BodyBytes, &apiKeyAlert)
	if err != nil {
		return nil, err
	}

	return &apiKeyAlert, nil
}

// UpdateApiKeyAlert - Enables API key alerting and set the advanced_notice_days.
func (c *PropelAuthClient) UpdateApiKeyAlert(advancedNoticeDays int32) error {
	updateReq := ApiKeyAlert{
		Enabled:           true,
		AdvanceNoticeDays: advancedNoticeDays,
	}

	body, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	_, err = c.put("end_user_api_key_alerts", body)
	if err != nil {
		return err
	}

	return nil
}

// DeleteApiKeyAlert - Disables API key alerting.
func (c *PropelAuthClient) DeleteApiKeyAlert() error {
	deleteReq := ApiKeyAlert{
		Enabled:           false,
		AdvanceNoticeDays: 1, // Need to provide a dummy value since the go default of 0 doesn't pass validation.
	}

	body, err := json.Marshal(deleteReq)
	if err != nil {
		return err
	}

	_, err = c.put("end_user_api_key_alerts", body)
	if err != nil {
		return err
	}

	return nil
}
