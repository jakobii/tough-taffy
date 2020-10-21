package db

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jakobii/tough-taffy/user"
)

type mapdb = *map[string][]interface{}

// MemDB is an in-memory driver
// used for testing only.
type MemDB struct {
	Requests chan MemReq
	DB       mapdb
	ctx      context.Context
	stop     context.CancelFunc
}

// NewMemDB creates a new MemDB
func NewMemDB() MemDB {
	db := make(map[string][]interface{})
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return MemDB{
		Requests: make(chan MemReq, 100),
		DB:       &db,
		ctx:      ctx,
		stop:     cancel,
	}
}

// Run the in-memory sever
func (m *MemDB) Run() {
	//log.Println("Starting server")
	for {
		select {
		case <-m.ctx.Done():
			//log.Println("stoping server")
			return

		case request := <-m.Requests:
			result := request.Query(m.DB)
			request.Response <- result

		}
	}
}

// Start starts the MemDB in its own process and returns a handler to stop it.
// the handler should be defered to stop the process.
//   example:
//      stopDB := db.Start()
//	    defer stopDB()
func (m *MemDB) Start() func() {
	var wg sync.WaitGroup

	//run the db sever in its own process to simluate network async.
	go func(db *MemDB) {
		wg.Add(1)
		defer wg.Done()
		db.Run()
	}(m)

	// stop server
	return func() {
		m.Stop()
		wg.Wait()
	}
}

// Stop the in-memory sever
func (m *MemDB) Stop() {
	m.stop()
}

// MemQuery a function for querying the MemDB
type MemQuery = func(mapdb) interface{}

// MemReq is An In-memory request
type MemReq struct {
	Query    MemQuery
	Response chan interface{}
}

// NewMemReq creates a new MemReq
func NewMemReq(q MemQuery) MemReq {
	return MemReq{
		Query:    q,
		Response: make(chan interface{}),
	}
}

// GetUserByName fulfills the UserDriver interface
func (m MemDB) GetUserByName(name string) (user.User, bool) {
	query := func(db mapdb) interface{} {
		if tb, ok := (*db)["Users"]; ok {
			for _, row := range tb {
				if u, ok := row.(user.User); ok {
					if u.Name == name {
						return u
					}
				}
			}
		}
		return nil
	}
	request := NewMemReq(query)
	m.Requests <- request
	response := <-request.Response
	if u, ok := response.(user.User); ok {
		return u, true
	}
	return user.User{}, false
}

// GetUser fulfills the UserDriver interface
func (m MemDB) GetUser(id uuid.UUID) (user.User, bool) {
	query := func(db mapdb) interface{} {
		if tb, ok := (*db)["Users"]; ok {
			for _, row := range tb {
				if u, ok := row.(user.User); ok {
					if u.ID == id {
						return u
					}
				}
			}
		}
		return nil
	}
	request := NewMemReq(query)
	m.Requests <- request
	response := <-request.Response
	if u, ok := response.(user.User); ok {
		return u, true
	}
	return user.User{}, false
}

// PutUser fulfills the UserDriver interface
func (m MemDB) PutUser(id uuid.UUID, u user.User) bool {
	query := func(db mapdb) interface{} {
		var index int
		var exists bool
		if _, ok := (*db)["Users"]; !ok {
			(*db)["Users"] = make([]interface{}, 0, 20)
		}
		if tb, ok := (*db)["Users"]; ok {
			for k, row := range tb {
				if u, ok := row.(user.User); ok {
					if u.ID == id {
						exists = true
						index = k
						break
					}
				}
			}
			if exists {
				(*db)["Users"][index] = u
				return true
			}
			(*db)["Users"] = append((*db)["Users"], u)
			return true
		}
		return false
	}
	request := NewMemReq(query)
	m.Requests <- request
	response := <-request.Response
	if r, ok := response.(bool); ok && r {
		return true
	}
	return false
}

// DelUser fulfills the UserDriver interface
func (m MemDB) DelUser(id uuid.UUID) bool {
	var index int
	var exists bool
	query := func(db mapdb) interface{} {
		if tb, ok := (*db)["Users"]; ok {
			for i, row := range tb {
				if u, ok := row.(user.User); ok {
					if u.ID == id {
						exists = true
						index = i
					}
				}
			}
		}
		if exists {
			(*db)["Users"][index] = (*db)["Users"][len((*db)["Users"])-1]
			(*db)["Users"][len((*db)["Users"])-1] = nil
			(*db)["Users"] = (*db)["Users"][:len((*db)["Users"])-1]
			return true
		}
		return false
	}
	request := NewMemReq(query)
	m.Requests <- request
	response := <-request.Response
	if r, ok := response.(bool); ok && r {
		return true
	}
	return false
}
