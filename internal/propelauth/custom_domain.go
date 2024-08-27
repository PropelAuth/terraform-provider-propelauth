package propelauth

import (
	"encoding/json"
	"fmt"
)

// GetCustomDomainInfo - Returns the custom domain info for the requested environment
func (c *PropelAuthClient) GetCustomDomainInfo(environment string, isSwitching bool) (*CustomDomainInfo, error) {
	res, err := c.get(fmt.Sprintf("custom_domain?environment=%v&is_switching=%v", environment, isSwitching), nil)
	if err != nil {
		return nil, err
	}

	customDomainInfo := CustomDomainInfo{}
	err = json.Unmarshal(res.BodyBytes, &customDomainInfo)
	if err != nil {
		return nil, err
	}

	return &customDomainInfo, nil
}

// UpdateCustomDomainInfo - Updates the custom domain info for the requested environment
func (c *PropelAuthClient) UpdateCustomDomainInfo(environment string, domain string, subdomain *string, isSwitching bool) (*CustomDomainInfo, error) {
	request := customDomainUpdateRequest{
		Domain: domain,
		Subdomain: subdomain,
		Environment: environment,
		IsSwitching: isSwitching,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	res, err := c.put("custom_domain", body)
	if err != nil {
		return nil, err
	}

	customDomainInfo := CustomDomainInfo{}
	err = json.Unmarshal(res.BodyBytes, &customDomainInfo)
	if err != nil {
		return nil, err
	}

	return &customDomainInfo, nil
}

// VerifyCustomDomainInfo - Verifies the custom domain info for the requested environment
func (c *PropelAuthClient) VerifyCustomDomainInfo(environment string, isSwitching bool) error {
	request := customDomainVerifyRequest{
		Environment: environment,
		IsSwitching: isSwitching,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.post("custom_domain/verify", body)
	if err != nil {
		return err
	}
	return nil
}