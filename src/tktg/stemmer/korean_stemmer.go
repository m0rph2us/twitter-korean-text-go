package stemmer

import (
	"reflect"
	"tktg/tokenizer"
	"tktg/util"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var Endings data_structure.Set
var Predicates data_structure.Set
var EndingsForNouns data_structure.Set

func init() {
	Endings = data_structure.NewLinkedSet(util.Eomi, util.PreEomi)
	Predicates = data_structure.NewLinkedSet(util.Verb, util.Adjective)
	EndingsForNouns = data_structure.NewLinkedSet("하다", "되다", "없다")
}

func Stem(tokens []tokenizer.KoreanToken) []tokenizer.KoreanToken {
	exists := scala.Exists(tokens, func(value interface{}) bool {
		t := value.(tokenizer.KoreanToken)
		return t.Pos == util.Verb || t.Pos == util.Adjective
	})

	if !exists {
		return tokens
	}

	t := scala.FoldLeft(tokens, []tokenizer.KoreanToken{},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			tokens := folded.([]tokenizer.KoreanToken)
			token := value.(tokenizer.KoreanToken)

			l := len(tokens)

			if l > 0 && Endings.Contains(token.Pos) {
				if Predicates.Contains(tokens[0].Pos) {
					prevToken := tokens[0]
					d := []tokenizer.KoreanToken{}

					d = append(d, tokenizer.KoreanToken{
						prevToken.Text,
						prevToken.Pos,
						prevToken.Offset,
						prevToken.Length + token.Length,
						prevToken.Unknown,
					})

					if l > 1 {
						d = append(d, tokens[1:]...)
					}

					return d
				} else {
					return tokens
				}
			} else if Predicates.Contains(token.Pos) {
				text := util.PredicateStems[token.Pos][token.Text]

				d := []tokenizer.KoreanToken{}

				d = append(d, tokenizer.KoreanToken{
					text,
					token.Pos,
					token.Offset,
					token.Length,
					token.Unknown,
				})

				d = append(d, tokens...)

				return d
			} else {
				d := []tokenizer.KoreanToken{}

				d = append(d, token)
				d = append(d, tokens...)

				return d
			}

		})

	stemmed := tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanToken(scala.Reverse(t).([]interface{}))

	validNounHeading := func(token tokenizer.KoreanToken) bool {
		validLength := utf8.RuneCountInString(token.Text) > 2
		validPos := token.Pos == util.Verb
		validEndings := false
		if validLength {
			r := []rune(token.Text)
			validEndings = EndingsForNouns.Contains(string(r[len(r)-2:]))
		}
		validNouns := false
		if validLength {
			r := []rune(token.Text)
			heading := string(r[0 : len(r)-2])
			validNouns = util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains(heading)
		}

		return validLength && validPos && validEndings && validNouns
	}

	flatMap := func(value []tokenizer.KoreanToken, f func(value interface{}) interface{}) []tokenizer.KoreanToken {
		s := []tokenizer.KoreanToken{}
		for _, v := range value {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				s = append(s, x.([]tokenizer.KoreanToken)...)
			}
		}
		return s
	}

	return flatMap(stemmed, func(value interface{}) interface{} {
		token := value.(tokenizer.KoreanToken)
		if validNounHeading(token) {
			r := []rune(token.Text)
			heading := string(r[0 : len(r)-2])
			ending := string(r[len(r)-2:])

			hl := utf8.RuneCountInString(heading)

			return []tokenizer.KoreanToken{
				{heading, util.Noun, token.Offset, hl, false},
				{ending, token.Pos, token.Offset + hl, token.Length - hl, false},
			}
		} else {
			return []tokenizer.KoreanToken{token}
		}
	})
}
