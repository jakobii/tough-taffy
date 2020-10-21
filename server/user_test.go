package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jakobii/rather/db"

	"github.com/jakobii/rather/user"
)

const wantgot string = "want: %v, got: %v"

// test the http method for general funtionality
func TestHandleGetUserAsync(t *testing.T) {

	memDB := db.NewMemDB()
	DB = memDB

	// start the in memory database
	stopDB := memDB.Start()
	defer stopDB()

	var wg sync.WaitGroup

	TestHandleGetUser := func(t *testing.T) {
		u1 := user.New("jacob")

		DB.PutUser(u1.ID, u1)
		//e.g.: /GetUser?id=<uuid>
		request := httptest.NewRequest("GET", "/GetUser?id="+u1.ID.String(), nil)
		recorder := httptest.NewRecorder()

		// the function being tested.
		handleGetUser(recorder, request)

		response := recorder.Result()
		body, err := ioutil.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		wantStatusCode := http.StatusOK
		if response.StatusCode != wantStatusCode {
			t.Errorf(wantgot, wantStatusCode, response.StatusCode)
		}

		wantContentType := "application/json"
		if contentType := response.Header.Get("Content-Type"); contentType != wantContentType {
			t.Errorf(wantgot, wantContentType, contentType)
		}

		wantBody := u1.JSON()
		if gotBody := string(body); gotBody != u1.JSON() {
			t.Errorf(wantgot, wantBody, gotBody)
		}
	}

	i := 0
	for i < 10 {
		go func(fn func(*testing.T)) {
			wg.Add(1)
			defer wg.Done()
			fn(t)
		}(TestHandleGetUser)
		i++
	}

	wg.Wait()

}
