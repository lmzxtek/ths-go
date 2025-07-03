package db

import (
	"testing"
)

func TestOpen(t *testing.T) {
	db, err := Connect("test.db")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
}
