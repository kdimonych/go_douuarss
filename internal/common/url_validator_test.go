package common

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		rawurl   string
		expected string
	}{
		{"http://example.com", ""},
		{"https://example.com", ""},
		{"ftp://example.com", "URL scheme must be http or https: ftp://example.com"},
		{"http://", "URL must have a host: http://"},
		{"invalid-url", "invalid URL: invalid-url, error: parse \"invalid-url\": invalid URI for request"},
	}

	for _, test := range tests {
		t.Run(test.rawurl, func(t *testing.T) {
			err := ValidateURL(test.rawurl)
			if err != nil && err.Error() != test.expected {
				t.Errorf("ValidateURL(%q) = %v, want %v", test.rawurl, err.Error(), test.expected)
			} else if err == nil && test.expected != "" {
				t.Errorf("ValidateURL(%q) = nil, want %v", test.rawurl, test.expected)
			}
		})
	}
}
