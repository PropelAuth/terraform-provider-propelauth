package propelauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
)

// UploadImage - Uploads an image to the project and returns the the new image_id.
func (c *PropelAuthClient) UploadImage(imageType string, pathToLocalImage string) (*ImageUploadResponse, error) {
	path := fmt.Sprintf("image/%s", imageType)
	url := c.assembleURL(path, nil)

	var requestBody bytes.Buffer
    w := multipart.NewWriter(&requestBody)

	// write image data to the request
	fw, err := w.CreateFormFile("file", pathToLocalImage)
	if err != nil {
		return nil, fmt.Errorf("error on creating form field: %w", err)
	}
	r, err := os.Open(pathToLocalImage)
	if err != nil {
		return nil, fmt.Errorf("error on opening image file: %w", err)
	}

	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, fmt.Errorf("error on writing image data to upload request: %w", err)
	}

	// close the writer or the image data will be incomplete
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("error on closing writer for image upload: %w", err)
	}
	
	// create request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("error on creating upload image request: %w", err)
	}

	// add headers
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer " + c.apiKey)
	req.Header.Set("User-Agent", "terraform-provider-propelauth/0.0 go/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH)

	// send request
	resp, err := c.httpClient.Do(req)
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
	imageUploadResponse := ImageUploadResponse{}
	err = json.Unmarshal(respBytes, &imageUploadResponse)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error on response: %s", string(respBytes[:]))
	}

	return &imageUploadResponse, nil
}
