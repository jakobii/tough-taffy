package db

import (
	"github.com/google/uuid"
	"github.com/jakobii/tough-taffy/user"
)

// Driver defines all methods that can be performed on the database.
// table methods should be seperated into there own interfaces, and
// then included into the main driver interface, e.g. UserDriver.
// child driver interfaces should abstract network and database errors,
// and simplify them to look and feel like working with an in-memory
// object.
//
// All drivers must be safe for async. and should be thoughly tested to
// prove they are safe.
type Driver interface {
	UserDriver
}

// UserDriver defines the methods for CRUD on users in a database
type UserDriver interface {

	// GetUserByName is the main means of discovering a user
	// without knowing its user id.
	GetUserByName(string) (user.User, bool)

	// returns true when the user was found and returned.
	GetUser(uuid.UUID) (user.User, bool)

	// creates and updates a user. replacings the resource. similar
	// to the http PUT method rules. return true when the user was,
	// and created/updated.
	PutUser(uuid.UUID, user.User) bool

	// completely removes user from the database. this can not be undone.
	// returns true if the user was found and deleted.
	DelUser(uuid.UUID) bool
}
