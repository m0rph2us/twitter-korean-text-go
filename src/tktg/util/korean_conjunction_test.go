package util

import (
	"testing"
	"tktg/data_structure"
)

func TestAddPreEomi(t *testing.T) {
	r := []rune{
		'x', 'y',
	}

	x := addPreEomi('1', r)

	if !(x[0] == "1x" && x[1] == "1y") {
		t.Error("Expected result should match.")
	}
}

func TestStringSetToRuneSet(t *testing.T) {

	expended := data_structure.NewLinkedSet()
	expended.Add("가")
	expended.Add("나")
	expended.Add("다")

	s := data_structure.NewLinkedSet()

	for v := range expended.Iter() {
		n := data_structure.NewLinkedSet()
		for _, v1 := range []rune(v.(string)) {
			n.Add(v1)
		}
		s.Add(n)
	}
}

func TestGetInit(t *testing.T) {
	s := "string"
	v := getInit(s)

	if v != "strin" {
		t.Error("Expected result should match.")
	}
}

func TestGetLast(t *testing.T) {
	s := "string"
	v := getLast(s)

	if v != 'g' {
		t.Error("Expected result should match.")
	}
}

func TestConjugatePredicatesToSet(t *testing.T) {
	d := data_structure.NewLinkedSet("귀엽")
	x := conjugatePredicatesToSet(d, true)

	e := data_structure.NewLinkedSet("귀여워", "귀여웠", "귀여운", "귀여", "귀엽")

	if x.Size() != e.Size() {
		t.Error("Expected result should match.")
	}

	for v := range x.Iter() {
		if !e.Contains(v) {
			t.Error("Expected result should match.")
		}
	}
}
