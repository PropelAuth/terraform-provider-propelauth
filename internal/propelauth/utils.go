package propelauth

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func Contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func FlipBoolRef(b *bool) *bool {
	if b == nil {
		return nil
	} else if *b {
		new_b := false
		return &new_b
	} else {
		new_b := true
		return &new_b
	}
}

func IsValidUrlWithoutTrailingSlash(inputUrl string) (bool, error) {
	// Parse the inputUrl string into a URL
	parsedURL, err := url.Parse(inputUrl)
	if err != nil {
		return false, err
	}

	// Check if the URL has a scheme (http, https, etc.)
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false, fmt.Errorf("invalid URL: missing scheme or host")
	}

	// Check if the URL has a trailing slash in the path
	if strings.HasSuffix(parsedURL.Path, "/") {
		return false, fmt.Errorf("URL has a trailing slash")
	}

	return true, nil
}

func IsValidUrl(inputUrl string) (bool, error) {
	// Parse the inputUrl string into a URL
	parsedURL, err := url.Parse(inputUrl)
	if err != nil {
		return false, err
	}

	// Check if the URL has a scheme (http, https, etc.)
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false, fmt.Errorf("invalid URL: missing scheme or host")
	}

	return true, nil
}

func GetPortFromLocalhost(inputUrl string) (bool, int) {
	// Parse the inputUrl string into a URL
	parsedURL, err := url.Parse(inputUrl)
	if err != nil {
		return false, 0
	}

	fmt.Println(parsedURL.Hostname())

	// Check if the URL is localhost
	if parsedURL.Hostname() == "localhost" && parsedURL.Scheme == "http" {
		fmt.Println(parsedURL.Port())
		port, err := strconv.Atoi(parsedURL.Port())
		if err == nil {
			return true, port
		} else {
			return false, 0
		}
	}

	return false, 0
}
