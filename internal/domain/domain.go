package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	Name  string
	Email string
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      Name      `json:"name" validate:"required, min=3, max=40"`
	Email     Email     `json:"email" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func New(name Name, email Email) (User, error) {
	now := time.Now()

	user := User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
	}

	// validate

	return user, nil
}
