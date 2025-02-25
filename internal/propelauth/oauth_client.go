package propelauth

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetOauthClientInfo - Returns the oauth client info for the requested environment.
func (c *PropelAuthClient) GetOauthClientInfo(environment string, oauthClientId string) (*OauthClientInfo, error) {
	res, err := c.get(
		fmt.Sprintf("%v/oauth_client/%v", strings.ToLower(environment), oauthClientId),
	)
	if err != nil {
		return nil, err
	}

	oauthClient := OauthClientInfo{}
	err = json.Unmarshal(res.BodyBytes, &oauthClient)
	if err != nil {
		return nil, err
	}

	return &oauthClient, nil
}

// CreateOauthClient - Creates a new oauth client and returns the client id and secret.
func (c *PropelAuthClient) CreateOauthClient(environment string, redirectUris []string) (*OauthClientCreationResponse, error) {
	request := OauthClientRequest{
		RedirectUris: redirectUris,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	res, err := c.post(
		fmt.Sprintf("%v/oauth_client", strings.ToLower(environment)),
		body,
	)
	if err != nil {
		return nil, err
	}

	oauthClient := OauthClientCreationResponse{}
	err = json.Unmarshal(res.BodyBytes, &oauthClient)
	if err != nil {
		return nil, err
	}

	return &oauthClient, nil
}

// UpdateOauthClient - Updates an existing oauth client and returns the result.
func (c *PropelAuthClient) UpdateOauthClient(environment string, oauthClientId string, redirectUris []string) error {
	request := OauthClientRequest{
		RedirectUris: redirectUris,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	res, err := c.patch(
		fmt.Sprintf("%v/oauth_client/%v", strings.ToLower(environment), oauthClientId),
		body,
	)

	if err != nil {
		return err
	}

	oauthClient := OauthClientInfo{}
	err = json.Unmarshal(res.BodyBytes, &oauthClient)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOauthClient - Deletes an existing oauth client.
func (c *PropelAuthClient) DeleteOauthClient(environment string, oauthClientId string) error {
	_, err := c.delete(
		fmt.Sprintf("%v/oauth_client/%v", strings.ToLower(environment), oauthClientId),
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
