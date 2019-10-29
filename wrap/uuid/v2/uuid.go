package v2

import (
	"fmt"

	"github.com/google/uuid"
)

// Represents of.UUIDGen
type UUID struct {
}

func (u *UUID) UUID() string {
	return fmt.Sprintf("%s", uuid.New())
}
