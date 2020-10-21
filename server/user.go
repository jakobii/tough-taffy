package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	var err error

	// get ID from Query
	var id uuid.UUID
	if x, ok := r.URL.Query()["id"]; ok {
		id, err = uuid.Parse(x[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid user id")
			return
		}
	}

	// get user
	if u, ok := DB.GetUser(id); ok {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, u.JSON())
		return
	}

	// if no user
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "user id not found")
	return
}
