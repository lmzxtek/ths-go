package math

import "testing"

func TestAdd(t *testing.T) {
	if Add(2, 3) != 5 { // Test case 1
		t.Error("Add(2, 3) should be 5")
	}
	if Add(10, -5) != 5 { // Test case 2
		t.Error("Add(10, -5) should be 5")
	}
}

func TestAdd2(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Error("Add(2, 3) should be 5")
	}
	if Add(10, -5) != 5 {
		t.Error("Add(10, -5) should be 5")
	}
}
