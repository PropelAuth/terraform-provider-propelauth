package propelauth

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// TODO: FOR THE LOVE OF GOD, DON'T FORGET TO CHANGE THIS TO THE REAL URL
const BaseURLTemplate string = "https://api.propelauth.localhost/iac/%s/project/%s"

// PropelAuthClient - Client for the PropelAuth API to manage an existing project and all its resources
type PropelAuthClient struct {
	BaseURL    string
	HTTPClient *http.Client
	ApiKey      string
}

type StandardResponse struct {
	StatusCode   int
	ResponseText string
	BodyBytes    []byte
	BodyText     string
}

func NewClient(tenant_id, project_id, api_key *string) (*PropelAuthClient, error) {
	c := PropelAuthClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		BaseURL: fmt.Sprintf(BaseURLTemplate, *tenant_id, *project_id),
		ApiKey: *api_key,
	}

	return &c, nil
}

// public http methods

func (c *PropelAuthClient) get(urlPostfix string, queryParams url.Values) (*StandardResponse, error) {
	url := c.assembleURL(urlPostfix, queryParams)

	return c.requestHelper("GET", url, nil)
}

func (c *PropelAuthClient) patch(urlPostfix string, body []byte) (*StandardResponse, error) {
	url := c.assembleURL(urlPostfix, nil)

	return c.requestHelper("PATCH", url, body)
}

func (c *PropelAuthClient) post(urlPostfix string, body []byte) (*StandardResponse, error) {
	url := c.assembleURL(urlPostfix, nil)

	return c.requestHelper("POST", url, body)
}

func (c *PropelAuthClient) put(urlPostfix string, body []byte) (*StandardResponse, error) {
	url := c.assembleURL(urlPostfix, nil)

	return c.requestHelper("PUT", url, body)
}

func (c *PropelAuthClient) delete(urlPostfix string, body []byte) (*StandardResponse, error) {
	url := c.assembleURL(urlPostfix, nil)

	return c.requestHelper("DELETE", url, body)
}


func (c *PropelAuthClient) requestHelper(method string, url string, body []byte) (*StandardResponse, error) {
	requestBody := bytes.NewBuffer(body)

	// create request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error on creating request: %w", err)
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + c.ApiKey)
	req.Header.Set("User-Agent", "terraform-provider-propelauth/0.0 go/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH)

	// send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	// convert the response body to a stream of bytes
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading response body: %w", err)
	}

	respBytes := buf.Bytes()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error on response: %s", string(respBytes[:]))
	}

	// return the response
	queryResponse := StandardResponse{
		StatusCode:   resp.StatusCode,
		ResponseText: resp.Status,
		BodyBytes:    respBytes,
		BodyText:     string(respBytes[:]),
	}

	return &queryResponse, nil
}

func (c *PropelAuthClient) assembleURL(urlPostfix string, queryParams url.Values) string {
	url := c.BaseURL + "/" + urlPostfix
	if queryParams != nil {
		url += "?" + queryParams.Encode()
	}

	return url
}