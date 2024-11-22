package propelauth

import (
	"encoding/json"
	"fmt"
)

// GetAllSocialLoginInfo - Returns all the social login info for environments and providers.
func (c *PropelAuthClient) GetAllSocialLoginInfo() (*AllSocialLoginInfoResponse, error) {
	res, err := c.get("social")
	if err != nil {
		return nil, err
	}

	allSocialLoginInfo := AllSocialLoginInfoResponse{}
	err = json.Unmarshal(res.BodyBytes, &allSocialLoginInfo)
	if err != nil {
		return nil, err
	}

	return &allSocialLoginInfo, nil
}

// GetSocialLoginInfo - Returns the social login redirect info for the requested environment + provider.
func (c *PropelAuthClient) GetSocialLoginInfo(provider string) (*SocialLoginInfo, error) {
	allSocialLoginInfo, err := c.GetAllSocialLoginInfo()
	if err != nil {
		return nil, err
	}

	switch provider {
	case "Google":
		return &allSocialLoginInfo.Google, nil
	case "Microsoft":
		return &allSocialLoginInfo.Microsoft, nil
	case "GitHub":
		return &allSocialLoginInfo.GitHub, nil
	case "Slack":
		return &allSocialLoginInfo.Slack, nil
	case "LinkedIn":
		return &allSocialLoginInfo.LinkedIn, nil
	case "Atlassian":
		return &allSocialLoginInfo.Atlassian, nil
	case "Apple":
		return &allSocialLoginInfo.Apple, nil
	case "Salesforce":
		return &allSocialLoginInfo.Salesforce, nil
	case "QuickBooks":
		return &allSocialLoginInfo.QuickBooks, nil
	case "Xero":
		return &allSocialLoginInfo.Xero, nil
	case "Salesloft":
		return &allSocialLoginInfo.Salesloft, nil
	case "Outreach":
		return &allSocialLoginInfo.Outreach, nil
	default:
		return nil, fmt.Errorf("invalid social login provider: %s", provider)
	}
}

// GetSocialLoginRedirectUrl - Returns the authorized redirect for the requested environment and provider.
func (c *PropelAuthClient) GetSocialLoginRedirectUrl(environment string, provider string) (*string, error) {
	socialLoginInfo, err := c.GetSocialLoginInfo(provider)
	if err != nil {
		return nil, err
	}

	switch environment {
	case "Test":
		return &socialLoginInfo.TestRedirectUrl, nil
	case "Staging":
		return &socialLoginInfo.StagingRedirectUrl, nil
	case "Prod":
		return &socialLoginInfo.ProdRedirectUrl, nil
	default:
		return nil, fmt.Errorf("invalid environment when fetching social login redirect URL: %s", environment)
	}
}
