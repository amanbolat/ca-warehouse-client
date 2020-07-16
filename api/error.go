package api

import (
	"fmt"
)

type Error struct {
	Message       string `json:"message"`
	Hint          string `json:"hint"`
	InternalError error  `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s. InternalErr: %v", e.Message, e.InternalError)
}

func NewError(err error, message string, hint string) error {
	apiErr, ok := err.(Error)
	if ok {
		return apiErr
	}

	return Error{
		Message:       message,
		Hint:          hint,
		InternalError: err,
	}
}
