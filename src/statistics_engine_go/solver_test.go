package main

import (
	"testing"
)

func TestCheckCheck(t *testing.T) {
	result := CheckCheck()
	want := "test"
	if result != want {
		t.Errorf("incorrect return")
	}
}
