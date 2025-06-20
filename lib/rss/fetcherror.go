package rss

import (
	"fmt"
)

type ErrorCode int

const (
	ErrorCodeUnreachable = iota
	ErrorCodeHttpError
	ErrorCodeNoData
)

func (e ErrorCode) String() string {
	switch e {
	case ErrorCodeUnreachable:
		return "Unreachable URL"
	case ErrorCodeHttpError:
		return "HTTP error"
	case ErrorCodeNoData:
		return "No data received"
	default:
		return "Unknown error"
	}
}

type FetchError struct {
	Code    ErrorCode // Error code for the error
	Details error     // Additional details about the error
}

func (e FetchError) Error() string {
	return fmt.Sprintf("%s (%v)", e.Code, e.Details)
}
