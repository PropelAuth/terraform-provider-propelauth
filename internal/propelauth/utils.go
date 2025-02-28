package propelauth

import (
	"fmt"
	"net/url"
	"strconv"
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
