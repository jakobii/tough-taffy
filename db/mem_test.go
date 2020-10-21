package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jakobii/tough-taffy/user"
)

// TestMemDriver prevents compilation, if MemDB does not impletement all the
// methods of the Driver interface.
func TestMemDriver(t *testing.T) {
	func(d Driver) {}(NewMemDB())
}

func TestMemDBUserCRUD(t *testing.T) {

	// instanciat a memdb
	db := NewMemDB()

	stopDB := db.Start()
	defer stopDB()

	// create user
	id := uuid.New()
	u := user.New("bob")
	u.ID = id

	// check get returns nil
	u2, ok := db.GetUser(id)
	if ok {
		t.Errorf("want: %v, got: %v", new(user.User), u2)
	}

	// create user
	ok = db.PutUser(id, u)
	if !ok {
		t.Errorf("want: %v, got: %v", true, ok)
	}

	// get newly created user
	u2, ok = db.GetUser(id)
	if !ok && u2.ID != u.ID {
		t.Errorf("want: %v, got: %v", u, u2)
	}

	// get newly created user
	u2, ok = db.GetUserByName(u.Name)
	if !ok && u2.Name != u.Name {
		t.Errorf("want: %v, got: %v", u, u2)
	}

	// update user and check
	u.Alias = "Robert"
	db.PutUser(id, u)
	u2, ok = db.GetUser(id)
	if !ok && u.Alias != u2.Alias {
		t.Errorf("want: %v, got: %v", u, u2)
	}

	// delete
	ok = db.DelUser(id)
	if !ok {
		t.Errorf("want: %v, got: %v", true, ok)
	}

	// check that deletion actually worked
	u2, ok = db.GetUser(id)
	if ok {
		t.Errorf("want: %v, got: %v", new(user.User), u2)
	}

}
