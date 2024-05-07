package client

import (
	"fmt"

	uuid "github.com/google/uuid"
)

type UserNotFoundError struct {
	client_id uuid.UUID
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user with id %s not found", e.client_id.String())
}
