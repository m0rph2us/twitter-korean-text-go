package data_structure

import (
	"testing"
)

func TestLinkedList_Add(t *testing.T) {
	l := NewLinkedList()

	l.Add("1")

	if !(l.ToSlice()[0] == "1") {
		t.Error("The likedlist should have an item.")
	}
}

func TestLinkedList_Remove00(t *testing.T) {
	l := NewLinkedList()

	l.Add("1")

	l.Remove("1")

	if !(l.Size() == 0) {
		t.Error("The linkedlist should have no item.")
	}
}

func TestLinkedList_Remove01(t *testing.T) {
	l := NewLinkedList()

	l.Add("1")
	l.Add("2")
	l.Add("3")

	l.Remove("2")

	if !(l.ToSlice()[0] == "1" && l.ToSlice()[1] == "3") {
		t.Error("The linkedlist should have two items.")
	}
}
