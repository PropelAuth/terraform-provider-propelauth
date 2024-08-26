package propelauth

import (
	"encoding/json"
	"fmt"
)

// GetCustomDomainInfo - Returns the custom domain info for the requested environment
func (c *PropelAuthClient) GetCustomDomainInfo(environment string) (*CustomDomainInfo, error) {
	res, err := c.get(fmt.Sprintf("custom_domain?environment=%v", environment), nil)
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
func (c *PropelAuthClient) UpdateCustomDomainInfo(environment string, domain string, subdomain *string) (*CustomDomainInfo, error) {
	request := customDomainUpdateRequest{
		Domain: domain,
		Subdomain: subdomain,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	res, err := c.put(fmt.Sprintf("%v/custom_domain", environment), body)
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