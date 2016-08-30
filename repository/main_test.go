package repository

import (
	"fmt"
	"testing"
)

func TestRepositoryOpen(t *testing.T) {
	db := Open(":memory:")
	if fmt.Sprintf("%T", db) != "*repository.Buntdb" {
		t.Errorf("Expect *repository.Buntdb but got %T\n", db)
	}
}

func TestRepositoryFind(t *testing.T) {
	db := Open(":memory:")
	err := db.CreateOrUpdate("test", "Hello World")
	handleErr(err, t)

	value, err := db.Find("test")
	handleErr(err, t)

	if value != "Hello World" {
		t.Errorf("Expect Hello World but got %s\n", value)
	}
}

func TestRepositoryCreateOrUpdate(t *testing.T) {
	db := Open(":memory:")
	err := db.CreateOrUpdate("test", "Hello World")
	handleErr(err, t)
}

func TestRepositoryIsAvailable(t *testing.T) {
	db := Open(":memory:")
	err := db.CreateOrUpdate("test", "Hello World")
	handleErr(err, t)

	if db.IsAvailable("test") != false {
		t.Fatalf("Expect test to be unavailable, got available")
	}

	if db.IsAvailable("test1") != true {
		t.Fatalf("Expect test1 to be available, got unavailable")
	}
}

func handleErr(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("An error ocurred: %s", err)
	}
}
