package stemmer

import (
	"fmt"
	"testing"
	"tktg/tokenizer"
	"tktg/util"
)

func TestStemmer(t *testing.T) {
	v := Stem(tokenizer.Tokenize("새로운 스테밍을 추가했었다."))

	r := []tokenizer.KoreanToken{
		{"새롭다", util.Adjective, 0, 3, false},
		{" ", util.Space, 3, 1, false},
		{"스테밍", util.ProperNoun, 4, 3, false},
		{"을", util.Josa, 7, 1, false},
		{" ", util.Space, 8, 1, false},
		{"추가", util.Noun, 9, 2, false},
		{"하다", util.Verb, 11, 3, false},
		{".", util.Punctuation, 14, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", r)) {
		t.Error("Expected result should match.")
	}

	v = Stem(tokenizer.Tokenize("그런 사람 없습니다.."))

	r = []tokenizer.KoreanToken{
		{"그렇다", util.Adjective, 0, 2, false},
		{" ", util.Space, 2, 1, false},
		{"사람", util.Noun, 3, 2, false},
		{" ", util.Space, 5, 1, false},
		{"없다", util.Adjective, 6, 4, false},
		{"..", util.Punctuation, 10, 2, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", r)) {
		t.Error("Expected result should match.")
	}

	v = Stem(tokenizer.Tokenize("라고만"))

	r = []tokenizer.KoreanToken{
		{"라고만", util.Eomi, 0, 3, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", r)) {
		t.Error("Expected result should match.")
	}

	v = Stem(tokenizer.Tokenize("하...나는 아이유가 좋아요."))

	r = []tokenizer.KoreanToken{
		{"하", util.Exclamation, 0, 1, false},
		{"...", util.Punctuation, 1, 3, false},
		{"나", util.Noun, 4, 1, false},
		{"는", util.Josa, 5, 1, false},
		{" ", util.Space, 6, 1, false},
		{"아이유", util.ProperNoun, 7, 3, false},
		{"가", util.Josa, 10, 1, false},
		{" ", util.Space, 11, 1, false},
		{"좋다", util.Adjective, 12, 3, false},
		{".", util.Punctuation, 15, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", r)) {
		t.Error("Expected result should match.")
	}
}
