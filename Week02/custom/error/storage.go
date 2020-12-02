package error

import (
	"database/sql"
)

var (
	// DefaultError is an error for non-platform-specific errors.
	DefaultError = &customErrorCode{"DefaultError", "Default custom error."}
	// EmptyResultError is used when the storage layer returns a platform specific error which means empty result, such as sql.ErrNoRows.
	EmptyResultError = &customErrorCode{"EmptyResultError", "Cannot find any results by given conditions."}
	// StorageConnectionError is used when the storage layer returns a platform specific error which stands for connection error, such as sql.ErrConnDone.
	StorageConnectionError = &customErrorCode{"StorageConnectionError", "Failed to connect storages."}
)

// IsPlatformSpecific checks an error is a platform specific error or not.
// We should protect the higher layer from depending on a platform specific error.
func IsPlatformSpecific(err error) bool {
	switch err {
	case sql.ErrConnDone, sql.ErrNoRows:
		return true
	default:
		return false
	}
}

// StorageErrToCustomErr transforms the error to custom error. When the given error is not platform specific, it is transformed into a default error.
func StorageErrToCustomErr(err error) *CustomError {
	switch err {
	case sql.ErrConnDone:
		return &CustomError{StorageConnectionError, err}
	case sql.ErrNoRows:
		return &CustomError{EmptyResultError, err}
	default:
		return &CustomError{DefaultError, err}
	}
}
