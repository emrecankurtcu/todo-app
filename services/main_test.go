package main

import (
	"testing"
)

func TestAddTodo(t *testing.T) {
	resp := addToDoFunc("Test")
	if !resp {
		t.Error("Error while adding")
	}
}
