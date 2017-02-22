package data_structure

import "testing"

func TestMap01(t *testing.T) {
	s := NewLinkedMap()

	s.Put("1", "hi")

	if !(s.Size() == 1 && s.Get("1") == "hi") {
		t.Error("Expected result should match.")
	}
}

func TestMap02(t *testing.T) {
	s := NewLinkedMap()

	s.Put("1", "hi")
	s.Put("2", "hello")

	if !(s.Size() == 2 && s.Get("1") == "hi" && s.Get("2") == "hello") {
		t.Error("Expected result should match.")
	}
}

func TestMap03(t *testing.T) {
	s := NewLinkedMap()

	s.Put("1", "hi")
	s.Put("2", "hello")

	s.Remove("2")

	if !(s.Size() == 1 && s.Get("1") == "hi") {
		t.Error("Expected result should match.")
	}
}

func TestMap04(t *testing.T) {
	s := NewLinkedMap()

	s.Put("1", "hi")
	s.Put("2", "hello")

	s.Remove("1")

	if !(s.Size() == 1 && s.Get("2") == "hello") {
		t.Error("Expected result should match.")
	}
}

func TestMap05(t *testing.T) {
	s := NewLinkedMap()

	s.Put("1", "hi")
	s.Put("2", "hello")

	if !(s.Size() == 2 && s.Contains("1") && s.Contains("2")) {
		t.Error("Expected result should match.")
	}
}