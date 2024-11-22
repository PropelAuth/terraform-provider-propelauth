package propelauth

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetAllSocialLoginInfo - Returns all the social login info for environments and sso providers.
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

// GetSocialLoginInfo - Returns the social login redirect info for the requested environment + sso provider.
func (c *PropelAuthClient) GetSocialLoginInfo(sso_provider string) (*SocialLoginInfo, error) {
	allSocialLoginInfo, err := c.GetAllSocialLoginInfo()
	if err != nil {
		return nil, err
	}

	switch sso_provider {
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
		return nil, fmt.Errorf("invalid social login sso_provider: %s", sso_provider)
	}
}

// GetSocialLoginRedirectUrl - Returns the authorized redirect for the requested environment and sso provider.
func (c *PropelAuthClient) GetSocialLoginRedirectUrl(environment string, sso_provider string) (*string, error) {
	socialLoginInfo, err := c.GetSocialLoginInfo(sso_provider)
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

// UpsertSocialLoginInfo - Upserts the social login info for the requested social sso provider.
func (c *PropelAuthClient) UpsertSocialLoginInfo(sso_provider string, clientId string, clientSecret string) error {
	request := SocialLoginUpdateRequest{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Enabled:      true,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.put(fmt.Sprintf("social/%s", strings.ToLower(sso_provider)), body)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSocialLogin - Deletes the social login info for the requested social sso provider and disables the integration.
func (c *PropelAuthClient) DeleteSocialLogin(sso_provider string) error {
	request := SocialLoginUpdateRequest{
		ClientId:     "DELETED",
		ClientSecret: "DELETED",
		Enabled:      false,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.put(fmt.Sprintf("social/%s", strings.ToLower(sso_provider)), body)
	if err != nil {
		return err
	}

	return nil
}
