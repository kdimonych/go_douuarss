package common

import (
	"fmt"
	"net/url"
)

func ValidateURL(rawurl string) error {
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return fmt.Errorf("invalid URL: %s, error: %w", rawurl, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https: %s", rawurl)
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host: %s", rawurl)
	}
	return nil
}
