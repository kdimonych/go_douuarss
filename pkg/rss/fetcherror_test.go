package rss

import (
	"testing"
)

func TestErrorCode_String(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected string
	}{
		{ErrorCodeUnreachable, "Unreachable URL"},
		{ErrorCodeHttpError, "HTTP error"},
		{ErrorCodeNoData, "No data received"},
		{999, "Unknown error"}, // Test for an unknown error code
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.code.String(); got != test.expected {
				t.Errorf("ErrorCode.String() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestFetchError_Error(t *testing.T) {
	tests := []struct {
		err      *FetchError
		expected string
	}{
		{&FetchError{ErrorCodeUnreachable, "Failed to fetch data"}, "Unreachable URL (Failed to fetch data)"},
		{&FetchError{ErrorCodeHttpError, "HTTP error ocured, Status: 404"}, "HTTP error (HTTP error ocured, Status: 404)"},
		{&FetchError{ErrorCodeNoData, "Server returned no data"}, "No data received (Server returned no data)"},
	}

	for _, test := range tests {
		t.Run(test.err.Code.String(), func(t *testing.T) {
			if got := test.err.Error(); got != test.expected {
				t.Errorf("FetchError.Error() = %v, want %v", got, test.expected)
			}
		})
	}
}
