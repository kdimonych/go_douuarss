package storage

type ErrorCode int

const (
	ErrorUnableToInsertFeed = iota
	ErrorUnableToInsertChannel
	ErrorUnableToInsertItem
)

func (e ErrorCode) String() string {
	switch e {
	case ErrorUnableToInsertFeed:
		return "Unable to insert feed"
	case ErrorUnableToInsertChannel:
		return "Unable to insert channel"
	case ErrorUnableToInsertItem:
		return "Unable to insert item"
	default:
		return "Unknown error"
	}
}

type StorageError struct {
	Code        ErrorCode // Error code for the error
	Description string    // Description of the error
	Details     error     // Additional details about the error
}

func (e *StorageError) Error() string {
	if e.Description == "" {
		return e.Code.String()
	}
	return e.Code.String() + ": " + e.Description
}

func (e *StorageError) Unwrap() error {
	return e.Details
}
