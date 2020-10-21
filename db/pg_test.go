package db

import "testing"

func createTestPG() PG {
	return *new(PG)
}

// requires postgres server to be running
func TestPGConn(t *testing.T) {

}
