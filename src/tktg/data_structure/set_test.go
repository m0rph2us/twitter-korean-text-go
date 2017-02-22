package data_structure

import (
	"testing"
)

func TestSet01(t *testing.T) {
	s := NewLinkedSet()

	s.Add("1")

	if !(s.Size() == 1 && s.ToSlice()[0] == "1") {
		t.Error("Expected result should match.")
	}
}

func TestSet02(t *testing.T) {
	s := NewLinkedSet("1", "2")

	if !(s.Size() == 2 && s.ToSlice()[0] == "1" && s.ToSlice()[1] == "2") {
		t.Error("Expected result should match.")
	}
}

func TestSet03(t *testing.T) {
	s := NewLinkedSet("1", "2")

	s.Remove("2")

	if !(s.Size() == 1 && s.ToSlice()[0] == "1") {
		t.Error("Expected result should match.")
	}
}

func TestSet04(t *testing.T) {
	s := NewLinkedSet("1", "2")

	s.Remove("1")

	if !(s.Size() == 1 && s.ToSlice()[0] == "2") {
		t.Error("Expected result should match.")
	}
}

func TestSet05(t *testing.T) {
	s := NewLinkedSet("1", "2")

	if !(s.Size() == 2 && s.Contains("1") && s.Contains("2")) {
		t.Error("Expected result should match.")
	}
}
