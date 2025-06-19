package rss

import (
	"fmt"
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
				t.Errorf("ErrorCode.String() = \"%v\", want \"%v\"", got, test.expected)
			}
		})
	}
}

func TestFetchError_Error(t *testing.T) {
	tests := []struct {
		err      *FetchError
		expected string
	}{
		{&FetchError{ErrorCodeUnreachable, fmt.Errorf("failed to fetch data")}, "Unreachable URL"},
		{&FetchError{ErrorCodeHttpError, fmt.Errorf("HTTP error ocured, Status: 404")}, "HTTP error"},
		{&FetchError{ErrorCodeNoData, fmt.Errorf("server returned no data")}, "No data received"},
	}

	for _, test := range tests {
		t.Run(test.err.Code.String(), func(t *testing.T) {
			if got := test.err.Error(); got != test.expected {
				t.Errorf("FetchError.Error() = \"%v\", want \"%v\"", got, test.expected)
			}
		})
	}
}

func TestFetchError_Unwrap(t *testing.T) {
	tests := []struct {
		err      *FetchError
		expected string
	}{
		{&FetchError{ErrorCodeUnreachable, fmt.Errorf("failed to fetch data")}, "failed to fetch data"},
		{&FetchError{ErrorCodeHttpError, fmt.Errorf("HTTP error ocured, Status: 404")}, "HTTP error ocured, Status: 404"},
		{&FetchError{ErrorCodeNoData, fmt.Errorf("server returned no data")}, "server returned no data"},
	}

	for _, test := range tests {
		t.Run(test.err.Code.String(), func(t *testing.T) {
			if got := test.err.Unwrap().Error(); got != test.expected {
				t.Errorf("FetchError.Unwrap() = \"%v\", want \"%v\"", got, test.expected)
			}
		})
	}
}
