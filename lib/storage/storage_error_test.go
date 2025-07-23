package storage

import (
	"fmt"
	"testing"
)

func TestErrorCode_String(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected string
	}{
		{ErrorUnableToInsertFeed, "Unable to insert feed"},
		{ErrorUnableToInsertChannel, "Unable to insert channel"},
		{ErrorUnableToInsertItem, "Unable to insert item"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.code.String(); got != test.expected {
				t.Errorf("ErrorCode.String() = \"%v\", want \"%v\"", got, test.expected)
			}
		})
	}
}

func TestStorageError_Error(t *testing.T) {
	tests := []struct {
		err      *StorageError
		expected string
	}{
		{&StorageError{ErrorUnableToInsertFeed, "1", fmt.Errorf("failed to insert feed")}, "Unable to insert feed: 1"},
		{&StorageError{ErrorUnableToInsertChannel, "2", fmt.Errorf("failed to insert channel")}, "Unable to insert channel: 2"},
		{&StorageError{ErrorUnableToInsertItem, "3", fmt.Errorf("failed to insert item")}, "Unable to insert item: 3"},
	}

	for _, test := range tests {
		t.Run(test.err.Code.String(), func(t *testing.T) {
			if got := test.err.Error(); got != test.expected {
				t.Errorf("StorageError.Error() = \"%v\", want \"%v\"", got, test.expected)
			}
		})
	}
}

func TestStorageError_Unwrap(t *testing.T) {
	tests := []struct {
		err      *StorageError
		expected string
	}{
		{&StorageError{ErrorUnableToInsertFeed, "1", fmt.Errorf("failed to insert feed")}, "failed to insert feed"},
		{&StorageError{ErrorUnableToInsertChannel, "2", fmt.Errorf("failed to insert channel")}, "failed to insert channel"},
		{&StorageError{ErrorUnableToInsertItem, "3", fmt.Errorf("failed to insert item")}, "failed to insert item"},
	}

	for _, test := range tests {
		t.Run(test.err.Code.String(), func(t *testing.T) {
			if got := test.err.Unwrap().Error(); got != test.expected {
				t.Errorf("StorageError.Unwrap() = \"%v\", want \"%v\"", got, test.expected)
			}
		})
	}
}
