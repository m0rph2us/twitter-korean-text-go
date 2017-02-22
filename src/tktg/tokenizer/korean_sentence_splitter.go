package tokenizer

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"unicode/utf8"
)

type Sentence struct {
	Text  string
	Start int
	End   int
}

func (t Sentence) toString() string {
	return t.Text + "(" + string(t.Start) + "," + string(t.End) + ")"
}

var re pcre.Regexp

func init() {
	re = pcre.MustCompile(
		`[^.!?…\s][^.!?…]*(?:[.!?…](?!['\"]?\s|$)[^.!?…]*)*[.!?…]?['\"]?(?=\s|$)`, pcre.MULTILINE|pcre.UTF8)
}

func Split(text string) []Sentence {
	ret := []Sentence{}

	start := 0
	for {
		offset := re.FindIndex([]byte(text[start:]), 0)

		if offset == nil {
			break
		}

		s := start + offset[0]
		e := start + offset[1]

		ret = append(ret,
			Sentence{
				string(text[s:e]),
				utf8.RuneCountInString(text[0:s]),
				utf8.RuneCountInString(text[0:e]),
			},
		)

		start = e
	}

	return ret
}
