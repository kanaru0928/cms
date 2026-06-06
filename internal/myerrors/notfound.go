package myerrors

import "fmt"

type NotFoundError struct {
	ID string
	ResourceType string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.ResourceType, e.ID)
}
