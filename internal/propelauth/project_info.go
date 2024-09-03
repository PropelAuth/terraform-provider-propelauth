package propelauth

import (
	"encoding/json"
)

// GetProjectInfo - Returns a project metadata.
func (c *PropelAuthClient) GetProjectInfo() (*ProjectInfoResponse, error) {
	res, err := c.get("info")
	if err != nil {
		return nil, err
	}

	projectInfo := ProjectInfoResponse{}
	err = json.Unmarshal(res.BodyBytes, &projectInfo)
	if err != nil {
		return nil, err
	}

	return &projectInfo, nil
}

// UpdateProjectInfo - Updates the project's metadata -- principally the name.
func (c *PropelAuthClient) UpdateProjectInfo(name *string) (*ProjectInfoResponse, error) {
	projectInfo := ProjectInfoUpdateRequest{
		Name: *name,
	}

	body, err := json.Marshal(projectInfo)
	if err != nil {
		return nil, err
	}

	_, err = c.patch("info", body)
	if err != nil {
		return nil, err
	}

	return c.GetProjectInfo()
}
