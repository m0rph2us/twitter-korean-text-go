package util

import (
	"testing"
)

func TestDecomposeHangul(t *testing.T) {
	x := DecomposeHangul('각')

	if !(x.Onset == 'ㄱ' && x.Vowel == 'ㅏ' && x.Coda == 'ㄱ') {
		t.Error("Expected result should match.")
	}

	x = DecomposeHangul('관')

	if !(x.Onset == 'ㄱ' && x.Vowel == 'ㅘ' && x.Coda == 'ㄴ') {
		t.Error("Expected result should match.")
	}

	x = DecomposeHangul('꼃')

	if !(x.Onset == 'ㄲ' && x.Vowel == 'ㅕ' && x.Coda == 'ㅀ') {
		t.Error("Expected result should match.")
	}
}

func TestDecomposeHangulHasNoCoda(t *testing.T) {
	x := DecomposeHangul('가')

	if !(x.Onset == 'ㄱ' && x.Vowel == 'ㅏ' && x.Coda == ' ') {
		t.Error("Expected result should match.")
	}

	x = DecomposeHangul('과')

	if !(x.Onset == 'ㄱ' && x.Vowel == 'ㅘ' && x.Coda == ' ') {
		t.Error("Expected result should match.")
	}

	x = DecomposeHangul('껴')

	if !(x.Onset == 'ㄲ' && x.Vowel == 'ㅕ' && x.Coda == ' ') {
		t.Error("Expected result should match.")
	}
}

func TestComposeHangulFull(t *testing.T) {
	x := ComposeHangulFull('ㄱ', 'ㅏ', 'ㄷ')

	if !(x == '갇') {
		t.Error("Expected result should match.")
	}

	x = ComposeHangulFull('ㄲ', 'ㅑ', 'ㅀ')

	if !(x == '꺓') {
		t.Error("Expected result should match.")
	}

	x = ComposeHangulFull('ㅊ', 'ㅘ', 'ㄴ')

	if !(x == '촨') {
		t.Error("Expected result should match.")
	}
}

func TestComposeHangulFullHasNoCoda(t *testing.T) {
	x := ComposeHangulFull('ㄱ', 'ㅏ', ' ')

	if !(x == '가') {
		t.Error("Expected result should match.")
	}

	x = ComposeHangulFull('ㄲ', 'ㅑ', ' ')

	if !(x == '꺄') {
		t.Error("Expected result should match.")
	}

	x = ComposeHangulFull('ㅊ', 'ㅘ', ' ')

	if !(x == '촤') {
		t.Error("Expected result should match.")
	}
}

func TestComposeHangulFullException(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected result should match.")
			}
		}()

		ComposeHangulFull(' ', 'ㅏ', ' ')
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected result should match.")
			}
		}()

		ComposeHangulFull('ㄲ', ' ', ' ')
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected result should match.")
			}
		}()

		ComposeHangulFull(' ', ' ', 'ㄴ')
	}()
}

func TestHasCoda(t *testing.T) {
	if !hasCoda('갈') {
		t.Error("Expected result should match.")
	}

	if !hasCoda('갉') {
		t.Error("Expected result should match.")
	}
}

func TestHasNoCoda(t *testing.T) {
	if hasCoda('가') {
		t.Error("Expected result should match.")
	}

	if hasCoda('ㅘ') {
		t.Error("Expected result should match.")
	}

	if hasCoda('ㄱ') {
		t.Error("Expected result should match.")
	}

	if hasCoda(' ') {
		t.Error("Expected result should match.")
	}
}
