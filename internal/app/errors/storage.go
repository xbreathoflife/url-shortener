package errors

import (
	"fmt"
)

type ULRDuplicateError struct{
	BaseURL string
	ShortURL string
}

func (e *ULRDuplicateError) Error() string {
	return fmt.Sprintf("URL %s was already shortened to %s", e.BaseURL, e.ShortURL)
}

func NewULRDuplicateError(baseURL string, shortURL string) *ULRDuplicateError {
	return &ULRDuplicateError{
		BaseURL: baseURL,
		ShortURL: shortURL,
	}
}

type EmptyStorageError struct{
	User string
}

func (e *EmptyStorageError) Error() string {
	return fmt.Sprintf("No URLs for user %s were found", e.User)
}

func NewEmptyStorageError(user string) *EmptyStorageError {
	return &EmptyStorageError{
		User: user,
	}
}