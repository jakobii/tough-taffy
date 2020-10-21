package user

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// ErrorNotFound for when a user can not be found
var ErrorNotFound = errors.New("user: not found")

// ErrorName is a general error regarding the username
var ErrorName = errors.New("user: name must be filled and unique")

// User is a user.
type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Alias string    `json:"alias"`
}

// New creates a new user
func New(name string) User {
	return User{
		ID:   uuid.New(),
		Name: name,
	}
}

// JSON marshalls the user to JSON
func (u *User) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}
