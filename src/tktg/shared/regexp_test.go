package shared

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"testing"
)

func TestReplaceAll00(t *testing.T) {
	re := pcre.MustCompile(`((ab)?[a-z]+)`, pcre.MULTILINE|pcre.UTF8)
	c := ReplaceAll(re, "abcd efg hijk", func(m *pcre.Matcher) string {
		return "<" + string(m.Group(0)) + ">"
	})

	if !(c == "<abcd> <efg> <hijk>") {
		t.Error("Expected result should match.")
	}
}

func TestReplaceAll01(t *testing.T) {
	re := pcre.MustCompile(`((ab)?[a-z]+)`, pcre.MULTILINE|pcre.UTF8)
	c := ReplaceAll(re, "가나다라마바사", func(m *pcre.Matcher) string {
		return "<" + string(m.Group(0)) + ">"
	})

	if !(c == "가나다라마바사") {
		t.Error("Expected result should match.")
	}
}
