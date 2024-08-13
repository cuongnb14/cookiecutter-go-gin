package vo

import (
	"github.com/google/uuid"
)

type UserVO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
